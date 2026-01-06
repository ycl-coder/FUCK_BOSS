# validator - 验证工具包

基于 `github.com/go-playground/validator/v10` 的统一验证工具。

## 功能

- 统一的验证规则定义
- 自定义验证器
- 验证错误消息格式化

## 使用示例

```go
import "fuck_boss/backend/pkg/validator"

// 验证结构体
type CreatePostRequest struct {
    Company  string `validate:"required,min=1,max=100"`
    CityCode string `validate:"required"`
    Content  string `validate:"required,min=10,max=5000"`
}

err := validator.Validate(req)
if err != nil {
    // 处理验证错误
}
```

## 验证规则

- `required` - 必填
- `min=n` - 最小长度
- `max=n` - 最大长度
- `email` - 邮箱格式（如需要）

