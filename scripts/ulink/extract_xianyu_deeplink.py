import sys
# 关键：Go 通过 pipe 调用本脚本时，Windows Python 默认 stdout 退到 GBK，
# 任何 emoji（✅/❌）会 UnicodeEncodeError 直接崩溃。强制 UTF-8。
try:
    sys.stdout.reconfigure(encoding='utf-8', errors='replace')
    sys.stderr.reconfigure(encoding='utf-8', errors='replace')
except Exception:
    pass

from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.support.ui import WebDriverWait
from selenium.common.exceptions import TimeoutException
import time
import re

# 配置 ChromeDriver 路径
#CHROME_DRIVER_PATH = "/opt/homebrew/bin/chromedriver"
CHROME_DRIVER_PATH = "D:\\150\\chromedriver-win64\\chromedriver.exe"

def get_xianyu_deeplink(short_url, driver=None, platform="android"):
    """
    从闲鱼短链接中提取 deeplink
    platform: "android" 或 "ios"
    """
    internal_driver = False
    if driver is None:
        internal_driver = True
        chrome_options = Options()

        # 根据平台选择不同的 User-Agent
        if platform.lower() == "android":
            mobile_emulation = {
                "deviceMetrics": {"width": 412, "height": 915, "pixelRatio": 3.5},
                "userAgent": "Mozilla/5.0 (Linux; Android 10; SM-G973F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Mobile Safari/537.36"
            }
        else:
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
            driver = webdriver.Chrome(
                service=service,
                options=chrome_options
            )
        except Exception as e:
            print(f"初始化 ChromeDriver 时出错: {e}")
            try:
                driver = webdriver.Chrome(
                    options=chrome_options
                )
            except Exception as e_path:
                print(f"从 PATH 初始化 ChromeDriver 时出错: {e_path}")
                return None

    try:
        print(f"正在访问闲鱼链接: {short_url}, 平台: {platform}")

        # Android 平台：使用 CDP 监听所有导航尝试
        if platform.lower() == "android":
            # 启用 CDP 的 Page 域来监听导航
            driver.execute_cdp_cmd('Page.enable', {})

            # 关键：在页面任何脚本执行前注入劫持代码
            # fleamarket:// 跳转通过 location.href / location.assign / a.click 触发
            # 浏览器对未知 scheme 直接拦截，不会产生网络请求，所以必须在 JS 层捕获
            hook_script = r"""
            (function () {
                window.__capturedDeeplinks = [];
                window.__hookStats = { fetch: 0, xhr: 0, locHref: 0, aClick: 0, hookInstalled: true };
                window.__capturedUrls = []; // 记录所有 fetch/xhr 的 URL，用于诊断
                var record = function (url, source) {
                    try {
                        if (typeof url === 'string' && url.indexOf('fleamarket://') === 0) {
                            window.__capturedDeeplinks.push({ url: url, source: source });
                            console.log('[DEEPLINK_CAPTURED]' + source + '|' + url);
                        }
                    } catch (e) {}
                };

                // 1. 劫持 location.href setter
                try {
                    var origLocation = window.location;
                    var hrefDesc = Object.getOwnPropertyDescriptor(Location.prototype, 'href')
                        || Object.getOwnPropertyDescriptor(window.Location.prototype, 'href');
                    if (hrefDesc && hrefDesc.set) {
                        var origSet = hrefDesc.set;
                        Object.defineProperty(Location.prototype, 'href', {
                            configurable: true,
                            enumerable: true,
                            get: hrefDesc.get,
                            set: function (v) {
                                window.__hookStats.locHref++;
                                record(v, 'location.href');
                                if (typeof v === 'string' && v.indexOf('fleamarket://') === 0) {
                                    return; // 阻止真正跳转，避免页面被打断
                                }
                                return origSet.call(this, v);
                            }
                        });
                    }
                } catch (e) {}

                // 2. 劫持 location.assign / replace
                try {
                    var origAssign = Location.prototype.assign;
                    Location.prototype.assign = function (v) {
                        record(v, 'location.assign');
                        if (typeof v === 'string' && v.indexOf('fleamarket://') === 0) return;
                        return origAssign.apply(this, arguments);
                    };
                    var origReplace = Location.prototype.replace;
                    Location.prototype.replace = function (v) {
                        record(v, 'location.replace');
                        if (typeof v === 'string' && v.indexOf('fleamarket://') === 0) return;
                        return origReplace.apply(this, arguments);
                    };
                } catch (e) {}

                // 3. 劫持 window.open
                try {
                    var origOpen = window.open;
                    window.open = function (v) {
                        record(v, 'window.open');
                        if (typeof v === 'string' && v.indexOf('fleamarket://') === 0) return null;
                        return origOpen.apply(this, arguments);
                    };
                } catch (e) {}

                // 4. 劫持 <a>.click —— 闲鱼会动态创建 <a href="fleamarket://..."> 再 click
                try {
                    var origClick = HTMLAnchorElement.prototype.click;
                    HTMLAnchorElement.prototype.click = function () {
                        window.__hookStats.aClick++;
                        record(this.href, 'a.click');
                        if (typeof this.href === 'string' && this.href.indexOf('fleamarket://') === 0) return;
                        return origClick.apply(this, arguments);
                    };
                } catch (e) {}

                // 5. 劫持 iframe.src —— 部分场景用隐藏 iframe 触发
                try {
                    var iframeDesc = Object.getOwnPropertyDescriptor(HTMLIFrameElement.prototype, 'src');
                    if (iframeDesc && iframeDesc.set) {
                        var origIframeSet = iframeDesc.set;
                        Object.defineProperty(HTMLIFrameElement.prototype, 'src', {
                            configurable: true,
                            enumerable: true,
                            get: iframeDesc.get,
                            set: function (v) {
                                record(v, 'iframe.src');
                                if (typeof v === 'string' && v.indexOf('fleamarket://') === 0) return;
                                return origIframeSet.call(this, v);
                            }
                        });
                    }
                } catch (e) {}

                // 6. 劫持 fetch —— 闲鱼通过 mtop API 异步拿 deeplink，必须扫响应体
                try {
                    var origFetch = window.fetch;
                    if (origFetch) {
                        window.fetch = function () {
                            window.__hookStats.fetch++;
                            var fetchArgs = arguments;
                            var fetchUrl = (fetchArgs[0] && fetchArgs[0].toString) ? fetchArgs[0].toString() : '';
                            try { window.__capturedUrls.push('fetch:' + fetchUrl.slice(0,200)); } catch(e){}
                            var p = origFetch.apply(this, arguments);
                            try {
                                p.then(function (res) {
                                    try {
                                        var clone = res.clone();
                                        clone.text().then(function (txt) {
                                            if (txt && txt.indexOf('fleamarket') !== -1) {
                                                // 兼容 / 编码、\\/ 转义
                                                var norm = txt.replace(/\\u002[fF]/g, '/').replace(/\\\//g, '/');
                                                var ms = norm.match(/fleamarket:\/\/[^\s"'<>\\)]+/g);
                                                if (ms) for (var j=0; j<ms.length; j++) {
                                                    record(ms[j], 'fetch:' + fetchUrl.slice(0,80));
                                                }
                                            }
                                        }).catch(function(){});
                                    } catch (e) {}
                                }).catch(function(){});
                            } catch (e) {}
                            return p;
                        };
                    }
                } catch (e) {}

                // 7. 劫持 XMLHttpRequest —— 老接口，闲鱼 mtop 走的也可能是 XHR
                try {
                    var OrigXHR = window.XMLHttpRequest;
                    if (OrigXHR) {
                        var XHRWrap = function () {
                            var xhr = new OrigXHR();
                            var openedUrl = '';
                            var origOpen = xhr.open;
                            xhr.open = function (method, url) {
                                openedUrl = url || '';
                                window.__hookStats.xhr++;
                                try { window.__capturedUrls.push('xhr:' + (openedUrl + '').slice(0,200)); } catch(e){}
                                return origOpen.apply(xhr, arguments);
                            };
                            xhr.addEventListener('load', function () {
                                try {
                                    var txt = xhr.responseText;
                                    if (txt && txt.indexOf('fleamarket') !== -1) {
                                        var norm = txt.replace(/\\u002[fF]/g, '/').replace(/\\\//g, '/');
                                        var ms = norm.match(/fleamarket:\/\/[^\s"'<>\\)]+/g);
                                        if (ms) for (var j=0; j<ms.length; j++) {
                                            record(ms[j], 'xhr:' + (xhr.responseURL || openedUrl).slice(0,80));
                                        }
                                    }
                                } catch (e) {}
                            });
                            return xhr;
                        };
                        XHRWrap.prototype = OrigXHR.prototype;
                        window.XMLHttpRequest = XHRWrap;
                    }
                } catch (e) {}
            })();
            """
            driver.execute_cdp_cmd('Page.addScriptToEvaluateOnNewDocument', {'source': hook_script})
            print("✅ 已注入 deeplink 劫持脚本（覆盖 location/open/a.click/iframe.src）")

        driver.get(short_url)

        # 使用显式等待页面加载完成
        time.sleep(1)
        print("等待页面加载...")
        try:
            WebDriverWait(driver, 5).until(
                lambda d: d.execute_script('return document.readyState') == 'complete'
            )
        except TimeoutException:
            print("页面加载等待超时，继续后续处理")

        current_url = driver.current_url
        print(f"当前 URL: {current_url}")

        # 落地后立刻验证 hook 是否被 Page.addScriptToEvaluateOnNewDocument 注入成功
        if platform.lower() == "android":
            try:
                hook_alive = driver.execute_script(
                    "return window.__hookStats ? window.__hookStats : null;"
                )
                print(f"Hook 状态（落地后立刻读）: {hook_alive}")
            except Exception as e:
                print(f"读取 hook 状态出错: {e}")

        # app=chrome / download=true 都会让页面进"非 APP/PC 浏览器"分支，去掉重新加载
        def _strip_param(u, key):
            u = re.sub(r'[?&]' + key + r'=[^&]*', '', u)
            # 修复 ?& → ?，?$ 没问题
            u = u.replace('?&', '?')
            if u.endswith('?'):
                u = u[:-1]
            return u

        if platform.lower() == "android" and ('download=true' in current_url or 'app=chrome' in current_url):
            clean_url = current_url
            for k in ('download', 'app'):
                clean_url = _strip_param(clean_url, k)
            print(f"检测到 download/app 参数，重新访问清洁 URL → {clean_url}")
            driver.get(clean_url)
            time.sleep(2)
            try:
                WebDriverWait(driver, 5).until(
                    lambda d: d.execute_script('return document.readyState') == 'complete'
                )
            except TimeoutException:
                pass
            current_url = driver.current_url
            print(f"清洁 URL 落地: {current_url}")

        # Android 平台：多种方法捕获 fleamarket 链接
        if platform.lower() == "android":
            print("=" * 80)
            print("Android 平台，开始搜索 fleamarket 链接...")
            print("=" * 80)

            # 排除明显是模板字符串/示例的伪命中
            def is_real_deeplink(u):
                if not u or not u.startswith('fleamarket://'):
                    return False
                # 至少要有一个非占位符字符
                if u.endswith('://') or u.endswith('://"') or u == 'fleamarket://':
                    return False
                # 排除模板字面量
                if '${' in u or '{{' in u:
                    return False
                return True

            # 诊断：打印 hook 状态和已经发出的请求 URL
            try:
                stats = driver.execute_script("return window.__hookStats || null;")
                urls = driver.execute_script("return window.__capturedUrls || [];")
                print(f"\n[诊断] hookStats = {stats}")
                print(f"[诊断] 已捕获 {len(urls or [])} 个 fetch/xhr 请求:")
                for u in (urls or [])[:30]:
                    print(f"  - {u}")
            except Exception as e:
                print(f"[诊断] 读取 hook 状态失败: {e}")

            # 方法1: 等 JS 触发跳转 / 异步 mtop 接口返回，从 hook 捕获
            # 闲鱼 callapp 是 mtop 异步接口，响应可能晚到，给足时间
            print("\n[方法1] 等待 JS 跳转/XHR 劫持 ...")
            for _ in range(20):
                try:
                    captured = driver.execute_script("return window.__capturedDeeplinks || [];")
                    if captured:
                        for item in captured:
                            url = item.get('url') if isinstance(item, dict) else None
                            if is_real_deeplink(url):
                                print(f"✅ 从 JS 劫持捕获 ({item.get('source')}): {url}")
                                return url
                except Exception as e:
                    print(f"读取劫持结果出错: {e}")

            # 方法2: 主动尝试调用页面里的 callapp 函数（闲鱼常用 callappReflow）
            print("\n[方法2] 尝试主动调用页面 callapp 函数 ...")
            try:
                driver.execute_script("""
                    try {
                        // 触摸文档/点击空白位置可能触发闲鱼的拉起逻辑
                        var ev = new Event('touchstart', {bubbles:true});
                        document.dispatchEvent(ev);
                        document.body && document.body.click();
                    } catch(e){}
                    // 扫描全局对象，找到名字带 callapp / openApp / launch 的函数并尝试调用
                    try {
                        var keys = Object.keys(window);
                        for (var i=0; i<keys.length; i++) {
                            var k = keys[i];
                            if (/callapp|openapp|launch|fleamarket|gotoapp/i.test(k)) {
                                try {
                                    var v = window[k];
                                    if (typeof v === 'function') v();
                                    else if (v && typeof v === 'object') {
                                        for (var sk in v) {
                                            if (typeof v[sk] === 'function' && /callapp|open|launch|jump/i.test(sk)) {
                                                try { v[sk](); } catch(e){}
                                            }
                                        }
                                    }
                                } catch(e){}
                            }
                        }
                    } catch(e){}
                """)
                time.sleep(2)
                captured = driver.execute_script("return window.__capturedDeeplinks || [];")
                for item in captured or []:
                    url = item.get('url') if isinstance(item, dict) else None
                    if is_real_deeplink(url):
                        print(f"✅ 主动调用后捕获 ({item.get('source')}): {url}")
                        return url
            except Exception as e:
                print(f"主动调用出错: {e}")

            # 方法3: 全文 + 全脚本正则扫描，挑出最完整的一个
            print("\n[方法3] 扫描全文 / script innerText / 全局对象 ...")
            try:
                blob = driver.execute_script("""
                    var parts = [document.documentElement.outerHTML];
                    try {
                        var scripts = document.scripts;
                        for (var i=0;i<scripts.length;i++) {
                            if (scripts[i].textContent) parts.push(scripts[i].textContent);
                        }
                    } catch(e){}
                    // 浅层扫描 window 上的字符串属性（避免遍历过深爆栈）
                    try {
                        var ks = Object.keys(window);
                        for (var i=0;i<ks.length;i++) {
                            try {
                                var v = window[ks[i]];
                                if (typeof v === 'string') parts.push(v);
                                else if (v && typeof v === 'object') {
                                    for (var sk in v) {
                                        try {
                                            var sv = v[sk];
                                            if (typeof sv === 'string') parts.push(sv);
                                        } catch(e){}
                                    }
                                }
                            } catch(e){}
                        }
                    } catch(e){}
                    return parts.join('\\n');
                """)
                # 找出所有命中，按长度优先（最长的最可能是带完整参数的真链接）
                matches = re.findall(r'fleamarket://[^\s"\'<>\\)]+', blob or '')
                # 去重保序
                seen = set()
                uniq = []
                for m in matches:
                    if m not in seen:
                        seen.add(m)
                        uniq.append(m)
                # 过滤掉明显伪命中
                cands = [m for m in uniq if is_real_deeplink(m) and len(m) > len('fleamarket://') + 5]
                # 优先选带 query / path 较丰富的（长度最大）
                cands.sort(key=len, reverse=True)
                print(f"全文扫描找到 {len(matches)} 处命中，{len(cands)} 个候选")
                for c in cands[:5]:
                    print(f"  候选: {c[:200]}")
                if cands:
                    return cands[0]
            except Exception as e:
                print(f"全文扫描出错: {e}")

            # 方法4: console 日志兜底
            try:
                logs = driver.get_log('browser')
                print(f"\n[方法4] 浏览器控制台日志 ({len(logs)} 条) ...")
                for log in logs:
                    message = log.get('message', '')
                    if 'fleamarket://' in message:
                        match = re.search(r'fleamarket://[^\s"\'<>\\)]+', message)
                        if match and is_real_deeplink(match.group(0)):
                            print(f"✅ 控制台日志找到: {match.group(0)}")
                            return match.group(0)
            except Exception as e:
                print(f"读取控制台日志出错: {e}")

            print("❌ 未能提取 deeplink，请检查链接是否有效")

        # iOS 平台：返回 pages.goofish.com 链接
        if platform.lower() == "ios":
            if 'pages.goofish.com' in current_url:
                print(f"iOS 平台，返回: {current_url}")
                return current_url
            for request in driver.requests:
                if 'pages.goofish.com' in request.url:
                    print(f"从网络请求找到: {request.url}")
                    return request.url

            return None

    except Exception as e:
        print(f"提取过程中发生错误: {e}")
        return None
    finally:
        if internal_driver and driver is not None:
            driver.quit()

    return None
