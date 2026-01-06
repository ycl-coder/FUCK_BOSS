# dto - 数据传输对象

应用层的数据传输对象（DTO），用于在不同层之间传输数据。

## 结构

- **content_dto.go** - 内容相关的 DTO
- **search_dto.go** - 搜索相关的 DTO

## DTOs

### PostDTO

Post 的数据传输对象，用于 API 响应。

**定义**:
```go
type PostDTO struct {
    ID        string      // Post ID (UUID)
    Company   string      // 公司名称
    CityCode  string      // 城市代码
    CityName  string      // 城市名称
    Content   string      // 内容
    OccurredAt *time.Time // 发生时间（可选）
    CreatedAt time.Time   // 创建时间
}
```

**使用场景**:
- gRPC API 响应
- 前端展示
- 跨层数据传输

**转换**:
- 从 Domain Entity (`Post`) 转换为 DTO
- 不包含业务逻辑，只包含数据

### PostsListDTO

Post 列表的数据传输对象，包含分页信息。

**定义**:
```go
type PostsListDTO struct {
    Posts    []*PostDTO  // Post 列表
    Total    int         // 总数（跨所有页面）
    Page     int         // 当前页码（1-based）
    PageSize int         // 每页数量
}
```

**使用场景**:
- 列表查询 API 响应
- 分页展示
- 前端分页控件

## 注意事项

- DTO 不包含业务逻辑
- 用于隔离 Domain Entity 和外部接口
- 可以包含格式化后的数据（如相对时间）

