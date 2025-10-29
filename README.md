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
port = 8010

[database]
type = "sqlite"
url = "./db/klog.db"

[jwt]
secret = "your-very-secret-key"  # 请修改为强密码(不低于32位)
expire_hour = 72

[redis]
addr = "localhost:6379"
password = ""
```

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
| GET | /media/i/:filename | 访问上传的文件 | 否 |

#### 用户接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /users | 获取用户列表 | 是（管理员）|
| GET | /users/:id | 获取用户信息 | 否 |
| PUT | /users/:id | 更新用户信息 | 是 |

### 请求示例

`todo ...` 

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
- file_hash: 文件内容hash
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

适合小型项目或开发测试：

```bash
docker-compose up -d
```

服务访问：`http://localhost:8010`

#### 场景2：后端 + Redis

启用Redis缓存提升性能：

```bash
docker-compose --profile full up -d redis klog-backend
```

需要在 `configs/config.toml` 中配置Redis：
```toml
[redis]
enabled = true
addr = "klog-redis:6379"
password = "klog123456"
db = 0
```

#### 场景3：后端 + MySQL

使用MySQL作为主数据库：

```bash
docker-compose --profile mysql up -d mysql klog-backend
```

需要在 `configs/config.toml` 中配置MySQL：
```toml
[database]
type = "mysql"
url = "klog:klogpassword@tcp(klog-mysql:3306)/klog?charset=utf8mb4&parseTime=True&loc=Local"
```

#### 场景4：完整服务栈

包含后端、Redis、MySQL和Nginx反向代理：

```bash
# 1. 创建nginx配置目录
mkdir -p nginx/conf.d nginx/ssl nginx/logs

# 2. 创建nginx配置文件（参考下方配置示例）
# 编辑 nginx/nginx.conf

# 3. 启动所有服务
docker-compose --profile full up -d

# 4. 查看服务状态
docker-compose ps
```

访问：
- API: `http://localhost:8010`（直接访问）
- Nginx: `http://localhost` 或 `http://localhost:80`（通过代理访问）

### 环境变量配置

可以通过环境变量覆盖默认配置。创建 `.env` 文件：

```bash
# 服务端口配置
SERVER_PORT=8010
NGINX_HTTP_PORT=80
NGINX_HTTPS_PORT=443
REDIS_PORT=6379
MYSQL_PORT=3306

# Gin运行模式
GIN_MODE=release

# 构建信息
VERSION=v1.0.0
BUILD_TIME=2025-10-29T00:00:00Z

# MySQL配置
MYSQL_ROOT_PASSWORD=your_strong_root_password
MYSQL_DATABASE=klog
MYSQL_USER=klog
MYSQL_PASSWORD=your_strong_password

# Redis配置
REDIS_PASSWORD=your_redis_password
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

### 常用Docker命令

```bash
# 启动服务
docker-compose up -d

# 重启服务
docker-compose restart klog-backend

# 停止服务
docker-compose stop

# 停止并删除容器
docker-compose down

# 重新构建并启动
docker-compose up -d --build

# 进入容器
docker exec -it klog-backend sh

# 查看容器资源使用
docker stats klog-backend

# 清理未使用的镜像和容器
docker system prune -a
```

## 生产环境部署

### 部署前检查清单

- [ ] 修改JWT密钥为强随机密码
- [ ] 修改数据库密码（如使用MySQL）
- [ ] 修改Redis密码
- [ ] 配置HTTPS证书
- [ ] 设置合理的日志轮转策略
- [ ] 配置防火墙规则
- [ ] 设置资源限制（CPU、内存）
- [ ] 配置定期数据备份
- [ ] 启用监控和告警

### 方式一：Docker Compose部署（推荐）

```bash
# 1. 克隆代码
git clone <repository-url>
cd klog/backend

# 2. 配置环境变量
cp .env.example .env
vim .env  # 修改敏感信息

# 3. 修改配置文件
vim configs/config.toml

# 4. 启动服务
docker-compose --profile full up -d

# 5. 验证服务
curl http://localhost:8010/health
```

### 方式二：传统部署

#### 1. 修改配置

编辑 `configs/config.toml`，修改为生产环境配置。

#### 2. 编译

```bash
make build-linux
# 或
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o klog-backend cmd/main.go
```

#### 3. 使用Systemd管理

创建 `/etc/systemd/system/klog-backend.service`:

```ini
[Unit]
Description=KLog Backend Service
After=network.target mysql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/klog-backend
ExecStart=/opt/klog-backend/klog-backend
Restart=on-failure
RestartSec=5s

# 资源限制
LimitNOFILE=65535
LimitNPROC=4096

[Install]
WantedBy=multi-user.target
```

启动服务：
```bash
sudo systemctl daemon-reload
sudo systemctl enable klog-backend
sudo systemctl start klog-backend
sudo systemctl status klog-backend
```

#### 4. Nginx反向代理

参考上方Nginx配置示例，配置反向代理和HTTPS。

### 性能优化建议

1. **数据库优化**
   - 生产环境使用MySQL或PostgreSQL
   - 配置数据库连接池
   - 添加必要的索引

2. **启用Redis缓存**
   - 缓存热点数据
   - 减少数据库查询

3. **静态资源**
   - 使用CDN加速上传的媒体文件
   - 启用Nginx gzip压缩

4. **负载均衡**
   - 多实例部署
   - 使用Nginx或云负载均衡器

5. **监控**
   - 配置应用性能监控（APM）
   - 设置日志收集和分析
   - 配置告警规则

### 安全建议

1. **网络安全**
   - 使用HTTPS（配置SSL证书）
   - 配置防火墙，只开放必要端口
   - 使用安全的密码和密钥

2. **应用安全**
   - 定期更新依赖包
   - 启用CORS限制
   - 配置请求速率限制
   - 验证和过滤用户输入

3. **数据安全**
   - 定期备份数据库
   - 加密敏感信息
   - 配置数据库访问权限

### Kubernetes部署（可选）

如需要部署到Kubernetes，可参考以下配置：

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: klog-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: klog-backend
  template:
    metadata:
      labels:
        app: klog-backend
    spec:
      containers:
      - name: klog-backend
        image: klog-backend:latest
        ports:
        - containerPort: 8010
        env:
        - name: GIN_MODE
          value: "release"
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8010
          initialDelaySeconds: 5
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8010
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

## 许可证

MIT License

