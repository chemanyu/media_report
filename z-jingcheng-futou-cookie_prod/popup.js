// 京橙复投接口的真实请求地址，Cookie 按这个 URL 取，理由见 background.js
const COOKIE_TARGET_URL = 'https://api.m.jd.com/api/';

const REQUIRED_COOKIES = ['thor', 'pin', 'light_key', '3AB9D23F7A4B3C9B'];

// 同 background.js：Chrome 按 RFC 6265 顺序返回，同名只留最贴近目标 URL 的那份
async function collectJdCookies() {
  const merged = new Map();

  const cookies = await chrome.cookies.getAll({ url: COOKIE_TARGET_URL });
  for (const cookie of cookies) {
    if (merged.has(cookie.name)) continue;
    merged.set(cookie.name, cookie.value);
  }

  return merged;
}

function buildCookieHeader(merged) {
  return Array.from(merged, ([name, value]) => `${name}=${value}`).join('; ');
}

document.addEventListener('DOMContentLoaded', function() {
  const sendButton = document.getElementById('sendButton');
  const statusDiv = document.getElementById('status');
  const resultDiv = document.getElementById('result');

  sendButton.addEventListener('click', async function() {
    sendButton.disabled = true;
    sendButton.textContent = 'Sending...';
    statusDiv.textContent = 'Fetching cookies...';

    try {
      const merged = await collectJdCookies();

      if (merged.size === 0) {
        showError('未找到Cookie，请先登录 https://jcheng.jd.com/ 并刷新后重试。');
        return;
      }

      const cookieString = buildCookieHeader(merged);
      const csrfToken = merged.get('csrftoken') || '';

      const missing = REQUIRED_COOKIES.filter(name => !merged.has(name));
      if (missing.length > 0) {
        console.warn('缺少关键Cookie，可能已退出登录:', missing.join(', '));
      }

      console.log(`Cookie获取成功，共 ${merged.size} 个 / ${cookieString.length} 字符`);
      console.log('所有Cookie名称:', Array.from(merged.keys()).join(', '));

      statusDiv.textContent = `已取到 ${merged.size} 个Cookie / ${cookieString.length} 字符，发送中...`;

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
      statusDiv.textContent = `Success! ${merged.size} 个Cookie / ${cookieString.length} 字符`;
      resultDiv.textContent = body;
      resultDiv.style.backgroundColor = '#d4edda';
      resultDiv.style.borderColor = '#c3e6cb';
      if (missing.length > 0) {
        resultDiv.textContent += `\n\n注意：缺少关键Cookie ${missing.join(', ')}，可能已退出登录。`;
      }
    } catch (error) {
      showError('Error sending request: ' + error.message);
    } finally {
      sendButton.disabled = false;
      sendButton.textContent = '更新Cookie';
    }
  });

  function showError(message) {
    statusDiv.textContent = 'Error:';
    resultDiv.textContent = message;
    resultDiv.style.backgroundColor = '#f8d7da';
    resultDiv.style.borderColor = '#f5c6cb';
  }
});
