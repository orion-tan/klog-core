.PHONY: build run clean dev help

# 变量定义
BINARY_NAME=klog-backend
GO=go
GOFLAGS=-v
LDFLAGS=-w -s

# 默认目标
.DEFAULT_GOAL := help

## help: 显示帮助信息
help:
	@echo "可用的make命令："
	@echo "  make build       - 编译应用程序"
	@echo "  make run         - 运行应用程序"
	@echo "  make dev         - 开发模式运行（热重载）"
	@echo "  make test        - 运行测试"
	@echo "  make clean       - 清理编译文件"
	@echo "  make deps        - 安装依赖"
	@echo "  make build-linux - 编译Linux版本"
	@echo "  make docker      - 构建Docker镜像"

## build: 编译应用程序
build:
	@echo "开始编译..."
	CGO_ENABLED=1 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) ./cmd/main.go
	@echo "编译完成: $(BINARY_NAME)"

## run: 运行应用程序
run: build
	@echo "启动应用..."
	./$(BINARY_NAME)

## dev: 开发模式运行
dev:
	@echo "开发模式启动..."
	$(GO) run ./cmd/main.go

## clean: 清理编译文件
clean:
	@echo "清理编译文件..."
	@rm -f $(BINARY_NAME)
	@rm -rf ./db/*.db
	@rm -rf ./uploads/*
	@rm -rf ./log/*
	@echo "清理完成"

## deps: 安装依赖
deps:
	@echo "安装依赖..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "依赖安装完成"

## build-linux: 编译Linux AMD64版本
build-linux:
	@echo "编译Linux版本..."
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-linux ./cmd/main.go
	@echo "编译完成: $(BINARY_NAME)-linux"

## build-linux-arm: 编译Linux ARM64版本
build-linux-arm:
	@echo "编译Linux ARM64版本..."
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-linux-arm64 ./cmd/main.go
	@echo "编译完成: $(BINARY_NAME)-linux-arm64"

## docker: 构建Docker镜像
docker:
	@echo "构建Docker镜像..."
	docker build -t klog-backend:latest .
	@echo "Docker镜像构建完成"

## docker-run: 运行Docker容器
docker-run:
	@echo "运行Docker容器..."
	docker-compose up -d

## docker-stop: 停止Docker容器
docker-stop:
	@echo "停止Docker容器..."
	docker-compose down

## init-db: 初始化数据库（会删除现有数据）
init-db:
	@echo "警告：此操作将删除现有数据库！"
	@read -p "确定继续吗？[y/N] " confirm && [ "$$confirm" = "y" ] || exit 1
	@rm -f ./db/klog.db
	@echo "数据库已重置"
	@$(MAKE) dev

## fmt: 格式化代码
fmt:
	@echo "格式化代码..."
	$(GO) fmt ./...
	@echo "代码格式化完成"

## lint: 代码检查
lint:
	@echo "代码检查..."
	golangci-lint run ./...

## mod-update: 更新依赖
mod-update:
	@echo "更新依赖..."
	$(GO) get -u ./...
	$(GO) mod tidy
	@echo "依赖更新完成"

