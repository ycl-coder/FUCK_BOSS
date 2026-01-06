# content - 内容用例

内容相关的应用用例（Use Cases）。

## 结构

- **create_post.go** - CreatePostUseCase（创建曝光内容）
- **list_posts.go** - ListPostsUseCase（列表查询）
- **get_post.go** - GetPostUseCase（详情查询）
- **dto.go** - 数据传输对象（DTO）

## Use Cases

### CreatePostUseCase

处理创建曝光内容的用例逻辑。

#### 创建 UseCase

```go
import (
    "fuck_boss/backend/internal/application/content"
    "fuck_boss/backend/internal/domain/content"
    "fuck_boss/backend/internal/application/cache"
    "fuck_boss/backend/internal/application/ratelimit"
)

uc := content.NewCreatePostUseCase(
    postRepo,      // content.PostRepository
    cacheRepo,     // cache.CacheRepository
    rateLimiter,   // ratelimit.RateLimiter
)
```

#### 使用示例

```go
dto, err := uc.Execute(ctx, content.CreatePostCommand{
    Company:   "测试公司",
    CityCode:  "beijing",
    CityName:  "北京",
    Content:   "这是一条测试内容，用于验证创建功能。内容应该足够长以满足最小长度要求。",
    ClientIP:  "127.0.0.1",
    OccurredAt: nil, // 可选
})
if err != nil {
    // 处理错误
    return err
}

// 使用返回的 DTO
fmt.Printf("Created post ID: %s\n", dto.ID)
```

#### 执行流程

1. **验证输入**: 检查必填字段（Company, CityCode, CityName, Content, ClientIP）
2. **检查限流**: 使用 RateLimiter 检查是否超过限制（3次/小时/IP）
3. **创建值对象**: 使用工厂方法创建 CompanyName, City, Content
4. **创建实体**: 使用 NewPost 创建 Post 聚合根
5. **保存到数据库**: 调用 Repository.Save 保存
6. **清除缓存**: 清除该城市相关的列表缓存
7. **返回 DTO**: 将 Post 实体转换为 PostDTO 返回

#### 错误处理

- **验证错误**: 返回 `VALIDATION_ERROR`
- **限流错误**: 返回 `RATE_LIMIT_EXCEEDED`
- **数据库错误**: 返回 `DATABASE_ERROR`
- **内部错误**: 返回 `INTERNAL_ERROR`

#### 限流策略

- **限制**: 每个 IP 每小时最多 3 条
- **Key 格式**: `rate_limit:post:{ip}:{hour}`
- **窗口**: 1 小时（滑动窗口）

### ListPostsUseCase

处理列表查询用例逻辑，包含缓存策略。

#### 创建 UseCase

```go
import (
    "fuck_boss/backend/internal/application/content"
    "fuck_boss/backend/internal/domain/content"
    "fuck_boss/backend/internal/application/cache"
)

uc := content.NewListPostsUseCase(
    postRepo,   // content.PostRepository
    cacheRepo,  // cache.CacheRepository
)
```

#### 使用示例

```go
result, err := uc.Execute(ctx, content.ListPostsQuery{
    CityCode: "beijing",
    Page:     1,
    PageSize: 20,
})
if err != nil {
    // 处理错误
    return err
}

// 使用返回的 DTO
fmt.Printf("Total: %d, Page: %d, Posts: %d\n", 
    result.Total, result.Page, len(result.Posts))
```

#### 执行流程

1. **验证输入**: 检查必填字段（CityCode），设置默认值（Page=1, PageSize=20）
2. **检查缓存**: 使用 Key `posts:city:{cityCode}:page:{page}` 查询缓存
3. **缓存命中**: 如果缓存存在，反序列化并返回
4. **缓存未命中**: 查询 Repository
5. **更新缓存**: 将查询结果序列化并存入缓存（TTL: 5-10 分钟）
6. **返回 DTO**: 将 Post 实体列表转换为 PostsListDTO 返回

#### 缓存策略

- **Key 格式**: `posts:city:{cityCode}:page:{page}`
- **TTL 策略**:
  - 热门城市（beijing, shanghai, guangzhou, shenzhen）: 5 分钟
  - 其他城市: 10 分钟
- **错误处理**: 缓存错误不影响主流程，自动回退到数据库查询

#### 错误处理

- **验证错误**: 返回 `VALIDATION_ERROR`
- **数据库错误**: 返回 `DATABASE_ERROR`
- **缓存错误**: 忽略，回退到数据库查询

### GetPostUseCase

处理详情查询用例逻辑，包含缓存策略。

#### 创建 UseCase

```go
import (
    "fuck_boss/backend/internal/application/content"
    "fuck_boss/backend/internal/domain/content"
    "fuck_boss/backend/internal/application/cache"
)

uc := content.NewGetPostUseCase(
    postRepo,   // content.PostRepository
    cacheRepo,  // cache.CacheRepository
)
```

#### 使用示例

```go
dto, err := uc.Execute(ctx, "post-uuid-here")
if err != nil {
    // 处理错误
    if apperrors.IsNotFoundError(err) {
        // Post 不存在
        return err
    }
    return err
}

// 使用返回的 DTO
fmt.Printf("Post ID: %s, Company: %s\n", dto.ID, dto.Company)
```

#### 执行流程

1. **验证输入**: 检查 Post ID 是否为空，验证 UUID 格式
2. **检查缓存**: 使用 Key `post:{postID}` 查询缓存
3. **缓存命中**: 如果缓存存在，反序列化并返回
4. **缓存未命中**: 查询 Repository
5. **处理 NotFound**: 如果 Post 不存在，返回 NotFound 错误
6. **更新缓存**: 将查询结果序列化并存入缓存（TTL: 10 分钟）
7. **返回 DTO**: 将 Post 实体转换为 PostDTO 返回

#### 缓存策略

- **Key 格式**: `post:{postID}`
- **TTL**: 10 分钟
- **错误处理**: 缓存错误不影响主流程，自动回退到数据库查询

#### 错误处理

- **验证错误**: 返回 `VALIDATION_ERROR`（空 ID 或无效 UUID）
- **NotFound 错误**: 返回 `NOT_FOUND`（Post 不存在）
- **数据库错误**: 返回 `DATABASE_ERROR`
- **缓存错误**: 忽略，回退到数据库查询

## DTOs

- **PostDTO** - Post 的数据传输对象
- **PostsListDTO** - Post 列表的数据传输对象

## 注意事项

- Use Case 只包含用例逻辑，不包含业务规则
- 业务规则在 Domain Layer
- 使用 DTO 进行数据传输，不直接暴露 Domain Entity

