# 多阶段构建 Dockerfile
# 第一阶段：构建阶段
FROM golang:1.22-alpine AS builder

# 设置工作目录
WORKDIR /build

# 安装必要的构建工具
RUN apk add --no-cache git ca-certificates tzdata

# 设置 Go 环境变量
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# 构建参数：可选配置 Go 代理（默认使用国内代理，如果网络正常可设置为空使用官方源）
# 使用方法: docker build --build-arg GOPROXY= .
ARG GOPROXY=https://goproxy.cn,direct
ARG GOSUMDB=sum.golang.google.cn

# 设置 Go 代理环境变量
ENV GOPROXY=${GOPROXY}
ENV GOSUMDB=${GOSUMDB}

# 复制项目源代码（先复制 go.mod 和 go.sum 以利用 Docker 缓存）
COPY go.mod go.sum ./

# 下载外部依赖（利用缓存层）
RUN echo "使用 Go 代理: $GOPROXY" && \
    go mod download

# 复制所有源代码
COPY . .

# 构建应用（此时所有源代码已就位，Go 可以找到内部包）
# 使用 server 作为输出文件名，避免与 app/ 目录冲突
RUN go build -ldflags="-w -s" -o server main.go && \
    ls -lh server

# 第二阶段：运行阶段
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata curl

# 设置时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 创建非 root 用户
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# 设置工作目录
WORKDIR /app

# 从构建阶段复制文件，并设置正确的所有者和权限
COPY --from=builder --chown=appuser:appuser /build/server /app/app
COPY --from=builder --chown=appuser:appuser /build/resource /app/resource
COPY --from=builder --chown=appuser:appuser /build/docs /app/docs

# 创建日志目录并设置权限
RUN mkdir -p /app/logs && \
    chmod +x /app/app && \
    ls -lh /app/ && \
    chown -R appuser:appuser /app/logs

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 17080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:17080/swagger/index.html || exit 1

# 启动应用
CMD ["./app"]

