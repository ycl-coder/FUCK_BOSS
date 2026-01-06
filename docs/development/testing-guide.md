# 测试指南

本文档介绍如何编写和运行 Fuck Boss 平台的测试。

## 测试策略

### 测试金字塔

```
        /\
       /  \      E2E Tests (5%)
      /----\
     /      \    Integration Tests (25%)
    /--------\
   /          \  Unit Tests (70%)
  /------------\
```

### 测试类型

1. **单元测试 (Unit Tests)**
   - 测试单个函数或方法
   - 使用 Mock 隔离依赖
   - 快速执行（毫秒级）
   - 覆盖率目标: >= 70%，核心逻辑 >= 90%

2. **集成测试 (Integration Tests)**
   - 测试组件与外部依赖的集成
   - 使用真实数据库和 Redis
   - 中等执行速度（秒级）
   - 验证技术实现正确性

3. **E2E 测试 (End-to-End Tests)**
   - 测试完整用户流程
   - 使用完整 Docker Compose 环境
   - 较慢执行（分钟级）
   - 验证系统整体功能

## 运行测试

### 单元测试

```bash
cd backend

# 运行所有单元测试
go test ./test/unit/...

# 运行特定包的测试
go test ./test/unit/domain/content/...

# 运行测试并显示详细输出
go test -v ./test/unit/application/content/...

# 运行测试并查看覆盖率
go test -cover ./test/unit/application/content/...

# 生成覆盖率报告
go test -coverprofile=coverage.out -coverpkg=./internal/application/content,./internal/application/dto ./test/unit/application/content/...
go tool cover -func=coverage.out

# 生成 HTML 覆盖率报告
go tool cover -html=coverage.out -o coverage.html
```

### 集成测试

```bash
# 1. 启动测试环境
make test-up

# 2. 等待服务就绪
docker-compose -f docker-compose.test.yml ps

# 3. 运行集成测试
make test-integration

# 或运行特定测试
make test-integration-repository
make test-integration-cache
make test-integration-usecase

# 4. 停止测试环境
make test-down
```

### E2E 测试

```bash
# 1. 启动测试环境
make test-up

# 2. 运行 E2E 测试
make test-e2e-grpc

# 3. 停止测试环境
make test-down
```

## 编写测试

### 单元测试

#### 使用 testify

```go
package content_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestNewPostID_ValidUUID(t *testing.T) {
    id, err := content.NewPostID("123e4567-e89b-12d3-a456-426614174000")
    
    require.NoError(t, err)
    assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", id.String())
}
```

#### 使用 Mock

```go
package content_test

import (
    "context"
    "testing"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/assert"
)

type MockPostRepository struct {
    mock.Mock
}

func (m *MockPostRepository) Save(ctx context.Context, post *content.Post) error {
    args := m.Called(ctx, post)
    return args.Error(0)
}

func TestCreatePostUseCase_Execute_Success(t *testing.T) {
    // 创建 Mock
    mockRepo := new(MockPostRepository)
    mockCache := new(MockCacheRepository)
    
    // 设置期望
    mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
    mockCache.On("DeleteByPattern", mock.Anything, "posts:city:*").Return(nil)
    
    // 创建 Use Case
    useCase := content.NewCreatePostUseCase(mockRepo, mockCache)
    
    // 执行
    cmd := content.CreatePostCommand{
        Company: "测试公司",
        CityCode: "beijing",
        CityName: "北京",
        Content: "这是一个测试内容，长度超过10个字符",
    }
    result, err := useCase.Execute(context.Background(), cmd)
    
    // 验证
    assert.NoError(t, err)
    assert.NotNil(t, result)
    mockRepo.AssertExpectations(t)
    mockCache.AssertExpectations(t)
}
```

### 集成测试

#### 使用 testify/suite

```go
package repository_test

import (
    "context"
    "database/sql"
    "testing"
    "github.com/stretchr/testify/suite"
    "fuck_boss/backend/internal/domain/content"
    "fuck_boss/backend/internal/infrastructure/persistence/postgres"
)

type PostRepositoryTestSuite struct {
    suite.Suite
    db   *sql.DB
    repo *postgres.PostRepository
    ctx  context.Context
}

func (s *PostRepositoryTestSuite) SetupSuite() {
    // 连接数据库
    s.db, _ = sql.Open("postgres", "postgres://test_user:test_password@localhost:5433/test_db?sslmode=disable")
    s.repo = postgres.NewPostRepository(s.db)
    s.ctx = context.Background()
}

func (s *PostRepositoryTestSuite) SetupTest() {
    // 清理数据
    s.db.Exec("TRUNCATE TABLE posts")
}

func (s *PostRepositoryTestSuite) TearDownTest() {
    // 清理数据
    s.db.Exec("TRUNCATE TABLE posts")
}

func (s *PostRepositoryTestSuite) TearDownSuite() {
    // 关闭连接
    s.db.Close()
}

func (s *PostRepositoryTestSuite) TestSave() {
    // 创建 Post
    post := content.NewPost(
        content.NewCompanyName("测试公司"),
        shared.NewCity("beijing", "北京"),
        content.NewContent("这是一个测试内容，长度超过10个字符"),
    )
    
    // 保存
    err := s.repo.Save(s.ctx, post)
    s.NoError(err)
    
    // 验证
    found, err := s.repo.FindByID(s.ctx, post.ID())
    s.NoError(err)
    s.NotNil(found)
    s.Equal(post.ID(), found.ID())
}

func TestPostRepositoryTestSuite(t *testing.T) {
    suite.Run(t, new(PostRepositoryTestSuite))
}
```

### E2E 测试

```go
package scenarios_test

import (
    "context"
    "testing"
    "github.com/stretchr/testify/suite"
    contentv1 "fuck_boss/backend/api/proto/content/v1"
)

type GRPCE2ETestSuite struct {
    suite.Suite
    client contentv1.ContentServiceClient
    ctx    context.Context
}

func (s *GRPCE2ETestSuite) SetupSuite() {
    // 连接到 gRPC 服务器
    conn, _ := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    s.client = contentv1.NewContentServiceClient(conn)
    s.ctx = context.Background()
}

func (s *GRPCE2ETestSuite) TestCreatePost_E2E() {
    // 创建 Post
    req := &contentv1.CreatePostRequest{
        Company:  "测试公司",
        CityCode: "beijing",
        CityName: "北京",
        Content:  "这是一个测试内容，长度超过10个字符",
    }
    
    resp, err := s.client.CreatePost(s.ctx, req)
    
    s.NoError(err)
    s.NotEmpty(resp.PostId)
}

func TestGRPCE2ETestSuite(t *testing.T) {
    suite.Run(t, new(GRPCE2ETestSuite))
}
```

## 测试最佳实践

### 1. 测试命名

```go
// 格式: Test<FunctionName>_<Scenario>_<ExpectedResult>
func TestNewPostID_ValidUUID_ReturnsPostID(t *testing.T) {}
func TestNewPostID_InvalidUUID_ReturnsError(t *testing.T) {}
func TestCreatePostUseCase_Execute_RateLimitExceeded_ReturnsError(t *testing.T) {}
```

### 2. 测试组织

- **表驱动测试**: 用于测试多个场景

```go
func TestNewCompanyName_Validation(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
        errMsg  string
    }{
        {
            name:    "valid company name",
            input:   "测试公司",
            wantErr: false,
        },
        {
            name:    "empty company name",
            input:   "",
            wantErr: true,
            errMsg:  "company name is required",
        },
        {
            name:    "too long company name",
            input:   strings.Repeat("a", 101),
            wantErr: true,
            errMsg:  "company name too long",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := content.NewCompanyName(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### 3. 测试隔离

- 每个测试应该独立运行
- 使用 `SetupTest` 和 `TearDownTest` 清理数据
- 不依赖测试执行顺序

### 4. 测试数据

- 使用工厂函数创建测试数据
- 使用随机数据避免冲突
- 清理测试数据

```go
func createTestPost(t *testing.T) *content.Post {
    company, _ := content.NewCompanyName("测试公司")
    city, _ := shared.NewCity("beijing", "北京")
    content, _ := content.NewContent("这是一个测试内容，长度超过10个字符")
    return content.NewPost(company, city, content)
}
```

### 5. 错误测试

- 测试所有错误场景
- 验证错误消息
- 测试边界条件

```go
func TestNewContent_TooShort_ReturnsError(t *testing.T) {
    _, err := content.NewContent("short")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "content too short")
}
```

## 测试覆盖率

### 查看覆盖率

```bash
# 生成覆盖率报告
go test -coverprofile=coverage.out -coverpkg=./internal/application/content,./internal/application/dto ./test/unit/application/content/...

# 查看覆盖率摘要
go tool cover -func=coverage.out

# 生成 HTML 报告
go tool cover -html=coverage.out -o coverage.html
```

### 覆盖率目标

- **整体覆盖率**: >= 70%
- **核心业务逻辑**: >= 90%
- **基础设施层**: >= 80%

### 排除不需要测试的代码

```go
// 在代码中添加注释排除
// +build !test

// 或使用 build tag
//go:build !test
```

## 测试工具

### 推荐工具

1. **testify**: 断言和 Mock
   ```bash
   go get github.com/stretchr/testify
   ```

2. **gomock**: 接口 Mock 生成
   ```bash
   go install github.com/golang/mock/mockgen@latest
   ```

3. **testcontainers**: 容器化测试（可选）
   ```bash
   go get github.com/testcontainers/testcontainers-go
   ```

## 常见问题

### Q: 如何测试私有方法？

**A**: 通过公开方法测试，或使用同包测试文件（`_test.go` 在同一包中）。

### Q: 集成测试太慢怎么办？

**A**: 
- 使用 `tmpfs` 加速数据库（测试环境）
- 并行运行测试（`go test -parallel`）
- 只运行必要的集成测试

### Q: 如何 Mock 外部服务？

**A**: 使用接口和 Mock 实现，在测试中注入 Mock。

### Q: 测试数据如何管理？

**A**: 
- 使用工厂函数创建测试数据
- 使用 `SetupTest` 清理数据
- 使用随机数据避免冲突

## 相关文档

- [环境设置指南](./setup-guide.md)
- [开发指南](./development-guide.md)
- [集成测试执行流程](../../backend/test/integration/EXECUTION_FLOW.md)

