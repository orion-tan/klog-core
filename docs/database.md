# 数据库设计与管理

## 数据库表结构

### 核心表

#### 1. users - 用户表

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | BIGINT | 主键，自增 |
| username | VARCHAR(50) | 用户名，唯一 |
| password | VARCHAR(255) | 密码哈希 |
| email | VARCHAR(100) | 邮箱，唯一 |
| nickname | VARCHAR(50) | 昵称 |
| avatar_url | VARCHAR(255) | 头像URL |
| role | VARCHAR(20) | 角色：author/admin |
| status | VARCHAR(20) | 状态：active/inactive |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

#### 2. posts - 文章表

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | BIGINT | 主键，自增 |
| category_id | BIGINT | 分类ID，外键 |
| author_id | BIGINT | 作者ID，外键 |
| title | VARCHAR(255) | 标题 |
| slug | VARCHAR(255) | URL别名，唯一 |
| content | LONGTEXT | 文章内容 |
| excerpt | TEXT | 摘要 |
| cover_image_url | VARCHAR(255) | 封面图片 |
| status | VARCHAR(20) | 状态：draft/published/archived |
| view_count | BIGINT | 浏览量 |
| published_at | TIMESTAMP | 发布时间 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

#### 3. categories - 分类表

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | BIGINT | 主键，自增 |
| name | VARCHAR(50) | 分类名称，唯一 |
| slug | VARCHAR(50) | URL别名，唯一 |
| description | VARCHAR(255) | 描述 |

#### 4. tags - 标签表

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | BIGINT | 主键，自增 |
| name | VARCHAR(50) | 标签名称，唯一 |
| slug | VARCHAR(50) | URL别名，唯一 |

#### 5. post_tags - 文章标签关联表

| 字段名 | 类型 | 说明 |
|--------|------|------|
| post_id | BIGINT | 文章ID，复合主键 |
| tag_id | BIGINT | 标签ID，复合主键 |

#### 6. comments - 评论表

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | BIGINT | 主键，自增 |
| post_id | BIGINT | 文章ID，外键 |
| user_id | BIGINT | 用户ID，外键（可为空） |
| name | VARCHAR(50) | 游客名称 |
| email | VARCHAR(100) | 游客邮箱 |
| content | TEXT | 评论内容 |
| ip | VARCHAR(100) | IP地址 |
| status | VARCHAR(20) | 状态：pending/approved/spam |
| parent_id | BIGINT | 父评论ID |
| created_at | TIMESTAMP | 创建时间 |

### 扩展表

#### 7. media - 媒体库表

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | BIGINT | 主键，自增 |
| uploader_id | BIGINT | 上传者ID，外键 |
| file_name | VARCHAR(255) | 文件名 |
| file_path | VARCHAR(255) | 文件路径 |
| mime_type | VARCHAR(100) | 文件类型 |
| size | BIGINT | 文件大小（字节） |
| created_at | TIMESTAMP | 创建时间 |

#### 8. settings - 系统设置表

| 字段名 | 类型 | 说明 |
|--------|------|------|
| key | VARCHAR(50) | 配置键，主键 |
| value | TEXT | 配置值 |
| type | VARCHAR(20) | 值类型：str/number/json |

#### 9. user_identities - 第三方关联表

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | BIGINT | 主键，自增 |
| user_id | BIGINT | 用户ID，外键 |
| provider | VARCHAR(50) | 第三方平台名称 |
| provider_id | VARCHAR(255) | 第三方平台用户ID |

## 数据库迁移

应用启动时会自动执行数据库迁移，创建所需的表结构。

迁移代码位于：`internal/database/migrate.go`

```go
func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &model.User{},
        &model.Post{},
        &model.Category{},
        &model.Tag{},
        &model.PostTag{},
        &model.Comment{},
        &model.Setting{},
        &model.Media{},
        &model.UserIdentity{},
    )
}
```

## 初始化数据

### 创建管理员账户

首次部署后，需要手动创建管理员账户：

#### 方法1：通过API创建后修改

1. 使用注册接口创建用户
2. 进入数据库修改角色为admin

SQLite:
```bash
sqlite3 ./db/klog.db
```

```sql
UPDATE users SET role = 'admin' WHERE username = 'admin';
.exit
```

MySQL:
```sql
USE klog;
UPDATE users SET role = 'admin' WHERE username = 'admin';
```

#### 方法2：直接在数据库中创建

SQLite:
```bash
sqlite3 ./db/klog.db
```

```sql
-- 注意：密码是 'admin123456' 的bcrypt哈希
INSERT INTO users (username, email, password, nickname, role, status, created_at, updated_at)
VALUES (
    'admin',
    'admin@example.com',
    '$2a$10$XYZ...', -- 需要生成bcrypt哈希
    '管理员',
    'admin',
    'active',
    datetime('now'),
    datetime('now')
);
```

#### 方法3：使用Go脚本生成密码哈希

创建 `tools/create_admin.go`：

```go
package main

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    password := "admin123456"
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Password hash: %s\n", string(hash))
}
```

运行：
```bash
go run tools/create_admin.go
```

### 创建初始分类和标签

```sql
-- 创建分类
INSERT INTO categories (name, slug, description) VALUES
('技术', 'tech', '技术相关文章'),
('生活', 'life', '生活随笔'),
('思考', 'thinking', '思考与感悟');

-- 创建标签
INSERT INTO tags (name, slug) VALUES
('Go', 'golang'),
('JavaScript', 'javascript'),
('Python', 'python'),
('Web开发', 'web-dev'),
('数据库', 'database'),
('Docker', 'docker');
```

## 数据库维护

### 备份数据库

#### SQLite备份

```bash
# 简单复制
cp ./db/klog.db ./backups/klog_$(date +%Y%m%d_%H%M%S).db

# 使用sqlite3备份
sqlite3 ./db/klog.db ".backup './backups/klog_$(date +%Y%m%d_%H%M%S).db'"
```

#### MySQL备份

```bash
# 导出数据
mysqldump -u klog -p klog > klog_backup_$(date +%Y%m%d_%H%M%S).sql

# 导入数据
mysql -u klog -p klog < klog_backup.sql
```

### 恢复数据库

#### SQLite恢复

```bash
# 停止应用
sudo systemctl stop klog

# 恢复数据库
cp ./backups/klog_20240101_120000.db ./db/klog.db

# 启动应用
sudo systemctl start klog
```

#### MySQL恢复

```bash
mysql -u klog -p klog < klog_backup.sql
```

### 清理旧数据

#### 删除已归档的文章

```sql
DELETE FROM posts WHERE status = 'archived' AND updated_at < DATE_SUB(NOW(), INTERVAL 1 YEAR);
```

#### 删除垃圾评论

```sql
DELETE FROM comments WHERE status = 'spam' AND created_at < DATE_SUB(NOW(), INTERVAL 30 DAY);
```

#### 清理未使用的标签

```sql
DELETE FROM tags WHERE id NOT IN (SELECT DISTINCT tag_id FROM post_tags);
```

## 数据库优化

### 索引优化

GORM会自动为外键和唯一字段创建索引，但你可以手动添加额外的索引：

```sql
-- 文章状态索引（如果还没有）
CREATE INDEX idx_posts_status ON posts(status);

-- 文章发布时间索引
CREATE INDEX idx_posts_published_at ON posts(published_at);

-- 评论状态索引
CREATE INDEX idx_comments_status ON comments(status);

-- 用户角色索引
CREATE INDEX idx_users_role ON users(role);
```

### 查询优化

使用EXPLAIN分析慢查询：

```sql
EXPLAIN SELECT * FROM posts 
WHERE status = 'published' 
ORDER BY published_at DESC 
LIMIT 10;
```

### SQLite特定优化

编辑 `cmd/main.go` 添加SQLite优化配置：

```go
db, err := gorm.Open(sqlite.Open(config.Cfg.Database.Url), &gorm.Config{})
if err != nil {
    fmt.Println("初始化数据库失败:", err)
    return
}

// SQLite优化
sqlDB, _ := db.DB()
sqlDB.SetMaxOpenConns(1) // SQLite只支持单个写入连接
sqlDB.Exec("PRAGMA journal_mode=WAL;") // 使用WAL模式提高性能
sqlDB.Exec("PRAGMA synchronous=NORMAL;") // 提高写入性能
sqlDB.Exec("PRAGMA cache_size=-64000;") // 设置缓存为64MB
```

### MySQL特定优化

```sql
-- 优化表
OPTIMIZE TABLE posts;
OPTIMIZE TABLE comments;
OPTIMIZE TABLE users;

-- 分析表
ANALYZE TABLE posts;
ANALYZE TABLE comments;
```

## 数据库监控

### SQLite监控

```bash
# 查看数据库大小
ls -lh ./db/klog.db

# 查看数据库统计
sqlite3 ./db/klog.db "
SELECT 
    name as table_name,
    (SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=m.name) as record_count
FROM sqlite_master m
WHERE type='table'
ORDER BY name;
"

# 检查数据库完整性
sqlite3 ./db/klog.db "PRAGMA integrity_check;"
```

### MySQL监控

```sql
-- 查看表大小
SELECT 
    table_name,
    ROUND(((data_length + index_length) / 1024 / 1024), 2) AS size_mb
FROM information_schema.TABLES
WHERE table_schema = 'klog'
ORDER BY size_mb DESC;

-- 查看记录数
SELECT 
    table_name,
    table_rows
FROM information_schema.TABLES
WHERE table_schema = 'klog'
ORDER BY table_rows DESC;
```

## 切换数据库

### 从SQLite切换到MySQL

1. 备份SQLite数据
2. 安装MySQL并创建数据库
3. 修改配置文件
4. 修改代码中的数据库驱动
5. 运行应用（自动创建表结构）
6. 迁移数据

数据迁移脚本示例：

```bash
# 导出SQLite数据为SQL
sqlite3 ./db/klog.db .dump > sqlite_dump.sql

# 转换并导入MySQL（需要手动处理一些语法差异）
# 或使用专门的迁移工具
```

## 常见问题

### 数据库锁定（SQLite）

如果遇到"database is locked"错误：

1. 确保只有一个进程访问数据库
2. 使用WAL模式：`PRAGMA journal_mode=WAL;`
3. 考虑切换到MySQL/PostgreSQL

### 外键约束错误

确保在删除数据前处理好外键关联：

```sql
-- 删除文章前先删除相关评论和标签关联
DELETE FROM comments WHERE post_id = 1;
DELETE FROM post_tags WHERE post_id = 1;
DELETE FROM posts WHERE id = 1;
```

### 字符编码问题

确保数据库使用UTF-8编码：

MySQL:
```sql
ALTER DATABASE klog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

## 最佳实践

1. **定期备份**：每天自动备份数据库
2. **监控性能**：定期检查慢查询
3. **优化索引**：根据查询模式添加合适的索引
4. **数据清理**：定期清理无用数据
5. **版本控制**：使用数据库迁移工具管理表结构变更
6. **测试恢复**：定期测试备份恢复流程

