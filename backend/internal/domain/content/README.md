# content - 内容领域

内容领域的核心业务逻辑，包含 Post 聚合根、值对象和 Repository 接口。

## 结构

- **entity.go** - Post 聚合根（Aggregate Root）
- **value_object.go** - 值对象（PostID, CompanyName, Content）
- **repository.go** - PostRepository 接口定义

## 核心概念

### Post（聚合根）

表示曝光内容的聚合根，包含业务规则和不变性约束。

#### 创建 Post

```go
import (
    "fuck_boss/backend/internal/domain/content"
    "fuck_boss/backend/internal/domain/shared"
)

// 创建值对象
company, _ := content.NewCompanyName("Example Company")
city, _ := shared.NewCity("beijing", "北京")
postContent, _ := content.NewContent("This is a detailed description of the company's misconduct...")

// 创建 Post
post, err := content.NewPost(company, city, postContent)
if err != nil {
    return err
}

// 发布 Post
err = post.Publish()
if err != nil {
    return err
}
```

#### 访问 Post 属性

```go
// 获取 ID
id := post.ID()

// 获取公司名称
company := post.Company()

// 获取城市
city := post.City()

// 获取内容
content := post.Content()

// 获取创建时间
createdAt := post.CreatedAt()
```

#### 业务规则

- **ID**: 自动生成 UUID，无需手动指定
- **创建时间**: 自动设置为当前时间
- **值对象验证**: 所有值对象（CompanyName、City、Content）在创建时已通过验证
- **不可变性**: Post 创建后，其值对象不可变（通过值对象的不变性保证）

#### 方法

- `NewPost(company, city, content)` - 创建新的 Post（工厂方法，自动生成 ID 和 createdAt）
- `NewPostFromDB(id, company, city, content, createdAt)` - 从数据库重建 Post（用于 Repository 层）
- `Publish()` - 发布内容（业务方法）
- `ID()` - 获取 Post ID
- `Company()` - 获取公司名称
- `City()` - 获取城市
- `Content()` - 获取内容
- `CreatedAt()` - 获取创建时间

### 值对象

#### PostID

UUID 格式的唯一标识符，用于标识 Post。

```go
// 从字符串创建
id, err := content.NewPostID("123e4567-e89b-12d3-a456-426614174000")
if err != nil {
    return err
}

// 从 UUID 创建
id := content.NewPostIDFromUUID(uuid.New())

// 生成新的 PostID
id := content.GeneratePostID()

// 使用
postID := id.String()
```

**验证规则**:
- 必须是有效的 UUID 格式
- 自动去除前后空白字符

**方法**:
- `String()` - 返回字符串表示
- `Value()` - 返回原始值
- `IsZero()` - 检查是否为零值
- `Equals(other PostID)` - 比较两个 PostID

#### CompanyName

公司名称值对象，用于封装公司名称的业务规则。

```go
// 创建公司名称
companyName, err := content.NewCompanyName("Example Company Ltd.")
if err != nil {
    return err
}

// 使用
name := companyName.String()
```

**验证规则**:
- 不能为空（trim 后）
- 长度必须在 1-100 字符之间（使用 rune 计数，支持 Unicode）
- 自动去除前后空白字符
- 内部空白字符保留

**方法**:
- `String()` - 返回字符串表示
- `Value()` - 返回原始值
- `IsZero()` - 检查是否为零值
- `Equals(other CompanyName)` - 比较两个 CompanyName

**常量**:
- `MinCompanyNameLength = 1` - 最小长度
- `MaxCompanyNameLength = 100` - 最大长度
#### Content

内容值对象，用于封装内容的业务规则。

```go
// 创建内容
content, err := content.NewContent("This is a detailed description of the company's misconduct...")
if err != nil {
    return err
}

// 使用
text := content.String()
summary := content.Summary() // 返回前 200 字符，如果更长则加 "..."
```

**验证规则**:
- 不能为空（trim 后）
- 长度必须在 10-5000 字符之间（使用 rune 计数，支持 Unicode）
- 自动去除前后空白字符
- 内部空白字符保留

**方法**:
- `String()` - 返回完整内容
- `Value()` - 返回原始值
- `Summary()` - 返回摘要（前 200 字符，如果更长则加 "..."）
- `IsZero()` - 检查是否为零值
- `Equals(other Content)` - 比较两个 Content

**常量**:
- `MinContentLength = 10` - 最小长度
- `MaxContentLength = 5000` - 最大长度
- `SummaryLength = 200` - 摘要长度
- **City**: 城市（code + name）

### Repository 接口

定义 Post 的持久化接口，遵循依赖倒置原则。

#### PostRepository

```go
import (
    "context"
    "fuck_boss/backend/internal/domain/content"
    "fuck_boss/backend/internal/domain/shared"
)

type PostRepository interface {
    // Save 保存 Post（如果已存在则更新）
    Save(ctx context.Context, post *content.Post) error
    
    // FindByID 根据 ID 查找 Post
    FindByID(ctx context.Context, id content.PostID) (*content.Post, error)
    
    // FindByCity 根据城市查找 Post 列表（分页）
    // page: 页码（从 1 开始）
    // pageSize: 每页数量
    // 返回: Posts 列表、总数、错误
    FindByCity(ctx context.Context, city shared.City, page, pageSize int) ([]*content.Post, int, error)
    
    // Search 搜索 Post（全文搜索，可选城市筛选，分页）
    // keyword: 搜索关键词
    // city: 城市筛选（nil 表示所有城市）
    // page: 页码（从 1 开始）
    // pageSize: 每页数量
    // 返回: Posts 列表、总数、错误
    Search(ctx context.Context, keyword string, city *shared.City, page, pageSize int) ([]*content.Post, int, error)
}
```

#### 设计原则

- **依赖倒置**: 接口定义在 Domain Layer，实现在 Infrastructure Layer
- **Context 支持**: 所有方法都接受 `context.Context` 作为第一个参数
- **错误处理**: 所有方法都返回 `error` 作为最后一个返回值
- **分页支持**: `FindByCity` 和 `Search` 方法支持分页（page 从 1 开始）

#### 使用示例

```go
// 在 Infrastructure Layer 实现
type PostgresPostRepository struct {
    db *sql.DB
}

func (r *PostgresPostRepository) Save(ctx context.Context, post *content.Post) error {
    // 实现保存逻辑
}

func (r *PostgresPostRepository) FindByID(ctx context.Context, id content.PostID) (*content.Post, error) {
    // 实现查找逻辑
}

// 在 Application Layer 使用
type CreatePostUseCase struct {
    repo content.PostRepository
}

func (uc *CreatePostUseCase) Execute(ctx context.Context, cmd CreatePostCommand) error {
    // 创建 Post
    post, err := content.NewPost(company, city, content)
    if err != nil {
        return err
    }
    
    // 保存到 Repository
    return uc.repo.Save(ctx, post)
}
```

## 业务规则

- 公司名称：1-100 字符，不能为空
- 内容：10-5000 字符，不能为空
- 城市：必须是有效城市
- Post ID：自动生成 UUID
- 创建时间：自动设置为当前时间

## 注意事项

- Domain Layer 不依赖任何其他层
- 所有业务规则都在 Domain Layer 中
- Repository 接口定义在这里，实现在 Infrastructure Layer

