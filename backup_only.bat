@echo off
REM 仅备份本地数据库
REM 使用方法: backup_only.bat

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
echo 数据库备份脚本
echo ========================================
echo.

REM 创建备份目录
if not exist "%BACKUP_DIR%" mkdir "%BACKUP_DIR%"

echo 正在备份数据库: %DB_NAME%
echo 备份文件: %BACKUP_FILE%
echo.

mysqldump -h localhost -P %LOCAL_PORT% -u %DB_USER% -p%DB_PASSWORD% --single-transaction --routines --triggers %DB_NAME% > "%BACKUP_FILE%"

if errorlevel 1 (
    echo 错误: 备份失败！
    pause
    exit /b 1
)

echo 备份成功！
echo 备份文件: %BACKUP_FILE%
echo.
pause





