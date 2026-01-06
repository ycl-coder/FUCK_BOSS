# 开发指南

本文档介绍 Fuck Boss 平台的开发规范、架构原则和最佳实践。

## 架构原则

### 领域驱动设计 (DDD)

项目采用 DDD 分层架构，将业务逻辑与基础设施分离：

```
Presentation Layer (表现层)
    ↓ 调用
Application Layer (应用层)
    ↓ 使用
Domain Layer (领域层)
    ↑ 实现
Infrastructure Layer (基础设施层)
```

### 分层职责

#### Domain Layer (领域层)

- **职责**: 核心业务逻辑，不依赖任何外部框架
- **包含**: 
  - Entities (聚合根): `Post`
  - Value Objects: `PostID`, `CompanyName`, `Content`, `City`
  - Repository Interfaces: `PostRepository`
  - Domain Services: 领域服务
- **规则**: 
  - 不依赖其他层
  - 只包含业务逻辑
  - 使用接口定义依赖

#### Application Layer (应用层)

- **职责**: 协调领域对象，处理用例
- **包含**:
  - Use Cases: `CreatePostUseCase`, `ListPostsUseCase`
  - DTOs: `PostDTO`, `PostsListDTO`
  - Commands/Queries: `CreatePostCommand`, `ListPostsQuery`
- **规则**:
  - 只依赖 Domain Layer
  - 不包含业务逻辑（业务逻辑在 Domain）
  - 处理事务边界

#### Infrastructure Layer (基础设施层)

- **职责**: 实现技术细节
- **包含**:
  - Repository Implementations: `PostRepository` (PostgreSQL)
  - Cache: `CacheRepository` (Redis)
  - Config: 配置管理
  - Logger: 日志记录
- **规则**:
  - 实现 Domain Layer 定义的接口
  - 可以依赖外部库
  - 处理技术细节（数据库、缓存等）

#### Presentation Layer (表现层)

- **职责**: 处理外部请求
- **包含**:
  - gRPC Handlers: `ContentService`
  - Middleware: `LoggingInterceptor`, `RecoveryInterceptor`
- **规则**:
  - 调用 Application Layer 的 Use Cases
  - 处理协议转换（gRPC ↔ Domain）
  - 错误处理和状态码转换

## 代码规范

### Go 代码风格

#### 格式化

- **[必须]** 使用 `gofmt` 格式化代码
- **[必须]** 使用 `goimports` 自动管理 import

```bash
# 格式化代码
go fmt ./...

# 格式化并整理 import
goimports -w .
```

#### 命名规范

- **包名**: `lowercase`, 简短有意义 (如 `content`, `postgres`)
- **类型**: `PascalCase` (如 `Post`, `PostRepository`)
- **函数/方法**: `PascalCase` (公开), `camelCase` (私有)
- **常量**: `UPPER_SNAKE_CASE` (如 `MAX_POST_LENGTH`)
- **变量**: `camelCase` (如 `postID`, `cityName`)

#### Import 顺序

```go
import (
    // 1. 标准库
    "context"
    "fmt"
    "time"
    
    // 2. 第三方库
    "github.com/lib/pq"
    "go.uber.org/zap"
    
    // 3. 项目内部包
    "fuck_boss/backend/internal/domain/content"
    "fuck_boss/backend/pkg/errors"
)
```

#### 错误处理

- **[必须]** 所有返回 error 的函数必须处理错误
- **[必须]** error 必须是最后一个返回值
- **[推荐]** 使用 `fmt.Errorf` 和 `%w` 进行错误包装 (Go 1.13+)
- **[必须]** 错误描述不需要标点结尾

```go
// ✅ 正确
func CreatePost(ctx context.Context, repo Repository) (*Post, error) {
    post, err := repo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to find post: %w", err)
    }
    return post, nil
}

// ❌ 错误
func CreatePost(ctx context.Context, repo Repository) (error, *Post) {
    // error 应该在最后
}
```

#### Panic 和 Recover

- **[必须]** 业务逻辑中禁止使用 panic
- **[必须]** 在 goroutine 顶层捕获 panic
- **[推荐]** 使用 recover 记录详细堆栈信息

```go
// ✅ 正确 - 在中间件中捕获 panic
func RecoveryInterceptor(log logger.Logger) grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
        defer func() {
            if r := recover(); r != nil {
                log.Error("panic recovered", zap.Any("panic", r), zap.String("stack", string(debug.Stack())))
                err = status.Error(codes.Internal, "internal server error")
            }
        }()
        return handler(ctx, req)
    }
}
```

### 代码组织

#### 文件大小

- **文件**: 建议 < 800 行
- **函数**: 建议 < 80 行
- **嵌套深度**: 建议 < 4 层

#### 单一职责

每个文件、函数只做一件事：

```go
// ✅ 正确 - 职责单一
// domain/content/entity.go - 只包含实体定义
// domain/content/repository.go - 只包含仓库接口
// domain/content/value_object.go - 只包含值对象

// ❌ 错误 - 职责混乱
// domain/content/all.go - 包含所有内容
```

#### 依赖倒置

依赖接口而非实现：

```go
// ✅ 正确 - 依赖接口
type CreatePostUseCase struct {
    repo content.PostRepository  // 接口
    cache cache.CacheRepository  // 接口
}

// ❌ 错误 - 依赖实现
type CreatePostUseCase struct {
    repo *postgres.PostRepository  // 具体实现
}
```

## 开发流程

### 1. 创建新功能

#### 步骤

1. **定义领域模型** (Domain Layer)
   - 创建 Entity 或 Value Object
   - 定义 Repository Interface

2. **实现基础设施** (Infrastructure Layer)
   - 实现 Repository
   - 编写集成测试

3. **实现用例** (Application Layer)
   - 创建 Use Case
   - 编写单元测试

4. **实现 Handler** (Presentation Layer)
   - 创建 gRPC Handler
   - 编写单元测试和 E2E 测试

#### 示例：添加新功能

```go
// 1. Domain Layer - 定义领域模型
// domain/content/entity.go
type Post struct {
    id        PostID
    company   CompanyName
    // ...
}

// domain/content/repository.go
type PostRepository interface {
    Save(ctx context.Context, post *Post) error
    FindByID(ctx context.Context, id PostID) (*Post, error)
}

// 2. Infrastructure Layer - 实现仓库
// infrastructure/persistence/postgres/post_repository.go
type PostRepository struct {
    db *sql.DB
}

func (r *PostRepository) Save(ctx context.Context, post *domain.Post) error {
    // 实现数据库操作
}

// 3. Application Layer - 实现用例
// application/content/create_post.go
type CreatePostUseCase struct {
    repo content.PostRepository
}

func (uc *CreatePostUseCase) Execute(ctx context.Context, cmd CreatePostCommand) (*dto.PostDTO, error) {
    // 实现用例逻辑
}

// 4. Presentation Layer - 实现 Handler
// presentation/grpc/content_handler.go
func (s *ContentService) CreatePost(ctx context.Context, req *contentv1.CreatePostRequest) (*contentv1.CreatePostResponse, error) {
    // 调用 Use Case
}
```

### 2. 编写测试

#### 测试金字塔

```
        /\
       /  \      E2E Tests (少量)
      /----\
     /      \    Integration Tests (中等)
    /--------\
   /          \  Unit Tests (大量)
  /------------\
```

#### 单元测试

- **位置**: `test/unit/`
- **特点**: 快速、隔离、使用 Mock
- **覆盖率**: >= 70%，核心逻辑 >= 90%

```go
// test/unit/application/content/create_post_test.go
func TestCreatePostUseCase_Execute_Success(t *testing.T) {
    // 创建 Mock
    mockRepo := new(MockPostRepository)
    mockCache := new(MockCacheRepository)
    
    // 设置期望
    mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
    
    // 执行测试
    useCase := NewCreatePostUseCase(mockRepo, mockCache)
    result, err := useCase.Execute(ctx, cmd)
    
    // 验证
    assert.NoError(t, err)
    assert.NotNil(t, result)
    mockRepo.AssertExpectations(t)
}
```

#### 集成测试

- **位置**: `test/integration/`
- **特点**: 使用真实数据库和 Redis
- **环境**: Docker Compose 测试环境

```go
// test/integration/repository/post_repository_test.go
func (s *PostRepositoryTestSuite) TestSave() {
    // 使用真实数据库
    post := content.NewPost(...)
    err := s.repo.Save(s.ctx, post)
    
    assert.NoError(s.T(), err)
}
```

#### E2E 测试

- **位置**: `test/e2e/`
- **特点**: 测试完整流程
- **环境**: 完整的 Docker Compose 环境

### 3. 代码审查清单

提交代码前检查：

- [ ] 代码通过 `gofmt` 和 `goimports`
- [ ] 代码通过 `go vet`
- [ ] 所有测试通过
- [ ] 测试覆盖率 >= 70%
- [ ] 代码注释完整
- [ ] README 更新（如需要）
- [ ] 遵循 DDD 架构原则
- [ ] 错误处理正确
- [ ] 没有 panic（业务逻辑中）

## 最佳实践

### 值对象 (Value Objects)

- **不可变**: 创建后不能修改
- **验证**: 在创建时验证
- **相等性**: 基于值而非引用

```go
// ✅ 正确
type PostID struct {
    value string
}

func NewPostID(value string) (PostID, error) {
    if _, err := uuid.Parse(value); err != nil {
        return PostID{}, fmt.Errorf("invalid post ID: %w", err)
    }
    return PostID{value: value}, nil
}

func (id PostID) String() string {
    return id.value
}
```

### 实体 (Entities)

- **有唯一标识**: 通过 ID 区分
- **业务方法**: 包含业务逻辑
- **封装**: 字段私有，通过方法访问

```go
// ✅ 正确
type Post struct {
    id        PostID
    company   CompanyName
    content   Content
    createdAt time.Time
}

func (p *Post) Publish() error {
    // 业务逻辑
    return nil
}

func (p *Post) ID() PostID {
    return p.id
}
```

### Repository 模式

- **接口在 Domain**: 定义在 Domain Layer
- **实现在 Infrastructure**: 实现在 Infrastructure Layer
- **依赖倒置**: Application Layer 依赖接口

```go
// Domain Layer - 定义接口
type PostRepository interface {
    Save(ctx context.Context, post *Post) error
    FindByID(ctx context.Context, id PostID) (*Post, error)
}

// Infrastructure Layer - 实现接口
type PostRepository struct {
    db *sql.DB
}

func (r *PostRepository) Save(ctx context.Context, post *domain.Post) error {
    // 实现
}
```

### Use Case 模式

- **单一职责**: 每个 Use Case 处理一个用例
- **协调作用**: 协调领域对象和基础设施
- **返回 DTO**: 不返回领域对象

```go
// ✅ 正确
type CreatePostUseCase struct {
    repo  content.PostRepository
    cache cache.CacheRepository
}

func (uc *CreatePostUseCase) Execute(ctx context.Context, cmd CreatePostCommand) (*dto.PostDTO, error) {
    // 1. 验证输入
    // 2. 创建领域对象
    // 3. 保存到仓库
    // 4. 更新缓存
    // 5. 返回 DTO
}
```

## 常见问题

### Q: 应该在哪个层处理验证？

**A**: 
- **输入验证**: Application Layer (Use Case)
- **业务规则验证**: Domain Layer (Value Objects, Entities)

### Q: 如何处理事务？

**A**: 在 Application Layer 的 Use Case 中处理事务边界。

### Q: 如何测试私有方法？

**A**: 通过公开方法测试，或使用 `_test.go` 文件（同包测试）。

### Q: 何时使用 Value Object vs Entity？

**A**: 
- **Value Object**: 没有唯一标识，通过值区分（如 `PostID`, `CompanyName`）
- **Entity**: 有唯一标识，通过 ID 区分（如 `Post`）

## 相关文档

- [环境设置指南](./setup-guide.md)
- [测试指南](./testing-guide.md)
- [API 文档](../api/README.md)
- [架构设计](../design/architecture.md)

