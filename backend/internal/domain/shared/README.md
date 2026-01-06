# shared - 共享领域概念（Shared Kernel）

跨有界上下文（Bounded Context）共享的值对象和领域概念。

## 什么是共享领域（Shared Kernel）？

在领域驱动设计（DDD）中，**共享领域（Shared Kernel）**是一个特殊的设计模式，用于处理多个有界上下文之间需要共享的通用概念。

### 核心特点

1. **跨上下文共享**: 被多个有界上下文共同使用的领域概念
2. **稳定且通用**: 这些概念相对稳定，不会频繁变化
3. **最小化共享**: 只共享真正需要共享的概念，避免过度耦合
4. **共同维护**: 所有使用它的上下文都需要参与维护和演化

### 为什么需要共享领域？

在我们的项目中，有多个有界上下文：

- **Content Context（内容领域）**: 处理曝光内容的发布、存储
- **Search Context（搜索领域）**: 处理内容搜索、筛选
- **User Context（用户领域）**: 未来可能添加的用户相关功能

这些上下文都需要使用**城市（City）**这个概念：
- Content Context: Post 需要关联城市
- Search Context: 需要按城市筛选内容
- User Context: 未来可能需要用户所在城市

如果每个上下文都自己实现 City，会导致：
- 代码重复
- 概念不一致（不同上下文对城市的理解可能不同）
- 维护困难（修改需要同步多个地方）

因此，我们将 City 放在共享领域（shared）中，让所有上下文共享使用。

## 设计原则

### 1. 最小化共享

**原则**: 只共享真正需要共享的概念，不要过度共享。

**示例**: 
- ✅ City - 多个上下文都需要，应该共享
- ❌ Post - 只在 Content Context 使用，不应该共享
- ❌ SearchQuery - 只在 Search Context 使用，不应该共享

### 2. 保持稳定

**原则**: 共享领域的概念应该相对稳定，避免频繁变化。

**原因**: 共享领域的变更会影响所有使用它的上下文，频繁变更会导致：
- 需要同步更新多个上下文
- 增加系统复杂度
- 提高出错风险

**实践**: 
- 共享领域的概念应该是核心业务概念，不是易变的业务规则
- 如果某个概念经常变化，考虑是否应该放在特定上下文中

### 3. 共同维护

**原则**: 所有使用共享领域的上下文都应该参与维护。

**实践**:
- 修改共享领域需要所有相关团队的 Review
- 变更需要评估对所有使用者的影响
- 保持向后兼容，或提供迁移方案

### 4. 明确边界

**原则**: 明确共享领域与各上下文的边界。

**实践**:
- 共享领域只包含值对象和基础概念
- 不包含业务逻辑（业务逻辑在各上下文中）
- 不包含基础设施代码（基础设施在 Infrastructure Layer）

## 当前实现

### City 值对象

表示城市的值对象，包含城市代码和名称。

#### 使用方式

```go
import "fuck_boss/backend/internal/domain/shared"

// 创建城市
city, err := shared.NewCity("beijing", "北京")
if err != nil {
    return err
}

// 使用
code := city.Code()  // "beijing"
name := city.Name()  // "北京"
```

#### 验证规则

- **code**: 不能为空（trim 后）
- **name**: 不能为空（trim 后）
- 自动去除前后空白字符
- 内部空白字符保留

#### 方法

- `Code()` - 返回城市代码
- `Name()` - 返回城市名称
- `String()` - 返回字符串表示（格式: "City{code: <code>, name: <name>}"）
- `IsZero()` - 检查是否为零值
- `Equals(other City)` - 比较两个 City（比较 code 和 name）

#### 使用场景

- **Content Domain**: Post 关联城市
- **Search Domain**: 按城市筛选
- **其他领域**: 需要城市信息的地方

#### 示例

```go
// 创建城市
beijing, err := shared.NewCity("beijing", "北京")
if err != nil {
    return err
}

// 比较城市
shanghai, _ := shared.NewCity("shanghai", "上海")
if beijing.Equals(shanghai) {
    // false
}

// 检查是否为零值
var city shared.City
if city.IsZero() {
    // true
}
```

## 未来可能扩展的共享概念

### 1. 时间相关概念

如果多个上下文都需要时间相关的概念，可以考虑：

```go
// 示例：时间范围
type TimeRange struct {
    Start time.Time
    End   time.Time
}
```

### 2. 分页概念

如果多个上下文都需要分页，可以考虑：

```go
// 示例：分页参数
type Pagination struct {
    Page     int
    PageSize int
}
```

### 3. 排序概念

如果多个上下文都需要排序，可以考虑：

```go
// 示例：排序参数
type SortOrder struct {
    Field string
    Order string // "asc" or "desc"
}
```

**注意**: 这些只是示例，是否加入共享领域需要根据实际需求评估。

## 使用注意事项

### ✅ 应该做的

1. **使用共享领域的概念**: 当多个上下文需要相同概念时，使用共享领域
2. **保持向后兼容**: 修改共享领域时，尽量保持向后兼容
3. **文档化变更**: 修改共享领域时，更新文档并通知所有使用者
4. **测试覆盖**: 确保共享领域的代码有充分的测试覆盖

### ❌ 不应该做的

1. **不要过度共享**: 不要把所有通用代码都放在共享领域
2. **不要包含业务逻辑**: 共享领域只包含值对象和基础概念，不包含业务逻辑
3. **不要频繁变更**: 避免频繁修改共享领域，影响所有使用者
4. **不要包含基础设施代码**: 基础设施代码应该在 Infrastructure Layer

## 与有界上下文的关系

```
┌─────────────────────────────────────────────────────────┐
│                    Shared Kernel                         │
│  ┌──────────────────────────────────────────────────┐   │
│  │  City (值对象)                                    │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
           │                    │                    │
           │                    │                    │
    ┌──────▼──────┐      ┌──────▼──────┐      ┌──────▼──────┐
    │   Content   │      │   Search    │      │    User     │
    │   Context   │      │   Context   │      │   Context   │
    │             │      │             │      │  (未来版本)  │
    │ 使用 City   │      │ 使用 City   │      │ 使用 City   │
    └─────────────┘      └─────────────┘      └─────────────┘
```

## 文件结构

```
backend/internal/domain/shared/
├── city.go              # City 值对象
└── README.md            # 本文档
```

## 相关文档

- [DDD 设计文档](../../../.spec-workflow/specs/content-management-v1/design.md)
- [技术栈文档](../../../.spec-workflow/steering/tech.md)
- [内容领域文档](../content/README.md)

