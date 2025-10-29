# API 测试示例

这个文档包含了所有API端点的测试示例。

## 环境变量

为了方便测试，先设置一些环境变量：

```bash
export BASE_URL="http://localhost:8010/api/v1"
export TOKEN="your_jwt_token_here"
```

## 1. 认证接口

### 1.1 用户注册

```bash
curl -X POST $BASE_URL/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "admin123456",
    "nickname": "管理员"
  }'
```

### 1.2 用户登录

```bash
curl -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "login": "admin",
    "password": "admin123456"
  }'
```

保存返回的token：
```bash
export TOKEN="eyJhbGc..."
```

### 1.3 获取当前用户信息

```bash
curl -X GET $BASE_URL/auth/me \
  -H "Authorization: Bearer $TOKEN"
```

## 2. 分类接口

### 2.1 创建分类

```bash
curl -X POST $BASE_URL/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "技术",
    "slug": "tech",
    "description": "技术相关文章"
  }'
```

```bash
curl -X POST $BASE_URL/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "生活",
    "slug": "life",
    "description": "生活随笔"
  }'
```

### 2.2 获取所有分类

```bash
curl -X GET $BASE_URL/categories
```

### 2.3 更新分类

```bash
curl -X PUT $BASE_URL/categories/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "技术分享",
    "description": "技术相关的分享文章"
  }'
```

### 2.4 删除分类

```bash
curl -X DELETE $BASE_URL/categories/2 \
  -H "Authorization: Bearer $TOKEN"
```

## 3. 标签接口

### 3.1 创建标签

```bash
curl -X POST $BASE_URL/tags \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Go语言",
    "slug": "golang"
  }'
```

```bash
curl -X POST $BASE_URL/tags \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Web开发",
    "slug": "web-development"
  }'
```

```bash
curl -X POST $BASE_URL/tags \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "后端",
    "slug": "backend"
  }'
```

### 3.2 获取所有标签

```bash
curl -X GET $BASE_URL/tags
```

### 3.3 更新标签

```bash
curl -X PUT $BASE_URL/tags/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Golang",
    "slug": "golang"
  }'
```

### 3.4 删除标签

```bash
curl -X DELETE $BASE_URL/tags/3 \
  -H "Authorization: Bearer $TOKEN"
```

## 4. 文章接口

### 4.1 创建文章（草稿）

```bash
curl -X POST $BASE_URL/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "我的第一篇博客",
    "slug": "my-first-blog",
    "content": "# 欢迎来到我的博客\n\n这是我的第一篇博客文章。",
    "excerpt": "这是我的第一篇博客文章",
    "status": "draft",
    "category_id": 1,
    "tags": ["golang", "web-development"]
  }'
```

### 4.2 创建文章（已发布）

```bash
curl -X POST $BASE_URL/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "Go Web开发入门",
    "slug": "go-web-development-intro",
    "content": "# Go Web开发入门\n\n本文将介绍如何使用Go语言开发Web应用...",
    "excerpt": "介绍如何使用Go语言开发Web应用",
    "cover_image_url": "/uploads/cover.jpg",
    "status": "published",
    "category_id": 1,
    "tags": ["golang", "web-development"]
  }'
```

### 4.3 获取文章列表

获取所有已发布的文章：
```bash
curl -X GET "$BASE_URL/posts?page=1&limit=10"
```

按分类筛选：
```bash
curl -X GET "$BASE_URL/posts?category=tech"
```

按标签筛选：
```bash
curl -X GET "$BASE_URL/posts?tag=golang"
```

获取自己的所有文章（需要认证）：
```bash
curl -X GET "$BASE_URL/posts?status=draft" \
  -H "Authorization: Bearer $TOKEN"
```

### 4.4 获取文章详情

```bash
curl -X GET $BASE_URL/posts/1
```

### 4.5 更新文章

```bash
curl -X PUT $BASE_URL/posts/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "我的第一篇博客（已更新）",
    "status": "published"
  }'
```

### 4.6 删除文章

```bash
curl -X DELETE $BASE_URL/posts/1 \
  -H "Authorization: Bearer $TOKEN"
```

## 5. 评论接口

### 5.1 发表评论（已登录用户）

注意：POST_ID 是文章的ID，比如 1、2、3 等

```bash
curl -X POST $BASE_URL/posts/1/comments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "content": "写得很好，学到了很多！"
  }'
```

### 5.2 发表评论（游客）

```bash
curl -X POST $BASE_URL/posts/1/comments \
  -H "Content-Type: application/json" \
  -d '{
    "content": "非常实用的教程！",
    "name": "访客",
    "email": "visitor@example.com"
  }'
```

### 5.3 发表回复评论

```bash
curl -X POST $BASE_URL/posts/1/comments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "content": "谢谢你的评论！",
    "parent_id": 1
  }'
```

### 5.4 获取文章评论列表

```bash
curl -X GET $BASE_URL/posts/1/comments
```

### 5.5 更新评论状态

```bash
curl -X PUT $BASE_URL/comments/2 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "status": "approved"
  }'
```

### 5.6 删除评论

```bash
curl -X DELETE $BASE_URL/comments/2 \
  -H "Authorization: Bearer $TOKEN"
```

## 6. 媒体库接口

### 6.1 上传文件（multipart）

```bash
curl -X POST $BASE_URL/media/upload \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@/path/to/image.jpg"
```

### 6.2 上传文件（base64）

```bash
curl -X POST $BASE_URL/media/upload \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "file_name": "test.jpg",
    "mime_type": "image/jpeg",
    "data": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="
  }'
```

### 6.3 获取媒体列表

```bash
curl -X GET "$BASE_URL/media?page=1&limit=20" \
  -H "Authorization: Bearer $TOKEN"
```

### 6.4 删除媒体文件

```bash
curl -X DELETE $BASE_URL/media/1 \
  -H "Authorization: Bearer $TOKEN"
```

### 6.5 访问上传的文件

```bash
curl -X GET "http://localhost:8010/uploads/filename.jpg"
```

## 7. 用户接口

### 7.1 获取用户列表（管理员）

```bash
curl -X GET "$BASE_URL/users?page=1&limit=20" \
  -H "Authorization: Bearer $TOKEN"
```

### 7.2 获取用户信息

```bash
curl -X GET $BASE_URL/users/1
```

### 7.3 更新用户信息

更新自己的信息：
```bash
curl -X PUT $BASE_URL/users/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "nickname": "新昵称",
    "avatar_url": "/uploads/avatar.jpg"
  }'
```

管理员禁用用户：
```bash
curl -X PUT $BASE_URL/users/2 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "status": "inactive"
  }'
```

## 测试脚本

创建一个完整的测试脚本 `test.sh`：

```bash
#!/bin/bash

BASE_URL="http://localhost:8010/api/v1"

echo "=== 1. 注册用户 ==="
REGISTER_RESPONSE=$(curl -s -X POST $BASE_URL/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "test123456",
    "nickname": "测试用户"
  }')
echo $REGISTER_RESPONSE | jq

echo -e "\n=== 2. 登录 ==="
LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "login": "testuser",
    "password": "test123456"
  }')
echo $LOGIN_RESPONSE | jq

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.token')
echo "Token: $TOKEN"

echo -e "\n=== 3. 获取当前用户信息 ==="
curl -s -X GET $BASE_URL/auth/me \
  -H "Authorization: Bearer $TOKEN" | jq

echo -e "\n=== 4. 创建分类 ==="
curl -s -X POST $BASE_URL/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "测试分类",
    "slug": "test-category",
    "description": "这是一个测试分类"
  }' | jq

echo -e "\n=== 5. 创建标签 ==="
curl -s -X POST $BASE_URL/tags \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "测试标签",
    "slug": "test-tag"
  }' | jq

echo -e "\n=== 6. 创建文章 ==="
curl -s -X POST $BASE_URL/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "测试文章",
    "slug": "test-post",
    "content": "这是一篇测试文章",
    "excerpt": "测试摘要",
    "status": "published",
    "category_id": 1,
    "tags": ["test-tag"]
  }' | jq

echo -e "\n=== 7. 获取文章列表 ==="
curl -s -X GET "$BASE_URL/posts?page=1&limit=10" | jq

echo -e "\n测试完成！"
```

使用方法：
```bash
chmod +x test.sh
./test.sh
```

注意：需要安装 `jq` 工具来格式化JSON输出：
```bash
# Ubuntu/Debian
sudo apt-get install jq

# macOS
brew install jq

# CentOS/RHEL
sudo yum install jq
```

