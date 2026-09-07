// 定时任务：每 3 分钟自动更新一次cookie
const UPDATE_INTERVAL = 3; // 分钟

// 京橙复投接口实际请求的域名，Cookie 必须按这个 URL 取。
// 不能按 domain: 'jd.com' 取：Chrome 的 domain 过滤会把所有 *.jd.com 子域
// （jcheng/www/passport…）的 host-only Cookie 一并返回，__jda / __jdc /
// 3AB9D23F7A4B3C9B 这类同名 Cookie 会出现多份，拼出的字符串带重复 key、值取错，
// 导致接口鉴权失败。
const COOKIE_TARGET_URL = 'https://api.m.jd.com/';

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

// 把 chrome.cookies 返回的列表拼成请求头用的 Cookie 字符串。
// Chrome 按 RFC 6265 的顺序返回（path 越长越靠前），同名 Cookie 只保留第一个，
// 也就是最贴近目标 URL 的那份，避免重复 key 覆盖出错误的值。
function buildCookieHeader(cookies) {
  const seen = new Set();
  return cookies
    .filter(cookie => {
      if (seen.has(cookie.name)) return false;
      seen.add(cookie.name);
      return true;
    })
    .map(cookie => `${cookie.name}=${cookie.value}`)
    .join('; ');
}

// 提取并发送Cookie的核心逻辑
function fetchAndSendCookies() {
  try {
    console.log('正在获取接口域名的Cookie:', COOKIE_TARGET_URL);

    // 按目标请求 URL 取 Cookie：等价于浏览器真正会发给 api.m.jd.com 的那一份
    chrome.cookies.getAll({ url: COOKIE_TARGET_URL }, function(cookies) {
      if (chrome.runtime.lastError) {
        console.error('获取Cookie失败:', chrome.runtime.lastError.message);
        return;
      }

      if (cookies.length === 0) {
        console.log('未找到Cookie，请先登录并刷新京橙页面:', COOKIE_TARGET_URL);
        return;
      }

      // 格式化cookies为字符串
      const cookieString = buildCookieHeader(cookies);

      // 提取 x-csrftoken（查找名为 X-Csrftoken 的 cookie）
      const csrfCookie = cookies.find(cookie => 
        cookie.name === 'csrftoken'
      );
      const csrfToken = csrfCookie ? csrfCookie.value : '';
      
      console.log('Cookie获取成功，准备发送到后端...');
      console.log('所有Cookie名称:', cookies.map(c => c.name).join(', '));
      if (csrfToken) {
        console.log('X-CSRF-Token已找到');
      } else {
        console.log('未找到X-CSRF-Token');
      }

      // 发送到后端API
      // Send POST request to the API
      fetch('http://127.0.0.1:8888/update/jingcheng/futou/cookie', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ 
          cookie: cookieString,
          csrfToken: csrfToken
        })
      })
      .then(response => {
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.text();
      })
      .then(body => {
        console.log('Cookie更新成功:', body);
        // 可以选择发送通知
        chrome.notifications.create({
          type: 'basic',
          iconUrl: 'icon48.png',
          title: '复投Cookie更新成功',
          message: `更新时间: ${new Date().toLocaleString()}`,
          priority: 1
        });
      })
      .catch(error => {
        console.error('发送Cookie失败:', error.message);
        chrome.notifications.create({
          type: 'basic',
          iconUrl: 'icon48.png',
          title: '复投Cookie更新失败',
          message: error.message,
          priority: 2
        });
      });
    });
  } catch (error) {
    console.error('处理URL时出错:', error.message);
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
