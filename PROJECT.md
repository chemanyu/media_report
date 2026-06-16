# Media Report 媒体报表服务

> 面向广告代理公司内部使用的媒体广告数据报表管理平台，覆盖多家广告媒体的数据抓取、报表生成、钉钉推送全流程自动化。

---

## 技术栈

| 层次 | 技术 |
|---|---|
| 后端框架 | go-zero v1.9.2（HTTP REST） |
| 编程语言 | Go 1.24 |
| 数据库 | MySQL + GORM |
| 定时任务 | robfig/cron v3 |
| Excel 生成 | excelize/v2 |
| 邮件发送 | gomail.v2 |
| 转链脚本 | Python 3 + Selenium + ChromeDriver |
| 前端 | Vue 3 + Element Plus（CDN，纯静态 HTML） |
| 浏览器插件 | Chrome Extension Manifest V3 × 4 套 |
| 部署 | Docker 多阶段构建 |

---

## 项目目录结构

```
media_report/
├── api/                          # go-zero API DSL 定义
│   └── media.api
├── common/                       # 公共基础模块
│   ├── database/mysql.go         # MySQL 连接池封装
│   └── httpclient/client.go      # 通用 HTTP 客户端
├── service/api/
│   ├── media.go                  # 服务入口（main）
│   ├── etc/media-api.yaml        # 全量配置文件
│   └── internal/
│       ├── config/               # 配置结构体
│       ├── handler/              # HTTP 处理器层（请求解析）
│       │   ├── routes.go         # 路由注册
│       │   ├── config/           # 配置管理接口
│       │   ├── download/         # 文件下载
│       │   ├── fz/               # 飞猪报表
│       │   ├── report/           # 快手/巨量/京橙报表
│       │   ├── tanx/             # 淘宝联盟
│       │   ├── ulink/            # 转链工具
│       │   ├── update/           # Cookie 更新
│       │   └── zfb/              # 支付宝下载
│       ├── logic/                # 业务逻辑层（核心实现）
│       ├── model/                # 数据库模型（GORM 实体）
│       ├── script/               # 定时任务调度器
│       ├── oppo/                 # OPPO 广告 API 客户端
│       ├── xiaomi/               # 小米营销 API 客户端
│       ├── svc/                  # 服务上下文（依赖注入）
│       └── types/                # 类型定义
├── scripts/ulink/                # Python 转链脚本
├── web/                          # 前端静态页面
├── z-juliang-cookie_prod/        # Chrome 插件：巨量 Cookie 上传
├── z-jingcheng-cookie_prod/      # Chrome 插件：京橙 Cookie 上传
├── z-tanx-cookie_prod/           # Chrome 插件：淘宝联盟 Cookie 上传
└── z-zfb-cookie_prod/            # Chrome 插件：支付宝 Cookie 上传
```

---

## 数据库模型

数据库名：`release_atd`

| 表名 | 说明 |
|---|---|
| `media_token` | 各平台 Cookie / Token 存储（字段：media, token, refresh_token, agent_id, advertiser_id）|
| `cainiao_advertiser` | 快手菜鸟广告主 ID 列表 |
| `cainiao_cardinality` | 菜鸟消耗倍率系数配置 |
| `rebate` | 返点配置（主体、端口、返点率、主体类型）|
| `service_fee` | 服务费配置（服务商、服务费率）|
| `task_type` | 任务类型配置（名称、编码、结算单价、适用媒体）|
| `fz_hourly_report` | 飞猪时报数据（OPPO / 小米 / ADN，按 media+adv_id+date 唯一）|
| `fz_media_advertiser` | 飞猪项目下的 OPPO / 小米 / ADN 媒体账户 |
| `elm_hc_performance_report` | 汇川饿了么绩效报表（客户-代理商-媒体账户映射）|
| `tanx_monitor` | 淘宝联盟广告位监控数据（按 pid+ds 唯一）|

**`media_token` 的 media 字段取值**：

| 值 | 说明 |
|---|---|
| `kuaishou` | 快手广告 OAuth Token |
| `juliang_pachong` | 巨量引擎业务后台 Cookie（插件维护）|
| `tanx_pachong` | 淘宝联盟后台 Cookie |
| `jingcheng_pachong` | 京橙后台 Cookie |
| `zfb_pachong` | 支付宝 Cookie |

---

## 功能模块

### 1. 快手（菜鸟）广告报表

**目标**：拉取菜鸟客户在快手的广告消耗数据，推送钉钉时报。

**执行时间**：每天 10 / 12 / 14 / 16 / 18 / 22 时

**流程**：
1. 每天 0:01 用 Refresh Token 刷新快手 Access Token，写回 `media_token`
2. 遍历 `cainiao_advertiser` 表所有广告主，并发调用快手 API `/rest/openapi/v1/report/account_report`
3. 读取 `cainiao_cardinality` 基数系数，换算实际计费消耗
4. 推送 Markdown 格式消息到"美数时报"钉钉群

**相关接口**：
- `POST /report/ks/account` — 手动触发账户报表查询

---

### 2. 巨量引擎（字节）广告报表

**目标**：汇总巨量代理商后台所有账户当日消耗，计算返后消耗、服务费、归因扣量、预估利润，生成 Excel 并推送钉钉。

**执行时间**：每天 10 / 12 / 14 / 16 / 18 / 22 时（每次 +3 分）

**流程**：

1. **Cookie 维护**：用 `z-juliang-cookie_prod` 插件提取巨量业务后台 Cookie，POST 到 `/update/juliang/cookie`，写入 `media_token`

2. **数据抓取**：
   - 调用巨量业务 API `https://business.oceanengine.com/nbs/api/bm/promotion/ad/get_account_list`
   - 分页 + 最大 10 并发拉取所有账户数据
   - 按备注格式 `主体-端口-服务商-任务`（如 `美数-优居-正好-APP`）解析账户归属

3. **多维指标计算**：

   | 指标 | 计算方式 |
   |---|---|
   | 返后消耗 | 消耗 / (1 + 返点率) |
   | 服务费成本 | 消耗 × 服务费率 |
   | 归因扣量数 | 从内部归因系统 `ad-ocpx.zhltech.net` 获取 |
   | 预估收入 | (转化数 + 扣量数) × 结算单价 |
   | 预估利润 | (预估收入 × 0.95) - 服务费成本 - 返后消耗 |
   | 预估利润率 | 预估利润 / 预估收入 |

4. **输出**：用 excelize 生成 20 列 Excel，保存到 `./download_files/`；将汇总数据和下载链接推送到"巨量时报"钉钉群

**相关接口**：
- `POST /report/juliang/trigger` — 手动触发巨量报表任务

---

### 3. 汇川饿了么报表

**目标**：通过巨量开放平台 OAuth API，拉取汇川饿了么客户的投放数据，推送到 ADX 外部接口。

**执行时间**：每天 11:00（日报）；每小时第 2 分（小时报）

**流程**：
- 每天 0:01 用 App Key/Secret 自动刷新巨量开放平台 Access Token
- 从 `elm_hc_performance_report` 表读取客户-媒体账户映射，逐账户拉取数据
- 单账户遇巨量限频（code=40110）时按退避表自动重试（时报预算 45min / 日报 40min），超预算记日志放弃，避免与下次任务重叠
- 调用 ADX 外部接口 `https://agent.kkforce.com/assistant-external` 写入数据

**相关接口**：
- `GET /api/elm_hc/report/list` — 查询报表数据

---

### 4. 飞猪外投数据报表（OPPO + 小米 + ADN）

**目标**：同步 OPPO、小米广告平台的飞猪 App 拉活数据，推送钉钉时报。

**执行时间**：每天 10 / 14 / 18 / 20 时；每天 23:58 更新当日数据

**OPPO API**：SHA1 签名（owner_id + api_id + timestamp），调用 `https://sapi.ads.heytapmobi.com/v3/data/common/summary/queryAdData`

**小米 API**：MD5 签名（参数 Key 排序拼接 + SecretKey），调用 `https://api.e.mi.com/openapi/v5/report/getData`

**数据写入**：Upsert 到 `fz_hourly_report`（幂等，按 media+adv_id+date 唯一键）

**钉钉推送**：区分"常规活动"和"集能量活动"，分别计算拉活成本、订单成本

**相关接口**：
- `GET /api/fz/sync_all_data` — 手动触发 OPPO+小米 数据同步
- `POST /api/fz/sync_adn_data` — 同步 ADN 媒体数据
- `GET /api/fz_hourly_report/list` — 查询时报数据列表
- `GET /api/fz_advertiser/list` / `POST` / `PUT` / `DELETE` — 账户 CRUD

---

### 5. 淘宝联盟（Tanx）监控

**目标**：爬取淘宝联盟后台 36 个广告位的监控数据（曝光、点击、预估收益等），存库并发邮件。

**执行时间**：每天 12:15

**流程**：
- Cookie 通过 `z-tanx-cookie_prod` 插件维护
- 逐 PID 调用 `https://tanx.alimama.com/api/media/debug/report/getReport.htm`
- Upsert 到 `tanx_monitor` 表（按 pid+ds 唯一）
- 通过 SMTP（Outlook）发送邮件到配置的接收人列表

**相关接口**：
- `POST /tanx/fetch_data` — 手动触发数据抓取
- `POST /tanx/export_data` — 手动触发导出并发邮件
- `POST /tanx/update_cookie` — 更新 Cookie

---

### 6. 京东广义转化归因数据

**目标**：查询内部归因系统的京东广义转化扣量数据，按账户和错误类型统计，支持多日导出。

**相关接口**：
- `GET /api/jd/attribution/data?date=YYYYMMDD` — 查询指定日期数据
- `GET /api/jd/attribution/export?num_days=10` — 导出近 N 天数据（Excel）

---

### 7. 支付宝（ZFB）文件下载

**目标**：代理下载支付宝广告报表文件（通过 Cookie 鉴权）。

- Cookie 通过 `z-zfb-cookie_prod` 插件维护
- 通过 Nginx `/zfb` 路径隔离，无需系统密码验证

**相关接口**：
- `GET /zfb/download?uid=xxx&start_date=xxx&end_date=xxx`

---

### 8. 转链工具（Ulink）

**目标**：将淘宝/闲鱼短链批量转换为 App Deeplink，用于广告投放落地页。

**技术实现**：Go 通过 `exec.Command` 调用 Python + Selenium 脚本，模拟手机端访问拦截 Deeplink；批量模式使用 3 个浏览器驱动池 + 线程池并发处理，结果写入 Excel 流式返回。

**子功能**：

| 功能 | 接口 | 说明 |
|---|---|---|
| 淘宝单链转链 | `POST /api/ulink/taobao/extract` | 单条短链 → Deeplink + H5 Deeplink |
| 淘宝批量转链 | `POST /api/ulink/taobao/batch` | 上传 TXT → 下载 Excel |
| 闲鱼单链转链 | `POST /api/ulink/xianyu/extract` | 单条短链 → Deeplink |
| 闲鱼批量转链 | `POST /api/ulink/xianyu/batch` | 上传 TXT → 下载 Excel |
| 活动短链批量 | `POST /api/ulink/activity/batch` | 淘宝客活动短链批量获取 |
| 活动报表查询 | `POST /api/ulink/activity/report` | 上传 Excel → 查询报表 → 下载 Excel |
| CPS 商品库 | `POST /api/ulink/cps/goods` | 高佣商品库批量导出 |

---

### 9. 配置管理

通过 Web 界面管理驱动核心计算逻辑的配置项：

| 配置 | 接口前缀 | 说明 |
|---|---|---|
| 返点配置 | `/config/rebate/*` | 主体-端口-返点率（用于巨量报表）|
| 服务费配置 | `/config/servicefee/*` | 服务商-服务费率 |
| 任务类型配置 | `/config/tasktype/*` | 任务-结算单价（用于预估收入）|
| 汇川饿了么报表 | `/config/elmhc/*` | 客户-代理商-媒体账户映射 |
| 菜鸟广告账户 | `/api/cainiao_advertiser/*` | 快手广告主 ID 管理 |
| 飞猪媒体账户 | `/api/fz_advertiser/*` | OPPO/小米/ADN 账户管理 |

---

### 10. Chrome 插件（Cookie 自动上传）

4 套 Chrome 扩展（Manifest V3），解决广告后台 Cookie 定期失效问题：

| 插件 | 平台 | 更新接口 |
|---|---|---|
| `z-juliang-cookie_prod` | 巨量引擎业务后台 | `POST /update/juliang/cookie` |
| `z-jingcheng-cookie_prod` | 京橙后台 | `POST /update/jingcheng/cookie` |
| `z-tanx-cookie_prod` | 淘宝联盟后台 | `POST /tanx/update_cookie` |
| `z-zfb-cookie_prod` | 支付宝 | `POST /update/zfb/cookie` |

**工作原理**：用户在浏览器登录对应平台 → 点击插件图标选择目标标签页 → 插件调用 `chrome.cookies.getAll()` → POST 到本地 `127.0.0.1:8888` → 写入 `media_token` 表。

---

## 定时任务调度

所有任务在 `service/api/internal/script/report_scheduler.go` 注册：

| 任务 | Cron 表达式 | 说明 |
|---|---|---|
| Token 刷新 | `1 0 * * *` | 刷新快手 + 巨量 DLS Access Token |
| 巨量报表 | `3 10,12,14,16,18,22 * * *` | 抓取数据、生成 Excel、推送钉钉 |
| 汇川饿了么日报 | `0 11 * * *` | 拉取日报，推送 ADX |
| 汇川饿了么小时报 | `2 * * * *` | 拉取小时报，推送 ADX |
| 淘宝联盟 Tanx | `15 12 * * *` | 抓取广告位数据，发邮件 |
| 飞猪时报 | `1 10,14,18,20 * * *` | 同步 OPPO+小米，推送钉钉 |
| 飞猪日报更新 | `58 23 * * *` | 更新当日飞猪全量数据 |

---

## 前端界面

Go 后端直接伺服静态 HTML 文件（`/web/:path` 路由），无需独立前端服务。

**主界面**（`web/index.html`）：左侧固定导航栏 + iframe 内容区，支持生产环境反向代理路径前缀自动适配。

**导航结构**：
```
巨量报表
  ├── 返点配置
  ├── 服务费配置
  └── 任务类型配置
汇川报表
  └── 汇川饿了么报表
菜鸟时报账户
京东广义数据统计
飞猪报表
  ├── 账户管理
  └── 数据报表
转链
  ├── 淘宝转链
  ├── 闲鱼转链
  ├── 活动转链
  └── CPS高佣商品库
```

---

## 部署

```bash
# 本地运行
go mod tidy
cd service/api && go run media.go -f etc/media-api.yaml

# Docker 构建与运行
make docker-build
make docker-run    # 监听 8888 端口
```

**生产环境**：`https://rta.zhltech.net/guangyixinmedia`，Nginx 反向代理。

---

## 关键设计

1. **Cookie 驱动的爬虫架构**：对无开放 API 的平台，Chrome 插件辅助获取 Cookie 存库，定时任务复用。
2. **Go + Python 混合**：浏览器模拟场景（Deeplink 提取）由 Go 调用 Python Selenium 脚本，结果通过 stdout JSON 传回。
3. **并发控制**：巨量分页拉取限制 10 并发（semaphore）；Ulink 批量使用 3 个浏览器驱动池。
4. **幂等写入**：GORM `clause.OnConflict` 实现 Upsert，多次同步相同日期不重复。
5. **多钉钉通道**：4 个独立 Webhook，对应菜鸟时报、巨量时报、飞猪时报（小时/日报）。
6. **go-zero 分层**：Handler（请求解析）→ Logic（业务逻辑）→ Model（数据库）严格三层分离。
7. **大响应强制定长**：前置 nginx 不正确转发 chunked 响应，body 超过 ~2KB 会被截断（表现为浏览器 axios `Network Error`、服务端却是 200）。返回列表类大响应的接口须用 `internal/response.OkJsonCtx`（显式 Content-Length），而非 `httpx.OkJsonCtx`。来源：session-log §2026-06-16 11:43 / commit 4a6c68c
