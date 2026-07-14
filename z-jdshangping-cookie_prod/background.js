// 定时任务：每 5 分钟自动刷新一次页面并抓取信息
const UPDATE_INTERVAL = 5; // 分钟

// 目标接口：京东联盟 api.m.jd.com/api 上任意带签名参数的请求即可
// 不挑 functionId（会变，如 union_orange_goods_api / union_orange_material_api ...）
const TARGET_URL_PATTERN = 'api.m.jd.com/api';
// 默认页面域名（顶级域名匹配）
const DEFAULT_PAGE_DOMAIN = 'jd.com';

// 只要 URL 命中目标接口且带上 h5st + x-api-eid-token 就算目标请求
function isTargetRequest(url) {
  return url.includes(TARGET_URL_PATTERN) &&
         url.includes('h5st=') &&
         url.includes('x-api-eid-token=');
}

// 监听定时器触发
chrome.alarms.onAlarm.addListener((alarm) => {
  if (alarm.name === 'updateJdCookie') {
    console.log('定时抓取京东Cookie任务触发:', new Date().toLocaleString());
    updateCookieAutomatically();
  }
});

// 扩展安装或更新时，创建定时器
chrome.runtime.onInstalled.addListener(() => {
  console.log('Chrome扩展已安装/更新（京东商品Cookie更新）');
  chrome.alarms.create('updateJdCookie', {
    periodInMinutes: UPDATE_INTERVAL
  });
});

// 扩展启动时，确保定时器存在
chrome.runtime.onStartup.addListener(() => {
  console.log('Chrome扩展启动（京东商品Cookie更新）');
  chrome.alarms.create('updateJdCookie', {
    periodInMinutes: UPDATE_INTERVAL
  });
});

// 自动更新：优先用已保存的域名，否则找 jd.com 相关页面
function updateCookieAutomatically() {
  chrome.storage.local.get(['savedDomain', 'savedUrl', 'savedTitle'], function(result) {
    if (result.savedDomain) {
      console.log('使用已保存的域名:', result.savedDomain);
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
          console.log('请打开', result.savedUrl, '或在弹窗手动选择标签页');
        }
      });
    } else {
      console.log('未保存域名，尝试查找京东联盟相关页面');
      chrome.tabs.query({}, function(allTabs) {
        const relevantTab = allTabs.find(tab =>
          tab.url && tab.url.includes('union.jd.com')
        );

        if (relevantTab) {
          reloadAndFetchCookies(relevantTab);
        } else {
          console.log('未找到京东联盟相关页面，跳过本次更新');
        }
      });
    }
  });
}

// 从 URL query string 中解析 x-api-eid-token / h5st / uuid
function parseTokensFromUrl(url) {
  let eidToken = '';
  let h5st = '';
  let uuid = '';
  try {
    const u = new URL(url);
    eidToken = u.searchParams.get('x-api-eid-token') || '';
    h5st = u.searchParams.get('h5st') || '';
    uuid = u.searchParams.get('uuid') || '';
  } catch (e) {
    console.warn('解析 URL query 失败:', e.message);
  }
  return { eidToken, h5st, uuid };
}

// 刷新标签页，attach debugger 监听目标接口，拿到真实 Cookie + URL 中的 token
function reloadAndFetchCookies(tab) {
  console.log('刷新标签页，准备用 debugger 抓取信息:', tab.title);
  const target = { tabId: tab.id };

  chrome.debugger.detach(target, () => {
    void chrome.runtime.lastError;
    chrome.debugger.attach(target, '1.3', () => {
      if (chrome.runtime.lastError) {
        console.error('attach debugger 失败:', chrome.runtime.lastError.message);
        return;
      }

      // requestWillBeSent 拿 URL(含 eid/h5st)，ExtraInfo 拿真实 Cookie，用 requestId 关联
      const pending = new Map(); // requestId -> { eidToken, h5st }

      chrome.debugger.sendCommand(target, 'Network.enable', {}, () => {
        const onEvent = (source, method, params) => {
          if (source.tabId !== tab.id) return;

          if (method === 'Network.requestWillBeSent') {
            const url = params.request.url || '';
            // 命中目标接口（含签名参数）即记录，functionId/方法不限
            if (isTargetRequest(url)) {
              pending.set(params.requestId, parseTokensFromUrl(url));
            }
            return;
          }

          if (method !== 'Network.requestWillBeSentExtraInfo') return;
          if (!pending.has(params.requestId)) return;
          const tokens = pending.get(params.requestId);
          pending.delete(params.requestId);

          const headers = params.headers || {};
          const cookie = headers['cookie'] || headers['Cookie'] || '';
          if (!cookie) {
            console.log('目标接口请求头中未找到 Cookie');
            return;
          }

          // 拿到后立即 detach，不再监听
          chrome.debugger.onEvent.removeListener(onEvent);
          chrome.debugger.detach(target, () => {});

          reportInfo(cookie, tokens.eidToken, tokens.h5st, tokens.uuid, tab.id);
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
  });
}

// 抓取到信息后的处理：先打印到 console，再上报
function reportInfo(cookie, eidToken, h5st, uuid, tabId) {
  console.log('==== 京东商品信息抓取成功 ====', new Date().toLocaleString());
  console.log('cookie 长度:', cookie.length);
  console.log('cookie:', cookie);
  console.log('x-api-eid-token:', eidToken);
  console.log('h5st:', h5st);
  console.log('uuid:', uuid);
  console.log('==============================');

  // 顺手把结果也打到目标页面自己的 console，方便直接查看
  if (tabId != null) {
    chrome.scripting.executeScript({
      target: { tabId },
      func: (c, e, h, u) => {
        console.log('%c==== [扩展] 京东商品信息抓取成功 ====', 'color:#d32f2f;font-weight:bold');
        console.log('[扩展] cookie:', c);
        console.log('[扩展] x-api-eid-token:', e);
        console.log('[扩展] h5st:', h);
        console.log('[扩展] uuid:', u);
      },
      args: [cookie, eidToken, h5st, uuid]
    }).catch(err => console.warn('注入页面 console 失败:', err.message));
  }

  chrome.notifications.create({
    type: 'basic', iconUrl: 'icon48.png',
    title: '京东信息抓取成功',
    message: `cookie/${cookie.length} eid/${eidToken ? 'Y' : 'N'} h5st/${h5st ? 'Y' : 'N'} uuid/${uuid ? 'Y' : 'N'}\n${new Date().toLocaleString()}`,
    priority: 1
  });

  // 上报到后端
  sendToServer({ cookie, eidToken, h5st, uuid });
}

// 上报接口：application/x-www-form-urlencoded，带 X-Secret 头
const REPORT_ENDPOINT = 'http://ad-ocpx.atd.com/index.php?r=tool%2Fjd-material-token%2Fsave';
const REPORT_SECRET = 'b8e04f21a7c93d65f018e2b4c7a95d3e61f0a8c2d94b7e35';

function sendToServer(payload) {
  const body = new URLSearchParams({
    h5st: payload.h5st || '',
    eid_token: payload.eidToken || '',
    cookie: payload.cookie || '',
    uuid: payload.uuid || ''
  });

  fetch(REPORT_ENDPOINT, {
    method: 'POST',
    headers: {
      'X-Secret': REPORT_SECRET,
      'Content-Type': 'application/x-www-form-urlencoded'
    },
    body: body.toString()
  })
    .then(r => { if (!r.ok) throw new Error(`HTTP ${r.status}`); return r.text(); })
    .then(text => console.log('上报成功:', text))
    .catch(err => console.error('上报失败:', err.message));
}

// 监听来自 popup 的消息（手动触发）
chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  if (request.action === 'updateCookieNow') {
    updateCookieAutomatically();
    sendResponse({ status: 'started' });
  }
  return true;
});
