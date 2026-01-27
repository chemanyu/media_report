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
      statusDiv.textContent = 'Fetching cookies...';

      // Get all cookies for the selected URL
      chrome.cookies.getAll({ domain: topLevelDomain }, function(cookies) {
        if (chrome.runtime.lastError) {
          showError('获取Cookie异常: ' + chrome.runtime.lastError.message);
          return;
        }
        // showError('123123==' + cookies);
        if (cookies.length === 0) {
          showError('当前页面未找到Cookie.');
          return;
        }

        // Format cookies as a semicolon-separated string (standard cookie format)
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
            showError('HTTP error! status.');
            throw new Error(`HTTP error! status: ${response.status}`);
          }
          return response.text();  // Get the body as text; use .json() if the response is JSON
        })
        .then(body => {
          statusDiv.textContent = 'Success!';
          resultDiv.textContent = body;  // Display the response body
          resultDiv.style.backgroundColor = '#d4edda';
          resultDiv.style.borderColor = '#c3e6cb';

          // 保存用户选择的域名，供定时任务使用
          chrome.storage.local.set({ 
            'savedDomain': topLevelDomain,
            'savedUrl': selectedTab.url,
            'savedTitle': selectedTab.title 
          }, function() {
            console.log('已保存域名:', topLevelDomain);
            showSavedDomain(); // 更新显示
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