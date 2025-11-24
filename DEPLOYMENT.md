# 项目部署指南

本文档详细说明如何将项目部署到服务器，包括 Docker 方式和非 Docker 方式。

## 📋 目录

1. [前置要求](#前置要求)
2. [Docker 方式部署](#docker-方式部署)
3. [非 Docker 方式部署](#非-docker-方式部署)
4. [生产环境配置](#生产环境配置)
5. [常见问题](#常见问题)

---

## 前置要求

### 服务器环境要求

- **操作系统**: Linux (推荐 Ubuntu 20.04+ / CentOS 7+)
- **CPU**: 2 核或以上
- **内存**: 4GB 或以上
- **磁盘**: 20GB 或以上可用空间

### 必需软件

#### Docker 方式
- Docker 20.10+
- Docker Compose 1.29+ (可选)

#### 非 Docker 方式
- Go 1.20+
- MySQL 8.0+
- Redis 6.0+
- Nginx (可选，用于反向代理)

---

## Docker 方式部署

### 方式一：使用 Docker Compose（推荐）

Docker Compose 可以一键启动完整的应用环境（包括 MySQL 和 Redis）。

#### 1. 准备配置文件

确保 `resource/static/config.yml` 存在并配置正确：

```yaml
app:
    env: prod  # 生产环境
    cache: true
    machine-id: 1
    server-port: 17080
    api-prefix: /activate
mysql:
    host: mysql  # Docker Compose 中使用服务名
    port: 3306
    user: root
    password: Cambridge#*DR
    dbname: activate_server
redis:
    addr: redis:6379  # Docker Compose 中使用服务名
    password:
    db: 0
```

#### 2. 部署步骤

```bash
# 1. 上传项目到服务器
scp -r ./go-viewer user@server:/opt/

# 2. SSH 登录服务器
ssh user@server

# 3. 进入项目目录
cd /opt/go-viewer

# 4. 使用 Docker Compose 启动所有服务
docker-compose up -d --build

# 5. 查看服务状态
docker-compose ps

# 6. 查看日志
docker-compose logs -f app
```

#### 3. 验证部署

```bash
# 检查服务健康状态
docker-compose ps

# 访问应用
curl http://localhost:17080/swagger/index.html

# 查看应用日志
docker-compose logs -f app
```

#### 4. 常用命令

```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose down

# 重启服务
docker-compose restart

# 查看日志
docker-compose logs -f

# 更新代码后重新构建
docker-compose up -d --build

# 停止并删除数据卷（⚠️ 会删除数据库数据）
docker-compose down -v
```

### 方式二：单独构建和运行 Docker 容器

如果已有外部 MySQL 和 Redis，可以只运行应用容器。

#### 1. 构建镜像

```bash
# 在项目根目录执行
docker build -t go-viewer:latest .

# 或者使用构建脚本
./build.sh -b  # Linux/Mac
build.bat -b   # Windows
```

#### 2. 运行容器

```bash
# 创建日志目录
mkdir -p logs

# 运行容器（需要外部 MySQL 和 Redis）
docker run -d \
  --name go-viewer-app \
  -p 17080:17080 \
  -v $(pwd)/resource/static/config.yml:/app/resource/static/config.yml:ro \
  -v $(pwd)/logs:/app/logs \
  --restart unless-stopped \
  go-viewer:latest
```

#### 3. 配置外部数据库

如果使用外部 MySQL 和 Redis，需要修改 `config.yml`：

```yaml
mysql:
    host: your-mysql-host  # 外部 MySQL 地址
    port: 3306
redis:
    addr: your-redis-host:6379  # 外部 Redis 地址
```

---

## 非 Docker 方式部署

### 1. 安装依赖

#### Ubuntu/Debian

```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装 Go
wget https://go.dev/dl/go1.20.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.20.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装 MySQL
sudo apt install mysql-server -y

# 安装 Redis
sudo apt install redis-server -y

# 安装 Nginx (可选)
sudo apt install nginx -y
```

#### CentOS/RHEL

```bash
# 更新系统
sudo yum update -y

# 安装 Go
wget https://go.dev/dl/go1.20.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.20.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装 MySQL
sudo yum install mysql-server -y
sudo systemctl start mysqld
sudo systemctl enable mysqld

# 安装 Redis
sudo yum install redis -y
sudo systemctl start redis
sudo systemctl enable redis

# 安装 Nginx (可选)
sudo yum install nginx -y
```

### 2. 配置数据库

```bash
# 登录 MySQL
sudo mysql -u root -p

# 创建数据库
CREATE DATABASE activate_server CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 创建用户（可选）
CREATE USER 'appuser'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON activate_server.* TO 'appuser'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

### 3. 配置 Redis

```bash
# 编辑 Redis 配置
sudo vim /etc/redis/redis.conf

# 修改以下配置（如果需要）
# bind 127.0.0.1  # 如果只允许本地访问
# requirepass your_password  # 设置密码

# 重启 Redis
sudo systemctl restart redis
```

### 4. 部署应用

```bash
# 1. 上传项目到服务器
scp -r ./go-viewer user@server:/opt/

# 2. SSH 登录服务器
ssh user@server

# 3. 进入项目目录
cd /opt/go-viewer

# 4. 配置 Go 环境
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=sum.golang.google.cn

# 5. 下载依赖
go mod download

# 6. 构建应用
go build -ldflags="-w -s" -o app main.go

# 7. 配置应用
vim resource/static/config.yml
# 修改数据库和 Redis 配置为实际地址

# 8. 创建日志目录
mkdir -p logs

# 9. 测试运行
./app
```

### 5. 使用 Systemd 管理服务

创建 systemd 服务文件：

```bash
sudo vim /etc/systemd/system/go-viewer.service
```

添加以下内容：

```ini
[Unit]
Description=Go Viewer Application
After=network.target mysql.service redis.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/go-viewer
ExecStart=/opt/go-viewer/app
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

# 环境变量
Environment="GIN_MODE=release"
Environment="TZ=Asia/Shanghai"

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
# 重载 systemd
sudo systemctl daemon-reload

# 启动服务
sudo systemctl start go-viewer

# 设置开机自启
sudo systemctl enable go-viewer

# 查看状态
sudo systemctl status go-viewer

# 查看日志
sudo journalctl -u go-viewer -f
```

### 6. 配置 Nginx 反向代理（可选）

创建 Nginx 配置：

```bash
sudo vim /etc/nginx/sites-available/go-viewer
```

添加以下内容：

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 日志
    access_log /var/log/nginx/go-viewer-access.log;
    error_log /var/log/nginx/go-viewer-error.log;

    # 反向代理到应用
    location / {
        proxy_pass http://127.0.0.1:17080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket 支持
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    # 静态文件缓存
    location ~* \.(jpg|jpeg|png|gif|ico|css|js|svg|woff|woff2|ttf|eot)$ {
        proxy_pass http://127.0.0.1:17080;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }
}
```

启用配置：

```bash
# 创建软链接
sudo ln -s /etc/nginx/sites-available/go-viewer /etc/nginx/sites-enabled/

# 测试配置
sudo nginx -t

# 重启 Nginx
sudo systemctl restart nginx
```

### 7. 配置防火墙

```bash
# Ubuntu/Debian (UFW)
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 17080/tcp  # 如果直接访问应用端口
sudo ufw enable

# CentOS/RHEL (firewalld)
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --permanent --add-port=17080/tcp
sudo firewall-cmd --reload
```

---

## 生产环境配置

### 1. 安全配置

#### 修改默认密码

```bash
# 修改 MySQL root 密码
sudo mysql_secure_installation

# 修改 Redis 密码
sudo vim /etc/redis/redis.conf
# 设置 requirepass your_strong_password
```

#### 配置 HTTPS

使用 Let's Encrypt 免费 SSL 证书：

```bash
# 安装 Certbot
sudo apt install certbot python3-certbot-nginx -y

# 获取证书
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo certbot renew --dry-run
```

### 2. 性能优化

#### 数据库优化

```sql
-- 创建索引
-- 根据实际业务需求创建索引

-- 配置 MySQL
sudo vim /etc/mysql/mysql.conf.d/mysqld.cnf

# 添加以下配置
[mysqld]
max_connections = 200
innodb_buffer_pool_size = 1G
innodb_log_file_size = 256M
```

#### Redis 优化

```bash
sudo vim /etc/redis/redis.conf

# 配置
maxmemory 512mb
maxmemory-policy allkeys-lru
```

#### 应用优化

在 `config.yml` 中配置：

```yaml
mysql:
    max-open-conns: 100
    max-idle-conns: 10
```

### 3. 监控和日志

#### 日志管理

```bash
# 配置日志轮转
sudo vim /etc/logrotate.d/go-viewer

# 添加内容
/opt/go-viewer/logs/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0644 www-data www-data
}
```

#### 监控工具

推荐使用：
- **Prometheus + Grafana**: 监控指标
- **ELK Stack**: 日志分析
- **Sentry**: 错误追踪

### 4. 备份策略

#### 数据库备份

```bash
# 创建备份脚本
sudo vim /opt/scripts/backup-mysql.sh

# 添加内容
#!/bin/bash
BACKUP_DIR="/opt/backups/mysql"
DATE=$(date +%Y%m%d_%H%M%S)
mkdir -p $BACKUP_DIR
mysqldump -u root -p'your_password' activate_server > $BACKUP_DIR/backup_$DATE.sql
find $BACKUP_DIR -name "backup_*.sql" -mtime +7 -delete

# 设置定时任务
sudo crontab -e
# 添加: 0 2 * * * /opt/scripts/backup-mysql.sh
```

#### 应用备份

```bash
# 备份配置和日志
tar -czf /opt/backups/app/app_backup_$(date +%Y%m%d).tar.gz \
  /opt/go-viewer/resource \
  /opt/go-viewer/logs
```

---

## 常见问题

### 1. 端口被占用

```bash
# 查找占用端口的进程
sudo lsof -i :17080
# 或
sudo netstat -tlnp | grep 17080

# 停止占用进程
sudo kill -9 <PID>
```

### 2. 数据库连接失败

- 检查 MySQL 服务是否运行：`sudo systemctl status mysql`
- 检查防火墙规则
- 验证配置文件中的数据库地址和密码
- 检查 MySQL 用户权限

### 3. Redis 连接失败

- 检查 Redis 服务：`sudo systemctl status redis`
- 检查 Redis 密码配置
- 验证网络连接

### 4. 应用无法启动

```bash
# 查看详细错误
./app

# 检查日志
tail -f logs/app.log

# 检查配置文件
cat resource/static/config.yml
```

### 5. WebSocket 连接失败

- 确保 Nginx 配置了 WebSocket 支持（见上方 Nginx 配置）
- 检查防火墙是否允许 WebSocket 连接
- 验证应用端口是否正确暴露

### 6. 内存不足

```bash
# 查看内存使用
free -h

# 优化应用配置
# 减少数据库连接数
# 优化缓存策略
```

---

## 快速部署脚本

### Docker 方式

```bash
#!/bin/bash
# 一键部署脚本 (Docker)

cd /opt/go-viewer
docker-compose down
docker-compose pull
docker-compose up -d --build
docker-compose logs -f
```

### 非 Docker 方式

```bash
#!/bin/bash
# 一键部署脚本 (非 Docker)

cd /opt/go-viewer
git pull  # 如果使用 Git
go mod download
go build -ldflags="-w -s" -o app main.go
sudo systemctl restart go-viewer
sudo journalctl -u go-viewer -f
```

---

## 更新部署

### Docker 方式

```bash
# 1. 停止服务
docker-compose down

# 2. 拉取最新代码（如果使用 Git）
git pull

# 3. 重新构建并启动
docker-compose up -d --build

# 4. 查看日志
docker-compose logs -f app
```

### 非 Docker 方式

```bash
# 1. 拉取最新代码
git pull

# 2. 重新构建
go build -ldflags="-w -s" -o app main.go

# 3. 重启服务
sudo systemctl restart go-viewer

# 4. 查看日志
sudo journalctl -u go-viewer -f
```

---

## 相关资源

- [Docker 官方文档](https://docs.docker.com/)
- [Docker Compose 文档](https://docs.docker.com/compose/)
- [Go 官方文档](https://go.dev/doc/)
- [Nginx 官方文档](https://nginx.org/en/docs/)
- [Systemd 服务管理](https://www.freedesktop.org/software/systemd/man/systemd.service.html)

