# 集成测试执行流程详解

本文档详细解释 PostgreSQL Repository 集成测试的完整执行过程。

## 测试架构

### 测试框架
- **testify/suite**: 提供测试套件（Test Suite）功能
- **testify/require**: 提供断言功能，失败时立即停止测试

### 测试文件结构
```
post_repository_test.go
├── PostRepositoryTestSuite (测试套件)
│   ├── SetupSuite()      - 套件级别初始化（只执行一次）
│   ├── SetupTest()       - 每个测试前执行
│   ├── TestXXX()         - 测试方法
│   ├── TearDownTest()    - 每个测试后执行
│   └── TearDownSuite()   - 套件级别清理（只执行一次）
└── TestPostRepositorySuite() - 测试入口
```

## 完整执行流程

### 阶段 1: 测试启动

#### 1.1 运行测试命令
```bash
make test-integration-repository
# 或
cd backend && go test -v ./test/integration/repository/...
```

**执行过程**:
1. Go 测试框架扫描 `test/integration/repository/` 目录
2. 找到 `post_repository_test.go` 文件
3. 发现 `TestPostRepositorySuite` 函数（测试入口）
4. 创建测试套件实例 `PostRepositoryTestSuite`

#### 1.2 调用测试入口
```go
func TestPostRepositorySuite(t *testing.T) {
    suite.Run(t, new(PostRepositoryTestSuite))
}
```

**执行过程**:
- `suite.Run()` 创建新的测试套件实例
- 开始执行测试套件的生命周期方法

---

### 阶段 2: 套件级别初始化 (SetupSuite)

#### 2.1 数据库连接配置
```go
func (s *PostRepositoryTestSuite) SetupSuite() {
    // 1. 获取数据库连接字符串
    dsn := os.Getenv("TEST_DATABASE_URL")
    if dsn == "" {
        dsn = "postgres://test_user:test_password@localhost:5433/test_db?sslmode=disable"
    }
```

**执行过程**:
- 检查环境变量 `TEST_DATABASE_URL`
- 如果未设置，使用默认连接字符串
- 默认连接: `localhost:5433` (测试数据库端口)

**数据流**:
```
环境变量 TEST_DATABASE_URL
    ↓ (未设置)
默认值: postgres://test_user:test_password@localhost:5433/test_db?sslmode=disable
```

#### 2.2 建立数据库连接
```go
    // 2. 连接到数据库
    s.db, err = sql.Open("postgres", dsn)
    require.NoError(s.T(), err, "Failed to connect to test database")
```

**执行过程**:
- `sql.Open()` 创建数据库连接对象（**不立即连接**）
- 返回 `*sql.DB` 对象，存储在 `s.db`
- 连接池会在第一次查询时建立实际连接

**重要**: `sql.Open()` 不会立即验证连接，只是创建连接对象。

#### 2.3 等待数据库就绪
```go
    // 3. 等待数据库就绪
    err = s.waitForDB()
    require.NoError(s.T(), err, "Database is not ready")
```

**执行过程** (`waitForDB` 方法):
```go
func (s *PostRepositoryTestSuite) waitForDB() error {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    for {
        err := s.db.PingContext(ctx)  // 尝试连接数据库
        if err == nil {
            return nil  // 连接成功
        }

        select {
        case <-ctx.Done():
            return fmt.Errorf("database not ready: %w", ctx.Err())
        case <-time.After(1 * time.Second):
            // 等待 1 秒后重试
        }
    }
}
```

**执行过程**:
1. 创建 30 秒超时的 context
2. 循环执行 `PingContext()` 尝试连接数据库
3. 如果连接失败，等待 1 秒后重试
4. 如果 30 秒内连接成功，返回 nil
5. 如果超时，返回错误

**时间线**:
```
t=0s:  第一次 PingContext() - 可能失败（数据库还在启动）
t=1s:  等待 1 秒
t=1s:  第二次 PingContext() - 可能失败
t=2s:  等待 1 秒
...
t=N:   第 N 次 PingContext() - 成功 ✓
```

#### 2.4 运行数据库迁移
```go
    // 4. 运行迁移
    err = s.runMigrations()
    require.NoError(s.T(), err, "Failed to run migrations")
```

**执行过程** (`runMigrations` 方法):
```go
func (s *PostRepositoryTestSuite) runMigrations() error {
    ctx := context.Background()  // 创建新的 context（s.ctx 还未初始化）
    
    migrationSQL := `
        CREATE TABLE IF NOT EXISTS posts (...);
        CREATE INDEX IF NOT EXISTS idx_posts_city_code ON posts(city_code);
        ...
    `
    
    _, err := s.db.ExecContext(ctx, migrationSQL)
    return err
}
```

**执行过程**:
1. 创建新的 context（因为 `s.ctx` 在 SetupSuite 最后才初始化）
2. 执行 SQL 语句创建表和索引
3. 使用 `IF NOT EXISTS` 确保幂等性（多次运行不会报错）

**数据库状态变化**:
```
迁移前: 数据库为空（或只有系统表）
    ↓
执行 CREATE TABLE IF NOT EXISTS posts
    ↓
迁移后: 
  - posts 表已创建
  - 所有索引已创建
  - 全文搜索索引已创建
```

#### 2.5 创建 Repository 实例
```go
    // 5. 创建 Repository
    s.repo = postgres.NewPostRepository(s.db)
```

**执行过程**:
- 创建 `PostRepository` 实例
- 传入数据库连接 `s.db`
- Repository 现在可以使用数据库进行 CRUD 操作

#### 2.6 初始化 Context
```go
    // 6. 创建 context
    s.ctx = context.Background()
}
```

**执行过程**:
- 创建根 context，用于所有测试方法
- 所有数据库操作都会使用这个 context

**SetupSuite 完成后的状态**:
```
s.db   = *sql.DB (已连接)
s.repo = *postgres.PostRepository (已创建)
s.ctx  = context.Background()
```

---

### 阶段 3: 每个测试的执行循环

对于每个测试方法（如 `TestPostRepository_Save`），执行以下流程：

#### 3.1 测试前准备 (SetupTest)
```go
func (s *PostRepositoryTestSuite) SetupTest() {
    // 清理任何现有的测试数据
    _, err := s.db.ExecContext(s.ctx, "TRUNCATE TABLE posts CASCADE")
    if err != nil {
        s.T().Logf("Failed to truncate posts table: %v", err)
    }
}
```

**执行过程**:
1. 在每个测试方法执行**之前**自动调用
2. 执行 `TRUNCATE TABLE posts CASCADE` 清空所有数据
3. 确保每个测试从干净的状态开始

**数据库状态**:
```
SetupTest 前: 可能有之前测试的数据
    ↓
执行 TRUNCATE TABLE posts CASCADE
    ↓
SetupTest 后: posts 表为空（但表结构保留）
```

**为什么使用 TRUNCATE 而不是 DELETE**:
- `TRUNCATE` 更快（不记录逐行删除）
- `TRUNCATE` 重置自增序列
- `TRUNCATE` 释放空间

#### 3.2 执行测试方法
```go
func (s *PostRepositoryTestSuite) TestPostRepository_Save() {
    // 1. 创建值对象
    company, err := content.NewCompanyName("测试公司")
    s.Require().NoError(err)
    
    city, err := shared.NewCity("beijing", "北京")
    s.Require().NoError(err)
    
    postContent, err := content.NewContent("这是一条测试内容...")
    s.Require().NoError(err)
    
    // 2. 创建 Post 实体
    post, err := content.NewPost(company, city, postContent)
    s.Require().NoError(err)
    
    // 3. 保存到数据库
    err = s.repo.Save(s.ctx, post)
    s.Require().NoError(err)
    
    // 4. 验证保存成功
    found, err := s.repo.FindByID(s.ctx, post.ID())
    s.Require().NoError(err)
    s.Equal(post.ID().String(), found.ID().String())
    // ... 更多断言
}
```

**执行过程**:
1. **创建值对象**: 使用工厂方法创建 `CompanyName`, `City`, `Content`
   - 每个值对象都会进行验证（长度、格式等）
2. **创建实体**: 使用 `NewPost` 创建 `Post` 聚合根
   - 自动生成 UUID
   - 设置 `createdAt` 为当前时间
3. **保存到数据库**: 调用 `repo.Save()`
   - 执行 `INSERT INTO posts ... ON CONFLICT DO UPDATE`
   - 数据写入数据库
4. **验证**: 从数据库读取并验证
   - 调用 `repo.FindByID()` 查询
   - 使用断言验证数据正确性

**数据库操作时间线**:
```
t=0ms:  TRUNCATE TABLE posts (SetupTest)
t=10ms: INSERT INTO posts ... (Save)
t=20ms: SELECT ... FROM posts WHERE id = ... (FindByID)
t=30ms: 测试完成
```

#### 3.3 测试后清理 (TearDownTest)
```go
func (s *PostRepositoryTestSuite) TearDownTest() {
    // 清理测试数据
    _, err := s.db.ExecContext(s.ctx, "TRUNCATE TABLE posts CASCADE")
    if err != nil {
        s.T().Logf("Failed to truncate posts table: %v", err)
    }
}
```

**执行过程**:
1. 在每个测试方法执行**之后**自动调用
2. 再次执行 `TRUNCATE TABLE posts CASCADE`
3. 确保测试数据不会影响下一个测试

**为什么测试前后都清理**:
- **SetupTest 清理**: 确保测试从干净状态开始
- **TearDownTest 清理**: 确保测试数据不会泄漏到下一个测试
- **双重保险**: 即使测试失败，TearDownTest 也会执行（defer 机制）

**数据库状态**:
```
测试执行中: posts 表有测试数据
    ↓
TearDownTest 执行 TRUNCATE
    ↓
测试完成后: posts 表为空
```

---

### 阶段 4: 所有测试完成后的清理 (TearDownSuite)

```go
func (s *PostRepositoryTestSuite) TearDownSuite() {
    if s.db != nil {
        s.db.Close()
    }
}
```

**执行过程**:
1. 所有测试方法执行完毕后调用
2. 关闭数据库连接
3. 释放资源

**资源清理**:
```
s.db.Close()
    ↓
关闭所有数据库连接
    ↓
释放连接池资源
```

---

## 完整执行时间线示例

假设运行 3 个测试方法：

```
时间    阶段                    操作
─────────────────────────────────────────────────
t=0s    SetupSuite             连接数据库
t=1s    SetupSuite             等待数据库就绪
t=2s    SetupSuite             运行迁移
t=3s    SetupSuite             创建 Repository
t=4s    SetupSuite             完成 ✓

t=5s    SetupTest (Test1)      TRUNCATE posts
t=6s    TestPostRepository_Save  执行测试
t=7s    TearDownTest (Test1)   TRUNCATE posts

t=8s    SetupTest (Test2)      TRUNCATE posts
t=9s    TestPostRepository_FindByID  执行测试
t=10s   TearDownTest (Test2)   TRUNCATE posts

t=11s   SetupTest (Test3)      TRUNCATE posts
t=12s   TestPostRepository_Search  执行测试
t=13s   TearDownTest (Test3)   TRUNCATE posts

t=14s   TearDownSuite          关闭数据库连接
t=15s   测试完成 ✓
```

---

## 数据生命周期

### 正常测试流程
```
SetupSuite
  ↓
  [数据库连接建立]
  ↓
  [表结构创建]
  ↓
循环: 对每个测试方法
  ↓
  SetupTest
    ↓
    [TRUNCATE posts] ← 清空数据
    ↓
  TestXXX()
    ↓
    [INSERT posts]   ← 创建测试数据
    ↓
    [SELECT posts]   ← 验证数据
    ↓
  TearDownTest
    ↓
    [TRUNCATE posts] ← 清空数据
    ↓
TearDownSuite
  ↓
  [关闭连接]
```

### 数据状态变化
```
时间点           posts 表状态
─────────────────────────────────
SetupSuite 后    表存在，但为空
SetupTest 后      空（刚清理）
TestXXX 执行中    有测试数据
TearDownTest 后   空（已清理）
下一个 SetupTest  空（再次清理）
```

---

## 关键设计决策

### 1. 为什么使用 TRUNCATE 而不是 DELETE？
- **性能**: TRUNCATE 更快，不记录逐行删除
- **重置**: 自动重置序列和索引
- **原子性**: 一次性清空所有数据

### 2. 为什么测试前后都清理？
- **隔离性**: 确保测试之间不相互影响
- **可重复性**: 每次运行测试结果一致
- **可靠性**: 即使测试失败，数据也会被清理

### 3. 为什么使用 tmpfs（内存存储）？
- **速度**: 内存操作比磁盘快
- **隔离**: 容器停止后数据自动消失
- **干净**: 每次启动都是全新环境

### 4. 为什么使用 testify/suite？
- **组织**: 将相关测试组织在一起
- **共享**: SetupSuite 中的资源可以被所有测试共享
- **生命周期**: 清晰的测试生命周期管理

---

## 常见问题

### Q: 为什么测试后看不到数据？
**A**: 因为 `TearDownTest` 在每个测试后执行 `TRUNCATE`，数据被立即清理。

**解决方案**: 
- 使用调试测试: `make test-integration-debug`
- 临时注释掉 `TearDownTest` 中的清理代码

### Q: 测试失败时数据会被清理吗？
**A**: 会的。`TearDownTest` 使用 defer 机制，即使测试失败也会执行。

### Q: 如何查看测试执行过程中的数据？
**A**: 
1. 在测试中添加 `time.Sleep()` 暂停执行
2. 在暂停期间连接数据库查看
3. 使用调试测试（不清理数据）

### Q: 多个测试并发运行会冲突吗？
**A**: 不会。每个测试前后都清理数据，确保隔离。但建议顺序运行以确保稳定性。

---

## 总结

集成测试的执行流程遵循严格的**生命周期管理**:

1. **一次性初始化** (SetupSuite): 连接数据库、创建表结构
2. **每个测试的隔离** (SetupTest → Test → TearDownTest): 确保测试独立性
3. **资源清理** (TearDownSuite): 释放数据库连接

这种设计确保了：
- ✅ 测试之间的完全隔离
- ✅ 可重复的测试结果
- ✅ 清晰的资源管理
- ✅ 易于调试和维护

