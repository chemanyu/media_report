// 定时任务：每 3 分钟把最近一次捕获的 Cookie 头重新上报一次
const UPDATE_INTERVAL = 3; // 分钟

// 复投接口的真实请求地址。京橙控制台刷新页面时会打这个接口，
// 我们用 webRequest 抓它请求头里的 Cookie 原文整条上报。
//
// 为什么不用 chrome.cookies 拼：
//   · domain: 'jd.com' —— Chrome 的 domain 过滤是后缀匹配，把所有 *.jd.com 子域的
//     host-only Cookie 混进来，__jda / __jdc / 3AB9D23F7A4B3C9B 出现多份，值取错。
//   · url: api.m.jd.com —— 范围对了，但仍然是「我们按规则重算一遍」，作用域、顺序、
//     值的编码都可能和浏览器实际发出去的那条有差异。
// 抓请求头 = 拿浏览器已经算好的最终结果，逐字节原样搬走，不做任何重组。
const TARGET_FILTER = { urls: ['https://api.m.jd.com/api/*'] };

// 只认京橙控制台发起的请求，避免抓到其它页面（可能是别的登录态）打的同一个接口
const ALLOWED_INITIATOR = 'https://jcheng.jd.com';

const REPORT_API = 'http://127.0.0.1:8888/update/jingcheng/futou/cookie';

// 关键 Cookie，缺了基本就是登录态失效，只告警不阻断（京东随时可能改名）
const REQUIRED_COOKIES = ['thor', 'pin', 'light_key', '3AB9D23F7A4B3C9B'];

// Cookie 没变时的最小重报间隔，防止一次页面刷新打十几个请求就刷爆后端
const MIN_REPOST_INTERVAL_MS = UPDATE_INTERVAL * 60 * 1000;

// 上报互斥：一次刷新会并发触发多个监听回调
let posting = false;

// ---------------------------------------------------------------- 捕获

chrome.webRequest.onBeforeSendHeaders.addListener(
  (details) => {
    // initiator 缺失时不拦（部分请求拿不到），有值就必须是京橙控制台
    if (details.initiator && !details.initiator.startsWith(ALLOWED_INITIATOR)) {
      return;
    }

    const cookieHeader = findCookieHeader(details.requestHeaders);
    if (!cookieHeader) return;

    captureCookieHeader(cookieHeader).catch(error => {
      console.error('处理捕获的Cookie失败:', error.message);
    });
  },
  TARGET_FILTER,
  // extraHeaders 必须加，否则 Chrome 不会把 Cookie 这类敏感头交给扩展
  ['requestHeaders', 'extraHeaders']
);

function findCookieHeader(headers) {
  if (!headers) return '';
  const header = headers.find(h => h.name.toLowerCase() === 'cookie');
  return header && header.value ? header.value : '';
}

// 存下最新的 Cookie 头；变了就立刻上报，没变则按 MIN_REPOST_INTERVAL_MS 节流
async function captureCookieHeader(cookieHeader) {
  const now = Date.now();
  const state = await chrome.storage.session.get(['cookieHeader', 'postedHeader', 'postedAt']);

  if (cookieHeader !== state.cookieHeader) {
    console.log(`捕获到新的Cookie头：${cookieHeader.length} 字符`);
  }
  await chrome.storage.session.set({ cookieHeader: cookieHeader, capturedAt: now });

  const changed = cookieHeader !== state.postedHeader;
  const stale = !state.postedAt || now - state.postedAt >= MIN_REPOST_INTERVAL_MS;
  if (changed || stale) {
    await reportCookieHeader(cookieHeader);
  }
}

// ---------------------------------------------------------------- 上报

// 原样上报，不解析、不重排、不重新编码；只额外数一下名字用于告警
async function reportCookieHeader(cookieHeader) {
  if (posting) {
    return { ok: false, message: '已有上报进行中，跳过本次' };
  }
  posting = true;

  try {
    const names = cookieNames(cookieHeader);
    const missing = REQUIRED_COOKIES.filter(name => !names.includes(name));
    if (missing.length > 0) {
      // 只告警不拦截：京东改 Cookie 名时不至于把整条链路卡死
      console.warn('缺少关键Cookie，可能已退出登录:', missing.join(', '));
    }

    console.log(`准备上报：${names.length} 个Cookie / ${cookieHeader.length} 字符`);
    console.log('所有Cookie名称:', names.join(', '));

    const response = await fetch(REPORT_API, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        cookie: cookieHeader,
        csrfToken: ''
      })
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const body = await response.text();
    await chrome.storage.session.set({ postedHeader: cookieHeader, postedAt: Date.now() });

    console.log('Cookie更新成功:', body);
    chrome.notifications.create({
      type: 'basic',
      iconUrl: 'icon48.png',
      title: '复投Cookie更新成功',
      message: `${names.length} 个Cookie / ${cookieHeader.length} 字符，${new Date().toLocaleString()}`,
      priority: 1
    });

    return { ok: true, count: names.length, length: cookieHeader.length, missing: missing, body: body };
  } catch (error) {
    console.error('发送Cookie失败:', error.message);
    chrome.notifications.create({
      type: 'basic',
      iconUrl: 'icon48.png',
      title: '复投Cookie更新失败',
      message: error.message,
      priority: 2
    });
    return { ok: false, message: error.message };
  } finally {
    posting = false;
  }
}

// 名字只用于「缺关键Cookie」告警和日志，不参与上报内容的拼装
function cookieNames(cookieHeader) {
  return cookieHeader
    .split(';')
    .map(part => part.trim().split('=')[0])
    .filter(Boolean);
}

// 定时重报：把已捕获的最新 Cookie 头再推一次，保持后端时间戳新鲜
async function updateCookieAutomatically() {
  const state = await chrome.storage.session.get(['cookieHeader', 'capturedAt']);

  if (!state.cookieHeader) {
    console.log('尚未捕获到Cookie头，请打开并刷新 https://jcheng.jd.com/');
    return;
  }

  const ageMin = Math.round((Date.now() - (state.capturedAt || 0)) / 60000);
  console.log(`重报已捕获的Cookie头（${ageMin} 分钟前捕获）`);
  await reportCookieHeader(state.cookieHeader);
}

// ---------------------------------------------------------------- 定时器与消息

chrome.alarms.onAlarm.addListener((alarm) => {
  if (alarm.name === 'updateCookie') {
    console.log('定时更新Cookie任务触发:', new Date().toLocaleString());
    updateCookieAutomatically();
  }
});

chrome.runtime.onInstalled.addListener(() => {
  console.log('Chrome扩展已安装/更新');
  chrome.alarms.create('updateCookie', { periodInMinutes: UPDATE_INTERVAL });
});

chrome.runtime.onStartup.addListener(() => {
  console.log('Chrome扩展启动');
  chrome.alarms.create('updateCookie', { periodInMinutes: UPDATE_INTERVAL });
});

// popup 的手动触发与状态查询
chrome.runtime.onMessage.addListener((request, _sender, sendResponse) => {
  if (request.action === 'updateCookieNow') {
    (async () => {
      const state = await chrome.storage.session.get(['cookieHeader']);
      if (!state.cookieHeader) {
        sendResponse({ ok: false, message: '尚未捕获到Cookie头，请打开并刷新 https://jcheng.jd.com/ 后重试' });
        return;
      }
      sendResponse(await reportCookieHeader(state.cookieHeader));
    })();
    return true;
  }

  if (request.action === 'getStatus') {
    (async () => {
      const state = await chrome.storage.session.get(['cookieHeader', 'capturedAt', 'postedAt']);
      sendResponse({
        captured: !!state.cookieHeader,
        count: state.cookieHeader ? cookieNames(state.cookieHeader).length : 0,
        length: state.cookieHeader ? state.cookieHeader.length : 0,
        capturedAt: state.capturedAt || 0,
        postedAt: state.postedAt || 0
      });
    })();
    return true;
  }

  return false;
});
