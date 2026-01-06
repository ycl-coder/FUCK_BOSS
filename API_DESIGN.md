# API 设计文档

## 📋 概述

本文档定义了"吐槽老板"平台的所有API接口规范。

**Base URL**: `https://api.fuckboss.com/v1`  
**认证方式**: JWT Bearer Token  
**数据格式**: JSON  
**字符编码**: UTF-8

---

## 🔐 认证相关

### 1. 用户注册
```
POST /auth/register
```

**请求体**:
```json
{
  "email": "user@example.com",
  "phone": "13800138000",  // 二选一
  "password": "password123",
  "username": "username",  // 可选
  "verification_code": "123456",  // 邮箱/手机验证码
  "agree_terms": true
}
```

**响应**:
```json
{
  "code": 200,
  "message": "注册成功",
  "data": {
    "user_id": 123,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 3600
  }
}
```

### 2. 用户登录
```
POST /auth/login
```

**请求体**:
```json
{
  "account": "user@example.com",  // 邮箱/手机/用户名
  "password": "password123",
  "remember_me": false
}
```

**响应**:
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "user_id": 123,
    "username": "username",
    "avatar": "https://...",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 3600
  }
}
```

### 3. 发送验证码
```
POST /auth/send-code
```

**请求体**:
```json
{
  "type": "email",  // email | phone
  "email": "user@example.com",
  "phone": "13800138000"
}
```

### 4. 刷新Token
```
POST /auth/refresh
```

**请求头**:
```
Authorization: Bearer {refresh_token}
```

---

## 👤 用户相关

### 1. 获取用户信息
```
GET /users/:id
```

**响应**:
```json
{
  "code": 200,
  "data": {
    "id": 123,
    "username": "username",
    "avatar": "https://...",
    "reputation_score": 150,
    "is_anonymous": false,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 2. 更新用户信息
```
PUT /users/me
Authorization: Bearer {token}
```

**请求体**:
```json
{
  "username": "newname",
  "avatar": "https://...",
  "is_anonymous": false
}
```

### 3. 修改密码
```
POST /users/me/password
Authorization: Bearer {token}
```

**请求体**:
```json
{
  "old_password": "oldpass",
  "new_password": "newpass"
}
```

---

## 📝 曝光相关

### 1. 发布曝光
```
POST /exposures
Authorization: Bearer {token}
```

**请求体**:
```json
{
  "company_name": "某某公司",
  "cities": ["beijing", "shanghai"],  // 城市代码数组
  "exposure_type": "拖欠工资",  // 曝光类型
  "title": "曝光标题",
  "content": "详细描述内容...",
  "boss_name": "张老板",  // 可选
  "department": "技术部",  // 可选
  "occurred_at": "2024-01-01",  // 可选
  "tags": ["加班", "996"],  // 可选，最多5个
  "evidence_urls": [  // 可选
    "https://oss.example.com/image1.jpg",
    "https://oss.example.com/video1.mp4"
  ]
}
```

**响应**:
```json
{
  "code": 200,
  "message": "发布成功",
  "data": {
    "exposure_id": 456,
    "status": "pending"  // pending | published
  }
}
```

### 2. 获取曝光列表
```
GET /exposures
```

**查询参数**:
- `page`: 页码 (默认: 1)
- `page_size`: 每页数量 (默认: 20, 最大: 100)
- `sort`: 排序方式 (latest | hot | credibility_high | credibility_low | controversial)
- `city`: 城市代码
- `exposure_type`: 曝光类型
- `time_range`: 时间范围 (today | week | month | all)
- `credibility_min`: 最小可信度 (0-100)
- `keyword`: 搜索关键词

**响应**:
```json
{
  "code": 200,
  "data": {
    "list": [
      {
        "id": 456,
        "company_name": "某某公司",
        "cities": ["北京", "上海"],
        "exposure_type": "拖欠工资",
        "title": "曝光标题",
        "content_preview": "内容预览...",
        "credibility_score": 75.5,
        "confirm_count": 10,
        "deny_count": 2,
        "view_count": 150,
        "comment_count": 5,
        "author": {
          "id": 123,
          "username": "username",
          "is_anonymous": false,
          "reputation_score": 150
        },
        "created_at": "2024-01-01T00:00:00Z",
        "tags": ["加班", "996"]
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 100,
      "total_pages": 5
    }
  }
}
```

### 3. 获取曝光详情
```
GET /exposures/:id
```

**响应**:
```json
{
  "code": 200,
  "data": {
    "id": 456,
    "company": {
      "id": 789,
      "name": "某某公司",
      "industry": "互联网",
      "exposure_count": 5,
      "avg_credibility": 70.0
    },
    "cities": ["北京", "上海"],
    "exposure_type": "拖欠工资",
    "title": "曝光标题",
    "content": "完整内容...",
    "boss_name": "张老板",
    "department": "技术部",
    "occurred_at": "2024-01-01",
    "tags": ["加班", "996"],
    "evidence_urls": [
      "https://oss.example.com/image1.jpg"
    ],
    "credibility_score": 75.5,
    "credibility_level": "待验证",  // 已验证 | 待验证 | 争议 | 不可信
    "verify_count": 12,
    "confirm_count": 10,
    "deny_count": 2,
    "view_count": 150,
    "comment_count": 5,
    "author": {
      "id": 123,
      "username": "username",
      "is_anonymous": false,
      "reputation_score": 150
    },
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "user_verification": null,  // 当前用户的验证状态
    "is_favorited": false  // 是否已收藏
  }
}
```

### 4. 更新曝光
```
PUT /exposures/:id
Authorization: Bearer {token}
```

### 5. 删除曝光
```
DELETE /exposures/:id
Authorization: Bearer {token}
```

---

## ✅ 验证相关

### 1. 证实曝光
```
POST /exposures/:id/verify/confirm
Authorization: Bearer {token}
```

**请求体**:
```json
{
  "reason": "我也遇到过类似情况，确实如此...",
  "evidence_urls": [  // 可选
    "https://oss.example.com/evidence.jpg"
  ]
}
```

**响应**:
```json
{
  "code": 200,
  "message": "证实成功",
  "data": {
    "verification_id": 789,
    "exposure_id": 456,
    "verification_type": "confirm",
    "credibility_score": 78.5  // 更新后的可信度
  }
}
```

### 2. 证伪曝光
```
POST /exposures/:id/verify/deny
Authorization: Bearer {token}
```

**请求体**: 同证实接口

### 3. 修改验证
```
PUT /verifications/:id
Authorization: Bearer {token}
```

**请求体**:
```json
{
  "reason": "修改后的理由...",
  "evidence_urls": []
}
```

**注意**: 只能在24小时内修改

### 4. 获取验证记录
```
GET /exposures/:id/verifications
```

**查询参数**:
- `type`: 验证类型 (confirm | deny | all)
- `sort`: 排序 (latest | popular | reputation)
- `page`: 页码
- `page_size`: 每页数量

**响应**:
```json
{
  "code": 200,
  "data": {
    "list": [
      {
        "id": 789,
        "verification_type": "confirm",
        "reason": "验证理由...",
        "evidence_urls": [],
        "like_count": 5,
        "dislike_count": 1,
        "user": {
          "id": 123,
          "username": "username",
          "reputation_score": 150
        },
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 12
    }
  }
}
```

---

## 💬 评论相关

### 1. 发表评论
```
POST /exposures/:id/comments
Authorization: Bearer {token}
```

**请求体**:
```json
{
  "content": "评论内容...",
  "parent_id": null  // 回复评论时填写父评论ID
}
```

### 2. 获取评论列表
```
GET /exposures/:id/comments
```

**查询参数**:
- `sort`: 排序 (latest | hot | earliest)
- `page`: 页码
- `page_size`: 每页数量

**响应**:
```json
{
  "code": 200,
  "data": {
    "list": [
      {
        "id": 101,
        "content": "评论内容...",
        "like_count": 10,
        "dislike_count": 2,
        "reply_count": 3,
        "user": {
          "id": 123,
          "username": "username",
          "avatar": "https://..."
        },
        "replies": [  // 子评论（最多3级）
          {
            "id": 102,
            "content": "回复内容...",
            "parent_id": 101,
            "user": {...}
          }
        ],
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {...}
  }
}
```

### 3. 点赞/点踩评论
```
POST /comments/:id/like
POST /comments/:id/dislike
Authorization: Bearer {token}
```

### 4. 删除评论
```
DELETE /comments/:id
Authorization: Bearer {token}
```

---

## 🔍 搜索相关

### 1. 基础搜索
```
GET /search
```

**查询参数**:
- `q`: 搜索关键词
- `page`: 页码
- `page_size`: 每页数量

### 2. 高级搜索
```
POST /search/advanced
```

**请求体**:
```json
{
  "keyword": "搜索关键词",
  "cities": ["beijing", "shanghai"],
  "exposure_types": ["拖欠工资", "职场霸凌"],
  "time_range": {
    "start": "2024-01-01",
    "end": "2024-12-31"
  },
  "credibility_range": {
    "min": 50,
    "max": 100
  },
  "has_evidence": true,
  "sort": "latest",
  "page": 1,
  "page_size": 20
}
```

### 3. 搜索建议
```
GET /search/suggestions
```

**查询参数**:
- `q`: 搜索关键词

**响应**:
```json
{
  "code": 200,
  "data": {
    "companies": ["某某公司", "另一公司"],
    "tags": ["加班", "996"],
    "cities": ["北京", "上海"]
  }
}
```

---

## 🏢 公司相关

### 1. 获取公司信息
```
GET /companies/:id
```

**响应**:
```json
{
  "code": 200,
  "data": {
    "id": 789,
    "name": "某某公司",
    "industry": "互联网",
    "website": "https://...",
    "description": "公司描述...",
    "exposure_count": 5,
    "avg_credibility": 70.0,
    "exposure_types": [
      {"type": "拖欠工资", "count": 2},
      {"type": "职场霸凌", "count": 3}
    ],
    "cities": ["北京", "上海"],
    "recent_exposures": [
      {
        "id": 456,
        "title": "曝光标题",
        "credibility_score": 75.5,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

### 2. 公司排行榜
```
GET /companies/rankings
```

**查询参数**:
- `type`: 排行类型 (exposure_count | low_credibility | controversial)
- `limit`: 返回数量 (默认: 20)

---

## 📍 城市相关

### 1. 获取城市曝光列表
```
GET /cities/:code/exposures
```

**查询参数**: 同曝光列表接口

### 2. 城市统计数据
```
GET /cities/:code/stats
```

**响应**:
```json
{
  "code": 200,
  "data": {
    "city_code": "beijing",
    "city_name": "北京",
    "exposure_count": 100,
    "avg_credibility": 65.5,
    "top_companies": [
      {"name": "某某公司", "count": 10}
    ],
    "exposure_types": [
      {"type": "拖欠工资", "count": 30}
    ]
  }
}
```

### 3. 城市热力图数据
```
GET /cities/heatmap
```

**响应**:
```json
{
  "code": 200,
  "data": [
    {
      "city_code": "beijing",
      "city_name": "北京",
      "exposure_count": 100,
      "avg_credibility": 65.5
    }
  ]
}
```

---

## ⭐ 收藏相关

### 1. 收藏曝光
```
POST /exposures/:id/favorite
Authorization: Bearer {token}
```

### 2. 取消收藏
```
DELETE /exposures/:id/favorite
Authorization: Bearer {token}
```

### 3. 获取收藏列表
```
GET /users/me/favorites
Authorization: Bearer {token}
```

---

## 📢 通知相关

### 1. 获取通知列表
```
GET /notifications
Authorization: Bearer {token}
```

**查询参数**:
- `type`: 通知类型
- `read`: 是否已读 (true | false | all)
- `page`: 页码
- `page_size`: 每页数量

### 2. 标记通知为已读
```
PUT /notifications/:id/read
Authorization: Bearer {token}
```

### 3. 标记全部为已读
```
PUT /notifications/read-all
Authorization: Bearer {token}
```

---

## 🚨 举报相关

### 1. 提交举报
```
POST /reports
Authorization: Bearer {token}
```

**请求体**:
```json
{
  "target_type": "exposure",  // exposure | comment | verification
  "target_id": 456,
  "report_type": "虚假信息",  // 举报类型
  "reason": "举报理由..."
}
```

---

## 👤 个人中心

### 1. 我的曝光
```
GET /users/me/exposures
Authorization: Bearer {token}
```

### 2. 我的验证
```
GET /users/me/verifications
Authorization: Bearer {token}
```

### 3. 我的评论
```
GET /users/me/comments
Authorization: Bearer {token}
```

---

## 📊 统计相关

### 1. 个人数据统计
```
GET /users/me/statistics
Authorization: Bearer {token}
```

**响应**:
```json
{
  "code": 200,
  "data": {
    "exposure_count": 5,
    "verification_count": 20,
    "comment_count": 50,
    "favorite_count": 10,
    "reputation_score": 150
  }
}
```

---

## 🔧 文件上传

### 1. 上传文件
```
POST /upload
Authorization: Bearer {token}
Content-Type: multipart/form-data
```

**请求体**: FormData
- `file`: 文件
- `type`: 文件类型 (image | video | document)

**响应**:
```json
{
  "code": 200,
  "data": {
    "url": "https://oss.example.com/file.jpg",
    "size": 1024000,
    "type": "image"
  }
}
```

---

## ⚠️ 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权（需要登录） |
| 403 | 禁止访问（权限不足） |
| 404 | 资源不存在 |
| 409 | 资源冲突（如重复验证） |
| 422 | 验证失败（如验证码错误） |
| 429 | 请求过于频繁（限流） |
| 500 | 服务器内部错误 |

**错误响应格式**:
```json
{
  "code": 400,
  "message": "错误描述",
  "errors": [
    {
      "field": "email",
      "message": "邮箱格式不正确"
    }
  ]
}
```

---

**文档版本**: v1.0  
**最后更新**: 2024年

