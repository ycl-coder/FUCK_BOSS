# Project Structure

## Directory Organization

```
fuck_boss/
├── docs/                           # 文档目录（统一管理）
│   ├── design/                     # 设计文档
│   │   ├── architecture.md         # 架构设计
│   │   ├── database-schema.md      # 数据库设计
│   │   └── api-design.md           # API 设计
│   ├── development/                # 开发文档
│   │   ├── setup-guide.md          # 环境搭建指南
│   │   ├── development-guide.md    # 开发指南
│   │   └── testing-guide.md        # 测试指南
│   ├── deployment/                 # 部署文档
│   │   ├── docker-deploy.md        # Docker 部署
│   │   └── production-deploy.md    # 生产环境部署
│   ├── fixes/                      # 修复文档（问题修复记录）
│   │   ├── bug-fixes/              # Bug 修复记录
│   │   ├── performance-fixes/      # 性能优化记录
│   │   └── security-fixes/         # 安全修复记录
│   └── changelog/                  # 变更日志
│       └── CHANGELOG.md            # 版本变更记录
│
├── backend/                        # 后端代码（Go + gRPC）
│   ├── cmd/                        # 应用程序入口
│   │   └── server/                 # gRPC 服务器入口
│   │       └── main.go
│   ├── internal/                   # 内部代码（不对外暴露）
│   │   ├── domain/                 # 领域层（DDD）
│   │   │   ├── content/            # 内容领域
│   │   │   │   ├── entity.go       # Post 实体
│   │   │   │   ├── value_object.go # 值对象
│   │   │   │   ├── repository.go   # Repository 接口
│   │   │   │   └── service.go      # 领域服务
│   │   │   ├── search/             # 搜索领域
│   │   │   │   ├── entity.go
│   │   │   │   ├── repository.go
│   │   │   │   └── service.go
│   │   │   └── shared/             # 共享领域概念
│   │   │       ├── city.go         # 城市值对象
│   │   │       └── company.go      # 公司值对象
│   │   │
│   │   ├── application/            # 应用层（Use Cases）
│   │   │   ├── content/            # 内容用例
│   │   │   │   ├── create_post.go  # 创建曝光用例
│   │   │   │   ├── get_post.go     # 获取曝光用例
│   │   │   │   └── list_posts.go   # 列表查询用例
│   │   │   ├── search/             # 搜索用例
│   │   │   │   └── search_posts.go
│   │   │   └── dto/                # 数据传输对象
│   │   │       ├── content_dto.go
│   │   │       └── search_dto.go
│   │   │
│   │   ├── infrastructure/         # 基础设施层
│   │   │   ├── persistence/        # 持久化实现
│   │   │   │   ├── postgres/       # PostgreSQL 实现
│   │   │   │   │   ├── post_repository.go
│   │   │   │   │   └── migrations/ # 数据库迁移
│   │   │   │   └── redis/          # Redis 实现
│   │   │   │       └── cache_repository.go
│   │   │   ├── config/             # 配置管理
│   │   │   │   └── config.go
│   │   │   └── logger/             # 日志
│   │   │       └── logger.go
│   │   │
│   │   └── presentation/           # 表现层
│   │       ├── grpc/               # gRPC 处理器
│   │       │   ├── content_handler.go
│   │       │   └── search_handler.go
│   │       └── middleware/         # 中间件
│   │           ├── logging.go
│   │           └── recovery.go
│   │
│   ├── pkg/                        # 可复用的公共包
│   │   ├── errors/                 # 错误定义
│   │   ├── validator/              # 验证工具
│   │   └── utils/                  # 工具函数
│   │
│   ├── api/                        # API 定义（protobuf）
│   │   └── proto/                  # .proto 文件
│   │       ├── content/
│   │       │   └── content.proto
│   │       └── search/
│   │           └── search.proto
│   │
│   ├── scripts/                    # 脚本文件
│   │   ├── migrate.sh              # 数据库迁移脚本
│   │   └── generate.sh             # 代码生成脚本
│   │
│   ├── test/                       # 测试文件
│   │   ├── unit/                   # 单元测试
│   │   ├── integration/            # 集成测试
│   │   │   ├── testdata/           # 测试数据
│   │   │   └── docker-compose.test.yml
│   │   └── e2e/                    # 端到端测试
│   │
│   ├── go.mod                      # Go 模块定义
│   ├── go.sum                      # 依赖校验
│   ├── Makefile                    # 构建脚本
│   └── Dockerfile                  # Docker 镜像
│
├── frontend/                       # 前端代码（React）
│   ├── public/                     # 静态资源
│   ├── src/
│   │   ├── api/                    # API 客户端
│   │   │   ├── grpc/               # gRPC 客户端
│   │   │   │   ├── contentClient.ts
│   │   │   │   └── searchClient.ts
│   │   │   └── types/              # API 类型定义
│   │   │
│   │   ├── domain/                 # 领域模型（前端）
│   │   │   ├── models/             # 数据模型
│   │   │   │   ├── Post.ts
│   │   │   │   └── City.ts
│   │   │   └── services/           # 领域服务
│   │   │       └── postService.ts
│   │   │
│   │   ├── features/               # 功能模块（按功能组织）
│   │   │   ├── post/               # 曝光内容功能
│   │   │   │   ├── components/    # 组件
│   │   │   │   │   ├── PostList.tsx
│   │   │   │   │   ├── PostDetail.tsx
│   │   │   │   │   └── PostForm.tsx
│   │   │   │   ├── hooks/         # React Hooks
│   │   │   │   │   └── usePosts.ts
│   │   │   │   └── index.ts       # 导出
│   │   │   │
│   │   │   ├── search/             # 搜索功能
│   │   │   │   ├── components/
│   │   │   │   │   └── SearchBar.tsx
│   │   │   │   └── hooks/
│   │   │   │       └── useSearch.ts
│   │   │   │
│   │   │   └── city/               # 城市选择功能
│   │   │       └── components/
│   │   │           └── CitySelector.tsx
│   │   │
│   │   ├── shared/                # 共享代码
│   │   │   ├── components/        # 通用组件
│   │   │   │   ├── Button.tsx
│   │   │   │   └── Card.tsx
│   │   │   ├── hooks/             # 通用 Hooks
│   │   │   ├── utils/             # 工具函数
│   │   │   └── constants/         # 常量
│   │   │
│   │   ├── app/                   # 应用入口
│   │   │   ├── App.tsx
│   │   │   ├── routes.tsx         # 路由配置
│   │   │   └── store.ts           # 状态管理
│   │   │
│   │   └── main.tsx               # 入口文件
│   │
│   ├── test/                       # 测试文件
│   │   ├── unit/                  # 单元测试
│   │   ├── integration/           # 集成测试
│   │   └── e2e/                   # 端到端测试
│   │
│   ├── package.json
│   ├── tsconfig.json
│   ├── vite.config.ts
│   └── Dockerfile
│
├── .spec-workflow/                 # 规范工作流文档
│   ├── steering/                   # 指导文档
│   │   ├── product.md
│   │   ├── tech.md
│   │   └── structure.md
│   └── specs/                     # 功能规范
│
├── docker-compose.yml              # 开发环境 Docker Compose
├── docker-compose.test.yml        # 测试环境 Docker Compose
├── .gitignore
├── README.md                       # 项目说明（指向 docs/）
└── Makefile                        # 项目级构建脚本
```

## Naming Conventions

### Files

#### 后端 (Go)
- **包名**: `snake_case` (小写，单词间用下划线)
- **文件名**: `snake_case.go` (与包名一致)
- **测试文件**: `snake_case_test.go`
- **接口文件**: `repository.go`, `service.go`
- **实现文件**: `postgres_repository.go`, `redis_cache.go`

#### 前端 (TypeScript/React)
- **组件文件**: `PascalCase.tsx` (如 `PostList.tsx`)
- **工具文件**: `camelCase.ts` (如 `formatDate.ts`)
- **类型文件**: `types.ts` 或 `*.types.ts`
- **测试文件**: `*.test.tsx` 或 `*.spec.tsx`
- **配置文件**: `kebab-case.config.ts` (如 `vite.config.ts`)

### Code

#### 后端 (Go)
- **包名**: `snake_case` (全小写)
- **类型/结构体**: `PascalCase` (如 `Post`, `PostRepository`)
- **接口**: `PascalCase` + `er` 后缀 (如 `PostRepository`, `CacheRepository`)
- **函数/方法**: `camelCase` (如 `CreatePost`, `GetPostByID`)
- **常量**: `UPPER_SNAKE_CASE` (如 `MAX_POST_LENGTH`)
- **变量**: `camelCase` (如 `postID`, `cityName`)
- **私有成员**: 小写开头 (如 `postRepository`)
- **公开成员**: 大写开头 (如 `PostRepository`)

#### 前端 (TypeScript)
- **组件**: `PascalCase` (如 `PostList`)
- **函数/变量**: `camelCase` (如 `getPosts`, `postList`)
- **类型/接口**: `PascalCase` (如 `Post`, `PostDTO`)
- **常量**: `UPPER_SNAKE_CASE` (如 `MAX_POST_LENGTH`)
- **枚举**: `PascalCase` (如 `PostStatus`)

## Import Patterns

### 后端 (Go)

#### Import Order
1. 标准库 (`fmt`, `context`, `time`)
2. 第三方库 (`google.golang.org/grpc`, `github.com/...`)
3. 项目内部包 (`internal/domain`, `pkg/utils`)

#### 示例
```go
package handler

import (
    "context"
    "time"
    
    "google.golang.org/grpc"
    "github.com/lib/pq"
    
    "fuck_boss/internal/domain/content"
    "fuck_boss/pkg/errors"
)
```

### 前端 (TypeScript)

#### Import Order
1. React 相关 (`react`, `react-dom`)
2. 第三方库 (`axios`, `zustand`)
3. 项目内部 (`@/api`, `@/domain`, `@/shared`)

#### 路径别名配置
```typescript
// tsconfig.json
{
  "compilerOptions": {
    "paths": {
      "@/*": ["src/*"],
      "@/api/*": ["src/api/*"],
      "@/domain/*": ["src/domain/*"],
      "@/shared/*": ["src/shared/*"]
    }
  }
}
```

## Code Structure Patterns

### 后端 (Go) - DDD 分层

#### Domain Layer 文件结构
```go
// domain/content/entity.go
package content

// Post 聚合根
type Post struct {
    id        PostID
    company   CompanyName
    city      City
    content   Content
    createdAt time.Time
}

// 业务方法
func (p *Post) Publish() error {
    // 领域逻辑
}

// domain/content/repository.go
package content

type Repository interface {
    Save(ctx context.Context, post *Post) error
    FindByID(ctx context.Context, id PostID) (*Post, error)
}
```

#### Application Layer 文件结构
```go
// application/content/create_post.go
package content

type CreatePostUseCase struct {
    repo domain.Repository
}

func (uc *CreatePostUseCase) Execute(ctx context.Context, cmd CreatePostCommand) (*PostDTO, error) {
    // 1. 验证输入
    // 2. 创建领域对象
    // 3. 调用 Repository
    // 4. 返回 DTO
}
```

#### Infrastructure Layer 文件结构
```go
// infrastructure/persistence/postgres/post_repository.go
package postgres

type PostRepository struct {
    db *sql.DB
}

func (r *PostRepository) Save(ctx context.Context, post *domain.Post) error {
    // 实现数据库操作
}
```

### 前端 (React) - 功能模块组织

#### 功能模块结构
```typescript
// features/post/components/PostList.tsx
export const PostList: React.FC<PostListProps> = ({ posts }) => {
    // 组件实现
}

// features/post/hooks/usePosts.ts
export const usePosts = () => {
    // Hook 实现
}

// features/post/index.ts
export { PostList } from './components/PostList'
export { usePosts } from './hooks/usePosts'
```

## Code Organization Principles

1. **单一职责**: 每个文件、函数只做一件事
2. **依赖倒置**: 依赖接口而非实现
3. **领域驱动**: 业务逻辑集中在 Domain Layer
4. **测试友好**: 代码结构便于测试
5. **一致性**: 遵循项目既定的模式和约定

## Module Boundaries

### 后端分层边界

```
Presentation Layer
    ↓ (依赖)
Application Layer
    ↓ (依赖)
Domain Layer
    ↑ (实现)
Infrastructure Layer
```

**规则**:
- Domain Layer 不依赖任何其他层（最纯净）
- Application Layer 只依赖 Domain Layer
- Infrastructure Layer 实现 Domain Layer 定义的接口
- Presentation Layer 调用 Application Layer 的 Use Cases

### 前端模块边界

```
Features (功能模块)
    ↓ (使用)
Domain (领域模型)
    ↓ (使用)
API (API 客户端)
    ↓ (使用)
Shared (共享组件)
```

**规则**:
- Features 之间相互独立，通过 Shared 共享
- Domain 定义业务模型，不依赖 UI
- API 层独立，可替换实现

## Code Size Guidelines

### 后端 (Go)
- **文件大小**: 建议 < 500 行
- **函数大小**: 建议 < 50 行
- **结构体复杂度**: 建议 < 10 个字段
- **嵌套深度**: 建议 < 4 层

### 前端 (TypeScript)
- **组件文件**: 建议 < 300 行
- **函数大小**: 建议 < 50 行
- **组件复杂度**: 建议 < 10 个 props
- **嵌套深度**: 建议 < 3 层 JSX

## Documentation Standards

### 文档组织原则

**所有文档统一放在 `docs/` 目录下，按类型分类**:

```
docs/
├── design/              # 设计文档（架构、数据库、API）
├── development/         # 开发文档（环境搭建、开发指南、测试）
├── deployment/          # 部署文档
├── fixes/               # 修复文档（问题修复记录）
│   ├── bug-fixes/       # Bug 修复
│   ├── performance-fixes/  # 性能优化
│   └── security-fixes/  # 安全修复
└── changelog/           # 变更日志
```

### 文档命名规范

- **设计文档**: `kebab-case.md` (如 `database-schema.md`)
- **开发文档**: `kebab-case-guide.md` (如 `setup-guide.md`)
- **修复文档**: `YYYY-MM-DD-issue-description.md` (如 `2026-01-06-redis-connection-fix.md`)

### 代码文档要求

#### 后端 (Go)
- **公开函数**: 必须有 `// FunctionName 描述` 注释
- **公开类型**: 必须有类型说明
- **复杂逻辑**: 必须有行内注释
- **包级别**: 每个包必须有 `package` 注释

```go
// Package content provides domain models and services for content management.
package content

// Post represents a company misconduct exposure post.
type Post struct {
    // ID is the unique identifier of the post.
    ID PostID
    
    // Company is the name of the company being exposed.
    Company CompanyName
}

// CreatePost creates a new post in the system.
// It validates the input and persists the post to the repository.
func CreatePost(ctx context.Context, repo Repository, cmd CreateCommand) (*Post, error) {
    // Implementation
}
```

#### 前端 (TypeScript)
- **组件**: 必须有 JSDoc 注释说明 props
- **函数**: 复杂函数必须有注释
- **类型**: 公开类型必须有注释

```typescript
/**
 * PostList component displays a list of company exposure posts.
 * 
 * @param posts - Array of posts to display
 * @param onPostClick - Callback when a post is clicked
 */
export const PostList: React.FC<PostListProps> = ({ posts, onPostClick }) => {
    // Implementation
}
```

### 修复文档模板

每次修复问题后，在 `docs/fixes/` 相应目录下创建修复文档：

```markdown
# [日期] - [问题描述]

## 问题描述
[详细描述问题]

## 原因分析
[分析问题根本原因]

## 解决方案
[描述解决方案]

## 测试验证
[描述如何验证修复]

## 相关文件
- `path/to/file.go`
- `path/to/file.ts`

## 影响范围
[描述影响范围]
```

## Testing Structure

### 测试文件组织

```
backend/test/
├── unit/                    # 单元测试
│   ├── domain/             # 领域层测试
│   ├── application/        # 应用层测试
│   └── infrastructure/     # 基础设施层测试
│
├── integration/            # 集成测试
│   ├── repository/        # Repository 集成测试
│   ├── cache/             # 缓存集成测试
│   └── grpc/              # gRPC 服务集成测试
│
└── e2e/                    # 端到端测试
    └── scenarios/          # 用户场景测试

frontend/test/
├── unit/                   # 单元测试
│   ├── components/        # 组件测试
│   ├── hooks/             # Hook 测试
│   └── utils/             # 工具函数测试
│
├── integration/           # 集成测试
│   └── api/               # API 集成测试
│
└── e2e/                    # 端到端测试
    └── flows/              # 用户流程测试
```

### 测试文件命名

- **单元测试**: `*_test.go` (Go), `*.test.tsx` (TypeScript)
- **集成测试**: `*_integration_test.go`, `*.integration.test.ts`
- **E2E 测试**: `*_e2e_test.go`, `*.e2e.test.ts`

## Build & Scripts

### Makefile 结构

```makefile
.PHONY: help build test lint

# 后端
backend-build:
backend-test:
backend-lint:

# 前端
frontend-build:
frontend-test:
frontend-lint:

# 数据库
db-migrate:
db-rollback:

# Docker
docker-up:
docker-down:

# 测试环境
test-up:        # 启动测试环境（PostgreSQL + Redis）
test-down:
```

### 脚本文件位置

- **构建脚本**: `backend/scripts/`, `frontend/scripts/`
- **部署脚本**: `scripts/deploy/`
- **工具脚本**: `scripts/tools/`

