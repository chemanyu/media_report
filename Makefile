.PHONY: api run-api build-api test clean help

# 项目变量
PROJECT_NAME=media_report
API_BINARY=media-api
API_DIR=service/api

# 颜色输出
GREEN=\033[0;32m
NC=\033[0m # No Color

help: ## 显示帮助信息
	@echo "$(GREEN)可用命令:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}'

init: ## 初始化项目依赖
	@echo "$(GREEN)初始化项目依赖...$(NC)"
	go mod tidy
	go mod download

build-api: ## 编译 API 服务
	@echo "$(GREEN)编译 API 服务...$(NC)"
	cd $(API_DIR) && go build -o ../../bin/$(API_BINARY) media.go

run-api: ## 运行 API 服务
	@echo "$(GREEN)启动 API 服务...$(NC)"
	cd $(API_DIR) && go run media.go -f etc/media-api.yaml

dev-api: ## 开发模式运行 API（使用 air 热重载）
	@echo "$(GREEN)开发模式启动 API 服务...$(NC)"
	@if command -v air > /dev/null; then \
		cd $(API_DIR) && air; \
	else \
		echo "air 未安装，使用普通模式运行"; \
		$(MAKE) run-api; \
	fi

test: ## 运行测试
	@echo "$(GREEN)运行测试...$(NC)"
	go test -v ./...

test-coverage: ## 运行测试并生成覆盖率报告
	@echo "$(GREEN)生成测试覆盖率报告...$(NC)"
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)覆盖率报告已生成: coverage.html$(NC)"

clean: ## 清理编译文件
	@echo "$(GREEN)清理编译文件...$(NC)"
	rm -rf bin/
	rm -f coverage.out coverage.html

fmt: ## 格式化代码
	@echo "$(GREEN)格式化代码...$(NC)"
	gofmt -w .
	@if command -v goimports > /dev/null; then \
		goimports -w .; \
	fi

lint: ## 代码检查
	@echo "$(GREEN)代码检查...$(NC)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint 未安装，请运行: brew install golangci-lint"; \
	fi

docker-build: ## 构建 Docker 镜像
	@echo "$(GREEN)构建 Docker 镜像...$(NC)"
	docker build -t $(PROJECT_NAME):latest .

docker-run: ## 运行 Docker 容器
	@echo "$(GREEN)运行 Docker 容器...$(NC)"
	docker run -p 8888:8888 $(PROJECT_NAME):latest

all: clean init build-api ## 完整构建流程
