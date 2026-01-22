# Media Report - 媒体报表服务

基于 go-zero 框架构建的媒体报表数据查询服务，集成快手广告 API。

## 项目简介

本项目提供了一个简单的 HTTP API 服务，用于查询快手广告账户的报表数据。支持单个账户查询和批量查询。

## 功能特性

- ✅ 单个账户报表查询
- ✅ 批量账户报表查询
- ✅ 支持多种时间粒度（日、周、月）
- ✅ 并发批量查询，提高性能
- ✅ 完整的错误处理和日志记录

## 技术栈

- **框架**: [go-zero](https://go-zero.dev/) - 高性能微服务框架
- **语言**: Go 1.23.12
- **外部 API**: 快手广告开放平台 API

## 项目结构

```
media_report/
├── api/                    # API 定义文件
│   └── media.api          # HTTP 接口定义
├── common/                # 公共模块
│   └── kuaishou/         # 快手 API 客户端
│       └── client.go     # HTTP 客户端封装
├── service/              # 服务实现
│   └── api/             # API 服务
│       ├── etc/         # 配置文件
│       ├── internal/    # 内部实现
│       │   ├── config/  # 配置结构
│       │   ├── handler/ # HTTP 处理器
│       │   ├── logic/   # 业务逻辑
│       │   ├── svc/     # 服务上下文
│       │   └── types/   # 类型定义
│       └── media.go     # 服务入口
├── go.mod               # Go 模块定义
├── Makefile            # 构建脚本
└── README.md           # 项目文档
```

## 快速开始

### 前置要求

- Go 1.23.12 或更高版本
- 快手广告平台 Access Token

### 安装依赖

```bash
make init
# 或
go mod tidy
```

### 配置

编辑配置文件 `service/api/etc/media-api.yaml`:

```yaml
Name: media-api
Host: 0.0.0.0
Port: 8888

Log:
  ServiceName: media-api
  Mode: console
  Level: info

# 快手广告 API 配置
Kuaishou:
  BaseUrl: https://ad.e.kuaishou.com
  AccessToken: YOUR_ACCESS_TOKEN_HERE  # 替换为你的 Access Token
  Timeout: 30
```

### 运行服务

```bash
# 使用 Makefile
make run-api

# 或直接运行
cd service/api && go run media.go -f etc/media-api.yaml
```

服务将在 `http://localhost:8888` 启动。

## 开发指南

### 编译项目

```bash
make build-api
```

编译后的二进制文件位于 `bin/media-api`。

### 运行测试

```bash
make test
```

### 代码格式化

```bash
make fmt
```

### 代码检查

```bash
make lint
```

### 查看所有命令

```bash
make help
```

## 部署

### Docker 部署

1. 构建镜像:
```bash
make docker-build
```

2. 运行容器:
```bash
make docker-run
```

### 生产环境建议

1. 使用环境变量管理敏感配置（Access Token）
2. 启用 HTTPS
3. 添加 API 限流和熔断机制
4. 配置日志收集和监控告警
5. 使用负载均衡器分发请求

## 常见问题

### 1. 如何获取快手广告平台 Access Token？

请访问 [快手广告开放平台](https://ad.e.kuaishou.com/) 注册账号并创建应用获取。

### 2. 支持哪些时间粒度？

支持以下时间粒度：
- `DAILY`: 按天
- `WEEKLY`: 按周
- `MONTHLY`: 按月

### 3. 批量查询的并发限制？

当前实现支持并发查询，但建议根据快手 API 的限流策略控制并发数量。

## 后续计划

- [ ] 添加数据缓存机制（Redis）
- [ ] 数据持久化（MySQL）
- [ ] 更多报表维度支持
- [ ] 数据可视化界面
- [ ] 定时任务自动拉取数据

## 许可证

MIT License

## 联系方式

如有问题，请提交 Issue 或联系项目维护者。
