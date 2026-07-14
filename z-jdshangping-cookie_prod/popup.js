document.addEventListener('DOMContentLoaded', function() {
  const sendButton = document.getElementById('sendButton');
  const statusDiv = document.getElementById('status');
  const resultDiv = document.getElementById('result');
  const tabSelect = document.getElementById('tabSelect');
  const refreshButton = document.getElementById('refreshButton');
  const savedDomainInfo = document.getElementById('savedDomainInfo');
  const savedDomainText = document.getElementById('savedDomainText');

  const TARGET_URL_PATTERN = 'api.m.jd.com/api';

  // 只要 URL 命中目标接口且带 h5st + x-api-eid-token 就算目标请求（functionId 不限）
  function isTargetRequest(url) {
    return url.includes(TARGET_URL_PATTERN) &&
           url.includes('h5st=') &&
           url.includes('x-api-eid-token=');
  }

  let allTabs = [];

  // 显示当前保存的域名
  function showSavedDomain() {
    chrome.storage.local.get(['savedDomain', 'savedTitle'], function(result) {
      if (result.savedDomain) {
        savedDomainText.textContent = result.savedDomain + (result.savedTitle ? ` (${result.savedTitle})` : '');
        savedDomainInfo.style.display = 'block';
      } else {
        savedDomainInfo.style.display = 'none';
      }
    });
  }
  showSavedDomain();

  // 加载标签页列表
  function loadTabs() {
    chrome.tabs.query({}, function(tabs) {
      allTabs = tabs;
      tabSelect.innerHTML = '';

      if (tabs.length === 0) {
        tabSelect.innerHTML = '<option value="">没有打开的标签页</option>';
        return;
      }

      tabs.forEach((tab, index) => {
        const option = document.createElement('option');
        option.value = tab.id;
        let title = tab.title || '无标题';
        if (title.length > 50) {
          title = title.substring(0, 47) + '...';
        }
        option.textContent = `[${index + 1}] ${title}`;
        if (tab.active && tab.windowId === chrome.windows.WINDOW_ID_CURRENT) {
          option.selected = true;
        }
        tabSelect.appendChild(option);
      });
    });
  }
  loadTabs();

  refreshButton.addEventListener('click', function() {
    loadTabs();
    statusDiv.textContent = '标签页列表已刷新';
    resultDiv.textContent = '';
    setTimeout(() => { statusDiv.textContent = ''; }, 2000);
  });

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
    } catch (e) {}
    return { eidToken, h5st, uuid };
  }

  // 上报接口
  const REPORT_ENDPOINT = 'http://ad-ocpx.atd.com/index.php?r=tool%2Fjd-material-token%2Fsave';
  const REPORT_SECRET = 'b8e04f21a7c93d65f018e2b4c7a95d3e61f0a8c2d94b7e35';

  function sendToServer(payload) {
    const body = new URLSearchParams({
      h5st: payload.h5st || '',
      eid_token: payload.eidToken || '',
      cookie: payload.cookie || '',
      uuid: payload.uuid || ''
    });
    return fetch(REPORT_ENDPOINT, {
      method: 'POST',
      headers: {
        'X-Secret': REPORT_SECRET,
        'Content-Type': 'application/x-www-form-urlencoded'
      },
      body: body.toString()
    })
      .then(r => { if (!r.ok) throw new Error(`HTTP ${r.status}`); return r.text(); });
  }

  sendButton.addEventListener('click', function() {
    const selectedTabId = parseInt(tabSelect.value);
    if (!selectedTabId) {
      showError('请选择一个标签页');
      return;
    }

    sendButton.disabled = true;
    sendButton.textContent = '抓取中...';
    statusDiv.textContent = '刷新页面中，等待接口请求...';

    const selectedTab = allTabs.find(tab => tab.id === selectedTabId);
    if (!selectedTab) {
      showError('所选标签页不存在，请刷新标签页列表');
      return;
    }

    try {
      const currentUrl = new URL(selectedTab.url);
      const topLevelDomain = currentUrl.hostname.split('.').slice(-2).join('.');
      const target = { tabId: selectedTab.id };

      chrome.debugger.detach(target, () => {
        void chrome.runtime.lastError;
        chrome.debugger.attach(target, '1.3', () => {
          if (chrome.runtime.lastError) {
            showError('attach debugger 失败: ' + chrome.runtime.lastError.message);
            return;
          }

          const pending = new Map(); // requestId -> { eidToken, h5st }

          chrome.debugger.sendCommand(target, 'Network.enable', {}, () => {
            const onEvent = (source, method, params) => {
              if (source.tabId !== selectedTab.id) return;

              if (method === 'Network.requestWillBeSent') {
                const url = params.request.url || '';
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
                showError('目标接口请求头中未找到 Cookie');
                return;
              }

              chrome.debugger.onEvent.removeListener(onEvent);
              chrome.debugger.detach(target, () => {});

              // 打印到 console
              console.log('==== 京东商品信息抓取成功 ====');
              console.log('cookie:', cookie);
              console.log('x-api-eid-token:', tokens.eidToken);
              console.log('h5st:', tokens.h5st);
              console.log('uuid:', tokens.uuid);

              statusDiv.textContent = '抓取成功，上报中...';
              resultDiv.textContent =
                `cookie(${cookie.length}):\n${cookie}\n\n` +
                `x-api-eid-token:\n${tokens.eidToken}\n\n` +
                `h5st:\n${tokens.h5st}\n\n` +
                `uuid:\n${tokens.uuid}`;
              resultDiv.style.backgroundColor = '#d4edda';
              resultDiv.style.borderColor = '#c3e6cb';

              chrome.storage.local.set({
                'savedDomain': topLevelDomain,
                'savedUrl': selectedTab.url,
                'savedTitle': selectedTab.title
              }, function() {
                console.log('已保存域名:', topLevelDomain);
                showSavedDomain();
              });

              sendToServer({ cookie, eidToken: tokens.eidToken, h5st: tokens.h5st, uuid: tokens.uuid })
                .then(text => { statusDiv.textContent = '上报成功: ' + text; })
                .catch(err => { statusDiv.textContent = '抓取成功，但上报失败: ' + err.message; });

              sendButton.disabled = false;
              sendButton.textContent = '立即抓取';
            };

            chrome.debugger.onEvent.addListener(onEvent);
            chrome.tabs.reload(selectedTab.id);

            // 超时保护：30秒
            setTimeout(() => {
              chrome.debugger.onEvent.removeListener(onEvent);
              chrome.debugger.detach(target, () => {});
              showError('超时未捕获到目标接口，请确认页面为京东联盟商品页');
              sendButton.disabled = false;
              sendButton.textContent = '立即抓取';
            }, 30000);
          });
        });
      });
    } catch (error) {
      showError('URL解析错误: ' + error.message);
    }
  });

  function showError(message) {
    statusDiv.textContent = 'Error:';
    resultDiv.textContent = message;
    resultDiv.style.backgroundColor = '#f8d7da';
    resultDiv.style.borderColor = '#f5c6cb';
    sendButton.disabled = false;
    sendButton.textContent = '立即抓取';
  }
});
