@echo off
REM 备份本地数据库并同步到 Docker MySQL
REM 使用方法: backup_and_sync.bat

setlocal enabledelayedexpansion

set DB_NAME=activate_server
set DB_USER=root
set DB_PASSWORD=Cambridge#*DR
set LOCAL_PORT=3388
set BACKUP_DIR=backups
set TIMESTAMP=%date:~0,4%%date:~5,2%%date:~8,2%_%time:~0,2%%time:~3,2%%time:~6,2%
set TIMESTAMP=!TIMESTAMP: =0!
set BACKUP_FILE=%BACKUP_DIR%\backup_%TIMESTAMP%.sql

echo ========================================
echo 数据库备份和同步脚本
echo ========================================
echo.

REM 创建备份目录
if not exist "%BACKUP_DIR%" mkdir "%BACKUP_DIR%"

echo [1/4] 备份本地数据库...
echo 备份文件: %BACKUP_FILE%
mysqldump -h localhost -P %LOCAL_PORT% -u %DB_USER% -p%DB_PASSWORD% --single-transaction --routines --triggers %DB_NAME% > "%BACKUP_FILE%"

if errorlevel 1 (
    echo 错误: 备份失败！请检查：
    echo   1. MySQL 服务是否运行
    echo   2. 端口 %LOCAL_PORT% 是否正确
    echo   3. 用户名和密码是否正确
    echo   4. mysqldump 命令是否可用
    pause
    exit /b 1
)

echo 备份成功！
echo.

REM 检查 Docker 容器是否运行
echo [2/4] 检查 Docker MySQL 容器状态...
docker ps | findstr "go-viewer-mysql" >nul
if errorlevel 1 (
    echo 警告: Docker MySQL 容器未运行，正在启动...
    docker-compose up -d mysql
    echo 等待 MySQL 启动...
    timeout /t 10 /nobreak >nul
)

echo Docker MySQL 容器运行中
echo.

REM 删除 Docker 中的旧数据库（可选）
echo [3/4] 准备同步数据到 Docker MySQL...
echo 注意: 这将覆盖 Docker 中的现有数据
set /p CONFIRM="是否继续? (y/n): "
if /i not "!CONFIRM!"=="y" (
    echo 操作已取消
    pause
    exit /b 0
)

echo 正在删除 Docker 中的旧数据库...
docker exec -i go-viewer-mysql mysql -uroot -p%DB_PASSWORD% -e "DROP DATABASE IF EXISTS %DB_NAME%; CREATE DATABASE %DB_NAME% CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

if errorlevel 1 (
    echo 错误: 无法删除/创建数据库
    pause
    exit /b 1
)

echo.

REM 导入数据到 Docker
echo [4/4] 导入数据到 Docker MySQL...
type "%BACKUP_FILE%" | docker exec -i go-viewer-mysql mysql -uroot -p%DB_PASSWORD% %DB_NAME%

if errorlevel 1 (
    echo 错误: 导入失败！
    pause
    exit /b 1
)

echo.
echo ========================================
echo 同步完成！
echo ========================================
echo 备份文件: %BACKUP_FILE%
echo Docker 数据库: %DB_NAME%
echo.
echo 验证数据:
echo   docker exec -it go-viewer-mysql mysql -uroot -p%DB_PASSWORD% -e "USE %DB_NAME%; SHOW TABLES;"
echo.
pause


