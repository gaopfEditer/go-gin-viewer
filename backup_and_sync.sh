#!/bin/bash

# 备份本地数据库并同步到 Docker MySQL
# 使用方法: ./backup_and_sync.sh

set -e

DB_NAME="activate_server"
DB_USER="root"
DB_PASSWORD="Cambridge#*DR"
LOCAL_PORT="3388"
BACKUP_DIR="backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/backup_${TIMESTAMP}.sql"

echo "========================================"
echo "数据库备份和同步脚本"
echo "========================================"
echo ""

# 创建备份目录
mkdir -p "${BACKUP_DIR}"

# 备份本地数据库
echo "[1/4] 备份本地数据库..."
echo "备份文件: ${BACKUP_FILE}"
mysqldump -h localhost -P "${LOCAL_PORT}" -u "${DB_USER}" -p"${DB_PASSWORD}" \
    --single-transaction --routines --triggers "${DB_NAME}" > "${BACKUP_FILE}"

if [ $? -ne 0 ]; then
    echo "错误: 备份失败！请检查："
    echo "  1. MySQL 服务是否运行"
    echo "  2. 端口 ${LOCAL_PORT} 是否正确"
    echo "  3. 用户名和密码是否正确"
    echo "  4. mysqldump 命令是否可用"
    exit 1
fi

echo "备份成功！"
echo ""

# 检查 Docker 容器是否运行
echo "[2/4] 检查 Docker MySQL 容器状态..."
if ! docker ps | grep -q "go-viewer-mysql"; then
    echo "警告: Docker MySQL 容器未运行，正在启动..."
    docker-compose up -d mysql
    echo "等待 MySQL 启动..."
    sleep 10
fi

echo "Docker MySQL 容器运行中"
echo ""

# 确认操作
echo "[3/4] 准备同步数据到 Docker MySQL..."
echo "注意: 这将覆盖 Docker 中的现有数据"
read -p "是否继续? (y/n): " CONFIRM
if [ "${CONFIRM}" != "y" ] && [ "${CONFIRM}" != "Y" ]; then
    echo "操作已取消"
    exit 0
fi

# 删除 Docker 中的旧数据库
echo "正在删除 Docker 中的旧数据库..."
docker exec -i go-viewer-mysql mysql -uroot -p"${DB_PASSWORD}" -e \
    "DROP DATABASE IF EXISTS ${DB_NAME}; CREATE DATABASE ${DB_NAME} CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

if [ $? -ne 0 ]; then
    echo "错误: 无法删除/创建数据库"
    exit 1
fi

echo ""

# 导入数据到 Docker
echo "[4/4] 导入数据到 Docker MySQL..."
cat "${BACKUP_FILE}" | docker exec -i go-viewer-mysql mysql -uroot -p"${DB_PASSWORD}" "${DB_NAME}"

if [ $? -ne 0 ]; then
    echo "错误: 导入失败！"
    exit 1
fi

echo ""
echo "========================================"
echo "同步完成！"
echo "========================================"
echo "备份文件: ${BACKUP_FILE}"
echo "Docker 数据库: ${DB_NAME}"
echo ""
echo "验证数据:"
echo "  docker exec -it go-viewer-mysql mysql -uroot -p${DB_PASSWORD} -e \"USE ${DB_NAME}; SHOW TABLES;\""
echo ""





