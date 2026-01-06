# config - 配置管理

基于 `github.com/spf13/viper` 的配置管理组件。

## 功能

- 从配置文件（YAML）读取配置
- 支持环境变量覆盖
- 配置验证和默认值
- 提供合理的默认配置

## 配置结构

```go
type Config struct {
    Database DatabaseConfig  // PostgreSQL 数据库配置
    Redis    RedisConfig     // Redis 缓存配置
    GRPC     GRPCConfig      // gRPC 服务器配置
    Log      LogConfig       // 日志配置
}
```

### DatabaseConfig

- `host`: 数据库主机地址（默认: localhost）
- `port`: 数据库端口（默认: 5432）
- `user`: 数据库用户名（默认: postgres）
- `password`: 数据库密码（默认: 空）
- `dbname`: 数据库名称（默认: fuck_boss）
- `sslmode`: SSL 模式（默认: disable）
- `max_open_conns`: 最大打开连接数（默认: 100）
- `max_idle_conns`: 最大空闲连接数（默认: 10）
- `conn_max_lifetime`: 连接最大生存时间（秒，默认: 3600）

### RedisConfig

- `host`: Redis 主机地址（默认: localhost）
- `port`: Redis 端口（默认: 6379）
- `password`: Redis 密码（默认: 空）
- `db`: Redis 数据库编号 0-15（默认: 0）
- `max_retries`: 最大重试次数（默认: 3）
- `pool_size`: 连接池大小（默认: 50）
- `min_idle_conns`: 最小空闲连接数（默认: 5）

### GRPCConfig

- `port`: gRPC 服务器端口（默认: 50051）
- `max_recv_msg_size`: 最大接收消息大小（字节，默认: 4MB）
- `max_send_msg_size`: 最大发送消息大小（字节，默认: 4MB）

### LogConfig

- `level`: 日志级别 debug/info/warn/error（默认: info）
- `format`: 日志格式 json/text（默认: json）
- `output_paths`: 日志输出路径列表（默认: ["stdout"]）
- `error_output_paths`: 错误日志输出路径列表（默认: ["stderr"]）

## 使用示例

```go
import "fuck_boss/backend/internal/infrastructure/config"

// 从指定文件加载配置
cfg, err := config.LoadConfig("config.yaml")
if err != nil {
    log.Fatal(err)
}

// 使用默认配置（不指定文件路径）
cfg, err := config.LoadConfig("")
if err != nil {
    log.Fatal(err)
}

// 使用配置
dsn := cfg.Database.GetDSN()
redisAddr := cfg.Redis.GetAddr()
```

## 配置文件示例

参考 `backend/config/config.example.yaml` 文件。

## 环境变量

所有配置项都可以通过环境变量覆盖，格式：`FUCK_BOSS_<SECTION>_<FIELD>`

示例：
- `FUCK_BOSS_DATABASE_HOST` - 覆盖 database.host
- `FUCK_BOSS_DATABASE_PORT` - 覆盖 database.port
- `FUCK_BOSS_REDIS_HOST` - 覆盖 redis.host
- `FUCK_BOSS_GRPC_PORT` - 覆盖 grpc.port
- `FUCK_BOSS_LOG_LEVEL` - 覆盖 log.level

环境变量的优先级高于配置文件。

## 配置验证

配置加载时会自动验证：
- 必填字段检查（database.host, database.user, database.dbname 等）
- 端口范围检查（1-65535）
- Redis DB 范围检查（0-15）
- 日志级别和格式验证
- 连接池参数验证

## 辅助方法

### DatabaseConfig.GetDSN()

返回 PostgreSQL 连接字符串（DSN）。

```go
dsn := cfg.Database.GetDSN()
// 输出: "host=localhost port=5432 user=postgres password=password dbname=fuck_boss sslmode=disable"
```

### RedisConfig.GetAddr()

返回 Redis 地址字符串。

```go
addr := cfg.Redis.GetAddr()
// 输出: "localhost:6379"
```
