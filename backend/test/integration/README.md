# integration - 集成测试

集成测试目录，用于测试组件与外部依赖（数据库、缓存等）的集成。

## 目录结构

```
integration/
├── repository/      # Repository 集成测试
├── cache/           # Cache 集成测试
└── usecase/         # Use Case 集成测试
```

## 运行集成测试

### 前置要求

1. **Docker 和 Docker Compose**: 用于启动测试数据库和缓存
2. **PostgreSQL**: 测试数据库（通过 Docker Compose 启动）

### 启动测试环境

```bash
# 启动测试数据库
docker-compose -f docker-compose.test.yml up -d

# 等待数据库就绪（健康检查）
docker-compose -f docker-compose.test.yml ps
```

### 运行测试

```bash
# 运行所有集成测试
cd backend && go test ./test/integration/...

# 运行特定测试
cd backend && go test ./test/integration/repository/...

# 运行测试并显示详细输出
cd backend && go test -v ./test/integration/repository/...
```

### 环境变量

可以通过环境变量配置测试数据库和 Redis 连接：

```bash
# 默认数据库连接字符串
TEST_DATABASE_URL="postgres://test_user:test_password@localhost:5433/test_db?sslmode=disable"

# 默认 Redis 地址
TEST_REDIS_ADDR="localhost:6380"

# 运行测试时指定
TEST_DATABASE_URL="postgres://..." TEST_REDIS_ADDR="localhost:6380" go test ./test/integration/...
```

### 停止测试环境

```bash
# 停止并删除容器
docker-compose -f docker-compose.test.yml down

# 停止并删除容器和卷（清理所有数据）
docker-compose -f docker-compose.test.yml down -v
```

## 测试原则

1. **使用真实数据库**: 所有集成测试使用真实的 PostgreSQL 数据库
2. **数据隔离**: 每个测试用例前后都会清理数据（TRUNCATE）
3. **独立性**: 测试用例之间相互独立，不依赖执行顺序
4. **错误处理**: 测试错误情况和边界条件

## Repository 集成测试

### PostRepository 测试

测试文件: `repository/post_repository_test.go`

**测试覆盖**:
- ✅ Save: 创建和更新 Post
- ✅ FindByID: 查找单个 Post（存在和不存在）
- ✅ FindByCity: 按城市查找，分页，排序
- ✅ Search: 全文搜索，城市过滤，分页
- ✅ 错误处理: 数据库错误、上下文取消等

**运行示例**:

```bash
# 运行 PostRepository 集成测试
cd backend && go test -v ./test/integration/repository/...
```

## Cache 集成测试

### Redis Cache 和 Rate Limiter 测试

测试文件: `cache/redis_cache_test.go`

**测试覆盖**:

#### CacheRepository
- ✅ Get: 获取缓存值（存在和不存在）
- ✅ Set: 设置缓存值
- ✅ Set TTL: 验证 TTL 过期
- ✅ Delete: 删除单个缓存键
- ✅ DeleteByPattern: 按模式删除多个缓存键
- ✅ 参数验证: 空键、负 TTL 等

#### RateLimiter
- ✅ Allow: 在限制内允许请求
- ✅ Allow: 超过限制拒绝请求
- ✅ Allow: 窗口过期后重新允许
- ✅ GetRemaining: 获取剩余可用请求数
- ✅ Reset: 重置限流计数器
- ✅ 参数验证: 空键、零/负 limit、零/负 window

**运行示例**:

```bash
# 运行 Cache 集成测试
cd backend && go test -v ./test/integration/cache/...

# 或使用 Makefile
make test-integration-cache
```

## 查看测试数据

### 方法 1: 使用调试测试（推荐）

运行调试测试，数据不会被清理：

```bash
# 运行调试测试（不会清理数据）
cd backend && go test -tags=debug -v ./test/integration/repository/... -run TestPostRepository_Save_Debug
```

然后连接到数据库查看：

```bash
# 连接到测试数据库
psql -h localhost -p 5433 -U test_user -d test_db

# 查看所有数据
SELECT * FROM posts;

# 查看特定数据
SELECT * FROM posts WHERE id = 'your-post-id';
```

### 方法 2: 临时禁用数据清理

在 `post_repository_test.go` 中，临时注释掉 `TearDownTest` 方法中的 `TRUNCATE` 语句：

```go
func (s *PostRepositoryTestSuite) TearDownTest() {
    // 临时注释掉清理代码
    // _, err := s.db.ExecContext(s.ctx, "TRUNCATE TABLE posts CASCADE")
}
```

### 方法 3: 在测试中添加延迟

在测试方法末尾添加延迟，给时间查看数据：

```go
func (s *PostRepositoryTestSuite) TestPostRepository_Save() {
    // ... 测试代码 ...
    
    // 添加延迟，有时间查看数据
    time.Sleep(30 * time.Second)
}
```

### 方法 4: 使用持久化存储

修改 `docker-compose.test.yml`，将 `tmpfs` 改为持久化卷：

```yaml
services:
  postgres-test:
    # ... 其他配置 ...
    volumes:
      - postgres_test_data:/var/lib/postgresql/data
    # 移除 tmpfs 配置

volumes:
  postgres_test_data:
```

**注意**: 使用持久化存储后，数据会保留，需要手动清理。

## 注意事项

1. **测试数据库端口**: 使用 5433 端口避免与生产数据库冲突
2. **数据清理**: 每个测试前后都会清理数据，确保测试独立性（调试测试除外）
3. **并发测试**: 测试可以并发运行，但需要确保数据库连接池足够
4. **迁移脚本**: 测试会自动运行迁移脚本创建表结构
5. **tmpfs 存储**: 默认使用内存存储，容器停止后数据消失

## 故障排查

### 数据库连接失败

```bash
# 检查容器是否运行
docker-compose -f docker-compose.test.yml ps

# 查看容器日志
docker-compose -f docker-compose.test.yml logs postgres-test

# 检查端口是否被占用
lsof -i :5433
```

### 迁移失败

确保迁移 SQL 语法正确，可以在 PostgreSQL 客户端中手动执行验证。

### 测试超时

如果测试超时，检查：
1. 数据库是否正常启动
2. 网络连接是否正常
3. 数据库连接池配置

