# Docker 镜像构建与部署指南

本文档详细说明如何为这个 Go 项目构建 Docker 镜像并运行容器。

## 📋 目录

1. [前置要求](#前置要求)
2. [构建镜像](#构建镜像)
3. [运行容器](#运行容器)
4. [使用 Docker Compose](#使用-docker-compose)
5. [镜像优化说明](#镜像优化说明)
6. [常见问题](#常见问题)

---

## 前置要求

### 必需软件

- **Docker**: 版本 20.10 或更高
- **Docker Compose**: 版本 1.29 或更高（可选，用于完整环境）

### 检查安装

```bash
# 检查 Docker 版本
docker --version

# 检查 Docker Compose 版本
docker-compose --version
```

### 项目准备

确保以下文件存在：

1. ✅ `resource/static/config.yml` - 配置文件
2. ✅ `resource/static/keys/private_key.pem` - RSA 私钥文件
3. ✅ `resource/static/web/` - 前端静态文件目录

---

## 构建镜像

### 方法一：使用 Dockerfile 直接构建

#### 基本构建命令

```bash
# 在项目根目录执行
docker build -t go-viewer:latest .
```

#### 带标签的构建

```bash
# 构建并打标签
docker build -t go-viewer:latest -t go-viewer:v1.0.0 .

# 查看构建的镜像
docker images | grep go-viewer
```

#### 构建参数说明

- `-t go-viewer:latest`: 指定镜像名称和标签
- `.`: 构建上下文（当前目录）

### 方法二：使用构建参数（高级）

如果需要自定义构建参数：

```bash
docker build \
  --build-arg GO_VERSION=1.20 \
  -t go-viewer:latest \
  .
```

---

## 运行容器

### 基本运行

```bash
# 运行容器（需要外部 MySQL 和 Redis）
docker run -d \
  --name go-viewer-app \
  -p 17080:17080 \
  go-viewer:latest
```

### 完整运行（带配置和日志）

```bash
# 创建日志目录
mkdir -p logs

# 运行容器
docker run -d \
  --name go-viewer-app \
  -p 17080:17080 \
  -v $(pwd)/resource/static/config.yml:/app/resource/static/config.yml:ro \
  -v $(pwd)/logs:/app/logs \
  --restart unless-stopped \
  go-viewer:latest
```

### 参数说明

- `-d`: 后台运行（detached mode）
- `--name go-viewer-app`: 容器名称
- `-p 17080:17080`: 端口映射（主机端口:容器端口）
- `-v`: 数据卷挂载
  - 配置文件挂载为只读（`:ro`）
  - 日志目录挂载为读写
- `--restart unless-stopped`: 自动重启策略

### 查看容器状态

```bash
# 查看运行中的容器
docker ps

# 查看容器日志
docker logs -f go-viewer-app

# 查看容器详细信息
docker inspect go-viewer-app
```

### 停止和删除容器

```bash
# 停止容器
docker stop go-viewer-app

# 删除容器
docker rm go-viewer-app

# 停止并删除
docker rm -f go-viewer-app
```

---

## 使用 Docker Compose

Docker Compose 可以一键启动完整的应用环境（包括 MySQL 和 Redis）。

### 启动所有服务

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f go-viewer
```

### 停止服务

```bash
# 停止所有服务
docker-compose down

# 停止并删除数据卷（⚠️ 会删除数据库数据）
docker-compose down -v
```

### 重新构建

```bash
# 重新构建镜像并启动
docker-compose up -d --build
```

### 环境变量配置

可以通过环境变量覆盖配置，编辑 `docker-compose.yml` 或在运行时指定：

```bash
# 使用环境变量文件
docker-compose --env-file .env up -d
```

---

## 镜像优化说明

### 多阶段构建的优势

本 Dockerfile 采用**多阶段构建**（Multi-stage Build），具有以下优势：

1. **减小镜像体积**
   - 构建阶段包含完整的 Go 编译工具链（~300MB）
   - 运行阶段只包含 Alpine Linux 和编译好的二进制（~20MB）
   - 最终镜像大小减少约 90%

2. **提高安全性**
   - 运行镜像不包含源代码和构建工具
   - 使用非 root 用户运行应用
   - 减少攻击面

3. **优化构建缓存**
   - `go mod download` 单独一层，依赖未变化时复用缓存
   - 加快后续构建速度

### 构建参数说明

```dockerfile
# 编译优化参数
-ldflags="-w -s"
```

- `-w`: 去除 DWARF 调试信息
- `-s`: 去除符号表和调试信息
- 可进一步减小二进制文件大小（约 30-40%）

### 健康检查

Dockerfile 中配置了健康检查：

```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:17080/swagger/index.html || exit 1
```

- `interval`: 每 30 秒检查一次
- `timeout`: 超时时间 3 秒
- `start-period`: 启动后 5 秒开始检查
- `retries`: 连续失败 3 次标记为不健康

查看健康状态：

```bash
docker ps  # 查看 STATUS 列
docker inspect go-viewer-app | grep -A 10 Health
```

---

## 常见问题

### 1. 构建失败：找不到依赖

**问题**: `go mod download` 失败

**解决方案**:
```bash
# 确保网络连接正常
# 如果使用代理，设置代理环境变量
docker build --build-arg HTTP_PROXY=http://proxy:port .
```

### 2. 运行时错误：配置文件不存在

**问题**: 容器内找不到 `config.yml`

**解决方案**:
- 确保配置文件路径正确
- 使用数据卷挂载配置文件：
  ```bash
  docker run -v $(pwd)/resource/static/config.yml:/app/resource/static/config.yml:ro ...
  ```

### 3. 数据库连接失败

**问题**: 应用无法连接到 MySQL

**解决方案**:
- 检查 `config.yml` 中的数据库配置
- 如果使用 Docker Compose，确保服务名称正确（`mysql`）
- 检查网络连接：
  ```bash
  docker network ls
  docker network inspect go-viewer-network
  ```

### 4. 端口被占用

**问题**: `Error: bind: address already in use`

**解决方案**:
```bash
# 查看占用端口的进程
netstat -ano | findstr :17080  # Windows
lsof -i :17080                  # Linux/Mac

# 或修改端口映射
docker run -p 18080:17080 ...
```

### 5. 权限问题

**问题**: 日志文件无法写入

**解决方案**:
- 确保日志目录权限正确
- 检查挂载的目录权限：
  ```bash
  chmod 755 logs
  ```

### 6. 镜像体积过大

**问题**: 镜像体积超过预期

**解决方案**:
- 检查 `.dockerignore` 是否正确配置
- 使用多阶段构建（已实现）
- 清理未使用的镜像：
  ```bash
  docker system prune -a
  ```

---

## 生产环境建议

### 1. 使用特定版本标签

```bash
# 不要使用 latest 标签
docker build -t go-viewer:v1.0.0 .
docker tag go-viewer:v1.0.0 registry.example.com/go-viewer:v1.0.0
```

### 2. 配置资源限制

在 `docker-compose.yml` 中添加：

```yaml
services:
  go-viewer:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

### 3. 使用私有镜像仓库

```bash
# 推送到私有仓库
docker tag go-viewer:latest registry.example.com/go-viewer:latest
docker push registry.example.com/go-viewer:latest

# 从私有仓库拉取
docker pull registry.example.com/go-viewer:latest
```

### 4. 配置日志轮转

在 `docker-compose.yml` 中配置日志驱动：

```yaml
services:
  go-viewer:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

---

## 快速参考

### 常用命令

```bash
# 构建镜像
docker build -t go-viewer:latest .

# 运行容器
docker run -d -p 17080:17080 --name go-viewer-app go-viewer:latest

# 查看日志
docker logs -f go-viewer-app

# 进入容器
docker exec -it go-viewer-app sh

# 停止容器
docker stop go-viewer-app

# 删除容器
docker rm go-viewer-app

# 使用 Docker Compose
docker-compose up -d
docker-compose logs -f
docker-compose down
```

---

## 相关资源

- [Docker 官方文档](https://docs.docker.com/)
- [Docker Compose 文档](https://docs.docker.com/compose/)
- [Go 官方 Docker 镜像](https://hub.docker.com/_/golang)
- [Alpine Linux 镜像](https://hub.docker.com/_/alpine)

