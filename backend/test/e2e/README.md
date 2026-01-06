# e2e - 端到端测试

端到端（End-to-End）测试，使用真实的 gRPC 服务器、PostgreSQL 和 Redis，测试完整的服务流程。

## 目录结构

```
test/e2e/
└── scenarios/         # 测试场景
    └── grpc_test.go   # gRPC E2E 测试
```

## gRPC E2E 测试

### 测试说明

gRPC E2E 测试使用真实的 gRPC 服务器、PostgreSQL 和 Redis，测试以下内容：

- 完整的服务流程（创建、查询、搜索）
- 错误处理（验证错误、NotFound 错误等）
- 中间件功能（日志和恢复）

### 运行测试

**重要**: E2E 测试需要先启动测试环境（PostgreSQL 和 Redis）。

```bash
# 1. 启动测试环境
make test-up

# 2. 运行 E2E 测试
make test-e2e-grpc

# 3. 停止测试环境（可选）
make test-down
```

或者直接使用 go test：

```bash
# 确保测试环境已启动
docker-compose -f docker-compose.test.yml up -d

# 运行 E2E 测试
cd backend && go test -v ./test/e2e/scenarios/...

# 停止测试环境
docker-compose -f docker-compose.test.yml down
```

### 测试覆盖

#### CreatePost

- ✅ `TestCreatePost_E2E`: 端到端创建帖子
- ✅ `TestCreatePost_E2E_ValidationError`: 验证错误处理

#### ListPosts

- ✅ `TestListPosts_E2E`: 端到端列出帖子

#### GetPost

- ✅ `TestGetPost_E2E`: 端到端获取帖子
- ✅ `TestGetPost_E2E_NotFound`: NotFound 错误处理

#### SearchPosts

- ✅ `TestSearchPosts_E2E`: 端到端搜索帖子

#### Middleware

- ✅ `TestGRPCMiddleware_E2E`: 中间件验证（日志和恢复）

### 测试环境

E2E 测试使用 Docker Compose 启动测试环境：

- **PostgreSQL**: `localhost:5433`
- **Redis**: `localhost:6380`
- **gRPC Server**: `localhost:50053`

测试环境配置在 `docker-compose.test.yml` 中。

### 测试流程

1. **SetupSuite**: 
   - 加载配置
   - 初始化日志
   - 连接数据库和 Redis
   - 运行数据库迁移
   - 创建 gRPC 服务器
   - 启动 gRPC 服务器
   - 创建 gRPC 客户端

2. **SetupTest**: 
   - 清理数据库（TRUNCATE）
   - 清理 Redis（删除所有键）

3. **TestXXX**: 
   - 执行测试用例

4. **TearDownTest**: 
   - 清理数据（为下一个测试做准备）

5. **TearDownSuite**: 
   - 关闭 gRPC 客户端连接
   - 停止 gRPC 服务器
   - 关闭数据库和 Redis 连接

### 注意事项

1. **测试环境**: 确保测试环境已启动，否则测试会失败。

2. **端口冲突**: 如果端口被占用，可以修改 `docker-compose.test.yml` 中的端口配置。

3. **数据清理**: 每个测试前都会清理数据，确保测试隔离。

4. **测试时间**: E2E 测试需要启动真实的服务器，因此运行时间较长。

5. **资源清理**: 测试完成后记得停止测试环境，释放资源。

## 相关文档

- [测试主文档](../README.md)
- [集成测试文档](../integration/README.md)
- [gRPC Handler 实现](../../internal/presentation/grpc/README.md)

