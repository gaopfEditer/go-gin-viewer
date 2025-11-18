#!/bin/bash

# Go Viewer Docker 构建脚本
# 使用方法: ./build.sh [选项]

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 默认值
IMAGE_NAME="go-viewer"
IMAGE_TAG="latest"
BUILD_TYPE="single"

# 显示帮助信息
show_help() {
    echo "Go Viewer Docker 构建脚本"
    echo ""
    echo "使用方法:"
    echo "  ./build.sh [选项]"
    echo ""
    echo "选项:"
    echo "  -n, --name NAME     镜像名称 (默认: go-viewer)"
    echo "  -t, --tag TAG       镜像标签 (默认: latest)"
    echo "  -c, --compose       使用 docker-compose 构建和启动"
    echo "  -b, --build-only    仅构建镜像，不运行"
    echo "  -h, --help          显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  ./build.sh                          # 构建并运行单个容器"
    echo "  ./build.sh -c                        # 使用 docker-compose 构建和启动"
    echo "  ./build.sh -n myapp -t v1.0.0        # 指定镜像名称和标签"
    echo "  ./build.sh -b                        # 仅构建镜像"
}

# 解析参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -n|--name)
            IMAGE_NAME="$2"
            shift 2
            ;;
        -t|--tag)
            IMAGE_TAG="$2"
            shift 2
            ;;
        -c|--compose)
            BUILD_TYPE="compose"
            shift
            ;;
        -b|--build-only)
            BUILD_TYPE="build-only"
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo -e "${RED}未知参数: $1${NC}"
            show_help
            exit 1
            ;;
    esac
done

# 检查必要文件
check_requirements() {
    echo -e "${YELLOW}检查必要文件...${NC}"
    
    if [ ! -f "Dockerfile" ]; then
        echo -e "${RED}错误: 未找到 Dockerfile${NC}"
        exit 1
    fi
    
    if [ ! -f "resource/static/config.yml" ]; then
        echo -e "${RED}错误: 未找到配置文件 resource/static/config.yml${NC}"
        exit 1
    fi
    
    if [ ! -f "resource/static/keys/private_key.pem" ]; then
        echo -e "${YELLOW}警告: 未找到私钥文件 resource/static/keys/private_key.pem${NC}"
        echo -e "${YELLOW}请先运行: cd resource/static/keys && go test -v -run TestGenerateRSA${NC}"
        read -p "是否继续? (y/n) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
    
    echo -e "${GREEN}文件检查完成${NC}"
}

# 构建镜像
build_image() {
    echo -e "${YELLOW}开始构建 Docker 镜像...${NC}"
    echo -e "镜像名称: ${GREEN}${IMAGE_NAME}:${IMAGE_TAG}${NC}"
    
    docker build -t "${IMAGE_NAME}:${IMAGE_TAG}" .
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}镜像构建成功!${NC}"
        docker images | grep "${IMAGE_NAME}" | head -1
    else
        echo -e "${RED}镜像构建失败!${NC}"
        exit 1
    fi
}

# 运行单个容器
run_container() {
    echo -e "${YELLOW}启动容器...${NC}"
    
    # 停止并删除已存在的容器
    if docker ps -a | grep -q "go-viewer-app"; then
        echo -e "${YELLOW}停止并删除已存在的容器...${NC}"
        docker rm -f go-viewer-app 2>/dev/null || true
    fi
    
    # 创建日志目录
    mkdir -p logs
    
    # 运行容器
    docker run -d \
        --name go-viewer-app \
        -p 17080:17080 \
        -v "$(pwd)/resource/static/config.yml:/app/resource/static/config.yml:ro" \
        -v "$(pwd)/logs:/app/logs" \
        --restart unless-stopped \
        "${IMAGE_NAME}:${IMAGE_TAG}"
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}容器启动成功!${NC}"
        echo -e "访问地址:"
        echo -e "  - API 文档: ${GREEN}http://localhost:17080/swagger/index.html${NC}"
        echo -e "  - 前端页面: ${GREEN}http://localhost:17080/${NC}"
        echo ""
        echo "查看日志: docker logs -f go-viewer-app"
    else
        echo -e "${RED}容器启动失败!${NC}"
        exit 1
    fi
}

# 使用 docker-compose
run_compose() {
    echo -e "${YELLOW}使用 Docker Compose 构建和启动...${NC}"
    
    if [ ! -f "docker-compose.yml" ]; then
        echo -e "${RED}错误: 未找到 docker-compose.yml${NC}"
        exit 1
    fi
    
    docker-compose up -d --build
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}所有服务启动成功!${NC}"
        echo -e "访问地址:"
        echo -e "  - API 文档: ${GREEN}http://localhost:17080/swagger/index.html${NC}"
        echo -e "  - 前端页面: ${GREEN}http://localhost:17080/${NC}"
        echo ""
        echo "查看日志: docker-compose logs -f"
        echo "停止服务: docker-compose down"
    else
        echo -e "${RED}服务启动失败!${NC}"
        exit 1
    fi
}

# 主流程
main() {
    check_requirements
    
    case $BUILD_TYPE in
        compose)
            run_compose
            ;;
        build-only)
            build_image
            ;;
        single)
            build_image
            run_container
            ;;
    esac
    
    echo -e "${GREEN}完成!${NC}"
}

# 执行主流程
main

