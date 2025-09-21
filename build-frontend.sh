#!/bin/bash

# 前端构建和部署脚本

echo "开始构建前端..."

# 确保在项目根目录
cd "$(dirname "$0")"

# 进入前端目录
cd frontend

# 安装依赖（如果需要）
if [ ! -d "node_modules" ]; then
    echo "安装前端依赖..."
    npm install
fi

# 构建前端
echo "构建前端应用..."
npm run build

# 检查构建是否成功
if [ $? -eq 0 ]; then
    echo "前端构建成功！"
    
    # 移动dist目录到项目根目录
    echo "移动构建文件到项目根目录..."
    rm -rf ../dist
    mv dist ../dist
    
    # 检查关键文件是否存在
    if [ -f "../dist/index.html" ]; then
        echo "✅ index.html 文件存在"
    else
        echo "❌ 警告：index.html 文件不存在！"
    fi
    
    if [ -d "../dist/static" ]; then
        echo "✅ static 目录存在"
    else
        echo "❌ 警告：static 目录不存在！"
    fi
    
    echo "前端部署完成！"
    echo "现在可以启动后端服务器了。"
    echo ""
    echo "访问方式："
    echo "- 根路径: http://localhost:8080/"
    echo "- 任何前端路由: http://localhost:8080/dashboard"
    echo "- API接口: http://localhost:8080/api/..."
else
    echo "前端构建失败！"
    exit 1
fi
