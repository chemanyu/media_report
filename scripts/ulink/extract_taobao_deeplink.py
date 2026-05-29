import sys
import os
import time
import json
from urllib.parse import parse_qs, urlparse, urlencode, quote

# Go 通过 pipe 调用本脚本时，Windows Python 默认 stdout 退到 GBK，
# 中文/emoji print 会 UnicodeEncodeError 直接崩溃，Go 拿不到 JSON。
try:
    sys.stdout.reconfigure(encoding='utf-8', errors='replace')
    sys.stderr.reconfigure(encoding='utf-8', errors='replace')
except Exception:
    pass

from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.common.exceptions import TimeoutException

# 回退路径用：uc 自带 chromedriver 管理，正常不会走到这里
CHROME_DRIVER_PATH = "D:\\148\\chromedriver-win64\\chromedriver.exe"

URL_PREFIXES = ('tbopen://', 'taobao://', 'tmall://')

MOBILE_UA = ("Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) "
             "AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Mobile Safari/537.36")

# 设置该环境变量后走 attach 模式：连到一个手动启动并登录过淘宝的 Chrome：
#   chrome.exe --remote-debugging-port=9222 --user-data-dir=D:\148\chrome-debug-profile
# 之后每次跑只在该 Chrome 里新开 tab，跑完关 tab、不退浏览器，登录态永远在
DEBUG_ADDR_ENV = 'TAOBAO_DEEPLINK_DEBUGGER_ADDR'

HOOK_SCRIPT = r"""
(function () {
    window.__capturedTaobao = [];
    var prefixes = ['tbopen://', 'taobao://', 'tmall://'];
    var record = function (url, source) {
        try {
            if (typeof url !== 'string') return;
            for (var i = 0; i < prefixes.length; i++) {
                if (url.indexOf(prefixes[i]) === 0) {
                    window.__capturedTaobao.push({ url: url, source: source });
                    return;
                }
            }
        } catch (e) {}
    };
    var isCustomScheme = function (v) {
        if (typeof v !== 'string') return false;
        for (var i = 0; i < prefixes.length; i++) if (v.indexOf(prefixes[i]) === 0) return true;
        return false;
    };
    try {
        var hrefDesc = Object.getOwnPropertyDescriptor(Location.prototype, 'href');
        if (hrefDesc && hrefDesc.set) {
            var origSet = hrefDesc.set;
            Object.defineProperty(Location.prototype, 'href', {
                configurable: true, enumerable: true, get: hrefDesc.get,
                set: function (v) {
                    record(v, 'location.href');
                    if (isCustomScheme(v)) return;
                    return origSet.call(this, v);
                }
            });
        }
    } catch (e) {}
    try {
        var winLocDesc = Object.getOwnPropertyDescriptor(window, 'location')
            || Object.getOwnPropertyDescriptor(Window.prototype, 'location');
        if (winLocDesc && winLocDesc.set) {
            var origWinLocSet = winLocDesc.set;
            Object.defineProperty(window, 'location', {
                configurable: true, enumerable: true, get: winLocDesc.get,
                set: function (v) {
                    record(typeof v === 'string' ? v : (v && v.href), 'window.location=');
                    if (typeof v === 'string' && isCustomScheme(v)) return;
                    return origWinLocSet.call(this, v);
                }
            });
        }
    } catch (e) {}
    try {
        var origAssign = Location.prototype.assign;
        Location.prototype.assign = function (v) {
            record(v, 'location.assign');
            if (isCustomScheme(v)) return;
            return origAssign.apply(this, arguments);
        };
        var origReplace = Location.prototype.replace;
        Location.prototype.replace = function (v) {
            record(v, 'location.replace');
            if (isCustomScheme(v)) return;
            return origReplace.apply(this, arguments);
        };
    } catch (e) {}
    try {
        var origOpen = window.open;
        window.open = function (v) {
            record(v, 'window.open');
            if (isCustomScheme(v)) return null;
            return origOpen.apply(this, arguments);
        };
    } catch (e) {}
    try {
        var origClick = HTMLAnchorElement.prototype.click;
        HTMLAnchorElement.prototype.click = function () {
            record(this.href, 'a.click');
            if (isCustomScheme(this.href)) return;
            return origClick.apply(this, arguments);
        };
    } catch (e) {}
})();
"""


def _build_chrome_options():
    opts = Options()

    # attach 模式：UA / 屏幕 / profile 全部由现有 Chrome 决定，mobile 在 attach 后用 CDP 动态加
    debug_addr = os.environ.get(DEBUG_ADDR_ENV)
    if debug_addr:
        opts.add_experimental_option("debuggerAddress", debug_addr)
        opts.set_capability("goog:loggingPrefs", {"performance": "ALL"})
        opts.set_capability("goog:perfLoggingPrefs", {"enableNetwork": True, "enablePage": True})
        return opts

    # 与浏览器 DevTools 手机模式一致：Android UA + 移动屏幕，
    # Android 下淘宝直接下发 tbopen://，iOS 走 universal link 走不通这条路径
    opts.add_experimental_option("mobileEmulation", {
        "deviceMetrics": {"width": 375, "height": 812, "pixelRatio": 3.0},
        "userAgent": MOBILE_UA,
    })
    opts.add_experimental_option("excludeSwitches", ["enable-automation"])
    opts.add_experimental_option("useAutomationExtension", False)
    opts.add_argument("--disable-gpu")
    opts.add_argument("--no-sandbox")
    opts.add_argument("--disable-dev-shm-usage")
    opts.add_argument("--disable-blink-features=AutomationControlled")
    opts.add_argument('log-level=3')

    # 持久化 profile：复用 cookie，绕开淘宝风控对全新 session 的拦截
    profile_dir = os.environ.get('TAOBAO_DEEPLINK_PROFILE') or os.path.join(
        os.path.dirname(os.path.abspath(__file__)), '.chrome-profile'
    )
    try:
        os.makedirs(profile_dir, exist_ok=True)
    except Exception:
        pass
    opts.add_argument(f"--user-data-dir={profile_dir}")

    # perf log：tbopen:// 顶层导航在某些场景只在 Page 域出事件
    opts.set_capability("goog:loggingPrefs", {"performance": "ALL"})
    opts.set_capability("goog:perfLoggingPrefs", {"enableNetwork": True, "enablePage": True})
    return opts


def _create_driver():
    """优先 undetected-chromedriver（绕风控），失败回退原生 selenium。"""
    chrome_options = _build_chrome_options()

    # attach 模式：接管已运行的 Chrome，uc 不适用（也不需要绕检测，本来就是真人会话）
    if os.environ.get(DEBUG_ADDR_ENV):
        try:
            return webdriver.Chrome(options=chrome_options)
        except Exception:
            return webdriver.Chrome(service=Service(CHROME_DRIVER_PATH), options=chrome_options)

    use_uc = os.environ.get('TAOBAO_DEEPLINK_USE_UC', '1') != '0'
    if use_uc:
        try:
            import undetected_chromedriver as uc
            uc_options = uc.ChromeOptions()
            for arg in chrome_options.arguments:
                uc_options.add_argument(arg)
            for k, v in chrome_options.experimental_options.items():
                try:
                    uc_options.add_experimental_option(k, v)
                except Exception:
                    pass
            return uc.Chrome(options=uc_options, use_subprocess=True)
        except Exception as e:
            print(f"uc 初始化失败，回退原生 selenium: {e}")
    try:
        return webdriver.Chrome(service=Service(CHROME_DRIVER_PATH), options=chrome_options)
    except Exception:
        return webdriver.Chrome(options=chrome_options)


def _scan_perf_logs(driver):
    """从 Chrome performance log 里捞 tbopen://，返回首个命中或 None。"""
    try:
        logs = driver.get_log('performance')
    except Exception:
        return None
    for entry in logs:
        try:
            msg = json.loads(entry.get('message', '{}')).get('message', {})
            method = msg.get('method', '')
            params = msg.get('params', {}) or {}
            candidates = []
            if method == 'Network.requestWillBeSent':
                candidates.append(params.get('request', {}).get('url', ''))
                candidates.append(params.get('documentURL', ''))
            elif method in ('Page.frameRequestedNavigation',
                            'Page.frameScheduledNavigation',
                            'Page.navigatedWithinDocument',
                            'Page.windowOpen'):
                candidates.append(params.get('url', ''))
            for url in candidates:
                if url and url.startswith(URL_PREFIXES):
                    return url
        except Exception:
            continue
    return None


def _scan_hook(driver):
    """从注入的 JS hook 里捞 tbopen://，返回首个命中或 None。"""
    try:
        captured = driver.execute_script("return window.__capturedTaobao || [];") or []
    except Exception:
        return None
    for item in captured:
        url = item.get('url') if isinstance(item, dict) else None
        if url and url.startswith(URL_PREFIXES):
            return url
    return None


def _replace_h5url(deeplink, short_url):
    """把 deeplink 里的 h5Url 替换成原始短链。"""
    try:
        parsed = urlparse(deeplink)
        params = parse_qs(parsed.query)
        params['h5Url'] = [short_url]
        new_query = urlencode(params, doseq=True)
        return f"{parsed.scheme}://{parsed.netloc}{parsed.path}?{new_query}"
    except Exception as e:
        print(f"替换 h5Url 时出错: {e}，使用原始 Deeplink")
        return deeplink


def get_taobao_deeplink(short_url, driver=None, platform="ios"):
    """从淘宝短链提取 tbopen:// deeplink。"""
    internal_driver = driver is None
    attach_mode = bool(os.environ.get(DEBUG_ADDR_ENV))
    if internal_driver:
        driver = _create_driver()
        if driver is None:
            return None, None

    # attach 模式：在已登录的 Chrome 中新开 tab，跑完关 tab，保留登录态给下次
    original_handle = None
    new_tab_opened = False
    if internal_driver and attach_mode:
        try:
            original_handle = driver.current_window_handle
            driver.switch_to.new_window('tab')
            new_tab_opened = True
        except Exception as e:
            print(f"新开 tab 失败: {e}")

    try:
        # attach 模式下 options.mobileEmulation 无效，用 CDP 在当前 target 上动态设
        if attach_mode:
            try:
                driver.execute_cdp_cmd('Emulation.setDeviceMetricsOverride', {
                    'width': 375, 'height': 812, 'deviceScaleFactor': 3.0, 'mobile': True
                })
                driver.execute_cdp_cmd('Network.setUserAgentOverride', {'userAgent': MOBILE_UA})
            except Exception as e:
                print(f"mobile emulation 设置失败: {e}")

        try:
            driver.execute_cdp_cmd('Network.enable', {})
            driver.execute_cdp_cmd('Page.enable', {})
            driver.execute_cdp_cmd('Page.addScriptToEvaluateOnNewDocument', {'source': HOOK_SCRIPT})
        except Exception as e:
            print(f"注入 hook 失败: {e}")

        driver.get(short_url)
        print(f"plat: {platform}")

        try:
            WebDriverWait(driver, 5).until(
                lambda d: d.execute_script('return document.readyState') == 'complete'
            )
        except TimeoutException:
            pass

        # 轮询 hook + perf log，最长 ~30s（活动页加载慢时 tbopen 在 5-10s 才发）
        deeplink = None
        for _ in range(60):
            deeplink = _scan_hook(driver) or _scan_perf_logs(driver)
            if deeplink:
                break
            time.sleep(0.5)

        # 兜底：从静态 <a> 标签找
        if not deeplink:
            try:
                els = driver.find_elements(
                    By.XPATH,
                    "//a[starts-with(@href,'taobao://') or starts-with(@href,'tbopen://') or starts-with(@href,'tmall://')]"
                )
                if els:
                    deeplink = els[0].get_attribute("href")
            except Exception:
                pass

        if deeplink:
            deeplink = _replace_h5url(deeplink, short_url)
            print(f"提取到 Deeplink: {deeplink}")
            return deeplink, process_deeplink(deeplink, platform)

        print(f"未能提取 Deeplink，落地 URL: {driver.current_url}")
    except Exception as e:
        print(f"提取 deeplink 过程中发生错误: {e}")
    finally:
        if internal_driver and driver is not None:
            if attach_mode:
                # 关 tab，不退浏览器：保留登录态 + 风控指纹给下次跑用
                try:
                    if new_tab_opened:
                        driver.close()
                    if original_handle and original_handle in driver.window_handles:
                        driver.switch_to.window(original_handle)
                except Exception:
                    pass
            else:
                try:
                    driver.quit()
                except Exception:
                    pass
    return None, None


def process_deeplink(deeplink, platform):
    """iOS 走 ace.tb.cn 通用链接二次包装，Android 直接返回 tbopen://。"""
    if platform.lower() == "ios":
        return f"https://ace.tb.cn/t?smburl={quote(deeplink, safe='')}"
    return deeplink
