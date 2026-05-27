import sys
# Go 通过 pipe 调用本脚本时，Windows Python 默认 stdout 退到 GBK，
# 任何中文/emoji print 会 UnicodeEncodeError 导致脚本崩溃，Go 拿不到 JSON。
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
import time
from urllib.parse import parse_qs, urlparse, quote

# 配置 ChromeDriver 路径 - 如果您的路径不同，请替换为您的 ChromeDriver 路径
# 对于 Linux，常见路径是 /usr/bin/chromedriver 或 /usr/local/bin/chromedriver
# 或者确保 chromedriver 在您的系统 PATH 环境变量中
#CHROME_DRIVER_PATH = "/opt/homebrew/bin/chromedriver" # <-- 请确保为 Linux 更新此路径
CHROME_DRIVER_PATH = "D:\\148\\chromedriver-win64\\chromedriver.exe" # <-- Windows 路径示例


# 添加一个参数 platform，表示选择的系统（安卓或 iOS）
def get_taobao_deeplink(short_url, driver=None, platform="ios"):
    """
    尝试从给定的淘宝短链接中提取 deeplink。
    模拟移动浏览器环境。
    如果提供了 driver 参数，则使用现有浏览器实例，否则创建新的。
    """
    internal_driver = False # 标记是否是内部创建的 driver
    if driver is None:
        internal_driver = True
        chrome_options = Options()
        mobile_emulation = {
            "deviceMetrics": {"width": 375, "height": 812, "pixelRatio": 3.0},
            "userAgent": "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1"
        }
        chrome_options.add_experimental_option("mobileEmulation", mobile_emulation)
        chrome_options.add_argument("--headless")
        chrome_options.add_argument("--disable-gpu")
        chrome_options.add_argument("--no-sandbox")
        chrome_options.add_argument("--disable-dev-shm-usage")
        chrome_options.add_argument('log-level=3')

        try:
            service = Service(CHROME_DRIVER_PATH)
            driver = webdriver.Chrome(service=service, options=chrome_options)  # 用 seleniumwire 的 webdriver
        except Exception as e:
            print(f"初始化 ChromeDriver 时出错 (路径: '{CHROME_DRIVER_PATH}'): {e}")
            print("尝试从系统 PATH 初始化 ChromeDriver...")
            try:
                driver = webdriver.Chrome(options=chrome_options)
            except Exception as e_path:
                print(f"从 PATH 初始化 ChromeDriver 时出错: {e_path}")
                print("请确保 Chrome/Chromium 和正确的 ChromeDriver 已安装并配置。")
                return None

    deeplink = None
    try:
        # 注入 hook：劫持 location.href / assign / replace / window.open / a.click
        # 淘宝/天猫的拉起方式是 location.href = "tbopen://..."，没有 <a> 标签，
        # 必须在 JS 层捕获并阻止真实跳转（避免页面被打断）。
        try:
            driver.execute_cdp_cmd('Page.enable', {})
            hook_script = r"""
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
            driver.execute_cdp_cmd('Page.addScriptToEvaluateOnNewDocument', {'source': hook_script})
        except Exception as e:
            print(f"注入 hook 失败: {e}")

        driver.get(short_url)
        current_url = driver.current_url
        print(f"plat: {platform}")

        time.sleep(1)
        try:
            WebDriverWait(driver, 5).until(
                lambda d: d.execute_script('return document.readyState') == 'complete'
            )
        except TimeoutException:
            print("页面加载等待超时，继续后续处理")

        # 策略 1：从 hook 捕获（最可靠 —— 包含 location.href / a.click / window.open 等）
        for _ in range(10):
            try:
                captured = driver.execute_script("return window.__capturedTaobao || [];") or []
                for item in captured:
                    url = item.get('url') if isinstance(item, dict) else None
                    if url and (url.startswith('tbopen://') or url.startswith('taobao://') or url.startswith('tmall://')):
                        deeplink = url
                        print(f"从 hook 捕获 ({item.get('source')}): {deeplink}")
                        break
            except Exception as e:
                print(f"读取 hook 出错: {e}")
            if deeplink:
                break
            time.sleep(0.5)

        # 策略 2：从 <a> 标签查找（旧逻辑，作为兜底）
        if not deeplink:
            try:
                deeplink_elements = driver.find_elements(
                    By.XPATH,
                    "//a[starts-with(@href, 'taobao://') or starts-with(@href, 'tbopen://') or starts-with(@href, 'tmall://')]"
                )
                if deeplink_elements:
                    deeplink = deeplink_elements[0].get_attribute("href")
                    print(f"从 <a> 标签找到: {deeplink}")
            except Exception as e:
                print(f"<a> 标签扫描出错: {e}")

        if deeplink:
            # 把 h5Url 替换成原始 short_url（保留旧行为）
            try:
                parsed = urlparse(deeplink)
                params = parse_qs(parsed.query)
                params['h5Url'] = [short_url]
                from urllib.parse import urlencode
                new_query = urlencode(params, doseq=True)
                deeplink = f"{parsed.scheme}://{parsed.netloc}{parsed.path}?{new_query}"
                print(f"替换 h5Url 后的 Deeplink: {deeplink}")
            except Exception as e:
                print(f"替换 h5Url 时出错: {e}，使用原始 Deeplink")
            return deeplink, process_deeplink(deeplink, platform)

        print(f"未能提取 Deeplink，落地 URL: {current_url}")

    except Exception as e:
        print(f"提取deeplink过程中发生错误: {e}")
    finally:
        if internal_driver and driver is not None:
            print("关闭内部创建的浏览器实例。")
            driver.quit()

    return None, None

def process_deeplink(deeplink, platform):
    """
    根据平台处理 Deeplink。
    如果提供了 short_url，则替换 deeplink 中的 h5Url 参数为 short_url 的 URL 编码值。
    如果是 iOS 平台，进行 URL 编码并拼接。
    如果是安卓平台，直接返回原始 Deeplink。
    """
    # print(f"处理 Deeplink: {deeplink}，平台: {platform}")

    # 如果提供了 short_url，替换 deeplink 中的 h5Url 参数

    if platform.lower() == "ios":
        # 对提取到的 Deeplink 进行 URL 编码并拼接
        encoded_deeplink = quote(deeplink, safe='')
        final_url = f"https://ace.tb.cn/t?smburl={encoded_deeplink}"
        print(f"最终拼接的 URL: {final_url}")
        return final_url
    else:
        print(f"安卓平台，返回处理后的 Deeplink: {deeplink}")
        return deeplink
