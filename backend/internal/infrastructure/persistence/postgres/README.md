# postgres - PostgreSQL 持久化实现

PostgreSQL Repository 的实现，实现 Domain Layer 定义的 Repository 接口。

## 结构

- **post_repository.go** - PostRepository 的 PostgreSQL 实现
- **migrations/** - 数据库迁移脚本

## 实现

### PostRepository

实现 `domain.PostRepository` 接口，使用 PostgreSQL 存储数据。

#### 创建 Repository

```go
import (
    "database/sql"
    "fuck_boss/backend/internal/infrastructure/persistence/postgres"
)

// db 是 *sql.DB 实例
repo := postgres.NewPostRepository(db)
```

#### 使用示例

```go
// 保存 Post
err := repo.Save(ctx, post)
if err != nil {
    return err
}

// 根据 ID 查找
post, err := repo.FindByID(ctx, postID)
if err != nil {
    return err
}

// 根据城市查找（分页）
posts, total, err := repo.FindByCity(ctx, city, page, pageSize)
if err != nil {
    return err
}

// 搜索（支持城市过滤和分页）
posts, total, err := repo.Search(ctx, keyword, &city, page, pageSize)
if err != nil {
    return err
}
```

#### 方法说明

- **Save**: 保存或更新 Post（使用 `ON CONFLICT` 实现 upsert）
- **FindByID**: 根据 ID 查找单个 Post
- **FindByCity**: 根据城市查找 Posts，支持分页，按创建时间倒序
- **Search**: 全文搜索，支持可选的城市过滤和分页

#### 全文搜索

使用 PostgreSQL 的全文搜索功能：
- 使用 `to_tsvector('simple', ...)` 创建搜索向量
- 使用 `plainto_tsquery('simple', ...)` 处理搜索关键词
- 搜索范围：`company_name` 和 `content` 字段

## 数据库 Schema

### posts 表

```sql
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_name VARCHAR(100) NOT NULL,
    city_code VARCHAR(50) NOT NULL,
    city_name VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    occurred_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

**字段说明**:
- `id` - UUID 主键，自动生成
- `company_name` - 公司名称（VARCHAR(100)，对应 CompanyName 值对象）
- `city_code` - 城市代码（VARCHAR(50)，对应 City.Code）
- `city_name` - 城市名称（VARCHAR(50)，对应 City.Name）
- `content` - 内容（TEXT，对应 Content 值对象）
- `occurred_at` - 发生时间（TIMESTAMP，可选，未来版本使用）
- `created_at` - 创建时间（TIMESTAMP，自动设置）
- `updated_at` - 更新时间（TIMESTAMP，自动设置）

### cities 表

```sql
CREATE TABLE cities (
    code VARCHAR(50) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    pinyin VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

**字段说明**:
- `code` - 城市代码（主键）
- `name` - 城市名称
- `pinyin` - 城市拼音（可选，用于搜索）
- `created_at` - 创建时间

### 索引

#### posts 表索引

- `idx_posts_city_code` - 城市代码索引（用于 FindByCity 查询）
- `idx_posts_created_at` - 创建时间索引（倒序，用于按时间排序）
- `idx_posts_company_name` - 公司名称索引（用于筛选和搜索）
- `idx_posts_search` - 全文搜索索引（GIN，用于全文搜索）

**全文搜索索引说明**:
- 当前使用 PostgreSQL 内置的 `simple` 配置
- 搜索范围：`company_name` 和 `content`
- 未来可以升级为 `pg_jieba` 扩展以获得更好的中文分词支持

#### cities 表索引

- `idx_cities_name` - 城市名称索引（用于搜索城市）

## 迁移

使用 `golang-migrate/migrate` 管理数据库迁移。

```bash
# 运行迁移
migrate -path migrations -database "postgres://..." up

# 回滚迁移
migrate -path migrations -database "postgres://..." down
```

## 技术细节

### 错误处理

- 使用统一的错误处理包 (`pkg/errors`)
- 数据库错误包装为 `DATABASE_ERROR`
- 未找到资源返回 `NOT_FOUND` 错误

### 查询优化

- 所有查询使用参数化查询，防止 SQL 注入
- 使用索引优化查询性能
- 分页查询使用 `LIMIT` 和 `OFFSET`
- 全文搜索使用 GIN 索引

### 数据转换

- 从数据库读取时，重建 Domain 层的值对象和实体
- 使用 `NewPostFromDB` 方法从数据库数据重建 Post 实体
- 所有值对象在重建时进行验证

## 注意事项

- 使用参数化查询，防止 SQL 注入
- 所有方法必须使用 context.Context
- 支持事务处理（通过 context 传递事务）
- 必须通过集成测试验证
- 分页参数：page 从 1 开始，pageSize 最小为 1

