#!/bin/bash

# 评价系统停止脚本

echo "=== 停止评价系统服务 ==="

# 获取当前脚本目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 读取PID文件并停止服务
if [ -f "$SCRIPT_DIR/backend.pid" ]; then
    BACKEND_PID=$(cat "$SCRIPT_DIR/backend.pid")
    if ps -p $BACKEND_PID > /dev/null; then
        echo "停止后端服务 (PID: $BACKEND_PID)..."
        kill $BACKEND_PID
        echo "✅ 后端服务已停止"
    else
        echo "后端服务未运行"
    fi
    rm -f "$SCRIPT_DIR/backend.pid"
else
    echo "未找到后端PID文件，尝试通过端口停止..."
    # 尝试通过端口号停止
    PID=$(lsof -ti:8080)
    if [ ! -z "$PID" ]; then
        kill $PID
        echo "✅ 通过端口8080停止了后端服务"
    fi
fi

if [ -f "$SCRIPT_DIR/frontend.pid" ]; then
    FRONTEND_PID=$(cat "$SCRIPT_DIR/frontend.pid")
    if ps -p $FRONTEND_PID > /dev/null; then
        echo "停止前端服务 (PID: $FRONTEND_PID)..."
        kill $FRONTEND_PID
        echo "✅ 前端服务已停止"
    else
        echo "前端服务未运行"
    fi
    rm -f "$SCRIPT_DIR/frontend.pid"
else
    echo "未找到前端PID文件，尝试通过端口停止..."
    # 尝试通过端口号停止
    PID=$(lsof -ti:3000)
    if [ ! -z "$PID" ]; then
        kill $PID
        echo "✅ 通过端口3000停止了前端服务"
    fi
fi

# 清理日志文件
if [ -f "$SCRIPT_DIR/backend/backend.log" ]; then
    echo "清理后端日志文件..."
    rm -f "$SCRIPT_DIR/backend/backend.log"
fi

echo "=== 所有服务已停止 ==="
