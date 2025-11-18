# 多阶段构建 Dockerfile
# 第一阶段：构建阶段
FROM golang:1.20-alpine AS builder

# 设置工作目录
WORKDIR /build

# 安装必要的构建工具
RUN apk add --no-cache git ca-certificates tzdata

# 设置 Go 环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖（利用 Docker 缓存层）
RUN go mod download

# 复制项目源代码
COPY . .

# 构建应用
# -ldflags 参数用于减小二进制文件大小并移除调试信息
RUN go build -ldflags="-w -s" -o app main.go

# 第二阶段：运行阶段
FROM alpine:latest

# 安装必要的运行时依赖（包括 curl 用于健康检查）
RUN apk --no-cache add ca-certificates tzdata curl

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非 root 用户运行应用
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /build/app .

# 从构建阶段复制必要的资源文件
COPY --from=builder /build/resource ./resource

# 创建日志目录
RUN mkdir -p /app/logs && \
    chown -R appuser:appuser /app

# 切换到非 root 用户
USER appuser

# 暴露端口（根据配置文件中的 server-port 设置）
EXPOSE 17080

# 健康检查（使用 curl 检查服务是否正常）
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:17080/swagger/index.html || exit 1

# 启动应用
CMD ["./app"]

