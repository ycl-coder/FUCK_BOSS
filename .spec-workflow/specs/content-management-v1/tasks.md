# Tasks Document

## 任务说明

本文档将设计文档拆分为可执行的开发任务。每个任务都包含：
- 文件路径和实现内容
- 依赖关系和可复用代码
- 对应的需求编号
- 详细的实现指引（_Prompt）
- 验证要求

**重要原则**：
- 每个任务完成后必须完成测试验证
- 三方组件（PostgreSQL、Redis）必须调通后再进入下一步
- 遵循 Go 代码规范和 DDD 架构原则

---

## Phase 1: 基础设施和基础组件

### 任务 1.1: 创建项目基础结构

- [x] 1.1. 初始化 Go 模块和项目目录结构
  - Files: `backend/go.mod`, `backend/go.sum`, 创建目录结构
  - 初始化 Go 模块：`go mod init fuck_boss/backend`
  - 创建 DDD 分层目录结构（domain, application, infrastructure, presentation）
  - 创建测试目录结构（unit, integration, e2e）
  - Purpose: 建立项目基础结构
  - _Leverage: 无（新项目）
  - _Requirements: 基础设施要求
  - _Prompt: Role: DevOps Engineer specializing in Go project setup | Task: Initialize Go module and create complete DDD directory structure following structure.md guidelines. Create all necessary directories: backend/internal/domain/, backend/internal/application/, backend/internal/infrastructure/, backend/internal/presentation/, backend/pkg/, backend/test/. Ensure go.mod is properly initialized with module name. | Restrictions: Follow Go modules best practices, do not create unnecessary files, maintain clean directory structure | Success: go.mod initialized correctly, all directories created, project structure matches structure.md specification_

### 任务 1.2: 创建错误处理包

- [x] 1.2. 实现统一的错误处理包
  - File: `backend/pkg/errors/errors.go`
  - 定义错误码类型和错误结构
  - 实现错误创建和转换函数
  - Purpose: 提供统一的错误处理机制
  - _Leverage: 无（基础组件）
  - _Requirements: 错误处理要求
  - _Prompt: Role: Go Developer specializing in error handling patterns | Task: Create comprehensive error handling package following Go best practices and tech.md code standards. Define ErrorCode type with constants (VALIDATION_ERROR, NOT_FOUND, RATE_LIMIT_EXCEEDED, INTERNAL_ERROR, DATABASE_ERROR). Create AppError struct with Code, Message, Details, and Cause fields. Implement error creation functions and error wrapping using fmt.Errorf with %w verb (Go 1.13+). Follow Go code standards: error messages without punctuation, proper error wrapping. | Restrictions: Must follow Go error handling conventions, error messages must be clear and user-friendly, support error wrapping for debugging | Success: Error package compiles without errors, all error types defined, error wrapping works correctly, follows Go code standards_

### 任务 1.3: 创建配置管理

- [x] 1.3. 实现配置管理组件
  - File: `backend/internal/infrastructure/config/config.go`
  - 使用 viper 读取配置文件
  - 定义配置结构体（数据库、Redis、gRPC 等）
  - Purpose: 统一管理应用配置
  - _Leverage: `github.com/spf13/viper`
  - _Requirements: 配置管理要求
  - _Prompt: Role: Backend Developer specializing in configuration management | Task: Implement configuration management using viper following tech.md standards. Create Config struct with fields for Database (host, port, user, password, dbname), Redis (host, port, password, db), gRPC (port), and Log (level, format). Implement LoadConfig function that reads from environment variables and config file (YAML). Support default values and validation. Follow Go code standards: proper error handling, clear function names. | Restrictions: Must use viper, support environment variable overrides, validate required fields, provide sensible defaults | Success: Configuration loads correctly from file and environment, validation works, all required fields are validated, follows Go code standards_

### 任务 1.4: 创建日志工具

- [x] 1.4. 实现结构化日志组件
  - File: `backend/internal/infrastructure/logger/logger.go`
  - 使用 zap 实现结构化日志
  - 提供不同级别的日志函数
  - Purpose: 提供统一的日志记录功能
  - _Leverage: `go.uber.org/zap`
  - _Requirements: 日志要求
  - _Prompt: Role: Backend Developer specializing in logging and observability | Task: Implement structured logging using zap following tech.md standards. Create Logger interface and zap implementation. Support different log levels (Debug, Info, Warn, Error). Implement context-aware logging with request ID support. Provide both development (console) and production (JSON) log formats. Follow Go code standards: proper error handling, clear function signatures. | Restrictions: Must use zap, support structured logging with fields, provide both dev and prod formats, handle errors gracefully | Success: Logger initializes correctly, all log levels work, structured logging with fields works, JSON format for production, follows Go code standards_

---

## Phase 2: Domain Layer（领域层）

### 任务 2.1: 创建值对象（Value Objects）

- [x] 2.1. 实现 PostID 值对象
  - File: `backend/internal/domain/content/value_object.go`
  - 定义 PostID 类型和验证逻辑
  - 实现 NewPostID 工厂方法
  - Purpose: 封装 Post ID 的业务规则
  - _Leverage: `github.com/google/uuid`
  - _Requirements: Requirement 1, Design - Value Objects
  - _Prompt: Role: Domain-Driven Design specialist with Go expertise | Task: Create PostID value object following DDD principles and tech.md code standards. Define PostID struct with value field (string, UUID format). Implement NewPostID factory function that validates UUID format and returns error if invalid. Implement String() method. Follow Go code standards: proper error handling, clear validation logic, value objects are immutable. | Restrictions: Must validate UUID format, value objects must be immutable, follow Go naming conventions | Success: PostID validates correctly, factory method works, String() method returns correct value, follows DDD principles and Go code standards_

- [x] 2.2. 实现 CompanyName 值对象
  - File: `backend/internal/domain/content/value_object.go` (继续)
  - 定义 CompanyName 类型和验证逻辑（1-100 字符）
  - 实现 NewCompanyName 工厂方法
  - Purpose: 封装公司名称的业务规则
  - _Leverage: 无
  - _Requirements: Requirement 1, Design - Value Objects
  - _Prompt: Role: Domain-Driven Design specialist with Go expertise | Task: Create CompanyName value object following DDD principles. Define CompanyName struct with value field (string). Implement NewCompanyName factory function that validates: non-empty, length 1-100 characters, trim whitespace. Return error if validation fails. Implement String() method. Follow Go code standards: proper error handling, clear validation messages. | Restrictions: Must validate length (1-100), trim whitespace, value objects are immutable | Success: CompanyName validates correctly, rejects invalid inputs, String() works, follows DDD principles_

- [x] 2.3. 实现 Content 值对象
  - File: `backend/internal/domain/content/value_object.go` (继续)
  - 定义 Content 类型和验证逻辑（10-5000 字符）
  - 实现 NewContent 工厂方法和 Summary 方法
  - Purpose: 封装内容的业务规则
  - _Leverage: 无
  - _Requirements: Requirement 1, Design - Value Objects
  - _Prompt: Role: Domain-Driven Design specialist with Go expertise | Task: Create Content value object following DDD principles. Define Content struct with value field (string). Implement NewContent factory function that validates: non-empty, length 10-5000 characters. Implement String() method and Summary() method that returns first 200 characters with ellipsis if longer. Follow Go code standards: proper error handling, efficient string operations. | Restrictions: Must validate length (10-5000), Summary must handle edge cases, value objects are immutable | Success: Content validates correctly, Summary works for all cases, follows DDD principles_

- [x] 2.4. 实现 City 值对象
  - File: `backend/internal/domain/shared/city.go`
  - 定义 City 类型（code 和 name）
  - 实现 NewCity 工厂方法
  - Purpose: 封装城市的业务规则
  - _Leverage: 无
  - _Requirements: Requirement 5, Design - Value Objects
  - _Prompt: Role: Domain-Driven Design specialist with Go expertise | Task: Create City value object following DDD principles. Define City struct with code (string) and name (string) fields. Implement NewCity factory function that validates both fields are non-empty. Implement Code() and Name() methods. Follow Go code standards: proper error handling, clear validation. | Restrictions: Must validate both code and name, value objects are immutable | Success: City validates correctly, accessor methods work, follows DDD principles_

- [x] 2.5. 值对象单元测试
  - File: `backend/test/unit/domain/content/value_object_test.go`, `backend/test/unit/domain/shared/city_test.go`
  - 测试所有值对象的验证逻辑
  - 测试边界情况和错误情况
  - Purpose: 确保值对象正确性
  - _Leverage: `github.com/stretchr/testify`
  - _Requirements: Testing Strategy - Unit Testing
  - _Prompt: Role: QA Engineer specializing in unit testing with Go | Task: Create comprehensive unit tests for all value objects (PostID, CompanyName, Content, City) covering requirements. Test valid inputs, invalid inputs, boundary cases (min/max length), edge cases. Use testify for assertions. Follow Go testing standards: test file naming (_test.go), test function naming (TestXxx), use table-driven tests where appropriate. | Restrictions: Must test all validation rules, cover edge cases, tests must be independent, use testify assertions | Success: All value object tests pass, coverage >= 90%, edge cases covered, tests are maintainable_

### 任务 2.2: 创建 Post Entity（聚合根）

- [x] 2.6. 实现 Post Entity
  - File: `backend/internal/domain/content/entity.go`
  - 定义 Post 结构体和业务方法
  - 实现 NewPost 工厂方法和 Publish 方法
  - Purpose: 实现 Post 聚合根，包含业务规则
  - _Leverage: Value Objects (PostID, CompanyName, City, Content)
  - _Requirements: Requirement 1, Design - Post Entity
  - _Prompt: Role: Domain-Driven Design specialist with Go expertise | Task: Create Post entity (aggregate root) following DDD principles and tech.md code standards. Define Post struct with private fields: id (PostID), company (CompanyName), city (City), content (Content), createdAt (time.Time). Implement NewPost factory function that generates UUID for id, sets createdAt to now, validates all value objects. Implement Publish() method (business method, can add validation logic). Implement getter methods: ID(), Company(), City(), Content(), CreatedAt(). Follow Go code standards: proper error handling, clear business logic, entities are mutable but controlled. | Restrictions: Must use value objects, id must be generated (UUID), createdAt must be set automatically, follow Go naming conventions | Success: Post entity works correctly, factory method validates inputs, getter methods work, follows DDD principles and Go code standards_

- [x] 2.7. Post Entity 单元测试
  - File: `backend/test/unit/domain/content/entity_test.go`
  - 测试 Post 创建、验证、业务方法
  - Purpose: 确保 Post Entity 正确性
  - _Leverage: `github.com/stretchr/testify`
  - _Requirements: Testing Strategy - Unit Testing
  - _Prompt: Role: QA Engineer specializing in unit testing with Go | Task: Create comprehensive unit tests for Post entity covering requirements. Test NewPost factory with valid/invalid inputs, test getter methods, test Publish() method, test business rules. Use testify for assertions. Follow Go testing standards. | Restrictions: Must test all business rules, cover edge cases, tests must be independent | Success: All Post entity tests pass, coverage >= 90%, business rules validated_

### 任务 2.3: 创建 Repository 接口

- [x] 2.8. 定义 PostRepository 接口
  - File: `backend/internal/domain/content/repository.go`
  - 定义所有 Repository 方法签名
  - Purpose: 定义持久化接口，遵循依赖倒置原则
  - _Leverage: Domain Entities (Post)
  - _Requirements: Design - PostRepository Interface
  - _Prompt: Role: Domain-Driven Design specialist with Go expertise | Task: Create PostRepository interface following DDD principles and tech.md code standards. Define interface with methods: Save(ctx context.Context, post *Post) error, FindByID(ctx context.Context, id PostID) (*Post, error), FindByCity(ctx context.Context, city City, page, pageSize int) ([]*Post, int, error), Search(ctx context.Context, keyword string, city *City, page, pageSize int) ([]*Post, int, error). All methods must use context.Context as first parameter. Return error as last parameter. Follow Go code standards: interface naming (Repository suffix), clear method signatures, proper error handling. | Restrictions: Must follow Go interface conventions, use context.Context, error as last return value, interface must be in domain layer | Success: Interface compiles correctly, all methods properly defined, follows DDD principles and Go code standards_

---

## Phase 3: Infrastructure Layer（基础设施层）

### 任务 3.1: PostgreSQL Repository 实现

- [x] 3.1. 创建数据库迁移脚本
  - Files: `backend/internal/infrastructure/persistence/postgres/migrations/000001_create_posts_table.up.sql`, `000001_create_posts_table.down.sql`
  - 创建 posts 表和索引
  - 创建 cities 表
  - Purpose: 定义数据库结构
  - _Leverage: `golang-migrate/migrate`
  - _Requirements: Design - Database Schema
  - _Prompt: Role: Database Engineer specializing in PostgreSQL and migrations | Task: Create database migration scripts following design.md schema specification. Create posts table with all fields (id UUID PRIMARY KEY, company_name VARCHAR(100), city_code VARCHAR(50), city_name VARCHAR(50), content TEXT, occurred_at TIMESTAMP, created_at TIMESTAMP, updated_at TIMESTAMP). Create indexes: idx_posts_city_code, idx_posts_created_at, idx_posts_company_name. Create full-text search index using pg_jieba (to_tsvector). Create cities table. Create up and down migrations. Follow PostgreSQL best practices: proper data types, indexes for performance, constraints. | Restrictions: Must match design.md schema exactly, include all indexes, support rollback (down migration), use proper PostgreSQL types | Success: Migrations run successfully, schema matches design, indexes created, rollback works_

- [x] 3.2. 实现 PostgreSQL Repository
  - File: `backend/internal/infrastructure/persistence/postgres/post_repository.go`
  - 实现 PostRepository 接口的所有方法
  - 使用 database/sql 和 lib/pq
  - Purpose: 实现 Post 的持久化
  - _Leverage: `database/sql`, `github.com/lib/pq`, Domain Repository Interface
  - _Requirements: Design - PostgreSQL Repository
  - _Prompt: Role: Backend Developer specializing in PostgreSQL and Go database operations | Task: Implement PostRepository interface using PostgreSQL following design.md and tech.md code standards. Use database/sql and lib/pq. Implement all methods: Save (INSERT with RETURNING), FindByID (SELECT by id), FindByCity (SELECT with WHERE city_code, ORDER BY created_at DESC, LIMIT/OFFSET for pagination, COUNT for total), Search (SELECT with full-text search using to_tsquery, support city filter, pagination). Use parameterized queries (prevent SQL injection). Handle errors properly. Follow Go code standards: proper error handling, use context for cancellation, transactions where needed, file length < 800 lines, function length < 80 lines. | Restrictions: Must use parameterized queries, handle all errors, support pagination correctly, implement full-text search, follow Go code standards | Success: Repository implements interface correctly, all methods work, SQL injection prevented, error handling proper, follows Go code standards_

- [x] 3.3. PostgreSQL Repository 集成测试 ⚠️ **必须完成**
  - File: `backend/test/integration/repository/post_repository_test.go`
  - 使用真实 PostgreSQL 数据库测试
  - 测试所有 CRUD 操作和搜索
  - Purpose: 验证 Repository 与数据库集成
  - _Leverage: Docker Compose (test database), `github.com/stretchr/testify`
  - _Requirements: Testing Strategy - Integration Testing, Design - Implementation Verification
  - _Prompt: Role: QA Engineer specializing in integration testing with Go and PostgreSQL | Task: Create comprehensive integration tests for PostgreSQL Repository using real PostgreSQL database. Use Docker Compose to start test database. Test all methods: Save (create new post, verify saved), FindByID (find existing, not found), FindByCity (pagination, ordering), Search (full-text search, city filter, pagination). Each test should use transaction and rollback after test. Test error cases (database errors). Follow Go testing standards. Use testify for assertions. **CRITICAL: All tests must pass before proceeding to next task.** | Restrictions: Must use real PostgreSQL, each test uses transaction with rollback, test all methods, test error cases, tests must be independent | Success: All integration tests pass, Repository works correctly with real database, full-text search works, pagination works, error handling verified, **ready to proceed to next task**_

### 任务 3.2: Redis Cache 实现

- [x] 3.4. 实现 Redis Cache Repository
  - File: `backend/internal/infrastructure/persistence/redis/cache_repository.go`
  - 实现 CacheRepository 接口
  - 使用 go-redis/v9
  - Purpose: 实现缓存功能
  - _Leverage: `github.com/redis/go-redis/v9`
  - _Requirements: Design - Redis Cache
  - _Prompt: Role: Backend Developer specializing in Redis and caching strategies | Task: Implement CacheRepository interface using Redis following design.md and tech.md code standards. Use go-redis/v9. Implement methods: Get (GET command), Set (SET with TTL), Delete (DEL), DeleteByPattern (SCAN + DEL). Handle Redis errors gracefully (return errors, don't panic). Support context for cancellation. Follow Go code standards: proper error handling, use context, file length < 800 lines, function length < 80 lines. | Restrictions: Must handle Redis errors, support TTL, implement pattern deletion, use context, follow Go code standards | Success: Cache repository implements interface correctly, all methods work, error handling proper, follows Go code standards_

- [x] 3.5. 实现 Redis Rate Limiter
  - File: `backend/internal/infrastructure/persistence/redis/rate_limiter.go`
  - 实现 RateLimiter 接口
  - 使用 Redis 实现滑动窗口限流
  - Purpose: 实现限流功能（防刷）
  - _Leverage: `github.com/redis/go-redis/v9`
  - _Requirements: Requirement 1 (限流要求), Design - Rate Limiter
  - _Prompt: Role: Backend Developer specializing in rate limiting and Redis | Task: Implement RateLimiter interface using Redis following design.md. Use sliding window algorithm with Redis. Implement Allow method: check if key (e.g., "rate_limit:post:{ip}:{hour}") count < limit within window. Use Redis INCR and EXPIRE. Return (allowed bool, error). Handle Redis errors. Follow Go code standards. | Restrictions: Must use sliding window algorithm, handle Redis errors, support configurable limit and window, use context | Success: Rate limiter works correctly, sliding window algorithm implemented, error handling proper, follows Go code standards_

- [x] 3.6. Redis Cache 集成测试 ⚠️ **必须完成**
  - File: `backend/test/integration/cache/redis_cache_test.go`
  - 使用真实 Redis 测试
  - 测试所有缓存操作和限流
  - Purpose: 验证 Cache 与 Redis 集成
  - _Leverage: Docker Compose (test Redis), `github.com/stretchr/testify`
  - _Requirements: Testing Strategy - Integration Testing, Design - Implementation Verification
  - _Prompt: Role: QA Engineer specializing in integration testing with Go and Redis | Task: Create comprehensive integration tests for Redis Cache and Rate Limiter using real Redis. Use Docker Compose to start test Redis. Test CacheRepository: Get/Set (verify value, TTL), Delete, DeleteByPattern. Test RateLimiter: Allow (within limit, exceed limit, window expiration). Test error handling (Redis connection failure - graceful degradation). Each test should clean up after itself. Follow Go testing standards. Use testify for assertions. **CRITICAL: All tests must pass before proceeding to next task.** | Restrictions: Must use real Redis, test all methods, test TTL, test rate limiting, test error handling, tests must be independent | Success: All integration tests pass, Cache works correctly with real Redis, Rate Limiter works, error handling verified, **ready to proceed to next task**_

---

## Phase 4: Application Layer（应用层）

### 任务 4.1: CreatePost UseCase

- [x] 4.1. 实现 CreatePost UseCase
  - File: `backend/internal/application/content/create_post.go`
  - 实现创建帖子的用例逻辑
  - 集成限流、验证、Repository
  - Purpose: 处理创建曝光内容的业务逻辑
  - _Leverage: Domain Repository, Cache Repository, Rate Limiter, Validator
  - _Requirements: Requirement 1, Design - CreatePost UseCase
  - _Prompt: Role: Backend Developer specializing in application layer and use cases | Task: Implement CreatePostUseCase following design.md and tech.md code standards. Create CreatePostCommand struct and CreatePostUseCase struct. Implement Execute method: 1) Validate input using validator, 2) Check rate limit (1 hour, 3 requests per IP), 3) Create domain entities (CompanyName, City, Content, Post), 4) Save to repository, 5) Clear related cache (city list cache), 6) Return PostDTO. Handle all errors properly. Follow Go code standards: proper error handling, use context, file length < 800 lines, function length < 80 lines. | Restrictions: Must validate input, check rate limit, use domain entities, handle all errors, clear cache, follow Go code standards | Success: UseCase works correctly, validation works, rate limiting works, cache cleared, error handling proper, follows Go code standards_

- [x] 4.2. CreatePost UseCase 单元测试
  - File: `backend/test/unit/application/content/create_post_test.go`
  - 使用 Mock Repository 测试用例逻辑
  - Purpose: 确保 UseCase 业务逻辑正确
  - _Leverage: `github.com/stretchr/testify/mock`, `github.com/golang/mock/gomock`
  - _Requirements: Testing Strategy - Unit Testing
  - _Prompt: Role: QA Engineer specializing in unit testing with Go and mocking | Task: Create comprehensive unit tests for CreatePostUseCase using mocked dependencies. Mock PostRepository, CacheRepository, RateLimiter. Test scenarios: valid input (success), invalid input (validation error), rate limit exceeded, repository error, cache error. Use gomock or testify/mock for mocking. Follow Go testing standards. | Restrictions: Must mock all dependencies, test all scenarios, tests must be independent | Success: All unit tests pass, all scenarios covered, mocking works correctly_

- [x] 4.3. CreatePost UseCase 集成测试
  - File: `backend/test/integration/usecase/create_post_test.go`
  - 使用真实 Repository 和 Cache 测试
  - Purpose: 验证 UseCase 与基础设施集成
  - _Leverage: Real PostgreSQL Repository, Real Redis Cache
  - _Requirements: Testing Strategy - Integration Testing
  - _Prompt: Role: QA Engineer specializing in integration testing | Task: Create integration tests for CreatePostUseCase using real PostgreSQL Repository and Redis Cache. Test complete flow: create post, verify saved in database, verify cache cleared, verify rate limiting works. Use Docker Compose for test environment. Follow Go testing standards. **CRITICAL: All tests must pass.** | Restrictions: Must use real infrastructure, test complete flow, verify all side effects | Success: All integration tests pass, complete flow works, **ready to proceed**_

### 任务 4.2: ListPosts UseCase

- [x] 4.4. 实现 ListPosts UseCase
  - File: `backend/internal/application/content/list_posts.go`
  - 实现列表查询用例逻辑
  - 集成缓存策略
  - Purpose: 处理列表查询业务逻辑
  - _Leverage: Domain Repository, Cache Repository
  - _Requirements: Requirement 2, Design - ListPosts UseCase
  - _Prompt: Role: Backend Developer specializing in application layer | Task: Implement ListPostsUseCase following design.md. Create ListPostsQuery struct and ListPostsUseCase struct. Implement Execute method: 1) Check cache first (key: "posts:city:{cityCode}:page:{page}"), 2) If cache miss, query repository, 3) Update cache (TTL: 5-10 minutes based on city popularity), 4) Return PostsListDTO. Handle cache errors gracefully (fallback to database). Follow Go code standards. | Restrictions: Must implement caching strategy, handle cache errors, support pagination, follow Go code standards | Success: UseCase works correctly, caching works, pagination works, error handling proper_

- [x] 4.5. ListPosts UseCase 测试
  - Files: `backend/test/unit/application/content/list_posts_test.go`, `backend/test/integration/usecase/list_posts_test.go`
  - 单元测试和集成测试
  - Purpose: 确保列表查询正确性
  - _Leverage: Mock dependencies, Real infrastructure
  - _Requirements: Testing Strategy
  - _Prompt: Role: QA Engineer | Task: Create unit and integration tests for ListPostsUseCase. Test caching (hit, miss), pagination, error handling. Use mocks for unit tests, real infrastructure for integration tests. | Restrictions: Must test caching, pagination, error handling | Success: All tests pass, caching verified, pagination verified_

### 任务 4.3: GetPost UseCase

- [x] 4.6. 实现 GetPost UseCase
  - File: `backend/internal/application/content/get_post.go`
  - 实现详情查询用例逻辑
  - Purpose: 处理详情查询业务逻辑
  - _Leverage: Domain Repository, Cache Repository
  - _Requirements: Requirement 3, Design - GetPost UseCase
  - _Prompt: Role: Backend Developer | Task: Implement GetPostUseCase following design.md. Create GetPostUseCase struct. Implement Execute method: 1) Check cache (key: "post:{postID}"), 2) If cache miss, query repository, 3) Update cache (TTL: 10 minutes), 4) Return PostDTO or NotFound error. Follow Go code standards. | Restrictions: Must implement caching, handle NotFound, follow Go code standards | Success: UseCase works correctly, caching works, NotFound handled_

- [x] 4.7. GetPost UseCase 测试
  - Files: `backend/test/unit/application/content/get_post_test.go`, `backend/test/integration/usecase/get_post_test.go`
  - 单元测试和集成测试
  - Purpose: 确保详情查询正确性
  - _Leverage: Mock dependencies, Real infrastructure
  - _Requirements: Testing Strategy
  - _Prompt: Role: QA Engineer | Task: Create unit and integration tests for GetPostUseCase. Test caching, NotFound case, error handling. | Restrictions: Must test all scenarios | Success: All tests pass_

### 任务 4.4: SearchPosts UseCase

- [x] 4.8. 实现 SearchPosts UseCase
  - File: `backend/internal/application/search/search_posts.go`
  - 实现搜索用例逻辑
  - Purpose: 处理搜索业务逻辑
  - _Leverage: Domain Repository, Cache Repository
  - _Requirements: Requirement 4, Design - SearchPosts UseCase
  - _Prompt: Role: Backend Developer | Task: Implement SearchPostsUseCase following design.md. Create SearchPostsQuery struct and SearchPostsUseCase struct. Implement Execute method: 1) Validate keyword (min 2 chars), 2) Check cache (key: "search:{keyword}:city:{cityCode}:page:{page}"), 3) If cache miss, query repository with full-text search, 4) Update cache (TTL: 5 minutes), 5) Return PostsListDTO. Handle cache errors gracefully. Follow Go code standards. | Restrictions: Must validate keyword, implement caching, support city filter, follow Go code standards | Success: UseCase works correctly, search works, caching works_

- [x] 4.9. SearchPosts UseCase 测试
  - Files: `backend/test/unit/application/search/search_posts_test.go`, `backend/test/integration/usecase/search_posts_test.go`
  - 单元测试和集成测试
  - Purpose: 确保搜索正确性
  - _Leverage: Mock dependencies, Real infrastructure
  - _Requirements: Testing Strategy
  - _Prompt: Role: QA Engineer | Task: Create unit and integration tests for SearchPostsUseCase. Test keyword validation, search functionality, city filter, caching, error handling. | Restrictions: Must test all scenarios | Success: All tests pass, search verified_

---

## Phase 5: Presentation Layer（表现层）

### 任务 5.1: Protocol Buffers 定义

- [x] 5.1. 定义 gRPC API（Protocol Buffers）
  - Files: `backend/api/proto/content/v1/content.proto`
  - 定义所有 gRPC 服务和消息
  - Purpose: 定义 API 契约
  - _Leverage: Protocol Buffers
  - _Requirements: Design - Protocol Buffers Schema
  - _Prompt: Role: API Designer specializing in gRPC and Protocol Buffers | Task: Create Protocol Buffers definition following design.md schema. Define ContentService with methods: CreatePost, ListPosts, GetPost, SearchPosts. Define all request/response messages. Use proper field numbers, follow protobuf best practices. Set go_package option correctly. | Restrictions: Must match design.md schema exactly, use proper field numbers, follow protobuf conventions | Success: proto file compiles correctly, all services and messages defined_

- [x] 5.2. 生成 gRPC 代码
  - Files: `backend/scripts/generate.sh`, 生成的 Go 代码
  - 使用 protoc 生成 Go 代码
  - Purpose: 生成 gRPC 客户端和服务器代码
  - _Leverage: `protoc`, `protoc-gen-go`, `protoc-gen-go-grpc`
  - _Requirements: gRPC 代码生成
  - _Prompt: Role: DevOps Engineer specializing in code generation | Task: Create script to generate gRPC Go code from proto files. Use protoc with plugins: protoc-gen-go, protoc-gen-go-grpc. Generate code to correct output directory. Create Makefile target or script. | Restrictions: Must use correct plugins, generate to correct directories | Success: Code generation works, all Go files generated correctly_

### 任务 5.2: gRPC Handlers

- [x] 5.3. 实现 gRPC Content Service
  - File: `backend/internal/presentation/grpc/content_handler.go`
  - 实现所有 gRPC 方法
  - 调用 Application Layer UseCases
  - Purpose: 处理 gRPC 请求
  - _Leverage: Application UseCases, Generated gRPC code
  - _Requirements: Design - gRPC Handlers
  - _Prompt: Role: Backend Developer specializing in gRPC | Task: Implement ContentService gRPC handler following design.md and tech.md code standards. Create ContentService struct with UseCase dependencies. Implement all methods: CreatePost (convert request to command, call usecase, convert response), ListPosts, GetPost, SearchPosts. Handle errors (convert to gRPC status codes). Use context for cancellation. Follow Go code standards: proper error handling, use context, file length < 800 lines. | Restrictions: Must handle all errors, convert to gRPC status, use context, follow Go code standards | Success: All gRPC methods work, error handling proper, follows Go code standards_

- [x] 5.4. 实现 gRPC Middleware
  - Files: `backend/internal/presentation/middleware/logging.go`, `recovery.go`
  - 实现日志和恢复中间件
  - Purpose: 提供请求日志和错误恢复
  - _Leverage: Logger infrastructure
  - _Requirements: Design - Middleware
  - _Prompt: Role: Backend Developer specializing in middleware | Task: Implement gRPC interceptors: LoggingInterceptor (log request/response, duration), RecoveryInterceptor (recover from panics, log stack trace). Use zap logger. Follow Go code standards. | Restrictions: Must log all requests, recover from panics, follow Go code standards | Success: Middleware works correctly, logging works, panic recovery works_

- [x] 5.5. gRPC Handler 测试
  - Files: `backend/test/unit/presentation/grpc/content_handler_test.go`, `backend/test/e2e/scenarios/grpc_test.go`
  - 单元测试和 E2E 测试
  - Purpose: 确保 gRPC 服务正确性
  - _Leverage: Mock UseCases, Real gRPC server
  - _Requirements: Testing Strategy
  - _Prompt: Role: QA Engineer | Task: Create unit tests (mock UseCases) and E2E tests (real gRPC server) for ContentService. Test all methods, error handling, middleware. | Restrictions: Must test all methods, error cases | Success: All tests pass, gRPC service works correctly_

### 任务 5.3: gRPC Server 启动

- [x] 5.6. 实现 gRPC Server 主程序
  - File: `backend/cmd/server/main.go`
  - 初始化所有依赖，启动 gRPC 服务器
  - Purpose: 应用入口点
  - _Leverage: 所有基础设施组件
  - _Requirements: Server startup
  - _Prompt: Role: Backend Developer specializing in application bootstrap | Task: Create main.go that initializes all dependencies: config, logger, database connection, Redis connection, repositories, usecases, gRPC service. Use dependency injection (wire or manual). Start gRPC server. Handle graceful shutdown. Follow Go code standards. | Restrictions: Must initialize all dependencies, handle graceful shutdown, follow Go code standards | Success: Server starts correctly, all dependencies initialized, graceful shutdown works_

---

## Phase 6: 前端开发（React）

### 任务 6.1: 前端基础设置

- [x] 6.1. 初始化 React 项目
  - Files: `frontend/package.json`, `frontend/vite.config.ts`, `frontend/tsconfig.json`
  - 使用 Vite + TypeScript + React
  - Purpose: 建立前端项目基础
  - _Leverage: Vite, React, TypeScript
  - _Requirements: Frontend setup
  - _Prompt: Role: Frontend Developer specializing in React and Vite | Task: Initialize React project with Vite, TypeScript, and React. Configure Vite for development and production builds. Set up TypeScript configuration. Install dependencies: React, React Router, gRPC Web client, UI library (antd or MUI), state management (zustand). Follow frontend best practices. | Restrictions: Must use Vite, TypeScript, React 18+, configure properly | Success: Project initializes correctly, all dependencies installed, build works_

- [x] 6.2. 设置 gRPC Web 客户端
  - Files: `frontend/src/api/grpc/client.ts`, `frontend/src/api/grpc/contentClient.ts`
  - 配置 gRPC Web 客户端
  - Purpose: 连接后端 gRPC 服务
  - _Leverage: `@grpc/grpc-js`, `@grpc-web/protoc-gen-grpc-web`
  - _Requirements: gRPC Web integration
  - _Prompt: Role: Frontend Developer specializing in gRPC Web | Task: Set up gRPC Web client. Generate TypeScript types from proto files. Create client wrapper for ContentService. Handle errors and connection. | Restrictions: Must use gRPC Web, handle errors properly | Success: gRPC Web client works, can connect to backend_

### 任务 6.2: 前端组件开发

- [x] 6.3. 创建发布表单组件
  - File: `frontend/src/features/post/components/PostForm.tsx`
  - 实现发布曝光内容的表单
  - Purpose: 用户发布界面
  - _Leverage: React, UI library, gRPC client
  - _Requirements: Requirement 1 (前端部分)
  - _Prompt: Role: Frontend Developer specializing in React forms | Task: Create PostForm component with fields: company (input), city (select), content (textarea), occurredAt (date picker, optional). Add validation (client-side). Handle form submission, loading state, error state. Use UI library components. Follow React best practices. | Restrictions: Must validate inputs, handle all states, use UI library | Success: Form works correctly, validation works, submission works_

- [x] 6.4. 创建列表组件
  - File: `frontend/src/features/post/components/PostList.tsx`
  - 实现内容列表展示
  - Purpose: 显示曝光内容列表
  - _Leverage: React, gRPC client, UI library
  - _Requirements: Requirement 2 (前端部分)
  - _Prompt: Role: Frontend Developer | Task: Create PostList component that displays list of posts. Show: company name, city, content summary, created time. Support pagination. Handle loading and error states. Use UI library components. | Restrictions: Must support pagination, handle all states | Success: List displays correctly, pagination works_

- [x] 6.5. 创建详情组件
  - File: `frontend/src/features/post/components/PostDetail.tsx`
  - 实现内容详情展示
  - Purpose: 显示内容详情
  - _Leverage: React, gRPC client
  - _Requirements: Requirement 3 (前端部分)
  - _Prompt: Role: Frontend Developer | Task: Create PostDetail component that displays full post details. Show all fields. Handle loading and error states (404). | Restrictions: Must handle 404, loading states | Success: Detail displays correctly, 404 handled_

- [x] 6.6. 创建搜索组件
  - File: `frontend/src/features/search/components/SearchBar.tsx`, `SearchResults.tsx`
  - 实现搜索功能
  - Purpose: 搜索曝光内容
  - _Leverage: React, gRPC client
  - _Requirements: Requirement 4 (前端部分)
  - _Prompt: Role: Frontend Developer | Task: Create SearchBar and SearchResults components. SearchBar: input field, city filter (optional), search button. SearchResults: display search results, highlight keywords, pagination. Handle debouncing. | Restrictions: Must support debouncing, keyword highlighting, pagination | Success: Search works correctly, debouncing works, highlighting works_

### 任务 6.3: 前端路由和集成

- [x] 6.7. 设置前端路由
  - File: `frontend/src/app/routes.tsx`
  - 配置 React Router
  - Purpose: 页面路由
  - _Leverage: `react-router-dom`
  - _Requirements: Frontend routing
  - _Prompt: Role: Frontend Developer | Task: Set up React Router with routes: / (home/list), /post/:id (detail), /create (create form), /search (search page). Create layout component. | Restrictions: Must use React Router, proper route structure | Success: Routing works correctly_

- [x] 6.8. 前端 E2E 测试
  - Files: `frontend/test/e2e/flows/create-and-view.spec.ts`, `search.spec.ts`, `navigation.spec.ts`
  - 使用 Playwright
  - Purpose: 验证完整用户流程
  - _Leverage: Playwright/Cypress
  - _Requirements: Testing Strategy - E2E Testing
  - _Prompt: Role: QA Automation Engineer | Task: Create E2E tests for complete user flows: create post → view list → view detail, search posts. Use Playwright or Cypress. Test all user interactions. | Restrictions: Must test complete flows, all interactions | Success: All E2E tests pass, user flows work correctly_

---

## Phase 7: 部署和文档

### 任务 7.1: Docker 和部署

- [x] 7.1. 创建 Docker 配置
  - Files: `docker-compose.yml`, `backend/Dockerfile`
  - 配置开发和生产环境
  - Purpose: 容器化部署
  - _Leverage: Docker, Docker Compose
  - _Requirements: Deployment
  - _Prompt: Role: DevOps Engineer specializing in Docker | Task: Create Docker configuration for development and production. docker-compose.yml: PostgreSQL, Redis, backend service, frontend service. Dockerfiles for backend and frontend. Configure networking, volumes, environment variables. | Restrictions: Must support development and production, proper networking | Success: Docker setup works, all services start correctly_

- [x] 7.2. 创建部署文档
  - Files: `docs/deployment/docker-deploy.md`, `docs/deployment/production-deploy.md`
  - 编写部署指南
  - Purpose: 部署文档
  - _Leverage: 无
  - _Requirements: Documentation
  - _Prompt: Role: Technical Writer | Task: Create deployment documentation. Include: prerequisites, setup steps, configuration, troubleshooting. | Restrictions: Must be clear and complete | Success: Documentation is complete and clear_

### 任务 7.2: 开发文档

- [x] 7.3. 创建开发指南
  - Files: `docs/development/setup-guide.md`, `docs/development/development-guide.md`, `docs/development/testing-guide.md`
  - 编写开发环境搭建和开发指南
  - Purpose: 帮助开发者快速上手
  - _Leverage: 无
  - _Requirements: Documentation
  - _Prompt: Role: Technical Writer | Task: Create development documentation: setup guide (how to set up dev environment), development guide (coding standards, architecture), testing guide (how to run tests). | Restrictions: Must be clear and complete | Success: Documentation is complete and helpful_

---

## 任务完成检查清单

每个任务完成后，必须完成以下检查：

### ⚠️ 代码提交（必须）

- [ ] **代码已提交到 Git 仓库**
  - 提交信息清晰描述完成的工作
  - 格式：`feat: [任务ID] 任务描述` 或 `fix: [任务ID] 修复描述`
  - 示例：`feat: [1.2] 实现统一错误处理包`
  - 已推送到远程仓库（如需要）
  - 确保敏感文件（如 `config.yaml`）不被提交

**重要**：未提交代码的任务视为未完成，不得进入下一个任务。

### 代码质量
- [ ] 8.1. 代码通过 `gofmt` 格式化
- [ ] 8.2. 代码通过 `go vet` 检查
- [ ] 8.3. 代码通过 `golangci-lint` 检查
- [ ] 8.4. 遵循 Go 代码规范（tech.md）
- [ ] 8.5. 文件长度 < 800 行
- [ ] 8.6. 函数长度 < 80 行
- [ ] 8.7. 嵌套深度 < 4 层

### 测试
- [ ] 8.8. 单元测试通过（覆盖率 >= 70%，核心逻辑 >= 90%）
- [ ] 8.9. 集成测试通过（使用真实数据库和 Redis）
- [ ] 8.10. E2E 测试通过（完整流程）
- [ ] 8.11. 所有测试独立运行

### 文档
- [ ] 8.12. 代码注释完整（包、类型、函数）
- [ ] 8.13. README 更新（如需要）
- [ ] 8.14. 设计文档更新（如需要）

### 三方组件验证 ⚠️ **关键**
- [ ] 8.15. PostgreSQL Repository: 集成测试通过，与真实数据库调通
- [ ] 8.16. Redis Cache: 集成测试通过，与真实 Redis 调通
- [ ] 8.17. 只有验证通过后，才能继续开发依赖该组件的功能

---

## 任务依赖关系

```
Phase 1 (基础设施) 
  ↓
Phase 2 (Domain Layer)
  ↓
Phase 3 (Infrastructure Layer) ⚠️ 必须调通
  ↓
Phase 4 (Application Layer)
  ↓
Phase 5 (Presentation Layer)
  ↓
Phase 6 (Frontend)
  ↓
Phase 7 (部署和文档)
```

**重要提醒**：
- Phase 3 完成后，必须验证 PostgreSQL 和 Redis 集成测试全部通过
- 只有 Phase 3 验证通过后，才能开始 Phase 4
- 每个任务完成后，更新 tasks.md 中的状态：`[ ]` → `[-]` (进行中) → `[x]` (完成)

