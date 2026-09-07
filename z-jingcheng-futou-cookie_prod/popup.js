// Cookie 采集顺序与优先级说明见 background.js（同名只取先出现的那份）
const COOKIE_SOURCES = [
  { url: 'https://api.m.jd.com/' },       // 复投接口本身，鉴权 Cookie 以这份为准
  { url: 'https://jcheng.jd.com/' },      // 京橙控制台 host-only：focus-*、me_saas_userInfo、switch_to_bpro
  { url: 'https://passport.jd.com/' },    // 登录态：pt_key、pt_pin、pt_token、pwdt_id、pt_st、sdtoken
  { url: 'https://www.jd.com/' },         // .jd.com 通用：__jdv、__jdb、visitkey、webp
  { domain: 'jd.com' },                   // 兜底：其余任意 *.jd.com 子域的 host-only Cookie
];

const REQUIRED_COOKIES = ['thor', 'pin', 'light_key', '3AB9D23F7A4B3C9B'];

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

    for (const cookie of cookies) {
      if (merged.has(cookie.name)) continue;
      merged.set(cookie.name, cookie.value);
    }
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
        showError('未找到 jd.com 的 Cookie，请先登录 https://jcheng.jd.com/ 并刷新后重试。');
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
