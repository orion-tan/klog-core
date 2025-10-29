# KLog 博客系统后端 API

这是一个完整的博客系统后端API，基于Go语言和Gin框架开发。

## 技术栈

- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: SQLite (可轻松切换到MySQL/PostgreSQL)
- **认证**: JWT
- **密码加密**: bcrypt

## 项目结构

```
backend/
├── cmd/                    # 应用入口
│   └── main.go
├── configs/                # 配置文件
│   └── config.toml
├── internal/               # 内部代码
│   ├── api/                # API请求/响应结构
│   ├── config/             # 配置管理
│   ├── database/           # 数据库迁移
│   ├── handler/            # 请求处理器（Controller）
│   ├── middleware/         # 中间件
│   ├── model/              # 数据模型
│   ├── repository/         # 数据访问层
│   ├── router/             # 路由配置
│   ├── services/           # 业务逻辑层
│   └── utils/              # 工具函数
└── uploads/                # 上传文件存储目录（自动创建）
```

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 配置

编辑 `configs/config.toml` 文件：

```toml
[server]
port = 8010

[database]
type = "sqlite"
url = "./db/klog.db"

[jwt]
secret = "your-very-secret-key"  # 请修改为强密码
expire_hour = 72
```

### 3. 运行

```bash
go run cmd/main.go
```

或编译后运行：

```bash
go build -o klog-backend cmd/main.go
./klog-backend
```

服务器将在 `http://localhost:8010` 启动。

## API 文档

### 基础URL

```
http://localhost:8010/api/v1
```

### 认证方式

在请求头中添加：
```
Authorization: Bearer <your_jwt_token>
```

### API 端点

#### 认证接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | /auth/register | 用户注册 | 否 |
| POST | /auth/login | 用户登录 | 否 |
| GET | /auth/me | 获取当前用户信息 | 是 |

#### 文章接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /posts | 获取文章列表 | 可选 |
| GET | /posts/:id | 获取文章详情 | 可选 |
| POST | /posts | 创建文章 | 是 |
| PUT | /posts/:id | 更新文章 | 是 |
| DELETE | /posts/:id | 删除文章 | 是 |

#### 分类接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /categories | 获取所有分类 | 否 |
| POST | /categories | 创建分类 | 是（管理员）|
| PUT | /categories/:id | 更新分类 | 是（管理员）|
| DELETE | /categories/:id | 删除分类 | 是（管理员）|

#### 标签接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /tags | 获取所有标签 | 否 |
| POST | /tags | 创建标签 | 是（管理员）|
| PUT | /tags/:id | 更新标签 | 是（管理员）|
| DELETE | /tags/:id | 删除标签 | 是（管理员）|

#### 评论接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /posts/:postId/comments | 获取文章评论列表 | 否 |
| POST | /posts/:postId/comments | 发表评论 | 可选 |
| PUT | /comments/:id | 更新评论状态 | 是（管理员）|
| DELETE | /comments/:id | 删除评论 | 是（管理员）|

#### 媒体库接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | /media/upload | 上传文件 | 是 |
| GET | /media | 获取媒体列表 | 是 |
| DELETE | /media/:id | 删除媒体文件 | 是 |
| GET | /uploads/:filename | 访问上传的文件 | 否 |

#### 用户接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /users | 获取用户列表 | 是（管理员）|
| GET | /users/:id | 获取用户信息 | 否 |
| PUT | /users/:id | 更新用户信息 | 是 |

### 请求示例

#### 用户注册

```bash
curl -X POST http://localhost:8010/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "nickname": "Test User"
  }'
```

#### 用户登录

```bash
curl -X POST http://localhost:8010/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "login": "testuser",
    "password": "password123"
  }'
```

#### 创建文章

```bash
curl -X POST http://localhost:8010/api/v1/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  -d '{
    "title": "My First Post",
    "slug": "my-first-post",
    "content": "This is my first blog post!",
    "excerpt": "A brief summary",
    "status": "published",
    "category_id": 1,
    "tags": ["tech", "golang"]
  }'
```

#### 上传文件（multipart）

```bash
curl -X POST http://localhost:8010/api/v1/media/upload \
  -H "Authorization: Bearer <your_token>" \
  -F "file=@/path/to/your/image.jpg"
```

### 响应格式

#### 成功响应

```json
{
  "success": true,
  "data": {
    // 返回的数据
  }
}
```

#### 错误响应

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "错误信息"
  }
}
```

## 数据模型

### 用户 (User)
- id: 用户ID
- username: 用户名
- email: 邮箱
- nickname: 昵称
- avatar_url: 头像URL
- role: 角色 (author/admin)
- status: 状态 (active/inactive)

### 文章 (Post)
- id: 文章ID
- category_id: 分类ID
- author_id: 作者ID
- title: 标题
- slug: URL别名
- content: 内容
- excerpt: 摘要
- cover_image_url: 封面图片
- status: 状态 (draft/published/archived)
- view_count: 浏览量
- published_at: 发布时间

### 分类 (Category)
- id: 分类ID
- name: 分类名称
- slug: URL别名
- description: 描述

### 标签 (Tag)
- id: 标签ID
- name: 标签名称
- slug: URL别名

### 评论 (Comment)
- id: 评论ID
- post_id: 文章ID
- user_id: 用户ID（可为空）
- name: 评论者名称（游客）
- email: 评论者邮箱（游客）
- content: 评论内容
- ip: IP地址
- status: 状态 (pending/approved/spam)
- parent_id: 父评论ID

### 媒体 (Media)
- id: 媒体ID
- uploader_id: 上传者ID
- file_name: 文件名
- file_path: 文件路径
- mime_type: 文件类型
- size: 文件大小

## 权限说明

### 角色类型
- **author**: 普通作者，可以创建和管理自己的文章
- **admin**: 管理员，拥有所有权限

### 权限规则
1. 未登录用户：只能查看已发布的文章、分类、标签、评论
2. 已登录用户：可以创建文章、发表评论、上传文件
3. 作者：可以管理自己的文章
4. 管理员：可以管理所有内容，包括用户、分类、标签等

## 开发说明

### 添加新功能

1. 在 `internal/model` 中定义数据模型
2. 在 `internal/repository` 中实现数据访问层
3. 在 `internal/services` 中实现业务逻辑
4. 在 `internal/handler` 中实现请求处理
5. 在 `internal/router` 中注册路由

### 数据库迁移

应用启动时会自动执行数据库迁移，创建或更新表结构。

### 日志

应用使用Gin默认的日志中间件，所有请求都会被记录。

## 生产环境部署

### 1. 修改配置

- 修改JWT密钥为强密码
- 如果使用MySQL/PostgreSQL，修改数据库配置
- 考虑使用环境变量存储敏感信息

### 2. 编译

```bash
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o klog-backend cmd/main.go
```

### 3. 运行

```bash
./klog-backend
```

### 4. 使用进程管理器

推荐使用 systemd 或 supervisord 管理进程。

### 5. 反向代理

使用 Nginx 作为反向代理：

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8010;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    location /uploads/ {
        alias /path/to/backend/uploads/;
    }
}
```

## 常见问题

### 1. 数据库文件在哪里？

SQLite数据库文件默认存储在 `./db/klog.db`

### 2. 上传的文件存储在哪里？

上传的文件存储在 `./uploads/` 目录

### 3. 如何创建管理员账户？

第一个注册的用户需要手动在数据库中将 role 修改为 'admin'

### 4. 如何切换到MySQL/PostgreSQL？

1. 修改 `configs/config.toml` 中的数据库配置
2. 在 `cmd/main.go` 中修改数据库驱动

## 许可证

MIT License

