#!/bin/bash

# Docker部署脚本

echo "=== 评价系统 Docker 部署脚本 ==="

# 检查Docker是否安装
if ! command -v docker &> /dev/null; then
    echo "错误: Docker未安装，请先安装Docker"
    exit 1
fi

# 检查Docker Compose是否安装
if ! command -v docker-compose &> /dev/null; then
    echo "错误: Docker Compose未安装，请先安装Docker Compose"
    exit 1
fi

# 获取当前脚本目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "项目目录: $SCRIPT_DIR"

# 创建必要的目录
echo "创建数据目录..."
mkdir -p data uploads

# 停止现有容器
echo "停止现有容器..."
docker-compose down

# 构建并启动服务
echo "构建并启动服务..."
docker-compose up --build -d

# 等待服务启动
echo "等待服务启动..."
sleep 10

# 检查服务状态
echo "检查服务状态..."
docker-compose ps

# 检查后端健康状态
echo "检查后端健康状态..."
if curl -f http://localhost:8080/api/health > /dev/null 2>&1; then
    echo "✅ 后端服务启动成功"
else
    echo "❌ 后端服务启动失败"
    echo "查看后端日志:"
    docker-compose logs backend
    exit 1
fi

# 检查前端状态
echo "检查前端状态..."
if curl -f http://localhost:80 > /dev/null 2>&1; then
    echo "✅ 前端服务启动成功"
else
    echo "❌ 前端服务启动失败"
    echo "查看前端日志:"
    docker-compose logs frontend
    exit 1
fi

echo ""
echo "=== 部署完成 ==="
echo "前端应用: http://localhost"
echo "后端API: http://localhost:8080"
echo ""
echo "默认管理员账户:"
echo "用户名: admin"
echo "密码: password"
echo ""
echo "管理命令:"
echo "查看日志: docker-compose logs -f"
echo "停止服务: docker-compose down"
echo "重启服务: docker-compose restart"
echo "查看状态: docker-compose ps"
echo ""
