# 评价系统

一个基于 Golang + Gin + SQLite + Vue2 的学生作品评价系统。

## 功能特性

### 用户角色
- **学生**: 可以上传作品、接受评价、给其他学生评价
- **专家**: 可以评价学生作品、关注学生
- **管理员**: 系统管理、用户管理、评价管理

### 核心功能
1. **登录系统**: 支持用户注册、登录、角色权限管理
2. **个人主页**: 用户信息展示、文件上传、作品展示
3. **评价系统**: 
   - 学生互相评价
   - 专家对学生评价
   - 评价反馈（有帮助/没帮助）
   - 评价者身份公开，被评价者身份匿名
4. **关注系统**: 专家可以关注学生
5. **后台管理**: 用户管理、评价管理、系统统计

## 技术栈

### 后端
- **语言**: Go 1.21+
- **框架**: Gin
- **数据库**: SQLite
- **ORM**: GORM
- **认证**: JWT
- **密码加密**: bcrypt

### 前端
- **框架**: Vue 2.6
- **UI库**: Element UI
- **状态管理**: Vuex
- **路由**: Vue Router
- **HTTP客户端**: Axios

## 项目结构

```
rating_system/
├── backend/                 # 后端代码
│   ├── config/             # 配置文件
│   ├── controllers/        # 控制器
│   ├── middleware/         # 中间件
│   ├── models/            # 数据模型
│   ├── routes/            # 路由配置
│   ├── go.mod             # Go模块文件
│   └── main.go            # 主程序入口
├── frontend/               # 前端代码
│   ├── public/            # 静态文件
│   ├── src/               # 源代码
│   │   ├── components/    # 组件
│   │   ├── views/         # 页面
│   │   ├── store/         # Vuex状态管理
│   │   ├── router/        # 路由配置
│   │   └── main.js        # 主程序入口
│   ├── package.json       # 依赖配置
│   └── vue.config.js      # Vue配置
└── README.md              # 项目说明
```

## 快速开始

### 环境要求
- Go 1.21+
- Node.js 16+
- npm 或 yarn

### 后端启动

1. 进入后端目录
```bash
cd backend
```

2. 安装依赖
```bash
go mod tidy
```

3. 启动服务
```bash
go run main.go
```

后端服务将在 `http://localhost:8080` 启动

### 前端启动

1. 进入前端目录
```bash
cd frontend
```

2. 安装依赖
```bash
npm install
```

3. 启动开发服务器
```bash
npm run serve
```

前端应用将在 `http://localhost:3000` 启动

## 默认账户

系统会自动创建一个默认管理员账户：
- 用户名: `admin`
- 密码: `password`

## API 接口

### 认证相关
- `POST /api/login` - 用户登录
- `POST /api/register` - 用户注册
- `GET /api/profile` - 获取用户信息
- `PUT /api/profile` - 更新用户信息

### 文件管理
- `POST /api/files` - 上传文件
- `GET /api/files/user/:user_id` - 获取用户文件列表
- `GET /api/files/:id` - 获取文件信息
- `GET /api/files/:id/download` - 下载文件
- `DELETE /api/files/:id` - 删除文件

### 评价管理
- `POST /api/ratings` - 创建评价
- `GET /api/ratings/user/:user_id` - 获取用户评价
- `GET /api/ratings/public` - 获取公开评价
- `POST /api/ratings/:rating_id/feedback` - 给评价反馈

### 关注系统
- `POST /api/follow/:user_id` - 关注用户
- `DELETE /api/follow/:user_id` - 取消关注
- `GET /api/following` - 获取关注列表
- `GET /api/followers` - 获取粉丝列表

### 管理员接口
- `GET /api/admin/users` - 获取用户列表
- `POST /api/admin/users` - 创建用户
- `PUT /api/admin/users/:id` - 更新用户
- `DELETE /api/admin/users/:id` - 删除用户
- `GET /api/admin/ratings` - 获取评价列表
- `DELETE /api/admin/ratings/:id` - 删除评价
- `GET /api/admin/stats` - 获取系统统计

## 数据库设计

### 用户表 (users)
- id, username, email, password, role, name, avatar, bio, links, is_active

### 文件表 (files)
- id, user_id, file_name, file_path, file_type, file_size, title, description, is_public

### 评价表 (ratings)
- id, rater_id, rated_id, file_id, content, score, is_anonymous

### 反馈表 (feedbacks)
- id, rating_id, user_id, is_helpful

### 关注表 (follows)
- id, follower_id, followed_id

## 部署说明

### 生产环境配置

1. 修改JWT密钥（backend/middleware/auth.go）
2. 配置数据库连接
3. 设置CORS策略
4. 构建前端应用：`npm run build`
5. 配置反向代理（如Nginx）

### Docker部署

可以创建Dockerfile来容器化部署应用。

## 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

MIT License
