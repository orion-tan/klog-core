# KLog API 接口文档

## 概述

KLog 是一个现代化的个人博客系统后端 API。本文档描述了所有可用的 REST API 端点。

### 基础信息

- **Base URL**: `http://your-domain.com/api/v1`
- **协议**: HTTP/HTTPS
- **数据格式**: JSON
- **认证方式**: JWT Bearer Token
- **字符编码**: UTF-8

### 认证说明

需要认证的接口需要在请求头中携带 JWT Token：

```
Authorization: Bearer <your-jwt-token>
```

### 统一响应格式

所有 API 响应采用统一的 JSON 格式：

**成功响应**:
```json
{
  "success": true,
  "data": {
    // 实际返回的数据
  }
}
```

**错误响应**:
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "错误描述信息"
  }
}
```

### 常见错误码

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| `INVALID_PARAMS` | 400 | 请求参数无效 |
| `INVALID_ID` | 400 | 无效的 ID |
| `UNAUTHORIZED` | 401 | 未认证或认证失败 |
| `FORBIDDEN` | 403 | 无权限访问 |
| `NOT_FOUND` | 404 | 资源不存在 |
| `RATE_LIMIT_EXCEEDED` | 429 | 请求频率超限 |
| `COMMENT_RATE_LIMIT_EXCEEDED` | 429 | 评论频率超限 |
| `REGISTER_FAILED` | 400 | 注册失败 |
| `LOGIN_FAILED` | 401 | 登录失败 |
| `CREATE_*_FAILED` | 500 | 创建资源失败 |
| `UPDATE_*_FAILED` | 500 | 更新资源失败 |
| `DELETE_*_FAILED` | 500 | 删除资源失败 |
| `UPLOAD_FAILED` | 500 | 文件上传失败 |
| `FILE_TOO_LARGE` | 400 | 文件大小超限 |
| `INVALID_FILE_TYPE` | 400 | 不支持的文件类型 |

#### 错误响应示例

**400 Bad Request - 参数验证失败**
```json
{
  "success": false,
  "error": {
    "code": "INVALID_PARAMS",
    "message": "标题不能为空"
  }
}
```

**401 Unauthorized - Token无效或过期**
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "未授权访问，请先登录"
  }
}
```

**403 Forbidden - 无权限操作**
```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "您没有权限执行此操作"
  }
}
```

**404 Not Found - 资源不存在**
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "文章不存在"
  }
}
```

**429 Too Many Requests - 限流**
```json
{
  "success": false,
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "请求过于频繁，请稍后再试"
  }
}
```

**500 Internal Server Error - 服务器错误**
```json
{
  "success": false,
  "error": {
    "code": "CREATE_POST_FAILED",
    "message": "创建文章失败，请稍后重试"
  }
}
```

---

## 健康检查 API

### 1. 健康检查

**端点**: `GET /health`

**描述**: 检查服务是否运行

**认证**: 无需认证

**响应**:
```json
{
  "success": true,
  "data": {
    "status": "ok"
  }
}
```

### 2. 就绪检查

**端点**: `GET /health/ready`

**描述**: 检查服务是否就绪（包括数据库连接等）

**认证**: 无需认证

### 3. 存活检查

**端点**: `GET /health/live`

**描述**: 检查服务是否存活

**认证**: 无需认证

### 4. 指标监控

**端点**: `GET /metrics`

**描述**: 获取服务监控指标

**认证**: 无需认证

---

## 认证 API

### 1. 用户注册

**端点**: `POST /api/v1/auth/register`

**描述**: 注册新用户（仅供首次设置管理员账号使用）

**认证**: 无需认证

**请求体**:
```json
{
  "username": "admin",
  "email": "admin@example.com",
  "password": "password123",
  "nickname": "管理员"
}
```

**字段说明**:
- `username` (string, 必需): 用户名，3-50 字符
- `email` (string, 必需): 邮箱地址，需符合邮箱格式
- `password` (string, 必需): 密码，8-30 字符
- `nickname` (string, 必需): 昵称，3-50 字符

**响应**: `201 Created`
```json
{
  "success": true,
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "nickname": "管理员"
  }
}
```

### 2. 用户登录

**端点**: `POST /api/v1/auth/login`

**描述**: 用户登录获取 JWT Token

**认证**: 无需认证

**请求体**:
```json
{
  "login": "admin",
  "password": "password123"
}
```

**字段说明**:
- `login` (string, 必需): 用户名或邮箱
- `password` (string, 必需): 密码

**响应**: `200 OK`
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 3. 获取当前用户信息

**端点**: `GET /api/v1/auth/me`

**描述**: 获取当前登录用户的详细信息

**认证**: 需要认证

**响应**: `200 OK`
```json
{
  "success": true,
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "nickname": "管理员",
    "bio": "这是个人简介",
    "avatar_url": "https://example.com/avatar.jpg"
  }
}
```

### 4. 用户登出

**端点**: `POST /api/v1/auth/logout`

**描述**: 登出当前用户（将 token 加入黑名单）

**认证**: 需要认证

**响应**: `200 OK`
```json
{
  "success": true,
  "data": {
    "message": "登出成功"
  }
}
```

---

## 文章 API

### 1. 获取文章列表

**端点**: `GET /api/v1/posts`

**描述**: 获取文章列表，支持分页、过滤和排序

**认证**: 可选（未认证用户只能看到已发布的文章）

**查询参数**:
| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `page` | integer | 1 | 页码（传统分页） |
| `limit` | integer | 10 | 每页数量 |
| `cursor` | string | - | 游标（Base64编码，用于游标分页） |
| `status` | string | - | 文章状态：`draft`/`published`/`archived` |
| `category` | string | - | 分类 slug |
| `tag` | string | - | 标签 slug |
| `sortBy` | string | `published_at` | 排序字段：`published_at`/`created_at`/`view_count`/`title` |
| `order` | string | `desc` | 排序方向：`asc`/`desc` |
| `detail` | integer | 0 | 是否包含文章内容：`0`/`1` |

**分页说明**:

系统支持两种分页方式：

1. **传统分页**（默认）
   - 使用 `page` 和 `limit` 参数
   - 适用于数据量较小的场景
   - 示例：`?page=1&limit=10`

2. **游标分页**（推荐用于大数据集）
   - 使用 `cursor` 和 `limit` 参数
   - 避免传统分页在大数据集时的性能问题
   - 游标值从上一页响应的 `next_cursor` 字段获取（Base64编码）
   - 示例：`?cursor=eyJpZCI6MTAsInB1Ymxpc2hlZF9hdCI6IjIwMjUtMDEtMDEifQ==&limit=10`

**注意**:
- `cursor` 和 `page` 参数互斥，同时提供时 `cursor` 优先
- 游标分页不支持跳页，只能顺序翻页
- 首次请求不需要提供 `cursor`，从响应中获取 `next_cursor` 用于下一页

**响应**: `200 OK`

传统分页响应：
```json
{
  "success": true,
  "data": {
    "total": 100,
    "page": 1,
    "limit": 10,
    "data": [
      {
        "id": 1,
        "category_id": 1,
        "author_id": 1,
        "title": "文章标题",
        "slug": "article-slug",
        "content": "",
        "excerpt": "文章摘要",
        "cover_image_url": "https://example.com/cover.jpg",
        "status": "published",
        "view_count": 100,
        "published_at": "2025-01-01T00:00:00Z",
        "created_at": "2025-01-01T00:00:00Z",
        "updated_at": "2025-01-01T00:00:00Z",
        "category": {
          "id": 1,
          "name": "技术",
          "slug": "tech",
          "description": "技术分类"
        },
        "author": {
          "id": 1,
          "username": "admin",
          "nickname": "管理员"
        },
        "tags": [
          {
            "id": 1,
            "name": "Go",
            "slug": "go"
          }
        ]
      }
    ]
  }
}
```

游标分页响应：
```json
{
  "success": true,
  "data": {
    "has_more": true,
    "next_cursor": "eyJpZCI6MTAsInB1Ymxpc2hlZF9hdCI6IjIwMjUtMDEtMDEifQ==",
    "data": [
      {
        "id": 1,
        "title": "文章标题",
        // ... 其他字段
      }
    ]
  }
}
```

**说明**:
- 当 `detail=0` 时，`content` 字段为空字符串以减少响应大小
- 游标分页时，`has_more` 表示是否有更多数据，`next_cursor` 用于获取下一页

### 2. 创建文章

**端点**: `POST /api/v1/posts`

**描述**: 创建新文章

**认证**: 需要认证

**请求体**:
```json
{
  "category_id": 1,
  "title": "文章标题",
  "slug": "article-slug",
  "content": "文章正文内容（Markdown格式）",
  "excerpt": "文章摘要",
  "cover_image_url": "https://example.com/cover.jpg",
  "status": "published",
  "tags": ["go", "backend"]
}
```

**字段说明**:
- `category_id` (integer, 可选): 分类 ID
- `title` (string, 必需): 文章标题
- `slug` (string, 必需): URL 友好的唯一标识符
- `content` (string, 必需): 文章内容
- `excerpt` (string, 可选): 摘要
- `cover_image_url` (string, 可选): 封面图 URL
- `status` (string, 必需): 状态，可选值：`draft`/`published`/`archived`
- `tags` (array, 可选): 标签 slug 数组

**响应**: `201 Created`
```json
{
  "success": true,
  "data": {
    "id": 1,
    "category_id": 1,
    "author_id": 1,
    "title": "文章标题",
    "slug": "article-slug",
    "content": "文章正文内容",
    "excerpt": "文章摘要",
    "cover_image_url": "https://example.com/cover.jpg",
    "status": "published",
    "view_count": 0,
    "published_at": "2025-01-01T00:00:00Z",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

### 3. 获取文章详情

**端点**: `GET /api/v1/posts/:id`

**描述**: 获取指定文章的详细信息

**认证**: 可选（未发布的文章需要认证）

**路径参数**:
- `id` (integer): 文章 ID

**响应**: `200 OK`
```json
{
  "success": true,
  "data": {
    "id": 1,
    "category_id": 1,
    "author_id": 1,
    "title": "文章标题",
    "slug": "article-slug",
    "content": "完整的文章内容",
    "excerpt": "文章摘要",
    "cover_image_url": "https://example.com/cover.jpg",
    "status": "published",
    "view_count": 100,
    "published_at": "2025-01-01T00:00:00Z",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z",
    "category": {
      "id": 1,
      "name": "技术",
      "slug": "tech",
      "description": "技术分类"
    },
    "author": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "nickname": "管理员"
    },
    "tags": [
      {
        "id": 1,
        "name": "Go",
        "slug": "go"
      }
    ]
  }
}
```

### 4. 更新文章

**端点**: `PUT /api/v1/posts/:id`

**描述**: 更新指定文章

**认证**: 需要认证

**路径参数**:
- `id` (integer): 文章 ID

**请求体**:
```json
{
  "category_id": 1,
  "title": "更新后的标题",
  "slug": "updated-slug",
  "content": "更新后的内容",
  "excerpt": "更新后的摘要",
  "cover_image_url": "https://example.com/new-cover.jpg",
  "status": "published",
  "tags": ["go", "api"]
}
```

**说明**: 所有字段都是可选的，只需提供要更新的字段。

**响应**: `200 OK`
```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "更新后的标题",
    // ... 其他字段
  }
}
```

### 5. 删除文章

**端点**: `DELETE /api/v1/posts/:id`

**描述**: 删除指定文章

**认证**: 需要认证

**路径参数**:
- `id` (integer): 文章 ID

**响应**: `204 No Content`

---

## 分类 API

### 1. 获取分类列表

**端点**: `GET /api/v1/categories`

**描述**: 获取所有分类

**认证**: 无需认证

**查询参数**:
| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `with_count` | string | - | 是否包含文章计数，传入任意值（如`1`或`true`）即启用 |

**响应**: `200 OK`

不带计数的响应：
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "技术",
      "slug": "tech",
      "description": "技术相关文章"
    },
    {
      "id": 2,
      "name": "生活",
      "slug": "life",
      "description": "生活随笔"
    }
  ]
}
```

带计数的响应（`?with_count=1`）：
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "技术",
      "slug": "tech",
      "description": "技术相关文章",
      "post_count": 15
    },
    {
      "id": 2,
      "name": "生活",
      "slug": "life",
      "description": "生活随笔",
      "post_count": 8
    }
  ]
}
```

**说明**: 当提供 `with_count` 参数时，响应中会包含 `post_count` 字段，显示该分类下的文章数量。

### 2. 创建分类

**端点**: `POST /api/v1/categories`

**描述**: 创建新分类

**认证**: 需要认证

**请求体**:
```json
{
  "name": "技术",
  "slug": "tech",
  "description": "技术相关文章"
}
```

**字段说明**:
- `name` (string, 必需): 分类名称
- `slug` (string, 必需): URL 友好的唯一标识符
- `description` (string, 可选): 分类描述

**响应**: `201 Created`
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "技术",
    "slug": "tech",
    "description": "技术相关文章"
  }
}
```

### 3. 更新分类

**端点**: `PUT /api/v1/categories/:id`

**描述**: 更新指定分类

**认证**: 需要认证

**路径参数**:
- `id` (integer): 分类 ID

**请求体**:
```json
{
  "name": "更新后的名称",
  "slug": "updated-slug",
  "description": "更新后的描述"
}
```

**说明**: 所有字段都是可选的。

**响应**: `200 OK`

### 4. 删除分类

**端点**: `DELETE /api/v1/categories/:id`

**描述**: 删除指定分类

**认证**: 需要认证

**路径参数**:
- `id` (integer): 分类 ID

**响应**: `204 No Content`

---

## 标签 API

### 1. 获取标签列表

**端点**: `GET /api/v1/tags`

**描述**: 获取所有标签

**认证**: 无需认证

**查询参数**:
| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `with_count` | string | - | 是否包含文章计数，传入任意值（如`1`或`true`）即启用 |

**响应**: `200 OK`

不带计数的响应：
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Go",
      "slug": "go"
    },
    {
      "id": 2,
      "name": "TypeScript",
      "slug": "typescript"
    }
  ]
}
```

带计数的响应（`?with_count=1`）：
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Go",
      "slug": "go",
      "post_count": 12
    },
    {
      "id": 2,
      "name": "TypeScript",
      "slug": "typescript",
      "post_count": 5
    }
  ]
}
```

**说明**: 当提供 `with_count` 参数时，响应中会包含 `post_count` 字段，显示使用该标签的文章数量。

### 2. 创建标签

**端点**: `POST /api/v1/tags`

**描述**: 创建新标签

**认证**: 需要认证

**请求体**:
```json
{
  "name": "Go",
  "slug": "go"
}
```

**字段说明**:
- `name` (string, 必需): 标签名称
- `slug` (string, 必需): URL 友好的唯一标识符

**响应**: `201 Created`
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Go",
    "slug": "go"
  }
}
```

### 3. 更新标签

**端点**: `PUT /api/v1/tags/:id`

**描述**: 更新指定标签

**认证**: 需要认证

**路径参数**:
- `id` (integer): 标签 ID

**请求体**:
```json
{
  "name": "Golang",
  "slug": "golang"
}
```

**说明**: 所有字段都是可选的。

**响应**: `200 OK`

### 4. 删除标签

**端点**: `DELETE /api/v1/tags/:id`

**描述**: 删除指定标签

**认证**: 需要认证

**路径参数**:
- `id` (integer): 标签 ID

**响应**: `204 No Content`

---

## 评论 API

### 1. 获取文章评论列表

**端点**: `GET /api/v1/posts/:id/comments`

**描述**: 获取指定文章的所有评论（包括嵌套回复）

**认证**: 无需认证

**路径参数**:
- `id` (integer): 文章 ID

**响应**: `200 OK`
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "post_id": 1,
      "user_id": 1,
      "name": "",
      "email": "",
      "content": "这是一条评论",
      "ip": "192.168.1.*",
      "status": "approved",
      "parent_id": null,
      "created_at": "2025-01-01T00:00:00Z",
      "user": {
        "id": 1,
        "username": "admin",
        "nickname": "管理员"
      },
      "replies": [
        {
          "id": 2,
          "post_id": 1,
          "user_id": null,
          "name": "游客",
          "email": "guest@example.com",
          "content": "这是回复",
          "ip": "192.168.1.*",
          "status": "approved",
          "parent_id": 1,
          "created_at": "2025-01-02T00:00:00Z"
        }
      ]
    }
  ]
}
```

**嵌套回复说明**:
- 评论分为**顶级评论**（`parent_id` 为 `null`）和**回复评论**（`parent_id` 指向父评论ID）
- 响应中顶级评论的 `replies` 数组包含所有子回复
- 目前仅支持**两层嵌套**（顶级评论 + 回复），不支持多层嵌套
- 游客评论时 `user_id` 为 `null`，使用 `name` 和 `email` 字段
- 登录用户评论时自动关联 `user_id`，`name` 和 `email` 为空
- IP地址返回时自动脱敏（最后一段用 `*` 替代）

### 2. 创建评论

**端点**: `POST /api/v1/posts/:id/comments`

**描述**: 为指定文章创建评论（支持游客和认证用户）

**认证**: 可选

**路径参数**:
- `id` (integer): 文章 ID

**请求体**:

认证用户：
```json
{
  "content": "这是一条评论",
  "parent_id": null
}
```

游客评论：
```json
{
  "content": "这是一条评论",
  "parent_id": null,
  "name": "游客",
  "email": "guest@example.com"
}
```

**字段说明**:
- `content` (string, 必需): 评论内容，1-1000 字符
- `parent_id` (integer, 可选): 父评论 ID（用于回复）
- `name` (string, 游客必需): 评论者姓名，2-50 字符
- `email` (string, 游客必需): 评论者邮箱

**响应**: `201 Created`
```json
{
  "success": true,
  "data": {
    "id": 1,
    "post_id": 1,
    "user_id": null,
    "name": "游客",
    "email": "guest@example.com",
    "content": "这是一条评论",
    "ip": "192.168.1.1",
    "status": "pending",
    "parent_id": null,
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

### 3. 更新评论状态

**端点**: `PUT /api/v1/comments/:id`

**描述**: 更新评论状态（审核）

**认证**: 需要认证

**路径参数**:
- `id` (integer): 评论 ID

**请求体**:
```json
{
  "status": "approved"
}
```

**字段说明**:
- `status` (string, 必需): 评论状态，可选值：`pending`/`approved`/`spam`

**响应**: `200 OK`

### 4. 删除评论

**端点**: `DELETE /api/v1/comments/:id`

**描述**: 删除指定评论

**认证**: 需要认证

**路径参数**:
- `id` (integer): 评论 ID

**响应**: `204 No Content`

---

## 媒体 API

### 1. 上传媒体文件

**端点**: `POST /api/v1/media/upload`

**描述**: 上传媒体文件（支持 multipart 和 base64 两种方式）

**认证**: 需要认证

**文件限制**:
- **最大文件大小**: 10MB（可通过配置 `media.max_file_size_mb` 调整）
- **支持的文件类型**:
  - 图片：`image/jpeg`, `image/png`, `image/gif`, `image/webp`, `image/svg+xml`
  - 视频：`video/mp4`, `video/webm`, `video/ogg`
  - 音频：`audio/mpeg`, `audio/ogg`, `audio/wav`
  - 文档：`application/pdf`, `text/plain`, `text/markdown`
- **文件名**: 上传后自动重命名为 `{hash}.{ext}` 格式，按日期分目录存储

**请求方式一：Multipart Form-Data**

```
Content-Type: multipart/form-data

file: [binary data]
```

**请求方式二：JSON (Base64)**

```json
{
  "file_name": "image.jpg",
  "data": "base64_encoded_data...",
  "mime_type": "image/jpeg"
}
```

**字段说明（Base64方式）**:
- `file_name` (string, 必需): 原始文件名（包含扩展名）
- `data` (string, 必需): Base64编码的文件数据
- `mime_type` (string, 必需): 文件的MIME类型

**响应**: `201 Created`
```json
{
  "success": true,
  "data": {
    "id": 1,
    "file_name": "image.jpg",
    "file_path": "2025/01/abc123.jpg",
    "file_hash": "abc123def456",
    "url": "/media/i/2025/01/abc123.jpg",
    "mime_type": "image/jpeg",
    "size": 102400,
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

### 2. 获取媒体列表

**端点**: `GET /api/v1/media`

**描述**: 获取媒体文件列表

**认证**: 需要认证

**查询参数**:
| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `page` | integer | 1 | 页码 |
| `limit` | integer | 20 | 每页数量（最大 100） |

**响应**: `200 OK`
```json
{
  "success": true,
  "data": {
    "total": 50,
    "page": 1,
    "limit": 20,
    "data": [
      {
        "id": 1,
        "file_name": "image.jpg",
        "file_path": "2025/01/abc123.jpg",
        "file_hash": "abc123def456",
        "mime_type": "image/jpeg",
        "size": 102400,
        "created_at": "2025-01-01T00:00:00Z",
        "updated_at": "2025-01-01T00:00:00Z"
      }
    ]
  }
}
```

### 3. 访问媒体文件

**端点**: `GET /api/v1/media/i/:filename`

**描述**: 访问上传的媒体文件

**认证**: 无需认证

**路径参数**:
- `filename` (string): 文件路径（如 `2025/01/abc123.jpg`）

**响应**: 返回文件二进制数据

### 4. 删除媒体文件

**端点**: `DELETE /api/v1/media/:id`

**描述**: 删除指定媒体文件（包括物理文件和数据库记录）

**认证**: 需要认证

**路径参数**:
- `id` (integer): 媒体 ID

**响应**: `204 No Content`

---

## 用户 API

### 1. 获取用户信息

**端点**: `GET /api/v1/users/:id`

**描述**: 获取指定用户的公开信息

**认证**: 无需认证

**路径参数**:
- `id` (integer): 用户 ID

**响应**: `200 OK`
```json
{
  "success": true,
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "nickname": "管理员",
    "bio": "这是个人简介",
    "avatar_url": "https://example.com/avatar.jpg"
  }
}
```

### 2. 更新用户信息

**端点**: `PUT /api/v1/users/:id`

**描述**: 更新用户信息（只能更新自己的信息）

**认证**: 需要认证

**路径参数**:
- `id` (integer): 用户 ID

**请求体**:
```json
{
  "nickname": "新昵称",
  "username": "newusername",
  "avatar_url": "https://example.com/new-avatar.jpg",
  "bio": "更新后的个人简介",
  "email": "newemail@example.com",
  "old_password": "oldpass123",
  "new_password": "newpass123"
}
```

**字段说明**:
- `nickname` (string, 可选): 昵称
- `username` (string, 可选): 用户名
- `avatar_url` (string, 可选): 头像 URL
- `bio` (string, 可选): 个人简介
- `email` (string, 可选): 邮箱
- `old_password` (string, 修改密码时必需): 旧密码
- `new_password` (string, 修改密码时必需): 新密码

**响应**: `200 OK`
```json
{
  "success": true,
  "data": {
    "id": 1,
    "username": "newusername",
    "email": "newemail@example.com",
    "nickname": "新昵称",
    "bio": "更新后的个人简介",
    "avatar_url": "https://example.com/new-avatar.jpg"
  }
}
```

---

## 设置 API

系统设置功能允许动态配置博客系统的各种参数，支持三种数据类型：字符串(str)、数字(number)和JSON对象(json)。

### 1. 获取所有设置

**端点**: `GET /api/v1/settings`

**描述**: 获取所有系统设置（返回值自动根据类型解析）

**认证**: 需要认证

**响应**: `200 OK`
```json
{
  "success": true,
  "data": {
    "site_title": "我的博客",
    "max_comments_per_hour": 10,
    "theme_config": {
      "color": "blue",
      "dark_mode": true
    },
    "site_description": "这是一个基于KLog的博客系统"
  }
}
```

**说明**:
- 字符串类型(`str`)直接返回字符串值
- 数字类型(`number`)自动转换为数值
- JSON类型(`json`)自动解析为对象/数组

---

### 2. 获取单个设置

**端点**: `GET /api/v1/settings/:key`

**描述**: 获取指定键的设置值

**认证**: 需要认证

**路径参数**:
- `key` (string): 设置键名

**响应**: `200 OK`
```json
{
  "success": true,
  "data": {
    "key": "site_title",
    "value": "我的博客",
    "type": "str"
  }
}
```

**错误响应**: `404 Not Found`
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "设置不存在"
  }
}
```

---

### 3. 创建/更新单个设置

**端点**: `PUT /api/v1/settings`

**描述**: 创建或更新单个设置（如果键已存在则更新，否则创建）

**认证**: 需要认证

**请求体**:
```json
{
  "key": "site_title",
  "value": "我的博客",
  "type": "str"
}
```

**字段说明**:
- `key` (string, 必需): 设置键名，建议使用下划线命名法（如`site_title`）
- `value` (string, 必需): 设置值（所有类型都以字符串形式传递）
- `type` (string, 必需): 值类型，可选值：`str`/`number`/`json`

**类型验证规则**:

1. **`str` (字符串类型)**
   - 直接存储，无需额外验证
   - 适用于：标题、描述、URL等

2. **`number` (数字类型)**
   - 必须是有效的数字字符串（如`"123"`、`"3.14"`）
   - 存储时保持字符串格式，读取时自动转换为float64
   - 适用于：阈值、计数、比例等
   - 验证示例：
     ```json
     {"key": "max_upload_size", "value": "10", "type": "number"}  // ✅ 有效
     {"key": "ratio", "value": "3.14", "type": "number"}          // ✅ 有效
     {"key": "count", "value": "abc", "type": "number"}           // ❌ 无效
     ```

3. **`json` (JSON类型)**
   - 必须是有效的JSON字符串
   - 支持对象和数组
   - 自动验证JSON格式
   - 适用于：配置对象、数组列表等
   - 验证示例：
     ```json
     {"key": "theme", "value": "{\"color\":\"blue\"}", "type": "json"}     // ✅ 有效
     {"key": "tags", "value": "[\"tech\",\"life\"]", "type": "json"}       // ✅ 有效
     {"key": "config", "value": "{invalid json}", "type": "json"}          // ❌ 无效
     ```

**响应**: `200 OK` (更新) 或 `201 Created` (创建)
```json
{
  "success": true,
  "data": {
    "key": "site_title",
    "value": "我的博客",
    "type": "str"
  }
}
```

**错误响应**: `400 Bad Request`
```json
{
  "success": false,
  "error": {
    "code": "INVALID_PARAMS",
    "message": "无效的数字格式"
  }
}
```

---

### 4. 批量创建/更新设置

**端点**: `PUT /api/v1/settings/batch`

**描述**: 批量创建或更新多个设置（事务操作，要么全部成功要么全部失败）

**认证**: 需要认证

**请求体**:
```json
{
  "settings": [
    {
      "key": "site_title",
      "value": "我的博客",
      "type": "str"
    },
    {
      "key": "max_comments_per_hour",
      "value": "10",
      "type": "number"
    },
    {
      "key": "theme_config",
      "value": "{\"color\":\"blue\",\"dark_mode\":true}",
      "type": "json"
    }
  ]
}
```

**字段说明**:
- `settings` (array, 必需): 设置数组，每个元素包含`key`、`value`、`type`字段

**响应**: `200 OK`
```json
{
  "success": true,
  "data": {
    "updated_count": 3,
    "settings": [
      {
        "key": "site_title",
        "value": "我的博客",
        "type": "str"
      },
      {
        "key": "max_comments_per_hour",
        "value": "10",
        "type": "number"
      },
      {
        "key": "theme_config",
        "value": "{\"color\":\"blue\",\"dark_mode\":true}",
        "type": "json"
      }
    ]
  }
}
```

**错误响应**: `400 Bad Request`
```json
{
  "success": false,
  "error": {
    "code": "INVALID_PARAMS",
    "message": "第2个设置的值类型验证失败: 无效的JSON格式"
  }
}
```

**说明**:
- 批量操作使用数据库事务，保证原子性
- 任何一个设置验证失败，整个操作回滚
- 建议批量操作的设置数量不超过100个

---

### 5. 删除设置

**端点**: `DELETE /api/v1/settings/:key`

**描述**: 删除指定的设置

**认证**: 需要认证

**路径参数**:
- `key` (string): 设置键名

**响应**: `204 No Content`

**错误响应**: `404 Not Found`
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "设置不存在"
  }
}
```

---

### 使用场景示例

**场景1: 配置站点信息**
```bash
# 设置站点标题
curl -X PUT http://localhost:8010/api/v1/settings \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "key": "site_title",
    "value": "技术博客",
    "type": "str"
  }'

# 设置站点描述
curl -X PUT http://localhost:8010/api/v1/settings \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "key": "site_description",
    "value": "分享技术与生活",
    "type": "str"
  }'
```

**场景2: 配置限流参数**
```bash
# 批量配置评论限流
curl -X PUT http://localhost:8010/api/v1/settings/batch \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "settings": [
      {"key": "comment_rate_limit_per_min", "value": "1", "type": "number"},
      {"key": "comment_rate_limit_per_hour", "value": "10", "type": "number"}
    ]
  }'
```

**场景3: 配置主题**
```bash
# 使用JSON类型存储主题配置
curl -X PUT http://localhost:8010/api/v1/settings \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "key": "theme_config",
    "value": "{\"primary_color\":\"#1890ff\",\"dark_mode\":true,\"font_size\":16}",
    "type": "json"
  }'
```

---

## 附录

### HTTP 状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 201 | 创建成功 |
| 204 | 删除成功（无内容返回） |
| 400 | 请求参数错误 |
| 401 | 未认证 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

### 限流说明

系统采用双重限流机制保护API：

#### 1. 全局限流
- **应用范围**: 所有 API 端点
- **限流规则**: 每个IP每秒最多10个请求
- **突发处理**: 允许短时间内最多20个请求
- **响应状态码**: 429 Too Many Requests

**触发限流时的响应**:
```json
{
  "success": false,
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "请求过于频繁，请稍后再试"
  }
}
```

#### 2. 评论限流（额外保护）
- **应用范围**: 仅评论创建接口 (`POST /api/v1/posts/:id/comments`)
- **限流规则**:
  - 每个IP每分钟最多1条评论
  - 每个IP每小时最多10条评论
- **目的**: 防止评论垃圾信息和恶意刷屏

**触发评论限流时的响应**:
```json
{
  "success": false,
  "error": {
    "code": "COMMENT_RATE_LIMIT_EXCEEDED",
    "message": "评论过于频繁，请稍后再试"
  }
}
```

#### 3. 其他限制
- **请求体大小**: 最大 30MB
- **超时时间**: API请求超时时间为30秒
- **并发连接**: 建议客户端合理控制并发请求数

### 数据类型说明

- 所有时间字段采用 RFC3339 格式（如 `2025-01-01T00:00:00Z`）
- 所有 ID 字段为正整数
- 可选字段在值为 `null` 时可能不出现在响应中

