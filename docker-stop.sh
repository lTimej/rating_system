#!/bin/bash

# Docker停止脚本

echo "=== 停止评价系统 Docker 服务 ==="

# 获取当前脚本目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# 停止并删除容器
echo "停止并删除容器..."
docker-compose down

# 可选：删除镜像（取消注释以启用）
# echo "删除镜像..."
# docker-compose down --rmi all

# 可选：删除数据卷（取消注释以启用，注意这会删除所有数据）
# echo "删除数据卷..."
# docker-compose down -v

echo "✅ 所有服务已停止"

# 显示剩余的容器和镜像
echo ""
echo "剩余容器:"
docker ps -a | grep rating-system || echo "无相关容器"

echo ""
echo "剩余镜像:"
docker images | grep rating-system || echo "无相关镜像"
