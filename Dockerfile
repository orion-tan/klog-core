# 构建阶段
FROM golang:1.23-alpine AS builder

# 设置Go环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# 安装编译所需的依赖
RUN apk add --no-cache gcc musl-dev git

# 设置工作目录
WORKDIR /app

# 复制go mod文件并下载依赖（利用Docker缓存层）
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# 复制源代码
COPY . .

# 编译应用（添加版本信息和构建时间）
ARG VERSION=dev
ARG BUILD_TIME
RUN go build -a -installsuffix cgo \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}" \
    -o klog-backend ./cmd/main.go

# 运行阶段
FROM alpine:latest

# 安装运行时依赖（添加wget用于健康检查）
RUN apk --no-cache add ca-certificates tzdata wget

# 设置时区为中国上海
ENV TZ=Asia/Shanghai

# 创建非root用户
RUN addgroup -g 1000 klog && \
    adduser -D -u 1000 -G klog klog

# 设置工作目录
WORKDIR /app

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/klog-backend .

# 复制配置文件
COPY --from=builder /app/configs ./configs

# 创建必要的目录并设置权限
RUN mkdir -p /app/db /app/uploads /app/logs && \
    chown -R klog:klog /app

# 切换到非root用户运行
USER klog

# 暴露端口
EXPOSE 8010

# 健康检查配置
# - interval: 每30秒检查一次
# - timeout: 检查超时时间为3秒
# - start-period: 启动后5秒开始检查
# - retries: 连续失败3次才标记为不健康
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8010/health/live || exit 1

# 运行应用
CMD ["./klog-backend"]

