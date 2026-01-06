# logger - 结构化日志

基于 `go.uber.org/zap` 的结构化日志组件。

## 功能

- 结构化日志记录
- 不同日志级别（Debug, Info, Warn, Error）
- 开发环境（Console）和生产环境（JSON）格式
- **准确的调用位置信息**：日志显示真实的代码调用位置，而不是包装层位置
- Context 感知日志（支持 Request ID、Trace ID、User ID）
- 字段追加功能

## 使用示例

### 基本使用

```go
import "fuck_boss/backend/internal/infrastructure/logger"

// 创建开发环境日志
logger, err := logger.NewLogger("development", nil)
if err != nil {
    log.Fatal(err)
}
defer logger.Sync()

// 记录日志
logger.Info("server started", zap.String("port", "50051"))
logger.Error("failed to connect", zap.Error(err))
```

### 使用配置

```go
import (
    "fuck_boss/backend/internal/infrastructure/config"
    "fuck_boss/backend/internal/infrastructure/logger"
)

// 从配置加载
cfg, err := config.LoadConfig("config.yaml")
if err != nil {
    log.Fatal(err)
}

// 创建日志器
logConfig := &logger.LogConfig{
    Level:            cfg.Log.Level,
    Format:           cfg.Log.Format,
    OutputPaths:       cfg.Log.OutputPaths,
    ErrorOutputPaths: cfg.Log.ErrorOutputPaths,
}

logger, err := logger.NewLoggerFromConfig(logConfig)
if err != nil {
    log.Fatal(err)
}
defer logger.Sync()
```

### Context 感知日志

```go
import (
    "context"
    "fuck_boss/backend/internal/infrastructure/logger"
    "go.uber.org/zap"
)

// 添加 Request ID 到 context
ctx := logger.WithRequestID(context.Background(), "req-123")

// 使用 context 记录日志
logger.WithContext(ctx).Info("request processed", zap.String("method", "POST"))
```

### 追加字段

```go
// 创建带字段的日志器
loggerWithFields := logger.WithFields(
    zap.String("service", "content-service"),
    zap.String("version", "1.0.0"),
)

loggerWithFields.Info("service started")
```

## 日志级别

- **Debug**: 调试信息，仅在开发环境使用
- **Info**: 一般信息，记录正常操作
- **Warn**: 警告信息，记录潜在问题
- **Error**: 错误信息，记录错误和异常

## 日志格式

### 开发环境（Console）

- 彩色输出，易于阅读
- 适合本地开发和调试
- 格式：`2024-01-01T12:00:00.000Z	INFO	message	{"key": "value"}`

### 生产环境（JSON）

- 结构化 JSON 格式
- 便于日志分析和聚合
- 格式：`{"level":"info","ts":1704110400.0,"msg":"message","key":"value"}`

## Context 支持

日志组件支持从 context 中提取以下信息：

- **Request ID**: 请求唯一标识
- **Trace ID**: 分布式追踪 ID
- **User ID**: 用户标识（未来版本）

### Context 辅助函数

```go
// 添加 Request ID
ctx := logger.WithRequestID(ctx, "req-123")

// 添加 Trace ID
ctx = logger.WithTraceID(ctx, "trace-456")

// 添加 User ID
ctx = logger.WithUserID(ctx, "user-789")

// 使用 context 记录日志
logger.WithContext(ctx).Info("operation completed")
```

## 配置选项

### LogConfig

```go
type LogConfig struct {
    Level            string   // 日志级别: debug, info, warn, error
    Format           string   // 日志格式: json, text, console
    OutputPaths      []string // 日志输出路径: ["stdout"], ["/var/log/app.log"]
    ErrorOutputPaths []string // 错误日志输出路径: ["stderr"], ["/var/log/error.log"]
}
```

## 最佳实践

1. **始终调用 Sync()**: 在程序退出前调用 `logger.Sync()` 确保所有日志都被刷新
2. **使用结构化字段**: 使用 `zap.Field` 而不是字符串拼接
3. **合理使用日志级别**: 
   - Debug: 详细的调试信息
   - Info: 重要的业务事件
   - Warn: 需要关注但不影响功能的问题
   - Error: 错误和异常
4. **Context 传递**: 在请求处理链中传递 context，自动记录 Request ID
5. **生产环境使用 JSON**: 生产环境使用 JSON 格式便于日志分析

## 示例

### 完整的服务启动示例

```go
package main

import (
    "fuck_boss/backend/internal/infrastructure/config"
    "fuck_boss/backend/internal/infrastructure/logger"
    "go.uber.org/zap"
)

func main() {
    // 加载配置
    cfg, err := config.LoadConfig("config.yaml")
    if err != nil {
        panic(err)
    }

    // 创建日志器
    logConfig := &logger.LogConfig{
        Level:            cfg.Log.Level,
        Format:           cfg.Log.Format,
        OutputPaths:      cfg.Log.OutputPaths,
        ErrorOutputPaths: cfg.Log.ErrorOutputPaths,
    }

    log, err := logger.NewLoggerFromConfig(logConfig)
    if err != nil {
        panic(err)
    }
    defer log.Sync()

    // 记录服务启动
    log.Info("service started",
        zap.String("service", "fuck_boss"),
        zap.String("version", "1.0.0"),
        zap.Int("grpc_port", cfg.GRPC.Port),
    )
}
```

### gRPC 中间件集成示例

```go
func LoggingInterceptor(log logger.Logger) grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        // 生成 Request ID
        requestID := generateRequestID()
        ctx = logger.WithRequestID(ctx, requestID)

        // 记录请求开始
        log := logger.WithContext(ctx)
        log.Info("request started",
            zap.String("method", info.FullMethod),
        )

        // 处理请求
        start := time.Now()
        resp, err := handler(ctx, req)
        duration := time.Since(start)

        // 记录请求完成
        if err != nil {
            log.Error("request failed",
                zap.Error(err),
                zap.Duration("duration", duration),
            )
        } else {
            log.Info("request completed",
                zap.Duration("duration", duration),
            )
        }

        return resp, err
    }
}
```
