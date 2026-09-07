// 定时任务：每 3 分钟自动更新一次cookie
const UPDATE_INTERVAL = 3; // 分钟

// 京橙复投接口的真实请求地址，Cookie 按这个 URL 取。
// 基准来自浏览器 Network 里一次正常成功的请求（functionId=union_orange_user_api_new），
// 它的 Cookie 头只有 18 个：shshshfpa/x/b、__jdu、__jda、__jdc、3AB9D23F7A4B3C9B/CSS、
// thor、flash、light_key、pinId、pin、unick、ceshi3.com、_tp、_pst、logining。
//
// 为什么必须按 url 取，两种错法都试过了：
//   · domain: 'jd.com' —— Chrome 的 domain 过滤是后缀匹配，把所有 *.jd.com 子域的
//     host-only Cookie 混进来，__jda / __jdc / 3AB9D23F7A4B3C9B 出现多份，值取错。
//   · 多域名合并 —— 会塞进 focus-* / me_saas_userInfo / pt_key / sdtoken 这类
//     jcheng 控制台和登录页专属的 Cookie，浏览器根本不会发给本接口，纯属噪声。
// getAll({url}) 让 Chrome 自己按 domain/path/secure 规则算，结果与真实请求头一致。
const COOKIE_TARGET_URL = 'https://api.m.jd.com/api/';

// 关键 Cookie，缺了基本就是登录态失效，只告警不阻断（京东随时可能改名）
const REQUIRED_COOKIES = ['thor', 'pin', 'light_key', '3AB9D23F7A4B3C9B'];

// 监听定时器触发
chrome.alarms.onAlarm.addListener((alarm) => {
  if (alarm.name === 'updateCookie') {
    console.log('定时更新Cookie任务触发:', new Date().toLocaleString());
    updateCookieAutomatically();
  }
});

// 扩展安装或更新时，创建定时器并立即执行一次
chrome.runtime.onInstalled.addListener(() => {
  console.log('Chrome扩展已安装/更新');
  
  // 创建定时器（每30分钟执行一次）
  chrome.alarms.create('updateCookie', {
    periodInMinutes: UPDATE_INTERVAL
  });
  
  // console.log('定时器已创建，立即执行一次Cookie更新');
  // updateCookieAutomatically();
});

// 扩展启动时，确保定时器存在
chrome.runtime.onStartup.addListener(() => {
  console.log('Chrome扩展启动');
  
  // 确保定时器存在
  chrome.alarms.create('updateCookie', {
    periodInMinutes: UPDATE_INTERVAL
  });
});

// 自动更新Cookie。
// Cookie 一律按 COOKIE_TARGET_URL 取，与打开的是哪个标签页无关；这里只用「有没有
// 打开京橙/京东页面」当作用户仍在登录态的信号，没有就跳过本次上报，免得把过期
// Cookie 覆盖上去。
function updateCookieAutomatically() {
  chrome.tabs.query({}, function(allTabs) {
    const relevantTab = allTabs.find(tab => tab.url && tab.url.includes('jd.com'));

    if (relevantTab) {
      console.log('找到京东相关标签页:', relevantTab.title);
      fetchAndSendCookies();
    } else {
      console.log('未找到京东相关页面，跳过本次更新');
    }
  });
}

// 取接口域名的 Cookie 并按 Cookie 头格式拼好。
// Chrome 按 RFC 6265 顺序返回（path 越长越靠前），同名只留第一个，即最贴近目标
// URL 的那份，避免重复 key 覆盖出错误的值。
async function collectJdCookies() {
  const merged = new Map();

  const cookies = await chrome.cookies.getAll({ url: COOKIE_TARGET_URL });
  for (const cookie of cookies) {
    if (merged.has(cookie.name)) continue;
    merged.set(cookie.name, cookie.value);
  }

  return merged;
}

// 拼成请求头用的 Cookie 字符串
function buildCookieHeader(merged) {
  return Array.from(merged, ([name, value]) => `${name}=${value}`).join('; ');
}

// 提取并发送Cookie的核心逻辑
async function fetchAndSendCookies() {
  const merged = await collectJdCookies();

  if (merged.size === 0) {
    console.log('未找到Cookie，请先登录京橙页面并刷新:', COOKIE_TARGET_URL);
    return;
  }

  const cookieString = buildCookieHeader(merged);

  const missing = REQUIRED_COOKIES.filter(name => !merged.has(name));
  if (missing.length > 0) {
    // 只告警不拦截：京东改 Cookie 名时不至于把整条链路卡死
    console.warn('缺少关键Cookie，可能已退出登录:', missing.join(', '));
  }

  // 提取 x-csrftoken（查找名为 csrftoken 的 cookie）
  const csrfToken = merged.get('csrftoken') || '';

  console.log(`Cookie获取成功，共 ${merged.size} 个 / ${cookieString.length} 字符，准备发送到后端...`);
  console.log('所有Cookie名称:', Array.from(merged.keys()).join(', '));

  try {
    const response = await fetch('http://127.0.0.1:8888/update/jingcheng/futou/cookie', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        cookie: cookieString,
        csrfToken: csrfToken
      })
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const body = await response.text();
    console.log('Cookie更新成功:', body);
    chrome.notifications.create({
      type: 'basic',
      iconUrl: 'icon48.png',
      title: '复投Cookie更新成功',
      message: `${merged.size} 个Cookie / ${cookieString.length} 字符，${new Date().toLocaleString()}`,
      priority: 1
    });
  } catch (error) {
    console.error('发送Cookie失败:', error.message);
    chrome.notifications.create({
      type: 'basic',
      iconUrl: 'icon48.png',
      title: '复投Cookie更新失败',
      message: error.message,
      priority: 2
    });
  }
}

// 监听来自popup的消息（可选，用于手动触发）
chrome.runtime.onMessage.addListener((request, _sender, sendResponse) => {
  if (request.action === 'updateCookieNow') {
    updateCookieAutomatically();
    sendResponse({ status: 'started' });
  }
  return true;
});
