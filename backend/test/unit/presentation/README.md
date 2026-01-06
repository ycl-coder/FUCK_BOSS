# presentation - 表现层单元测试

表现层（Presentation Layer）的单元测试，包括 gRPC Handler 测试。

## 目录结构

```
test/unit/presentation/
└── grpc/              # gRPC Handler 单元测试
    └── content_handler_test.go
```

## gRPC Handler 测试

### 测试说明

gRPC Handler 的单元测试使用 mock UseCases 来隔离 handler 逻辑，测试以下内容：

- 所有 gRPC 方法的成功场景
- 错误处理（验证错误、速率限制错误、NotFound 错误等）
- 错误转换（将应用错误转换为 gRPC 状态码）
- 客户端 IP 提取功能

### 运行测试

```bash
# 使用 Makefile（推荐）
make test-unit-grpc

# 或直接使用 go test
cd backend && go test -v ./test/unit/presentation/grpc/...
```

### 查看覆盖率

```bash
# 使用 Makefile（推荐）
make test-unit-grpc-coverage

# 或直接使用 go test
cd backend && go test -coverprofile=coverage-grpc.out \
  -coverpkg=./internal/presentation/grpc \
  ./test/unit/presentation/grpc/... \
  && go tool cover -func=coverage-grpc.out
```

### 生成 HTML 覆盖率报告

```bash
# 使用 Makefile（推荐）
make test-unit-grpc-coverage-html

# 或直接使用 go test
cd backend && go test -coverprofile=coverage-grpc.out \
  -coverpkg=./internal/presentation/grpc \
  ./test/unit/presentation/grpc/... \
  && go tool cover -html=coverage-grpc.out -o coverage-grpc.html
```

## 测试覆盖

### CreatePost

- ✅ 成功创建帖子
- ✅ 带发生时间的创建
- ✅ 验证错误处理
- ✅ 速率限制错误处理

### ListPosts

- ✅ 成功列出帖子

### GetPost

- ✅ 成功获取帖子
- ✅ NotFound 错误处理

### SearchPosts

- ✅ 成功搜索帖子
- ✅ 不带城市过滤的搜索

### 其他

- ✅ 错误转换为 gRPC 状态码
- ✅ 客户端 IP 提取（从 peer 和 metadata）

## Mock 类

测试使用以下 mock 类：

- `MockCreatePostUseCase`: CreatePostUseCase 的 mock 实现
- `MockListPostsUseCase`: ListPostsUseCase 的 mock 实现
- `MockGetPostUseCase`: GetPostUseCase 的 mock 实现
- `MockSearchPostsUseCase`: SearchPostsUseCase 的 mock 实现

## 相关文档

- [测试主文档](../../README.md)
- [gRPC Handler 实现](../../../internal/presentation/grpc/README.md)

