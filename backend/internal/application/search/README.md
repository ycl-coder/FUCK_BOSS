# search - 搜索用例

搜索相关的应用用例（Use Cases）。

## 结构

- **search_posts.go** - SearchPostsUseCase（搜索曝光内容）

## Use Cases

### SearchPostsUseCase

处理搜索用例逻辑，包含缓存策略。

#### 创建 UseCase

```go
import (
    "fuck_boss/backend/internal/application/search"
    "fuck_boss/backend/internal/domain/content"
    "fuck_boss/backend/internal/application/cache"
)

uc := search.NewSearchPostsUseCase(
    postRepo,   // content.PostRepository
    cacheRepo,  // cache.CacheRepository
)
```

#### 使用示例

```go
// 搜索所有城市
result, err := uc.Execute(ctx, search.SearchPostsQuery{
    Keyword:  "测试公司",
    CityCode: nil, // 搜索所有城市
    Page:     1,
    PageSize: 20,
})

// 搜索特定城市
cityCode := "beijing"
result, err := uc.Execute(ctx, search.SearchPostsQuery{
    Keyword:  "测试公司",
    CityCode: &cityCode,
    Page:     1,
    PageSize: 20,
})
```

#### 执行流程

1. **验证输入**: 检查关键词是否为空，验证最小长度（2 个字符）
2. **设置默认值**: Page=1, PageSize=20
3. **检查缓存**: 使用 Key `search:{keyword}:city:{cityCode}:page:{page}` 或 `search:{keyword}:page:{page}` 查询缓存
4. **缓存命中**: 如果缓存存在，反序列化并返回
5. **缓存未命中**: 查询 Repository（使用全文搜索）
6. **更新缓存**: 将查询结果序列化并存入缓存（TTL: 5 分钟）
7. **返回 DTO**: 将 Post 实体列表转换为 PostsListDTO 返回

#### 缓存策略

- **Key 格式**: 
  - 有城市过滤: `search:{keyword}:city:{cityCode}:page:{page}`
  - 无城市过滤: `search:{keyword}:page:{page}`
- **TTL**: 5 分钟
- **关键词规范化**: 转换为小写并去除首尾空格
- **错误处理**: 缓存错误不影响主流程，自动回退到数据库查询

#### 搜索策略

- 使用 PostgreSQL 全文搜索（tsvector/tsquery）
- 搜索范围：公司名称（company_name）和内容（content）
- 支持可选的城市过滤
- 支持分页

#### 错误处理

- **验证错误**: 返回 `VALIDATION_ERROR`（空关键词或长度不足）
- **数据库错误**: 返回 `DATABASE_ERROR`
- **缓存错误**: 忽略，回退到数据库查询

## 注意事项

- Use Case 只包含用例逻辑，不包含业务规则
- 业务规则在 Domain Layer
- 使用 DTO 进行数据传输，不直接暴露 Domain Entity

## 结构

- **search_posts.go** - SearchPostsUseCase（搜索曝光内容）

## SearchPostsUseCase

处理搜索用例逻辑，支持关键词搜索和城市筛选。

```go
uc := application.NewSearchPostsUseCase(repo, cacheRepo)
result, err := uc.Execute(ctx, application.SearchPostsQuery{
    Keyword:  "测试",
    CityCode: &cityCode, // 可选
    Page:     1,
    PageSize: 20,
})
```

## 搜索策略

- 使用 PostgreSQL 全文搜索
- 搜索范围：公司名称、内容
- 支持城市筛选
- 结果按相关性排序

