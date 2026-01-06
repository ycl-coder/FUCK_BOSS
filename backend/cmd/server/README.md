# server - gRPC 服务器入口

应用程序的主入口，负责初始化所有依赖并启动 gRPC 服务器。

## 功能

- 初始化配置（从配置文件或环境变量）
- 初始化结构化日志（zap）
- 连接数据库（PostgreSQL）
- 连接缓存（Redis）
- 运行数据库迁移
- 初始化 Repository
- 初始化 Use Cases
- 初始化 gRPC 服务
- 注册中间件（日志、恢复）
- 启动 gRPC 服务器
- 优雅关闭（Graceful Shutdown）

## 启动流程

1. **加载配置**: 从配置文件（`config/config.yaml`）或环境变量加载配置
2. **初始化日志**: 创建结构化日志器（支持 JSON/Console 格式）
3. **连接数据库**: 连接 PostgreSQL，设置连接池参数
4. **连接 Redis**: 连接 Redis，设置连接池参数
5. **运行迁移**: 检查并创建数据库表（如果不存在）
6. **创建 Repository**: 初始化 PostgreSQL Repository 和 Redis Cache/Rate Limiter
7. **创建 Use Case**: 初始化所有 Use Cases（CreatePost, ListPosts, GetPost, SearchPosts）
8. **创建 gRPC 服务**: 创建 ContentService 实例
9. **注册中间件**: 注册日志和恢复中间件
10. **启动服务器**: 在指定端口启动 gRPC 服务器
11. **等待信号**: 监听 SIGTERM/SIGINT 信号
12. **优雅关闭**: 停止接受新连接，关闭数据库和 Redis 连接

## 运行

### 开发环境

**重要**: 首次运行前，需要配置数据库连接信息。

#### 1. 创建配置文件

```bash
# 复制示例配置文件
cp config/config.example.yaml config/config.yaml

# 编辑配置文件，修改数据库用户名和密码
# 或者使用环境变量
```

#### 2. 配置数据库连接

**方式一：使用配置文件**

编辑 `config/config.yaml`，修改数据库连接信息：

```yaml
database:
  host: localhost
  port: 5432
  user: your_username  # 修改为你的 PostgreSQL 用户名
  password: your_password  # 修改为你的 PostgreSQL 密码
  dbname: fuck_boss
```

**方式二：使用环境变量**

```bash
export FUCK_BOSS_DATABASE_USER=your_username
export FUCK_BOSS_DATABASE_PASSWORD=your_password
export FUCK_BOSS_DATABASE_DBNAME=fuck_boss
```

#### 3. 启动服务器

```bash
# 使用配置文件
go run cmd/server/main.go

# 或指定配置文件路径
CONFIG_PATH=config/config.yaml go run cmd/server/main.go

# 或使用环境变量（无需配置文件）
FUCK_BOSS_DATABASE_USER=your_user \
FUCK_BOSS_DATABASE_PASSWORD=your_pass \
go run cmd/server/main.go
```

### 生产环境

```bash
# 构建
go build -o bin/server ./cmd/server

# 运行
./bin/server

# 或指定配置文件
CONFIG_PATH=/path/to/config.yaml ./bin/server
```

### 使用 Makefile

```bash
# 构建
make backend-build

# 运行（需要先配置数据库和 Redis）
./backend/bin/server
```

## 配置

### 配置文件

默认配置文件路径：`config/config.yaml`

可以通过环境变量 `CONFIG_PATH` 指定配置文件路径。

### 环境变量

所有配置项都可以通过环境变量覆盖，格式：`FUCK_BOSS_<SECTION>_<FIELD>`

示例：
- `FUCK_BOSS_DATABASE_HOST=localhost`
- `FUCK_BOSS_DATABASE_PORT=5432`
- `FUCK_BOSS_REDIS_HOST=localhost`
- `FUCK_BOSS_REDIS_PORT=6379`
- `FUCK_BOSS_GRPC_PORT=50051`
- `FUCK_BOSS_LOG_LEVEL=info`

### 配置项

#### 数据库配置

- `database.host`: 数据库主机（默认: localhost）
- `database.port`: 数据库端口（默认: 5432）
- `database.user`: 数据库用户名（默认: postgres）
- `database.password`: 数据库密码（必需）
- `database.dbname`: 数据库名称（默认: fuck_boss）
- `database.sslmode`: SSL 模式（默认: disable）
- `database.max_open_conns`: 最大打开连接数（默认: 100）
- `database.max_idle_conns`: 最大空闲连接数（默认: 10）
- `database.conn_max_lifetime`: 连接最大生存时间（秒，默认: 3600）

#### Redis 配置

- `redis.host`: Redis 主机（默认: localhost）
- `redis.port`: Redis 端口（默认: 6379）
- `redis.password`: Redis 密码（可选）
- `redis.db`: Redis 数据库编号（默认: 0）
- `redis.max_retries`: 最大重试次数（默认: 3）
- `redis.pool_size`: 连接池大小（默认: 50）
- `redis.min_idle_conns`: 最小空闲连接数（默认: 5）

#### gRPC 配置

- `grpc.port`: gRPC 服务器端口（默认: 50051）
- `grpc.max_recv_msg_size`: 最大接收消息大小（字节，默认: 4MB）
- `grpc.max_send_msg_size`: 最大发送消息大小（字节，默认: 4MB）

#### 日志配置

- `log.level`: 日志级别（debug, info, warn, error，默认: info）
- `log.format`: 日志格式（json, text, console，默认: json）
- `log.output_paths`: 日志输出路径（默认: ["stdout"]）
- `log.error_output_paths`: 错误日志输出路径（默认: ["stderr"]）

## 数据库迁移

服务器启动时会自动检查数据库表是否存在。如果表不存在，会自动创建：

- `cities` 表（城市配置）
- `posts` 表（曝光内容）

**注意**: 生产环境建议使用专门的迁移工具（如 `golang-migrate`）管理数据库迁移。

## 优雅关闭

服务器支持优雅关闭：

1. 收到 SIGTERM 或 SIGINT 信号后，停止接受新连接
2. 等待现有请求完成（最多 30 秒）
3. 关闭数据库连接
4. 关闭 Redis 连接
5. 退出程序

## gRPC Reflection

服务器启用了 gRPC Reflection，可以使用以下工具进行测试：

- **grpcurl**: 命令行工具
  ```bash
  grpcurl -plaintext localhost:50051 list
  ```

- **grpcui**: Web UI 工具
  ```bash
  grpcui -plaintext localhost:50051
  ```

## 日志输出

### 开发环境

建议使用 `console` 格式，输出到 `stdout`：

```yaml
log:
  level: debug
  format: console
  output_paths:
    - stdout
```

### 生产环境

建议使用 `json` 格式，输出到文件：

```yaml
log:
  level: info
  format: json
  output_paths:
    - /var/log/fuck_boss/app.log
  error_output_paths:
    - /var/log/fuck_boss/error.log
```

## 故障排查

### 数据库连接失败

1. 检查数据库是否运行
2. 检查连接配置（host, port, user, password, dbname）
3. 检查网络连接
4. 检查防火墙设置

### Redis 连接失败

1. 检查 Redis 是否运行
2. 检查连接配置（host, port, password）
3. 检查网络连接
4. 检查 Redis 认证配置

### 端口被占用

1. 检查端口是否被其他进程占用
2. 修改 `grpc.port` 配置使用其他端口

### 迁移失败

1. 检查数据库用户权限
2. 检查数据库是否已存在表
3. 查看日志获取详细错误信息

## 相关文档

- [配置管理](../../internal/infrastructure/config/README.md)
- [日志组件](../../internal/infrastructure/logger/README.md)
- [PostgreSQL Repository](../../internal/infrastructure/persistence/postgres/README.md)
- [Redis Cache](../../internal/infrastructure/persistence/redis/README.md)
- [gRPC Handler](../../internal/presentation/grpc/README.md)
- [Middleware](../../internal/presentation/middleware/README.md)

