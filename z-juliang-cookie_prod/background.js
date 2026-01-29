// 定时任务：每半天自动更新一次cookie
const UPDATE_INTERVAL = 1800; // 分钟

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
          fetchAndSendCookies(matchingTab.url);
        } else {
          console.log('未找到匹配域名的标签页:', result.savedDomain);
          console.log('请打开', result.savedUrl, '或手动选择标签页');
        }
      });
    } else {
      console.log('未保存域名，尝试查找京东联盟相关页面');
      // 如果没有保存的域名，查找京东联盟相关的标签页
      chrome.tabs.query({}, function(allTabs) {
        const relevantTab = allTabs.find(tab => 
          tab.url && (tab.url.includes('jd.com') || 
                     tab.url.includes('union.jd.com'))
        );
        
        if (relevantTab) {
          fetchAndSendCookies(relevantTab.url);
        } else {
          console.log('未找到京东联盟相关页面，跳过本次更新');
        }
      });
    }
  });
}

// 提取并发送Cookie的核心逻辑
function fetchAndSendCookies(url) {
  try {
    const currentUrl = new URL(url);
    const domain = currentUrl.hostname;
    const topLevelDomain = domain.split('.').slice(-2).join('.');

    console.log('正在获取域名的Cookie:', topLevelDomain);

    // 获取指定域名的所有cookies
    chrome.cookies.getAll({ domain: topLevelDomain }, function(cookies) {
      if (chrome.runtime.lastError) {
        console.error('获取Cookie失败:', chrome.runtime.lastError.message);
        return;
      }

      if (cookies.length === 0) {
        console.log('当前域名未找到Cookie:', topLevelDomain);
        return;
      }

      // 格式化cookies为字符串
      const cookieString = cookies.map(cookie => `${cookie.name}=${cookie.value}`).join('; ');
      
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
      fetch('http://127.0.0.1:8888/update/juliang/cookie', {
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
          title: 'Cookie更新成功',
          message: `更新时间: ${new Date().toLocaleString()}`,
          priority: 1
        });
      })
      .catch(error => {
        console.error('发送Cookie失败:', error.message);
        chrome.notifications.create({
          type: 'basic',
          iconUrl: 'icon48.png',
          title: 'Cookie更新失败',
          message: error.message,
          priority: 2
        });
      });
    });
  } catch (error) {
    console.error('处理URL时出错:', error.message);
  }
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
