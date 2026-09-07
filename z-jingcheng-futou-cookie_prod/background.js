// 定时任务：每 3 分钟自动更新一次cookie
const UPDATE_INTERVAL = 3; // 分钟

// Cookie 采集顺序（重要）：从「最贴近接口」到「最兜底」。
// 同名 Cookie 只取先出现的那一份，所以越靠前的来源优先级越高。
//
// 为什么要分层而不是单取一个来源：
//   · 只按 domain: 'jd.com' 取 —— Chrome 的 domain 过滤是后缀匹配，会把所有
//     *.jd.com 子域的 host-only Cookie 混进来，__jda / __jdc / 3AB9D23F7A4B3C9B
//     这类同名 Cookie 出现多份，值取错导致接口鉴权失败。
//   · 只按 url: api.m.jd.com 取 —— 值是对的，但只剩浏览器真会发给该接口的那份，
//     jcheng 控制台的 host-only Cookie（focus-* / me_saas_userInfo / switch_to_bpro）
//     和登录态 Cookie（pt_key / pt_pin / pt_token / pwdt_id / pt_st / sdtoken）会全丢，
//     入库字符串明显偏短。
// 分层合并 = 范围取全量，冲突时接口域名的值优先。
const COOKIE_SOURCES = [
  { url: 'https://api.m.jd.com/' },       // 复投接口本身，鉴权 Cookie 以这份为准
  { url: 'https://jcheng.jd.com/' },      // 京橙控制台 host-only：focus-*、me_saas_userInfo、switch_to_bpro
  { url: 'https://passport.jd.com/' },    // 登录态：pt_key、pt_pin、pt_token、pwdt_id、pt_st、sdtoken
  { url: 'https://www.jd.com/' },         // .jd.com 通用：__jdv、__jdb、visitkey、webp
  { domain: 'jd.com' },                   // 兜底：其余任意 *.jd.com 子域的 host-only Cookie
];

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

// 按 COOKIE_SOURCES 顺序合并，同名只留先出现的那份（前面的来源优先级更高）。
// 单个来源内部 Chrome 已按 RFC 6265 排序（path 越长越靠前），所以同来源内也是
// 最贴近目标 URL 的那份胜出。
async function collectJdCookies() {
  const merged = new Map();

  for (const filter of COOKIE_SOURCES) {
    let cookies;
    try {
      cookies = await chrome.cookies.getAll(filter);
    } catch (error) {
      console.warn('取Cookie失败，跳过该来源:', JSON.stringify(filter), error.message);
      continue;
    }

    let added = 0;
    for (const cookie of cookies) {
      if (merged.has(cookie.name)) continue;
      merged.set(cookie.name, cookie.value);
      added++;
    }
    console.log(`来源 ${filter.url || 'domain:' + filter.domain}：返回 ${cookies.length} 个，新增 ${added} 个`);
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
    console.log('未找到任何 jd.com Cookie，请先登录并刷新京橙页面');
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
