# middleware - gRPC 中间件

gRPC 拦截器（Interceptors），提供日志记录和错误恢复功能。

## 结构

- **logging.go** - 日志拦截器
- **recovery.go** - 恢复拦截器

## LoggingInterceptor

记录所有 gRPC 请求和响应，包括：

- 请求方法名（FullMethod）
- 客户端 IP 地址
- 请求元数据（metadata）
- 响应状态码
- 处理时长
- 错误信息（如果有）

### 功能特性

- **自动生成 Request ID**: 为每个请求生成唯一的 Request ID，并添加到 context 中
- **客户端 IP 提取**: 从 metadata 中提取客户端 IP（支持 X-Forwarded-For 和 X-Real-IP）
- **状态码提取**: 从 gRPC 错误中提取状态码
- **结构化日志**: 使用 zap 进行结构化日志记录

### 使用示例

```go
import (
    "fuck_boss/backend/internal/infrastructure/logger"
    "fuck_boss/backend/internal/presentation/middleware"
)

// 创建日志器
log, err := logger.NewLogger("production", nil)
if err != nil {
    panic(err)
}

// 创建 gRPC 服务器，使用日志拦截器
server := grpc.NewServer(
    grpc.UnaryInterceptor(middleware.LoggingInterceptor(log)),
)
```

## RecoveryInterceptor

捕获 panic 并恢复，记录堆栈信息。

### 功能特性

- **Panic 恢复**: 自动捕获和处理 panic
- **堆栈跟踪**: 记录完整的堆栈信息
- **错误转换**: 将 panic 转换为 gRPC Internal 错误
- **结构化日志**: 使用 zap 记录 panic 详情

### 使用示例

```go
import (
    "fuck_boss/backend/internal/infrastructure/logger"
    "fuck_boss/backend/internal/presentation/middleware"
)

// 创建日志器
log, err := logger.NewLogger("production", nil)
if err != nil {
    panic(err)
}

// 创建 gRPC 服务器，使用恢复拦截器
server := grpc.NewServer(
    grpc.UnaryInterceptor(middleware.RecoveryInterceptor(log)),
)
```

## 组合使用

中间件可以链式组合，建议的顺序是：

1. **RecoveryInterceptor** - 最外层，捕获所有 panic
2. **LoggingInterceptor** - 记录所有请求

```go
import (
    "fuck_boss/backend/internal/infrastructure/logger"
    "fuck_boss/backend/internal/presentation/middleware"
    "google.golang.org/grpc"
)

// 创建日志器
log, err := logger.NewLogger("production", nil)
if err != nil {
    panic(err)
}

// 创建 gRPC 服务器，链式组合中间件
server := grpc.NewServer(
    grpc.ChainUnaryInterceptor(
        middleware.RecoveryInterceptor(log),  // 最外层：恢复
        middleware.LoggingInterceptor(log),   // 内层：日志
    ),
)
```

## 注意事项

- **Recovery 必须在最外层**: 确保所有 panic 都能被捕获
- **Logging 应该在 Recovery 之后**: 这样即使发生 panic，也能记录日志
- **使用结构化日志**: 所有日志都使用 zap 的结构化字段
- **Request ID 自动传递**: LoggingInterceptor 会自动生成 Request ID 并添加到 context，后续的日志会自动包含 Request ID

## 日志格式示例

### 请求开始日志

```json
{
  "level": "info",
  "ts": 1704110400.0,
  "msg": "gRPC request started",
  "request_id": "20240106120000-123456-0010",
  "method": "/content.v1.ContentService/CreatePost",
  "client_ip": "192.168.1.100",
  "metadata": {...}
}
```

### 请求完成日志

```json
{
  "level": "info",
  "ts": 1704110401.0,
  "msg": "gRPC request completed",
  "request_id": "20240106120000-123456-0010",
  "method": "/content.v1.ContentService/CreatePost",
  "status_code": "OK",
  "duration": 0.123
}
```

### Panic 恢复日志

```json
{
  "level": "error",
  "ts": 1704110401.0,
  "msg": "gRPC panic recovered",
  "request_id": "20240106120000-123456-0010",
  "method": "/content.v1.ContentService/CreatePost",
  "panic": "runtime error: invalid memory address",
  "stack": "goroutine 1 [running]:\n..."
}
```

