import sys
import os
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
        # 与你浏览器 DevTools 手机模式一致：Android UA + 移动屏幕，
        # 淘宝对 iOS / Android 的拉端策略不同，Android 直接下发 tbopen://
        mobile_emulation = {
            "deviceMetrics": {"width": 375, "height": 812, "pixelRatio": 3.0},
            "userAgent": "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Mobile Safari/537.36"
        }
        chrome_options.add_experimental_option("mobileEmulation", mobile_emulation)
        chrome_options.add_experimental_option("excludeSwitches", ["enable-automation"])
        chrome_options.add_experimental_option("useAutomationExtension", False)
        # 默认走有头模式 —— 淘宝活动页有风控，--headless 会让 starlink SDK 卡死在 loading,
        # tbopen 请求根本不会发出。Linux 虚拟机请用 Xvfb 提供虚拟显示器:
        #   sudo apt install xvfb
        #   Xvfb :99 -screen 0 1280x1024x24 &
        #   export DISPLAY=:99
        # 或者用 xvfb-run 启动调用方进程。
        # 如果调试时想强制走 headless，设置环境变量 TAOBAO_DEEPLINK_HEADLESS=1。
        if os.environ.get('TAOBAO_DEEPLINK_HEADLESS') == '1':
            chrome_options.add_argument("--headless=new")
        chrome_options.add_argument("--disable-gpu")
        chrome_options.add_argument("--no-sandbox")
        chrome_options.add_argument("--disable-dev-shm-usage")
        chrome_options.add_argument("--disable-blink-features=AutomationControlled")
        chrome_options.add_argument('log-level=3')
        # 持久化 user-data-dir：第一次手动登录一次淘宝（在弹出的窗口里输手机号/扫码），
        # 之后所有调用都会复用这份 cookie，绕过风控弹登录。
        # 路径用环境变量 TAOBAO_DEEPLINK_PROFILE 覆盖；默认放在脚本同目录下 .chrome-profile/
        profile_dir = os.environ.get('TAOBAO_DEEPLINK_PROFILE') or os.path.join(
            os.path.dirname(os.path.abspath(__file__)), '.chrome-profile'
        )
        try:
            os.makedirs(profile_dir, exist_ok=True)
        except Exception:
            pass
        chrome_options.add_argument(f"--user-data-dir={profile_dir}")
        # 开启 performance log，让 Chrome 把 Network.requestWillBeSent 事件写入日志，
        # 兜底捕获 tbopen://（hook 不一定能拦到，但浏览器发起请求时一定有 Network 事件）
        chrome_options.set_capability("goog:loggingPrefs", {"performance": "ALL", "browser": "ALL"})
        # 必须同时开 Page —— tbopen:// 这种自定义 scheme 的顶层导航在 headless 下
        # 经常不触发 Network.requestWillBeSent（被渲染端在到达网络栈前就拒掉），
        # 但 Page.frameRequestedNavigation 一定会带着 URL 触发。
        chrome_options.set_capability("goog:perfLoggingPrefs", {"enableNetwork": True, "enablePage": True})

        # 默认用 undetected-chromedriver 绕开淘宝风控；
        # 设 TAOBAO_DEEPLINK_USE_UC=0 可以回退到原生 selenium 用于对比。
        use_uc = os.environ.get('TAOBAO_DEEPLINK_USE_UC', '1') != '0'
        try:
            if use_uc:
                import undetected_chromedriver as uc
                # uc 自己管理 chromedriver 二进制，CHROME_DRIVER_PATH 仅在回退路径用
                # 把已经构造好的 chrome_options 转成 uc.ChromeOptions（属性兼容）
                uc_options = uc.ChromeOptions()
                for arg in chrome_options.arguments:
                    uc_options.add_argument(arg)
                # mobileEmulation / excludeSwitches 必须重新加，experimental_options 不会自动迁移
                for k, v in chrome_options.experimental_options.items():
                    try:
                        uc_options.add_experimental_option(k, v)
                    except Exception:
                        pass
                # uc 不支持 set_capability 的 perfLog 设置走 desired_capabilities 路径，
                # 改用 enable_cdp_events=True，后续用 driver.execute_cdp_cmd 自取
                driver = uc.Chrome(options=uc_options, use_subprocess=True)
                print("[init] 使用 undetected-chromedriver")
            else:
                service = Service(CHROME_DRIVER_PATH)
                driver = webdriver.Chrome(service=service, options=chrome_options)
                print("[init] 使用原生 selenium webdriver")
        except Exception as e:
            print(f"初始化驱动时出错: {e}")
            print("回退到原生 selenium webdriver...")
            try:
                service = Service(CHROME_DRIVER_PATH)
                driver = webdriver.Chrome(service=service, options=chrome_options)
            except Exception as e2:
                print(f"回退也失败 (路径 '{CHROME_DRIVER_PATH}'): {e2}")
                try:
                    driver = webdriver.Chrome(options=chrome_options)
                except Exception as e_path:
                    print(f"从 PATH 初始化 ChromeDriver 时出错: {e_path}")
                    print("请确保 Chrome/Chromium 和正确的 ChromeDriver 已安装并配置。")
                    return None, None

    deeplink = None
    try:
        # 启用 CDP Network 域：浏览器发起的任何 request（含 tbopen://）都会被记录
        # 即使被浏览器拦截/无 handler 也会有 requestWillBeSent 事件
        try:
            driver.execute_cdp_cmd('Network.enable', {})
        except Exception as e:
            print(f"启用 Network 域失败: {e}")
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
                // Location.prototype.href
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
                // window.location = "..."  ←— 天猫活动页常用这条路径，必须 hook Window.prototype.location
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
                // document.location = "..."
                try {
                    var docLocDesc = Object.getOwnPropertyDescriptor(Document.prototype, 'location');
                    if (docLocDesc && docLocDesc.set) {
                        var origDocLocSet = docLocDesc.set;
                        Object.defineProperty(Document.prototype, 'location', {
                            configurable: true, enumerable: true, get: docLocDesc.get,
                            set: function (v) {
                                record(typeof v === 'string' ? v : (v && v.href), 'document.location=');
                                if (typeof v === 'string' && isCustomScheme(v)) return;
                                return origDocLocSet.call(this, v);
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
                // iframe 内部也可能 top.location / parent.location = "tbopen://..."
                // 这两个的 setter 在某些浏览器版本只读，try/catch 包住即可
                try {
                    ['top','parent'].forEach(function(name){
                        try {
                            var d = Object.getOwnPropertyDescriptor(window, name);
                            if (!d) return;
                        } catch(e) {}
                    });
                } catch (e) {}
                // 兜底：周期性扫描 document.location.href（某些场景拉端会先改 href 再被浏览器拦）
                try {
                    var lastHref = location.href;
                    setInterval(function () {
                        try {
                            var h = location.href;
                            if (h !== lastHref) {
                                record(h, 'href-poll');
                                lastHref = h;
                            }
                        } catch (e) {}
                    }, 100);
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

        # 策略 1：从 hook 捕获 + 同步轮询 performance log（最长 ~10s）
        # 活动页（如 pages.tmall.com/wow/...）的 tbopen 请求通常在 1.5-3s 才发出，
        # 必须真等，否则脚本会在请求发生之前就读完日志返回 None。
        import json as _json
        # 累计所有 perf log（get_log 是消费型的，每次读完就清空，必须自己累积），
        # 顺便累计所有 ping/method 用于失败时诊断到底见过哪些 URL。
        seen_urls = []
        url_scan_prefixes = ('tbopen://', 'taobao://', 'tmall://')

        def _scan_perf_logs():
            nonlocal deeplink
            try:
                logs = driver.get_log('performance')
            except Exception:
                return
            for entry in logs:
                try:
                    msg = _json.loads(entry.get('message', '{}')).get('message', {})
                    method = msg.get('method', '')
                    params = msg.get('params', {}) or {}
                    # 三类来源都看：Network 请求 / Page 顶层导航 / Page 计划导航
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
                        if not url:
                            continue
                        if url.startswith(url_scan_prefixes):
                            seen_urls.append((method, url))
                            if not deeplink:
                                deeplink = url
                                print(f"从 perf log 捕获 ({method}): {deeplink}")
                except Exception:
                    continue

        for i in range(60):  # 最长 ~30s，活动页有时加载非常慢
            try:
                captured = driver.execute_script("return window.__capturedTaobao || [];") or []
                for item in captured:
                    url = item.get('url') if isinstance(item, dict) else None
                    if url and url.startswith(url_scan_prefixes):
                        deeplink = url
                        print(f"从 hook 捕获 ({item.get('source')}): {deeplink}")
                        break
            except Exception as e:
                print(f"读取 hook 出错: {e}")
            if deeplink:
                break
            _scan_perf_logs()
            if deeplink:
                break
            # 每 5 次循环（~2.5s）做一次"是否还卡在 loading"诊断；
            # 若 body 仍是 loading 且时间还充裕，就继续等而不是过早放弃
            if i > 0 and i % 10 == 0:
                try:
                    body_text = driver.execute_script("return (document.body && document.body.innerText || '').slice(0,200);") or ''
                    print(f"[wait] {i*0.5:.1f}s 时 body 文字: {body_text[:80]!r}")
                except Exception:
                    pass
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

        # 策略 3：循环结束后再扫一次 perf log（兜底）
        _scan_perf_logs()

        # 策略 4：直接从 page source / 全量 JS 上下文里正则扫 tbopen://
        # 适用场景：活动页（如天猫红包封面）把 schemaUrl 直接写进 __INITIAL_DATA__ /
        # SDK 配置数据块，但需要用户点击才会真的拉端 —— headless 没人点，但数据已经在 DOM 里了。
        if not deeplink:
            print("[strategy4] 开始扫 page source / window 数据")
            try:
                import re as _re
                sources = []
                try:
                    src = driver.page_source or ''
                    sources.append(('page_source', src))
                    print(f"[strategy4] page_source 长度={len(src)}")
                except Exception as e:
                    print(f"[strategy4] page_source 读取失败: {e}")
                try:
                    js_blob = driver.execute_script(
                        "return Array.from(document.scripts).map(s=>s.textContent||'').join('\\n');"
                    ) or ''
                    sources.append(('scripts', js_blob))
                    print(f"[strategy4] <script> 文本总长={len(js_blob)}")
                except Exception as e:
                    print(f"[strategy4] 读 scripts 失败: {e}")
                try:
                    win_blob = driver.execute_script(
                        "try { return JSON.stringify(window.__INITIAL_DATA__||window.__INIT_DATA__||window.pageData||window.g_config||{}); } catch(e){ return ''; }"
                    ) or ''
                    sources.append(('window_data', win_blob))
                    print(f"[strategy4] window 数据块长度={len(win_blob)}")
                except Exception as e:
                    print(f"[strategy4] 读 window 失败: {e}")
                pattern = _re.compile(r'(tbopen://[^"\'\\\s<>]+|taobao://[^"\'\\\s<>]+|tmall://[^"\'\\\s<>]+)')
                for name, blob in sources:
                    m = pattern.search(blob)
                    if m:
                        deeplink = m.group(1)
                        try:
                            deeplink = deeplink.encode('utf-8').decode('unicode_escape')
                        except Exception:
                            pass
                        print(f"[strategy4] 命中 ({name}): {deeplink[:200]}")
                        break
                    else:
                        # 退而求其次：看看里面有没有 schemaUrl / schema_url / linkUrl 字样，便于诊断
                        for kw in ('schemaUrl', 'schema_url', 'tbopen', 'tmall://', 'taobao://'):
                            if kw in blob:
                                idx = blob.find(kw)
                                print(f"[strategy4] {name} 里见到 '{kw}' @ {idx}: {blob[max(0,idx-30):idx+200]!r}")
                                break
            except Exception as e:
                print(f"[strategy4] 异常: {e}")

        # 策略 5：尝试自动点击页面上看起来像「拉起 App」的按钮，然后再扫一遍
        if not deeplink:
            print("[strategy5] 开始尝试自动点击")
            # 先打印页面上所有可见文字 + 所有可点击元素轮廓，便于诊断
            try:
                body_text = driver.execute_script("return document.body && document.body.innerText || '';") or ''
                print(f"[strategy5] body 可见文字（前 500 字符）: {body_text[:500]!r}")
            except Exception:
                pass
            try:
                btn_xpaths = [
                    "//*[contains(text(),'立即打开')]",
                    "//*[contains(text(),'打开淘宝')]",
                    "//*[contains(text(),'打开APP')]",
                    "//*[contains(text(),'打开app')]",
                    "//*[contains(text(),'去淘宝')]",
                    "//*[contains(text(),'开红包')]",
                    "//*[normalize-space(text())='开']",
                ]
                clicked = False
                for xp in btn_xpaths:
                    try:
                        els = driver.find_elements(By.XPATH, xp)
                        if els:
                            print(f"[strategy5] xpath 命中 {xp} -> {len(els)} 个元素")
                        for el in els:
                            try:
                                driver.execute_script("arguments[0].click();", el)
                                print(f"[strategy5] 已点击: {xp}")
                                clicked = True
                                break
                            except Exception as e:
                                print(f"[strategy5] 点击失败: {e}")
                                continue
                        if clicked:
                            break
                    except Exception:
                        continue
                if not clicked:
                    print("[strategy5] 没有命中任何拉端按钮 xpath")
                else:
                    for _ in range(10):
                        try:
                            captured = driver.execute_script("return window.__capturedTaobao || [];") or []
                            for item in captured:
                                u = item.get('url') if isinstance(item, dict) else None
                                if u and u.startswith(url_scan_prefixes):
                                    deeplink = u
                                    print(f"[strategy5] 点击后 hook 捕获 ({item.get('source')}): {deeplink}")
                                    break
                        except Exception:
                            pass
                        if deeplink:
                            break
                        _scan_perf_logs()
                        if deeplink:
                            break
                        time.sleep(0.5)
            except Exception as e:
                print(f"[strategy5] 异常: {e}")

        # 策略 6：调用淘宝 starlink/mtop 的 evoke 接口获取 schemaUrl（兜底）
        # 活动页 starlink.sdk.evoke 请求 cancelled 的响应里其实带 schemaUrl，
        # 我们把已经出现的 starlink.sdk.evoke 请求拦截一下，从响应体里取
        if not deeplink:
            print("[strategy6] 尝试从 starlink.sdk.evoke 响应解析")
            try:
                import json as _json
                logs = driver.get_log('performance')
                req_ids = []
                for entry in logs:
                    try:
                        msg = _json.loads(entry.get('message', '{}')).get('message', {})
                        if msg.get('method') != 'Network.responseReceived':
                            continue
                        params = msg.get('params', {}) or {}
                        url = (params.get('response') or {}).get('url', '')
                        if 'starlink.sdk.evoke' in url or 'mtop.taobao.starlink' in url:
                            req_ids.append((params.get('requestId'), url))
                    except Exception:
                        continue
                print(f"[strategy6] 找到 {len(req_ids)} 条 evoke 响应")
                for rid, url in req_ids:
                    try:
                        body = driver.execute_cdp_cmd('Network.getResponseBody', {'requestId': rid})
                        text = body.get('body', '') if body else ''
                        import re as _re
                        m = _re.search(r'(tbopen://[^"\'\\\s<>]+|taobao://[^"\'\\\s<>]+|tmall://[^"\'\\\s<>]+)', text)
                        if m:
                            deeplink = m.group(1)
                            try:
                                deeplink = deeplink.encode('utf-8').decode('unicode_escape')
                            except Exception:
                                pass
                            print(f"[strategy6] 从 evoke 响应命中: {deeplink[:200]}")
                            break
                        else:
                            print(f"[strategy6] {url[:80]} 响应无 tbopen，体长={len(text)}, 前 200: {text[:200]!r}")
                    except Exception as e:
                        print(f"[strategy6] getResponseBody 失败: {e}")
            except Exception as e:
                print(f"[strategy6] 异常: {e}")

        if deeplink:
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
        # 失败时打印诊断信息：到底见过什么
        try:
            captured = driver.execute_script("return window.__capturedTaobao || [];") or []
            print(f"[diag] hook 捕获条目数={len(captured)}, 内容={captured}")
        except Exception:
            pass
        print(f"[diag] perf log 中见过的 tb/taobao/tmall URL（{len(seen_urls)} 条）: {seen_urls[:5]}")

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
