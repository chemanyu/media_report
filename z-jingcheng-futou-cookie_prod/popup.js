document.addEventListener('DOMContentLoaded', function() {
  const sendButton = document.getElementById('sendButton');
  const statusDiv = document.getElementById('status');
  const resultDiv = document.getElementById('result');
  const captureInfo = document.getElementById('captureInfo');

  // Cookie 头由 background 的 webRequest 监听器捕获，popup 只负责展示和手动触发
  function showStatus() {
    chrome.runtime.sendMessage({ action: 'getStatus' }, function(state) {
      if (!state) {
        captureInfo.textContent = '后台未响应，请重新加载扩展';
        return;
      }

      if (!state.captured) {
        captureInfo.textContent = '尚未捕获，请打开并刷新 https://jcheng.jd.com/';
        return;
      }

      captureInfo.textContent =
        `${state.count} 个Cookie / ${state.length} 字符，` +
        `${describeAge(state.capturedAt)}捕获` +
        (state.postedAt ? `，${describeAge(state.postedAt)}上报` : '，尚未上报');
    });
  }

  function describeAge(ts) {
    if (!ts) return '未知时间';
    const min = Math.round((Date.now() - ts) / 60000);
    if (min < 1) return '刚刚';
    if (min < 60) return `${min} 分钟前`;
    return `${Math.round(min / 60)} 小时前`;
  }

  showStatus();

  sendButton.addEventListener('click', function() {
    sendButton.disabled = true;
    sendButton.textContent = 'Sending...';
    statusDiv.textContent = '上报中...';

    chrome.runtime.sendMessage({ action: 'updateCookieNow' }, function(result) {
      sendButton.disabled = false;
      sendButton.textContent = '更新Cookie';

      if (!result) {
        showError('后台未响应，请重新加载扩展');
        return;
      }

      if (!result.ok) {
        showError(result.message || '上报失败');
        return;
      }

      statusDiv.textContent = `Success! ${result.count} 个Cookie / ${result.length} 字符`;
      resultDiv.textContent = result.body || '';
      resultDiv.style.backgroundColor = '#d4edda';
      resultDiv.style.borderColor = '#c3e6cb';
      if (result.missing && result.missing.length > 0) {
        resultDiv.textContent += `\n\n注意：缺少关键Cookie ${result.missing.join(', ')}，可能已退出登录。`;
      }
      showStatus();
    });
  });

  function showError(message) {
    statusDiv.textContent = 'Error:';
    resultDiv.textContent = message;
    resultDiv.style.backgroundColor = '#f8d7da';
    resultDiv.style.borderColor = '#f5c6cb';
  }
});
