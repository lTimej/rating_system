#!/bin/bash

# 单一镜像 Docker 部署脚本

echo "=== 评价系统单一镜像 Docker 部署脚本 ==="

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
docker-compose -f docker-compose.yml down

# 构建并启动服务
echo "构建并启动单一镜像服务..."
echo "注意: 首次构建可能需要较长时间，请耐心等待..."
docker-compose -f docker-compose.yml up --build -d

# 等待服务启动
echo "等待服务启动..."
sleep 15

# 检查服务状态
echo "检查服务状态..."
docker-compose -f docker-compose.yml ps

# 检查应用健康状态
echo "检查应用健康状态..."
if curl -f http://localhost:8080/api/health > /dev/null 2>&1; then
    echo "✅ 后端API服务启动成功"
else
    echo "❌ 后端API服务启动失败"
    echo "查看应用日志:"
    docker-compose -f docker-compose.yml logs
    exit 1
fi

# 检查前端状态
echo "检查前端状态..."
if curl -f http://localhost:8080 > /dev/null 2>&1; then
    echo "✅ 前端应用启动成功"
else
    echo "❌ 前端应用启动失败"
    echo "查看应用日志:"
    docker-compose -f docker-compose.yml logs
    exit 1
fi

# 显示镜像信息
echo ""
echo "=== 镜像信息 ==="
docker images | grep rating-system

echo ""
echo "=== 部署完成 ==="
echo "应用访问地址: http://localhost:8080"
echo "API接口地址: http://localhost:8080/api"
echo ""
echo "默认管理员账户:"
echo "用户名: admin"
echo "密码: password"
echo ""
echo "管理命令:"
echo "查看日志: docker-compose -f docker-compose.yml logs -f"
echo "停止服务: docker-compose -f docker-compose.yml down"
echo "重启服务: docker-compose -f docker-compose.yml restart"
echo "查看状态: docker-compose -f docker-compose.yml ps"
echo ""
echo "数据持久化:"
echo "数据库文件: ./data/rating_system.db"
echo "上传文件: ./uploads/"
echo ""
