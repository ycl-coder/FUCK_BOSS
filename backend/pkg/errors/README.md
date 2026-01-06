# errors - 错误处理包

统一的错误处理机制，提供错误码、错误类型和错误包装功能。

## 功能

- 定义统一的错误码类型
- 提供结构化的错误类型（AppError）
- 支持错误包装（Error Wrapping，Go 1.13+）
- 错误码和错误消息映射

## 使用示例

```go
import "fuck_boss/backend/pkg/errors"

// 创建错误
err := errors.NewValidationError("company name is required")

// 包装错误
err := fmt.Errorf("failed to save post: %w", originalErr)

// 检查错误类型
if errors.IsValidationError(err) {
    // 处理验证错误
}
```

## 错误码

- `VALIDATION_ERROR` - 验证错误
- `NOT_FOUND` - 资源未找到
- `RATE_LIMIT_EXCEEDED` - 限流错误
- `INTERNAL_ERROR` - 内部错误
- `DATABASE_ERROR` - 数据库错误

