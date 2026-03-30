document.addEventListener('DOMContentLoaded', function() {
  const sendButton = document.getElementById('sendButton');
  const statusDiv = document.getElementById('status');
  const resultDiv = document.getElementById('result');
  const tabSelect = document.getElementById('tabSelect');
  const refreshButton = document.getElementById('refreshButton');
  const savedDomainInfo = document.getElementById('savedDomainInfo');
  const savedDomainText = document.getElementById('savedDomainText');

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

  // 初始显示保存的域名
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
        // 显示标题，最多显示50个字符
        let title = tab.title || '无标题';
        if (title.length > 50) {
          title = title.substring(0, 47) + '...';
        }
        option.textContent = `[${index + 1}] ${title}`;
        // 默认选中当前活动的标签页
        if (tab.active && tab.windowId === chrome.windows.WINDOW_ID_CURRENT) {
          option.selected = true;
        }
        tabSelect.appendChild(option);
      });
    });
  }

  // 初始加载标签页列表
  loadTabs();

  // 刷新按钮点击事件
  refreshButton.addEventListener('click', function() {
    loadTabs();
    statusDiv.textContent = '标签页列表已刷新';
    resultDiv.textContent = '';
    setTimeout(() => {
      statusDiv.textContent = '';
    }, 2000);
  });

  sendButton.addEventListener('click', function() {
    const selectedTabId = parseInt(tabSelect.value);
    if (!selectedTabId) {
      showError('请选择一个标签页');
      return;
    }

    sendButton.disabled = true;
    sendButton.textContent = 'Sending...';
    statusDiv.textContent = 'Fetching selected tab...';

    // 查找选中的标签页
    const selectedTab = allTabs.find(tab => tab.id === selectedTabId);
    if (!selectedTab) {
      showError('所选标签页不存在，请刷新标签页列表');
      return;
    }

    try {
      const currentUrl = new URL(selectedTab.url);
      const domain = currentUrl.hostname;
      const topLevelDomain = domain.split('.').slice(-2).join('.');

      const doFetchCookies = () => {
        statusDiv.textContent = 'Fetching cookies...';
        chrome.cookies.getAll({ url: 'https://dhh.taobao.com/polystar/api/creative/material/forminfo' }, function(cookies) {
          if (chrome.runtime.lastError) {
            showError('获取Cookie异常: ' + chrome.runtime.lastError.message);
            return;
          }
          if (cookies.length === 0) {
            showError('当前页面未找到Cookie.');
            return;
          }

          const excludeKeys = ['_uab_collina', '__wpkreporterwid_'];
          const cookieString = cookies
            .filter(cookie => !excludeKeys.includes(cookie.name))
            .map(cookie => `${cookie.name}=${cookie.value}`)
            .join('; ');

          const xsrfToken = (cookies.find(cookie => cookie.name === 'XSRF-TOKEN') || {}).value || '';

          fetch('http://127.0.0.1:8888/update/dhh/cookie', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ cookie: cookieString, csrfToken: xsrfToken })
          })
          .then(response => {
            if (!response.ok) {
              showError('HTTP error! status.');
              throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.text();
          })
          .then(body => {
            statusDiv.textContent = 'Success!';
            resultDiv.textContent = body;
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
          })
          .catch(error => {
            showError('Error sending request: ' + error.message);
          })
          .finally(() => {
            sendButton.disabled = false;
            sendButton.textContent = '更新Cookie';
          });
        });
      };

      // 先刷新标签页，等加载完成后再获取 cookie
      statusDiv.textContent = '刷新页面中...';
      chrome.tabs.reload(selectedTab.id, function() {
        setTimeout(() => {
          const onUpdated = (tabId, changeInfo) => {
            if (tabId === selectedTab.id && changeInfo.status === 'complete') {
              chrome.tabs.onUpdated.removeListener(onUpdated);
              doFetchCookies();
            }
          };
          chrome.tabs.onUpdated.addListener(onUpdated);
        }, 2000); // 2000 毫秒 = 2 秒
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
    sendButton.textContent = 'Send Cookies';
  }
});