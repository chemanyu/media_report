# CLAUDE.md

本文件供 Claude Code 每次自动读取，记录 media_report 项目级规则与约定。
项目整体说明见 [PROJECT.md](PROJECT.md)，详细排查记录见 `docs/engineering/session-log.md`。

## 环境约束

- **`git diff` / `git status` 看不到刚做的改动 ≠ 没落盘**：本项目用户习惯在对话中途自行 `git commit`。收工或编译前发现工作区干净，先 `git log -3` 看最近提交是否已包含改动，不要误判为「没保存」。
  - 来源：session-log §2026-06-15、§2026-06-16（已第 3 次踩中）

- **`go build` 报 `version "go1.23.0" does not match go tool version "go1.25.0"`**：goenv 把 `GOROOT` 钉在 1.23.0，但 `go` 二进制是 1.25.0。临时绕过：
  ```bash
  GOROOT=/Users/chemanyu/.goenv/versions/1.25.0 GOTOOLCHAIN=local go build ./...
  ```
  - 来源：session-log §2026-06-15、§2026-06-16（已第 2 次踩中）

## 响应约定

- **大响应强制定长**：线上前置 nginx 不正确转发 chunked 响应，body 超过 ~2KB 会被截断（浏览器 axios 报 `Network Error`，服务端日志却是 200）。返回列表类大响应的接口须用 `service/api/internal/response.OkJsonCtx`（显式 Content-Length），而非 `httpx.OkJsonCtx`。
  - 来源：session-log §2026-06-16 / commit 4a6c68c

## 业务约定

- **飞猪媒体分两类接入方式**：① API 主动拉取（OPPO/小米/荣耀，logic 里有 `SyncXxxData`，按账户凭据调外部 API）；② 外部回传接收（ADN/华为，只有 `SaveXxxData` + 一个 POST handler，由外部 PHP 主动 POST 数据进来）。新增「外部回传型」媒体照搬 ADN 三件套即可：types `FzSyncXxxDataReq` + logic `SaveXxxData`（cost/成本字段 ×100 元转分）+ handler + 路由 + 前端 media 选项；**不要**套荣耀的 API 客户端逻辑。
  - 原因：两类形态不同，回传型不需要凭据/拉取，照搬错模板会做无用功。
  - 来源：session-log §2026-06-16 12:16 / commit 800e853（华为接口）

## 排查约定

- **「线上看不到刚改的功能」先怀疑未部署 / 页面缓存，而非代码漏改**：本项目前端是静态 HTML，改动只 commit 在本地，线上跑的是部署前的旧页面。遇到「某功能没生效」先 ① 核对代码是否真的已加（grep）② 确认部署状态（线上二进制/HTML 是否已更新）③ 让用户强刷浏览器（Ctrl+F5 清缓存），确认是部署/缓存问题后再考虑改代码。
  - 来源：session-log §2026-06-16 12:16（华为下拉「线上没有」实为页面缓存）
