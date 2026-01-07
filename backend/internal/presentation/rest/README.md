# REST API Handler - REST API 处理器

REST API 处理器，将 JSON HTTP 请求转换为应用层的 Use Cases 调用。

## 功能

- **CreatePost**: 创建新帖子
- **ListPosts**: 获取帖子列表（支持城市筛选和分页）
- **GetPost**: 获取帖子详情
- **SearchPosts**: 搜索帖子（支持关键词和城市筛选）

## 使用示例

```go
import (
    "fuck_boss/backend/internal/application/content"
    "fuck_boss/backend/internal/application/search"
    "fuck_boss/backend/internal/presentation/rest"
    "fuck_boss/backend/internal/infrastructure/logger"
)

// 创建 REST API 处理器
restHandler := rest.NewContentHandler(
    createUseCase,  // content.CreatePostUseCaseInterface
    listUseCase,    // content.ListPostsUseCaseInterface
    getUseCase,     // content.GetPostUseCaseInterface
    searchUseCase,  // search.SearchPostsUseCaseInterface
    logger,         // logger.Logger
)
```

## API 端点

### POST /api/posts
创建新帖子

**请求体**:
```json
{
  "company": "公司名称",
  "cityCode": "beijing",
  "cityName": "北京",
  "content": "曝光内容...",
  "occurredAt": 1767715620  // 可选，Unix 时间戳
}
```

**响应**:
```json
{
  "postId": "uuid",
  "createdAt": 1767715620
}
```

### GET /api/posts
获取帖子列表

**查询参数**:
- `cityCode` (可选): 城市代码，不传则返回所有城市
- `page` (可选): 页码，默认 1
- `pageSize` (可选): 每页数量，默认 20

**响应**:
```json
{
  "posts": [...],
  "total": 12,
  "page": 1,
  "pageSize": 20
}
```

### GET /api/posts/:id
获取帖子详情

**响应**:
```json
{
  "id": "uuid",
  "company": "公司名称",
  "cityCode": "beijing",
  "cityName": "北京",
  "content": "内容...",
  "occurredAt": 1767715620,  // 可选
  "createdAt": 1767715620
}
```

### POST /api/posts/search
搜索帖子

**请求体**:
```json
{
  "keyword": "搜索关键词",
  "cityCode": "beijing",  // 可选
  "page": 1,              // 可选
  "pageSize": 20          // 可选
}
```

**响应**:
```json
{
  "posts": [...],
  "total": 5,
  "page": 1,
  "pageSize": 20
}
```

## 错误处理

所有错误都会转换为标准的 HTTP 状态码：

- `400 Bad Request`: 验证错误（VALIDATION_ERROR）
- `404 Not Found`: 资源未找到（NOT_FOUND）
- `429 Too Many Requests`: 限流错误（RATE_LIMIT_EXCEEDED）
- `500 Internal Server Error`: 内部错误

错误响应格式：
```json
{
  "error": "错误消息"
}
```

## 主要组件

### ContentHandler
REST API 请求处理器，实现了所有 HTTP 端点。

### 请求/响应类型
- `CreatePostRequest` / `CreatePostResponse`
- `ListPostsRequest` / `ListPostsResponse`
- `PostResponse`
- `SearchPostsRequest` / `SearchPostsResponse`

## 注意事项

1. **CORS 支持**: 所有端点都支持 CORS，允许跨域请求
2. **客户端 IP**: CreatePost 会自动从请求头提取客户端 IP（X-Forwarded-For, X-Real-IP）
3. **错误转换**: 应用层错误会自动转换为对应的 HTTP 状态码
4. **JSON 格式**: 所有请求和响应都使用 JSON 格式

## 相关文档

- [设计文档](../../../../.spec-workflow/specs/content-management-v1/design.md)
- [gRPC 处理器](../grpc/README.md)
- [应用层用例](../../application/content/README.md)

