# 生产环境部署指南

本指南将帮助你将KLog博客系统后端部署到生产环境。

## 准备工作

### 系统要求

- Linux服务器（推荐Ubuntu 20.04+）
- Go 1.21+（如果需要在服务器上编译）
- 至少512MB内存
- 至少1GB磁盘空间

### 可选组件

- Nginx（作为反向代理）
- Systemd（进程管理）
- MySQL/PostgreSQL（如果不使用SQLite）

## 部署步骤

### 1. 编译应用

#### 在本地编译

如果你在本地开发环境编译，然后上传到服务器：

```bash
# Linux AMD64
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o klog-backend cmd/main.go

# Linux ARM64（树莓派等）
CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o klog-backend cmd/main.go
```

注意：SQLite需要CGO支持，所以 `CGO_ENABLED=1` 是必需的。

#### 在服务器上编译

```bash
# 安装Go（如果还没有）
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# 编译
cd /path/to/klog/backend
go build -o klog-backend cmd/main.go
```

### 2. 上传文件到服务器

创建应用目录：

```bash
sudo mkdir -p /opt/klog
sudo chown $USER:$USER /opt/klog
```

上传文件：

```bash
# 使用scp
scp -r klog-backend configs/ your-server:/opt/klog/

# 或使用rsync
rsync -avz --exclude='*.db' --exclude='uploads/' \
  klog-backend configs/ your-server:/opt/klog/
```

### 3. 配置应用

编辑配置文件：

```bash
cd /opt/klog
nano configs/config.toml
```

生产环境配置示例：

```toml
[server]
port = 8010

[database]
type = "sqlite"
url = "/opt/klog/db/klog.db"

# 如果使用MySQL
# [database]
# type = "mysql"
# url = "username:password@tcp(localhost:3306)/klog?charset=utf8mb4&parseTime=True&loc=Local"

[jwt]
# 使用强密码！
secret = "your-production-secret-key-change-this"
expire_hour = 72
```

### 4. 创建必要的目录

```bash
mkdir -p /opt/klog/db
mkdir -p /opt/klog/uploads
mkdir -p /opt/klog/logs
```

### 5. 使用Systemd管理进程

创建systemd服务文件：

```bash
sudo nano /etc/systemd/system/klog.service
```

内容：

```ini
[Unit]
Description=KLog Backend Service
After=network.target

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/opt/klog
ExecStart=/opt/klog/klog-backend
Restart=on-failure
RestartSec=5s

# 环境变量
Environment="GIN_MODE=release"

# 日志
StandardOutput=append:/opt/klog/logs/stdout.log
StandardError=append:/opt/klog/logs/stderr.log

[Install]
WantedBy=multi-user.target
```

设置权限：

```bash
sudo chown -R www-data:www-data /opt/klog
sudo chmod +x /opt/klog/klog-backend
```

启动服务：

```bash
# 重新加载systemd配置
sudo systemctl daemon-reload

# 启动服务
sudo systemctl start klog

# 查看状态
sudo systemctl status klog

# 设置开机自启
sudo systemctl enable klog

# 查看日志
sudo journalctl -u klog -f
```

常用命令：

```bash
# 停止服务
sudo systemctl stop klog

# 重启服务
sudo systemctl restart klog

# 查看日志
sudo tail -f /opt/klog/logs/stdout.log
sudo tail -f /opt/klog/logs/stderr.log
```

### 6. 配置Nginx反向代理

安装Nginx：

```bash
sudo apt update
sudo apt install nginx
```

创建Nginx配置：

```bash
sudo nano /etc/nginx/sites-available/klog
```

基础配置：

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 日志
    access_log /var/log/nginx/klog-access.log;
    error_log /var/log/nginx/klog-error.log;

    # API代理
    location /api/ {
        proxy_pass http://localhost:8010;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }

    # 上传文件
    location /uploads/ {
        alias /opt/klog/uploads/;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }

    # 客户端最大上传大小
    client_max_body_size 10M;
}
```

启用配置：

```bash
sudo ln -s /etc/nginx/sites-available/klog /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 7. 配置HTTPS（使用Let's Encrypt）

安装Certbot：

```bash
sudo apt install certbot python3-certbot-nginx
```

获取SSL证书：

```bash
sudo certbot --nginx -d your-domain.com
```

Certbot会自动修改Nginx配置并设置自动续期。

完整的HTTPS配置：

```nginx
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    # SSL证书
    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;
    
    # SSL配置
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # 日志
    access_log /var/log/nginx/klog-access.log;
    error_log /var/log/nginx/klog-error.log;

    # API代理
    location /api/ {
        proxy_pass http://localhost:8010;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }

    # 上传文件
    location /uploads/ {
        alias /opt/klog/uploads/;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }

    client_max_body_size 10M;
}
```

## 使用MySQL数据库

### 1. 安装MySQL

```bash
sudo apt install mysql-server
sudo mysql_secure_installation
```

### 2. 创建数据库和用户

```bash
sudo mysql
```

在MySQL命令行中：

```sql
CREATE DATABASE klog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'klog'@'localhost' IDENTIFIED BY 'strong_password_here';
GRANT ALL PRIVILEGES ON klog.* TO 'klog'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

### 3. 修改配置

修改 `configs/config.toml`：

```toml
[database]
type = "mysql"
url = "klog:strong_password_here@tcp(localhost:3306)/klog?charset=utf8mb4&parseTime=True&loc=Local"
```

### 4. 更新代码

修改 `cmd/main.go`，将SQLite驱动改为MySQL：

```go
import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func main() {
    // ...
    db, err := gorm.Open(mysql.Open(config.Cfg.Database.Url), &gorm.Config{})
    // ...
}
```

并在 `go.mod` 中添加MySQL驱动：

```bash
go get -u gorm.io/driver/mysql
```

## 安全建议

### 1. 防火墙配置

```bash
# 只开放必要的端口
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 22/tcp
sudo ufw enable
```

### 2. 修改默认配置

- 修改JWT密钥为强随机字符串
- 使用环境变量存储敏感信息
- 限制上传文件大小和类型

### 3. 定期备份

创建备份脚本 `/opt/klog/backup.sh`：

```bash
#!/bin/bash

BACKUP_DIR="/opt/klog/backups"
DATE=$(date +%Y%m%d_%H%M%S)

# 创建备份目录
mkdir -p $BACKUP_DIR

# 备份数据库
cp /opt/klog/db/klog.db $BACKUP_DIR/klog_$DATE.db

# 备份上传文件
tar -czf $BACKUP_DIR/uploads_$DATE.tar.gz /opt/klog/uploads

# 删除30天前的备份
find $BACKUP_DIR -type f -mtime +30 -delete

echo "Backup completed: $DATE"
```

设置定时任务：

```bash
chmod +x /opt/klog/backup.sh
crontab -e
```

添加：
```
0 2 * * * /opt/klog/backup.sh >> /opt/klog/logs/backup.log 2>&1
```

### 4. 日志轮转

创建 `/etc/logrotate.d/klog`：

```
/opt/klog/logs/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0640 www-data www-data
    postrotate
        systemctl reload klog > /dev/null 2>&1 || true
    endscript
}
```

## 监控和维护

### 1. 监控应用状态

```bash
# 查看服务状态
sudo systemctl status klog

# 查看实时日志
sudo journalctl -u klog -f

# 检查端口
sudo netstat -tlnp | grep 8010
```

### 2. 性能监控

安装htop监控系统资源：

```bash
sudo apt install htop
htop
```

### 3. 更新应用

```bash
# 编译新版本
go build -o klog-backend cmd/main.go

# 上传到服务器
scp klog-backend your-server:/opt/klog/

# 重启服务
sudo systemctl restart klog
```

## 故障排查

### 应用无法启动

1. 检查日志：
```bash
sudo journalctl -u klog -n 50
tail -f /opt/klog/logs/stderr.log
```

2. 检查配置文件：
```bash
cat /opt/klog/configs/config.toml
```

3. 检查权限：
```bash
ls -la /opt/klog/
```

### 502 Bad Gateway

1. 检查应用是否运行：
```bash
sudo systemctl status klog
```

2. 检查端口：
```bash
sudo netstat -tlnp | grep 8010
```

3. 检查Nginx配置：
```bash
sudo nginx -t
```

### 数据库连接失败

1. 检查数据库文件权限（SQLite）
2. 测试数据库连接（MySQL）
3. 查看数据库日志

## Docker部署（可选）

创建 `Dockerfile`：

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN apk add --no-cache gcc musl-dev
RUN go mod download
RUN CGO_ENABLED=1 go build -o klog-backend cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/klog-backend .
COPY --from=builder /app/configs ./configs

RUN mkdir -p /root/db /root/uploads

EXPOSE 8010
CMD ["./klog-backend"]
```

创建 `docker-compose.yml`：

```yaml
version: '3.8'

services:
  klog-backend:
    build: .
    ports:
      - "8010:8010"
    volumes:
      - ./db:/root/db
      - ./uploads:/root/uploads
      - ./configs:/root/configs
    environment:
      - GIN_MODE=release
    restart: unless-stopped
```

运行：

```bash
docker-compose up -d
```

## 总结

完成以上步骤后，你的KLog博客系统后端应该已经成功部署到生产环境了。记得：

- 定期备份数据
- 保持系统和依赖更新
- 监控应用性能和日志
- 使用HTTPS保护通信
- 设置强密码和安全配置

