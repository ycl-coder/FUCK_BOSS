# redis - Redis 缓存和限流实现

Redis 缓存和限流功能的实现。

## 结构

- **cache_repository.go** - CacheRepository 实现
- **rate_limiter.go** - RateLimiter 实现

## 实现

### CacheRepository

实现 `application/cache.CacheRepository` 接口，提供 Get/Set/Delete 操作。

#### 创建 Repository

```go
import (
    "github.com/redis/go-redis/v9"
    "fuck_boss/backend/internal/infrastructure/persistence/redis"
)

// client 是 *redis.Client 实例
cacheRepo := redis.NewCacheRepository(client)
```

#### 使用示例

```go
// 设置缓存（带 TTL）
err := cacheRepo.Set(ctx, "post:123", "post data", 10*time.Minute)
if err != nil {
    return err
}

// 获取缓存
value, err := cacheRepo.Get(ctx, "post:123")
if err != nil {
    // 处理错误（可能是缓存未命中）
    return err
}

// 删除缓存
err = cacheRepo.Delete(ctx, "post:123")
if err != nil {
    return err
}

// 按模式删除缓存
err = cacheRepo.DeleteByPattern(ctx, "posts:city:beijing:*")
if err != nil {
    return err
}
```

#### 方法说明

- **Get**: 获取缓存值，如果不存在返回 `NotFoundError`
- **Set**: 设置缓存值，支持 TTL（过期时间）
- **Delete**: 删除单个缓存键
- **DeleteByPattern**: 按模式删除多个缓存键（使用 SCAN 命令，避免阻塞 Redis）

#### 错误处理

- 使用统一的错误处理包 (`pkg/errors`)
- Redis 错误包装为 `DATABASE_ERROR`
- 缓存未命中返回 `NOT_FOUND` 错误
- 参数验证错误返回 `VALIDATION_ERROR`

### RateLimiter

实现 `application/ratelimit.RateLimiter` 接口，使用滑动窗口算法实现限流。

#### 创建 RateLimiter

```go
import (
    "github.com/redis/go-redis/v9"
    "fuck_boss/backend/internal/infrastructure/persistence/redis"
)

// client 是 *redis.Client 实例
limiter := redis.NewRateLimiter(client)
```

#### 使用示例

```go
// 检查是否允许请求
// key: 限流键（如 "rate_limit:post:127.0.0.1:2026-01-06-14"）
// limit: 限制数量（如 3）
// window: 时间窗口（如 1 小时）
allowed, err := limiter.Allow(ctx, "rate_limit:post:127.0.0.1:2026-01-06-14", 3, time.Hour)
if err != nil {
    // 处理错误（Redis 连接失败等）
    return err
}

if !allowed {
    // 限流，拒绝请求
    return apperrors.NewRateLimitError("rate limit exceeded")
}

// 允许请求，继续处理
```

#### 方法说明

- **Allow**: 检查是否允许请求，使用滑动窗口算法
  - 使用 Redis INCR 增加计数器
  - 使用 EXPIRE 设置过期时间（窗口大小）
  - 如果计数 <= limit，允许请求
  - 如果计数 > limit，拒绝请求

- **GetRemaining**: 获取剩余可用请求数（用于返回给客户端）

- **Reset**: 重置限流计数器（用于测试或手动重置）

#### 滑动窗口算法

**工作原理**:
1. 第一次请求时，Redis 键不存在，INCR 返回 1，设置 EXPIRE
2. 后续请求在窗口内，INCR 递增计数器
3. 如果计数 <= limit，允许请求
4. 如果计数 > limit，拒绝请求
5. 窗口过期后，键自动删除，重新开始计数

**示例**:
```
时间线（limit=3, window=1小时）:
t=0:   请求1 → INCR=1 → 允许 ✓
t=10m: 请求2 → INCR=2 → 允许 ✓
t=20m: 请求3 → INCR=3 → 允许 ✓
t=30m: 请求4 → INCR=4 → 拒绝 ✗ (超过限制)
t=1h:  键过期，重新开始
t=1h5m: 请求5 → INCR=1 → 允许 ✓
```

#### 错误处理

- 使用统一的错误处理包 (`pkg/errors`)
- Redis 错误包装为 `DATABASE_ERROR`
- 参数验证错误返回 `VALIDATION_ERROR`

## 缓存 Key 规范

- 列表缓存: `posts:city:{cityCode}:page:{page}`
- 详情缓存: `post:{postID}`
- 搜索缓存: `search:{keyword}:city:{cityCode}:page:{page}`
- 限流 Key: `rate_limit:post:{ip}:{hour}`

## TTL 策略

- 列表缓存: 5-10 分钟（根据城市热度）
- 详情缓存: 10 分钟
- 搜索缓存: 5 分钟
- 限流窗口: 1 小时

## 注意事项

- 必须处理 Redis 连接失败（降级到数据库）
- 支持 context 取消
- 必须通过集成测试验证
- 使用滑动窗口算法实现限流

