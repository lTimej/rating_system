# 多阶段构建 - 前端构建阶段
FROM node:16-alpine AS frontend-builder

WORKDIR /app/frontend

# 复制前端package文件
COPY frontend/package*.json ./

# 安装前端依赖
RUN npm ci --only=production
RUN npm install

# 复制前端源代码
COPY frontend/ ./

# 构建前端应用
RUN npm run build

# 后端构建阶段
FROM golang:1.21-alpine AS backend-builder
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOPROXY="https://goproxy.cn,direct" \
    GO111MODULE=on

WORKDIR /app/backend

# 安装必要的系统依赖 (只需要git，不需要C编译器)

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update && \
    apk add --no-cache git gcc musl-dev sqlite-dev build-base

# 复制后端go mod文件
COPY backend/go.mod backend/go.sum ./

# 下载后端依赖
RUN go mod tidy
RUN go mod download

# 复制后端源代码
COPY backend/ ./

# 设置环境变量 (禁用CGO使用纯Go SQLite驱动)
ENV CGO_ENABLED=0
ENV GOOS=linux

# 构建后端应用
RUN go build -o main .

# 最终运行阶段
FROM alpine

# 安装运行时依赖
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update && \
    apk add --no-cache ca-certificates

# 创建非root用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 设置工作目录
WORKDIR /app

# 从后端构建阶段复制二进制文件
COPY --from=backend-builder /app/backend/main .

# 从前端构建阶段复制构建结果
COPY --from=frontend-builder /app/frontend/dist ./static/dist

# 创建必要的目录
RUN mkdir -p uploads data && \
    chown -R appuser:appgroup /app

# 切换到非root用户
USER appuser

# 暴露端口
EXPOSE 8080

# 设置环境变量
ENV GIN_MODE=release
ENV DB_PATH=/app/data/rating_system.db

# 运行应用
CMD ["./main"]