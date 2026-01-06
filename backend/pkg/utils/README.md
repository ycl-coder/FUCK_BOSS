# utils - 工具函数包

通用的工具函数集合。

## 功能

- 字符串处理工具
- 时间处理工具
- 其他通用工具函数

## 使用示例

```go
import "fuck_boss/backend/pkg/utils"

// 字符串工具
str := utils.Truncate("long string", 100)

// 时间工具
relativeTime := utils.RelativeTime(timestamp)
```

## 注意事项

- 工具函数应该是纯函数（无副作用）
- 避免在 utils 中放置业务逻辑
- 保持函数职责单一

