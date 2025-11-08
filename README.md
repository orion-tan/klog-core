> 项目目前处于激进开发阶段，绝不能用于生产环境！！！

# KLog 博客系统后端 API

这是一个完整的博客系统后端API，基于Go语言和Gin框架开发。

## 技术栈

- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: SQLite (可轻松切换到MySQL/PostgreSQL)
- **认证**: JWT (golang-jwt/jwt/v5)
- **密码加密**: bcrypt
- **缓存**: Redis（可选）
- **日志**: Zap + Lumberjack（日志轮转）
- **配置管理**: Viper

## 项目结构

```
backend/
├── cmd/                    # 应用入口
│   └── main.go
├── configs/                # 配置文件
│   └── config.toml
├── internal/               # 内部代码
│   ├── api/                # API请求/响应结构
│   ├── cache               # redis cache
│   ├── config/             # 配置管理
│   ├── database/           # 数据库迁移
│   ├── handler/            # 请求处理器（Controller）
│   ├── middleware/         # 中间件
│   ├── model/              # 数据模型
│   ├── repository/         # 数据访问层
│   ├── router/             # 路由配置
│   ├── services/           # 业务逻辑层
│   └── utils/              # 工具函数
├── log/                    # 日志目录（默认位置）
├── db/                     # 数据库目录（默认位置）
└── uploads/                # 上传文件存储目录（默认位置）
```

## 快速开始

### 方式一：Docker Compose（推荐）

`todo...` 

### 方式二：本地开发

#### 1. 安装依赖

```bash
go mod tidy
```

#### 2. 配置

编辑 `configs/config.toml`(详细字段参考 `internal/config/config.go`) 文件：

```toml
[server]
port = 8010                              # 服务器端口
cors = ["http://localhost:3000"]         # CORS允许的源（支持通配符，如 "*.example.com"）

[database]
type = "sqlite"                          # 数据库类型：sqlite/mysql/postgres
url = "./db/klog.db"                     # 数据库连接URL

[jwt]
secret = "your-very-secret-key"          # JWT密钥（必须>=32字符）
expire_hour = 168                        # Token过期时间（小时，默认7天）

[redis]
addr = "localhost:6379"                  # Redis地址（可选，用于Token黑名单和缓存）
password = ""                            # Redis密码

[logger]
level = "info"                           # 日志级别：debug/info/warn/error
path = "./log/klog.log"                 # 日志文件路径
max_size = 100                          # 单个日志文件最大大小（MB）
max_backups = 3                         # 保留的备份数量
max_age = 30                            # 日志保留时间（天）

[media]
media_dir = "./uploads"                  # 上传文件存储目录
max_file_size_mb = 10                   # 单个文件最大大小（MB）

[scheduler]
enabled = true                           # 是否启用定时任务
cleanup_cron = "0 0 3 * * 0"            # 清理任务Cron表达式（每周日凌晨3点）
```

**配置说明**:
- `server.cors`: 支持精确匹配和通配符子域名（如`*.example.com`）
- `jwt.secret`: 生产环境请使用强随机密钥，长度不低于32字符
- `redis`: 可选配置，启用后支持Token黑名单和缓存功能
- `scheduler.cleanup_cron`: 使用标准Cron表达式，默认每周日凌晨3点清理过期文件

#### 3. 运行

开发模式：
```bash
make dev
# 或
go run cmd/main.go
```

编译后运行：
```bash
make build
./klog-backend
```

使用Makefile命令：
```bash
make help        # 查看所有可用命令
make build       # 编译应用程序
make run         # 编译并运行
make test        # 运行测试
make clean       # 清理编译文件
make docker      # 构建Docker镜像
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
| POST | /auth/register | 用户注册（仅供首次设置管理员） | 否 |
| POST | /auth/login | 用户登录 | 否 |
| POST | /auth/logout | 用户登出 | 是 |
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
| POST | /categories | 创建分类 | 是 |
| PUT | /categories/:id | 更新分类 | 是 |
| DELETE | /categories/:id | 删除分类 | 是 |

#### 标签接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /tags | 获取所有标签 | 否 |
| POST | /tags | 创建标签 | 是 |
| PUT | /tags/:id | 更新标签 | 是 |
| DELETE | /tags/:id | 删除标签 | 是 |

#### 评论接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /posts/:postId/comments | 获取文章评论列表 | 否 |
| POST | /posts/:postId/comments | 发表评论 | 可选 |
| PUT | /comments/:id | 更新评论状态 | 是 |
| DELETE | /comments/:id | 删除评论 | 是 |

#### 媒体库接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | /media/upload | 上传文件 | 是 |
| GET | /media | 获取媒体列表 | 是 |
| DELETE | /media/:id | 删除媒体文件 | 是 |
| GET | /media/i/:filename | 访问上传的文件 | 否 |

#### 用户接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /users/:id | 获取用户信息 | 否 |
| PUT | /users/:id | 更新用户信息（仅限本人） | 是 |

#### 设置接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /settings | 获取所有系统设置 | 是 |
| GET | /settings/:key | 获取指定设置 | 是 |
| PUT | /settings | 创建/更新单个设置 | 是 |
| PUT | /settings/batch | 批量创建/更新设置 | 是 |
| DELETE | /settings/:key | 删除指定设置 | 是 |

### 请求示例

详细的API请求示例请参考 [API文档](docs/API.md)，以下是常用操作的快速示例：

**登录获取Token**:
```bash
curl -X POST http://localhost:8010/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login":"admin","password":"password123"}'
```

**创建文章**:
```bash
curl -X POST http://localhost:8010/api/v1/posts \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title":"我的第一篇文章",
    "slug":"my-first-post",
    "content":"# Hello World\n这是内容",
    "status":"published",
    "tags":["tech"]
  }'
```

**上传图片**:
```bash
curl -X POST http://localhost:8010/api/v1/media/upload \
  -H "Authorization: Bearer <your_token>" \
  -F "file=@/path/to/image.jpg"
```

更多示例请查看 [完整API文档](docs/API.md)。

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
- username: 用户名（唯一）
- email: 邮箱（唯一）
- password: 密码（bcrypt加密）
- nickname: 昵称
- bio: 个人简介（可选）
- avatar_url: 头像URL（可选）
- created_at: 创建时间
- updated_at: 更新时间

### 文章 (Post)
- id: 文章ID
- category_id: 分类ID（可选）
- author_id: 作者ID
- title: 标题
- slug: URL别名（唯一）
- content: 内容（Markdown格式）
- excerpt: 摘要（可选）
- cover_image_url: 封面图片（可选）
- status: 状态 (draft/published/archived)
- view_count: 浏览量（uint64类型，默认0）
- published_at: 发布时间
- created_at: 创建时间
- updated_at: 更新时间

### 分类 (Category)
- id: 分类ID
- name: 分类名称（唯一）
- slug: URL别名（唯一）
- description: 描述（可选）

### 标签 (Tag)
- id: 标签ID
- name: 标签名称（唯一）
- slug: URL别名（唯一）

### 评论 (Comment)
- id: 评论ID
- post_id: 文章ID
- user_id: 用户ID（登录用户评论时使用）
- name: 评论者名称（游客评论时必填）
- email: 评论者邮箱（游客评论时必填）
- content: 评论内容
- ip: IP地址（自动记录）
- status: 状态 (pending/approved/spam)
- parent_id: 父评论ID（用于嵌套回复）
- created_at: 创建时间

### 媒体 (Media)
- id: 媒体ID
- file_name: 文件名
- file_hash: 文件内容hash
- file_path: 文件路径
- mime_type: 文件类型
- size: 文件大小
- created_at: 上传时间
- updated_at: 更新时间

### 设置 (Setting)
- key: 设置键（主键）
- value: 设置值
- type: 值类型（str/number/json）

**说明**: 系统设置支持三种数据类型：
- `str`: 字符串类型
- `number`: 数字类型（存储时为字符串，读取时自动转换为数值）
- `json`: JSON对象或数组（自动验证和解析）



### 权限规则
1. **未登录用户**：只能查看已发布的文章、分类、标签、评论
2. **管理员（登录用户）**：拥有所有权限
   - 创建、编辑、删除文章（包括草稿）
   - 管理分类、标签
   - 审核、删除评论
   - 上传、删除媒体文件
   - 修改个人信息

### 首次设置流程
1. 系统初次启动后，数据库中无用户
2. 访问注册接口 `POST /api/v1/auth/register` 创建管理员账号
3. 注册成功后，系统将拒绝任何后续注册请求
4. 使用管理员账号登录即可管理整个博客

## 安全特性

### 认证与授权
- **JWT认证**: 使用JWT Bearer Token进行用户认证
- **Token黑名单**: 登出时Token加入Redis黑名单（需启用Redis）
- **密码加密**: 使用bcrypt算法加密存储用户密码
- **权限控制**: 基于中间件的权限验证，区分公开接口和需认证接口

### 限流保护
1. **全局限流**
   - 每个IP：10请求/秒
   - 突发请求：最多20个
   - 适用于所有API端点

2. **评论限流**（额外保护）
   - 每个IP：1条评论/分钟
   - 每个IP：最多10条评论/小时
   - 防止评论垃圾信息

### CORS安全
- 支持配置允许的源（origins）
- 支持精确匹配：`http://localhost:3000`
- 支持通配符子域名：`*.example.com`
- 预检请求缓存：86400秒（1天）

### 内容安全
- **Markdown验证**: 自动检测和清理危险的HTML标签和脚本
  - 拒绝：`<script>`, `<iframe>`, `<embed>`, `<object>` 等标签
  - 拒绝：`javascript:`, `data:text/html` 等危险协议
  - 拒绝：事件处理器（`onclick`等）
- **IP脱敏**: 评论IP地址返回时自动脱敏（最后一段用*替代）
- **文件类型验证**: 上传文件时验证MIME类型和扩展名
- **文件大小限制**: 可配置的文件大小上限（默认10MB）

### 请求限制
- **请求体大小**: 最大30MB
- **超时控制**: 合理的请求超时时间
- **SQL注入防护**: 使用GORM参数化查询

## 后台任务

系统启动时会自动运行以下后台任务：

### 定时清理任务
- **执行周期**: 根据配置的Cron表达式（默认：每周日凌晨3点）
- **清理内容**:
  - 已删除的媒体文件（物理文件）
  - 孤立的媒体记录（数据库中存在但文件已丢失）
  - 过期的临时文件
- **配置项**: `scheduler.enabled` 和 `scheduler.cleanup_cron`

### 限流器清理
- **评论限流器**: 每10分钟清理2小时前的记录
- **全局限流器**: 每10分钟清理过期的IP限流记录
- **目的**: 防止内存无限增长

### 文件删除队列
- **异步处理**: 文件删除操作加入队列，异步执行
- **防止阻塞**: 避免API请求因文件IO操作而阻塞
- **自动重试**: 删除失败时可重试

### 日志轮转
- **自动轮转**: 基于配置的文件大小和保留策略
- **压缩归档**: 旧日志自动压缩
- **自动清理**: 超过保留期的日志自动删除

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

## Docker 部署说明

### Docker镜像构建

```bash
# 基础构建
docker build -t klog-backend:latest .

# 带版本信息构建
docker build \
  --build-arg VERSION=v1.0.0 \
  --build-arg BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  -t klog-backend:v1.0.0 .

# 使用Makefile构建
make docker
```

### Docker Compose 配置

项目提供了灵活的 Docker Compose 配置，支持多种部署场景。

#### 场景1：仅后端服务（SQLite）

```bash
docker-compose up -d
```

服务访问：`http://localhost:8010`

#### 场景2：后端 + Redis

启用Redis缓存提升性能：

需要在 `configs/config.toml` 中配置Redis：
```toml
[redis]
addr = "klog-redis:6379"
password = "klog123456"
```

#### 场景3：后端 + MySQL

使用MySQL作为主数据库：

需要在 `configs/config.toml` 中配置MySQL：
```toml
[database]
type = "mysql"
url = "klog:klogpassword@tcp(klog-mysql:3306)/klog?charset=utf8mb4&parseTime=True&loc=Local"
```






### 健康检查

所有服务都配置了健康检查：

```bash
# 查看容器健康状态
docker-compose ps

# 手动健康检查
curl http://localhost:8010/health
curl http://localhost:8010/health/live
curl http://localhost:8010/health/ready
curl http://localhost:8010/metrics
```

健康检查端点说明：
- `/health` - 完整健康检查（包含数据库和Redis）
- `/health/live` - 存活检查（用于Docker和K8s liveness probe）
- `/health/ready` - 就绪检查（用于K8s readiness probe）
- `/metrics` - 应用指标（内存、goroutine、数据库连接池等）

### 数据持久化

Docker Compose配置了以下数据卷：

```bash
# 查看数据卷
docker volume ls | grep klog

# 数据备份（示例：SQLite）
docker cp klog-backend:/app/db ./backup/

# 数据恢复
docker cp ./backup/db klog-backend:/app/

# MySQL数据备份
docker exec klog-mysql mysqldump -u root -p'rootpassword' klog > backup.sql

# MySQL数据恢复
docker exec -i klog-mysql mysql -u root -p'rootpassword' klog < backup.sql
```

### 日志管理

```bash
# 查看实时日志
docker-compose logs -f klog-backend

# 查看最近100行日志
docker-compose logs --tail=100 klog-backend

# 查看所有服务日志
docker-compose logs -f

# 应用内日志文件位置
./logs/app.log
```

### Nginx反向代理配置示例

创建 `nginx/nginx.conf`:

```nginx
user nginx;
worker_processes auto;
error_log /var/log/nginx/error.log warn;
pid /var/run/nginx.pid;

events {
    worker_connections 1024;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';

    access_log /var/log/nginx/access.log main;

    sendfile on;
    tcp_nopush on;
    keepalive_timeout 65;
    gzip on;

    include /etc/nginx/conf.d/*.conf;
}
```

创建 `nginx/conf.d/klog.conf`:

```nginx
upstream klog_backend {
    server klog-backend:8010;
}

server {
    listen 80;
    server_name localhost;
    client_max_body_size 30M;

    # 健康检查
    location /health {
        proxy_pass http://klog_backend;
        access_log off;
    }

    # API代理
    location /api/ {
        proxy_pass http://klog_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 静态文件（上传的文件）
    location /uploads/ {
        alias /usr/share/nginx/html/uploads/;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }

    # 其他路径转发到后端
    location / {
        proxy_pass http://klog_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## 许可证

MIT License

