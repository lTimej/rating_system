#!/bin/bash

# 评价系统启动脚本

echo "=== 评价系统启动脚本 ==="

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo "错误: Go未安装，请先安装Go 1.21+"
    exit 1
fi

# 检查Node.js是否安装
if ! command -v node &> /dev/null; then
    echo "错误: Node.js未安装，请先安装Node.js 16+"
    exit 1
fi

# 获取当前脚本目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$SCRIPT_DIR/backend"
FRONTEND_DIR="$SCRIPT_DIR/frontend"

echo "项目目录: $SCRIPT_DIR"

# 启动后端
echo "=== 启动后端服务 ==="
cd "$BACKEND_DIR"

# 安装Go依赖
echo "安装Go依赖..."
go mod tidy

# 后台启动Go服务
echo "启动后端服务 (端口: 8080)..."
nohup go run main.go > backend.log 2>&1 &
BACKEND_PID=$!
echo "后端服务PID: $BACKEND_PID"

# 等待后端启动
sleep 3

# 检查后端是否启动成功
if ps -p $BACKEND_PID > /dev/null; then
    echo "✅ 后端服务启动成功"
else
    echo "❌ 后端服务启动失败，请检查 backend.log"
    exit 1
fi

# 启动前端
echo "=== 启动前端服务 ==="
cd "$FRONTEND_DIR"

# 检查是否已安装依赖
if [ ! -d "node_modules" ]; then
    echo "安装前端依赖..."
    npm install
fi

# 启动前端开发服务器
echo "启动前端服务 (端口: 3000)..."
npm run serve &
FRONTEND_PID=$!
echo "前端服务PID: $FRONTEND_PID"

# 保存PID到文件
echo $BACKEND_PID > "$SCRIPT_DIR/backend.pid"
echo $FRONTEND_PID > "$SCRIPT_DIR/frontend.pid"

echo ""
echo "=== 服务启动完成 ==="
echo "后端服务: http://localhost:8080"
echo "前端应用: http://localhost:3000"
echo ""
echo "默认管理员账户:"
echo "用户名: admin"
echo "密码: password"
echo ""
echo "要停止服务，请运行: ./stop.sh"
echo ""

# 等待用户输入来保持脚本运行
echo "按 Ctrl+C 停止所有服务..."
trap 'echo "正在停止服务..."; kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; exit 0' INT

# 保持脚本运行
wait
