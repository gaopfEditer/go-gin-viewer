@echo off
REM Go Viewer Docker 构建脚本 (Windows)
REM 使用方法: build.bat [选项]

setlocal enabledelayedexpansion

set IMAGE_NAME=go-viewer
set IMAGE_TAG=latest
set BUILD_TYPE=single

REM 解析参数
:parse_args
if "%~1"=="" goto :check_requirements
if /i "%~1"=="-n" (
    set IMAGE_NAME=%~2
    shift
    shift
    goto :parse_args
)
if /i "%~1"=="--name" (
    set IMAGE_NAME=%~2
    shift
    shift
    goto :parse_args
)
if /i "%~1"=="-t" (
    set IMAGE_TAG=%~2
    shift
    shift
    goto :parse_args
)
if /i "%~1"=="--tag" (
    set IMAGE_TAG=%~2
    shift
    shift
    goto :parse_args
)
if /i "%~1"=="-c" (
    set BUILD_TYPE=compose
    shift
    goto :parse_args
)
if /i "%~1"=="--compose" (
    set BUILD_TYPE=compose
    shift
    goto :parse_args
)
if /i "%~1"=="-b" (
    set BUILD_TYPE=build-only
    shift
    goto :parse_args
)
if /i "%~1"=="--build-only" (
    set BUILD_TYPE=build-only
    shift
    goto :parse_args
)
if /i "%~1"=="-h" goto :show_help
if /i "%~1"=="--help" goto :show_help
shift
goto :parse_args

:show_help
echo Go Viewer Docker 构建脚本 (Windows)
echo.
echo 使用方法:
echo   build.bat [选项]
echo.
echo 选项:
echo   -n, --name NAME     镜像名称 (默认: go-viewer)
echo   -t, --tag TAG       镜像标签 (默认: latest)
echo   -c, --compose       使用 docker-compose 构建和启动
echo   -b, --build-only    仅构建镜像，不运行
echo   -h, --help          显示此帮助信息
echo.
echo 示例:
echo   build.bat                          # 构建并运行单个容器
echo   build.bat -c                        # 使用 docker-compose 构建和启动
echo   build.bat -n myapp -t v1.0.0        # 指定镜像名称和标签
echo   build.bat -b                        # 仅构建镜像
exit /b 0

:check_requirements
echo 检查必要文件...

if not exist "Dockerfile" (
    echo 错误: 未找到 Dockerfile
    exit /b 1
)

if not exist "resource\static\config.yml" (
    echo 错误: 未找到配置文件 resource\static\config.yml
    exit /b 1
)

if not exist "resource\static\keys\private_key.pem" (
    echo 警告: 未找到私钥文件 resource\static\keys\private_key.pem
    echo 请先运行: cd resource\static\keys ^&^& go test -v -run TestGenerateRSA
    set /p CONTINUE="是否继续? (y/n) "
    if /i not "!CONTINUE!"=="y" exit /b 1
)

echo 文件检查完成

:build_image
echo 开始构建 Docker 镜像...
echo 镜像名称: %IMAGE_NAME%:%IMAGE_TAG%

docker build -t %IMAGE_NAME%:%IMAGE_TAG% .

if errorlevel 1 (
    echo 镜像构建失败!
    exit /b 1
)

echo 镜像构建成功!
docker images | findstr %IMAGE_NAME%

if "%BUILD_TYPE%"=="build-only" (
    echo 完成!
    exit /b 0
)

if "%BUILD_TYPE%"=="compose" goto :run_compose
if "%BUILD_TYPE%"=="single" goto :run_container

:run_container
echo 启动容器...

REM 停止并删除已存在的容器
docker ps -a | findstr "go-viewer-app" >nul
if not errorlevel 1 (
    echo 停止并删除已存在的容器...
    docker rm -f go-viewer-app 2>nul
)

REM 创建日志目录
if not exist "logs" mkdir logs

REM 运行容器
docker run -d --name go-viewer-app -p 17080:17080 -v "%CD%\resource\static\config.yml:/app/resource/static/config.yml:ro" -v "%CD%\logs:/app/logs" --restart unless-stopped %IMAGE_NAME%:%IMAGE_TAG%

if errorlevel 1 (
    echo 容器启动失败!
    exit /b 1
)

echo 容器启动成功!
echo 访问地址:
echo   - API 文档: http://localhost:17080/swagger/index.html
echo   - 前端页面: http://localhost:17080/
echo.
echo 查看日志: docker logs -f go-viewer-app
goto :end

:run_compose
echo 使用 Docker Compose 构建和启动...

if not exist "docker-compose.yml" (
    echo 错误: 未找到 docker-compose.yml
    exit /b 1
)

docker-compose up -d --build

if errorlevel 1 (
    echo 服务启动失败!
    exit /b 1
)

echo 所有服务启动成功!
echo 访问地址:
echo   - API 文档: http://localhost:17080/swagger/index.html
echo   - 前端页面: http://localhost:17080/
echo.
echo 查看日志: docker-compose logs -f
echo 停止服务: docker-compose down

:end
echo 完成!
exit /b 0

