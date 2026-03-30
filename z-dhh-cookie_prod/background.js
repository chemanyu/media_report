// 定时任务：每半天自动更新一次cookie
const UPDATE_INTERVAL = 1; // 分钟

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

// 自动更新Cookie的函数
// 首先尝试从存储中获取用户上次选择的域名
function updateCookieAutomatically() {
  chrome.storage.local.get(['savedDomain', 'savedUrl', 'savedTitle'], function(result) {
    if (result.savedDomain) {
      console.log('使用已保存的域名:', result.savedDomain);
      // 查找该域名的标签页
      chrome.tabs.query({}, function(allTabs) {
        const matchingTab = allTabs.find(tab => {
          if (!tab.url) return false;
          try {
            const tabUrl = new URL(tab.url);
            const tabDomain = tabUrl.hostname.split('.').slice(-2).join('.');
            return tabDomain === result.savedDomain;
          } catch (e) {
            return false;
          }
        });
        
        if (matchingTab) {
          console.log('找到匹配的标签页:', matchingTab.title);
          reloadAndFetchCookies(matchingTab);
        } else {
          console.log('未找到匹配域名的标签页:', result.savedDomain);
          console.log('请打开', result.savedUrl, '或手动选择标签页');
        }
      });
    } else {
      console.log('未保存域名，尝试查大航海相关页面');
      // 如果没有保存的域名，查找大喊Ghia相关的标签页
      chrome.tabs.query({}, function(allTabs) {
        const relevantTab = allTabs.find(tab => 
          tab.url && (tab.url.includes('dhh.taobao.com'))
        );
        
        if (relevantTab) {
          reloadAndFetchCookies(relevantTab);
        } else {
          console.log('未找到大航海相关页面，跳过本次更新');
        }
      });
    }
  });
}

// 刷新标签页，attach debugger 监听目标接口请求头，拿到真实 Cookie
function reloadAndFetchCookies(tab) {
  console.log('刷新标签页，准备用 debugger 抓取 Cookie:', tab.title);
  const target = { tabId: tab.id };
  const TARGET_URL_PATTERN = 'dhh.taobao.com/polystar/api/creative/material/forminfo';

  chrome.debugger.detach(target, () => {
    void chrome.runtime.lastError;
    chrome.debugger.attach(target, '1.3', () => {
    if (chrome.runtime.lastError) {
      console.error('attach debugger 失败:', chrome.runtime.lastError.message);
      return;
    }

    const pendingRequestIds = new Set();

    chrome.debugger.sendCommand(target, 'Network.enable', {}, () => {
      const onEvent = (source, method, params) => {
        if (source.tabId !== tab.id) return;

        if (method === 'Network.requestWillBeSent') {
          if (params.request.url.includes(TARGET_URL_PATTERN)) {
            pendingRequestIds.add(params.requestId);
          }
          return;
        }

        if (method !== 'Network.requestWillBeSentExtraInfo') return;
        if (!pendingRequestIds.has(params.requestId)) return;
        pendingRequestIds.delete(params.requestId);

        const headers = params.headers || {};
        const cookie = headers['cookie'] || headers['Cookie'] || '';
        if (!cookie) {
          console.log('目标接口请求头中未找到 Cookie');
          return;
        }

        console.log('从接口请求头捕获到 Cookie，长度:', cookie.length);

        // 拿到后立即 detach，不再监听
        chrome.debugger.onEvent.removeListener(onEvent);
        chrome.debugger.detach(target, () => {});

        const xsrfToken = cookie.split('; ')
          .map(c => c.split('='))
          .find(([k]) => k === 'XSRF-TOKEN')?.[1] || '';

        sendCookieToServer(cookie, xsrfToken);
      };

      chrome.debugger.onEvent.addListener(onEvent);

      // 刷新页面，触发目标接口请求
      chrome.tabs.reload(tab.id);

      // 超时保护：30秒内没捕获到就 detach
      setTimeout(() => {
        chrome.debugger.onEvent.removeListener(onEvent);
        chrome.debugger.detach(target, () => {});
        console.log('超时未捕获到目标接口，已 detach debugger');
      }, 30000);
    });
  });
  }); // detach wrapper
}

// 发送 Cookie 到后端
function sendCookieToServer(cookieString, xsrfToken) {
  fetch('http://127.0.0.1:8888/update/dhh/cookie', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ cookie: cookieString, csrfToken: xsrfToken })
      })
      .then(response => {
        if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
        return response.text();
      })
      .then(body => {
        console.log('Cookie更新成功:', body);
        chrome.notifications.create({
          type: 'basic', iconUrl: 'icon48.png',
          title: 'Cookie更新成功',
          message: `更新时间: ${new Date().toLocaleString()}`,
          priority: 1
        });
      })
      .catch(error => {
        console.error('发送Cookie失败:', error.message);
        chrome.notifications.create({
          type: 'basic', iconUrl: 'icon48.png',
          title: 'Cookie更新失败',
          message: error.message, priority: 2
        });
      });
}

function logError(message) {
  console.error('自动更新Cookie错误:', message);
  chrome.notifications.create({
    type: 'basic',
    iconUrl: 'icon48.png',
    title: 'Cookie更新错误',
    message: message,
    priority: 2
  });
}

// 监听来自popup的消息（可选，用于手动触发）
chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  if (request.action === 'updateCookieNow') {
    updateCookieAutomatically();
    sendResponse({ status: 'started' });
  }
  return true;
});
