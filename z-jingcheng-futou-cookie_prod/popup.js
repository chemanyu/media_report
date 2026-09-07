// 京橙复投接口实际请求的域名，Cookie 必须按这个 URL 取，理由见 background.js
const COOKIE_TARGET_URL = 'https://api.m.jd.com/';

// 同 background.js：同名 Cookie 只留最贴近目标 URL 的那一份
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

document.addEventListener('DOMContentLoaded', function() {
  const sendButton = document.getElementById('sendButton');
  const statusDiv = document.getElementById('status');
  const resultDiv = document.getElementById('result');
  const cookieSourceText = document.getElementById('cookieSourceText');

  cookieSourceText.textContent = COOKIE_TARGET_URL;

  sendButton.addEventListener('click', function() {
    sendButton.disabled = true;
    sendButton.textContent = 'Sending...';
    statusDiv.textContent = 'Fetching cookies...';

    // 按接口域名取 Cookie，与当前打开哪个标签页无关
    chrome.cookies.getAll({ url: COOKIE_TARGET_URL }, function(cookies) {
      if (chrome.runtime.lastError) {
        showError('获取Cookie异常: ' + chrome.runtime.lastError.message);
        return;
      }
      if (cookies.length === 0) {
        showError('未找到Cookie，请先登录京橙页面（https://jcheng.jd.com/）并刷新后重试。');
        return;
      }

      // Format cookies as a semicolon-separated string (standard cookie format)
      const cookieString = buildCookieHeader(cookies);

      // 提取 x-csrftoken（查找名为 csrftoken 的 cookie）
      const csrfCookie = cookies.find(cookie => cookie.name === 'csrftoken');
      const csrfToken = csrfCookie ? csrfCookie.value : '';

      console.log('Cookie获取成功，准备发送到后端...');
      console.log('所有Cookie名称:', cookies.map(c => c.name).join(', '));
      if (csrfToken) {
        console.log('X-CSRF-Token已找到');
      } else {
        console.log('未找到X-CSRF-Token');
      }

      // 发送到后端API
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
        return response.text();  // Get the body as text; use .json() if the response is JSON
      })
      .then(body => {
        statusDiv.textContent = 'Success!';
        resultDiv.textContent = body;  // Display the response body
        resultDiv.style.backgroundColor = '#d4edda';
        resultDiv.style.borderColor = '#c3e6cb';
      })
      .catch(error => {
        showError('Error sending request: ' + error.message);
      })
      .finally(() => {
        sendButton.disabled = false;
        sendButton.textContent = '更新Cookie';
      });
    });
  });

  function showError(message) {
    statusDiv.textContent = 'Error:';
    resultDiv.textContent = message;
    resultDiv.style.backgroundColor = '#f8d7da';
    resultDiv.style.borderColor = '#f5c6cb';
    sendButton.disabled = false;
    sendButton.textContent = '更新Cookie';
  }
});
