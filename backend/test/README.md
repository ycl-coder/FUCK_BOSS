# test - 测试目录

测试文件组织目录，包含单元测试、集成测试和 E2E 测试。

## 目录结构

```
test/
├── unit/              # 单元测试
│   └── domain/        # 领域层单元测试
├── integration/       # 集成测试
└── e2e/              # 端到端测试
```

## 测试覆盖率查看

### 问题说明

由于测试文件在 `test/` 目录下，使用外部测试包（如 `content_test`），直接运行 `go test -cover` 会显示 `[no statements]`，因为 Go 默认只统计被测试包内的代码覆盖率。

### 解决方案

使用 `-coverpkg` 参数指定要统计覆盖率的包：

```bash
# 查看特定包的覆盖率
go test -coverprofile=coverage.out \
  -coverpkg=./internal/domain/content \
  ./test/unit/domain/content/...

# 查看覆盖率报告（文本格式）
go tool cover -func=coverage.out

# 查看覆盖率报告（HTML 格式）
go tool cover -html=coverage.out -o coverage.html
```

### 示例输出

```
ok  	fuck_boss/backend/test/unit/domain/content	0.006s	coverage: 100.0% of statements in ./internal/domain/content
fuck_boss/backend/internal/domain/content/value_object.go:22:	NewPostID		100.0%
fuck_boss/backend/internal/domain/content/value_object.go:35:	NewPostIDFromUUID	100.0%
fuck_boss/backend/internal/domain/content/value_object.go:40:	GeneratePostID		100.0%
...
total:								(statements)		100.0%
```

### 常用命令

#### 1. 运行测试并查看覆盖率

```bash
# 单个包的覆盖率
go test -coverprofile=coverage.out \
  -coverpkg=./internal/domain/content \
  ./test/unit/domain/content/... \
  && go tool cover -func=coverage.out
```

#### 2. 生成 HTML 覆盖率报告

```bash
go test -coverprofile=coverage.out \
  -coverpkg=./internal/domain/content \
  ./test/unit/domain/content/... \
  && go tool cover -html=coverage.out -o coverage.html
```

然后在浏览器中打开 `coverage.html` 查看详细的覆盖率报告。

#### 3. 查看多个包的覆盖率

```bash
# 查看所有领域层的覆盖率
go test -coverprofile=coverage.out \
  -coverpkg=./internal/domain/... \
  ./test/unit/domain/... \
  && go tool cover -func=coverage.out
```

#### 4. 在 Makefile 中添加覆盖率命令

可以在 `Makefile` 中添加：

```makefile
.PHONY: test-coverage
test-coverage:
	go test -coverprofile=coverage.out \
		-coverpkg=./internal/domain/content \
		./test/unit/domain/content/... \
		&& go tool cover -func=coverage.out

.PHONY: test-coverage-html
test-coverage-html:
	go test -coverprofile=coverage.out \
		-coverpkg=./internal/domain/content \
		./test/unit/domain/content/... \
		&& go tool cover -html=coverage.out -o coverage.html
```

## 测试覆盖率要求

根据项目规范（tech.md）：

- **整体覆盖率要求**: >= 70%
- **核心业务逻辑覆盖率**: >= 90%
- **Domain Layer**: 所有实体、值对象、领域服务必须 >= 90%

## 当前覆盖率状态

### PostID 值对象

- **覆盖率**: 100.0%
- **测试文件**: `test/unit/domain/content/value_object_test.go`
- **覆盖的方法**:
  - `NewPostID`: 100.0%
  - `NewPostIDFromUUID`: 100.0%
  - `GeneratePostID`: 100.0%
  - `String`: 100.0%
  - `Value`: 100.0%
  - `IsZero`: 100.0%
  - `Equals`: 100.0%

## gRPC Handler 测试

### 单元测试

gRPC Handler 的单元测试使用 mock UseCases 来隔离 handler 逻辑。

#### 运行单元测试

```bash
# 使用 Makefile
make test-unit-grpc

# 或直接使用 go test
cd backend && go test -v ./test/unit/presentation/grpc/...
```

#### 查看覆盖率

```bash
# 使用 Makefile
make test-unit-grpc-coverage

# 或直接使用 go test
cd backend && go test -coverprofile=coverage-grpc.out \
  -coverpkg=./internal/presentation/grpc \
  ./test/unit/presentation/grpc/... \
  && go tool cover -func=coverage-grpc.out
```

#### 生成 HTML 覆盖率报告

```bash
# 使用 Makefile
make test-unit-grpc-coverage-html

# 或直接使用 go test
cd backend && go test -coverprofile=coverage-grpc.out \
  -coverpkg=./internal/presentation/grpc \
  ./test/unit/presentation/grpc/... \
  && go tool cover -html=coverage-grpc.out -o coverage-grpc.html
```

### E2E 测试

gRPC E2E 测试使用真实的 gRPC 服务器、PostgreSQL 和 Redis，测试完整的服务流程。

#### 运行 E2E 测试

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

#### E2E 测试覆盖

- `CreatePost_E2E`: 端到端创建帖子
- `CreatePost_E2E_ValidationError`: 验证错误处理
- `ListPosts_E2E`: 端到端列出帖子
- `GetPost_E2E`: 端到端获取帖子
- `GetPost_E2E_NotFound`: NotFound 错误处理
- `SearchPosts_E2E`: 端到端搜索帖子
- `TestGRPCMiddleware_E2E`: 中间件验证（日志和恢复）

## 注意事项

1. **外部测试包**: 测试文件使用 `package content_test` 而不是 `package content`，这样可以测试包的公开 API，但不能访问私有成员。

2. **覆盖率统计**: 使用 `-coverpkg` 参数时，需要指定被测试包的完整路径。

3. **HTML 报告**: 生成的 HTML 报告会显示哪些行被覆盖（绿色），哪些行未被覆盖（红色）。

4. **CI/CD 集成**: 在 CI/CD 流程中，可以使用覆盖率报告来确保代码质量。

5. **E2E 测试环境**: E2E 测试需要真实的数据库和 Redis，确保在运行前启动测试环境。
