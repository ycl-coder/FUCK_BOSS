# proto - Protocol Buffers 定义

gRPC API 的 Protocol Buffers 定义文件。

## 结构

- **content/v1/content.proto** - 内容服务的 API 定义
- **search/v1/search.proto** - 搜索服务的 API 定义（如需要）

## 使用

### 定义服务

在 `.proto` 文件中定义服务和消息：

```protobuf
syntax = "proto3";

package content.v1;

service ContentService {
  rpc CreatePost(CreatePostRequest) returns (CreatePostResponse);
}
```

### 生成代码

使用 `protoc` 生成 Go 代码：

```bash
protoc --go_out=. --go-grpc_out=. api/proto/content/v1/content.proto
```

## 版本管理

- 使用 `/v1` 路径进行版本管理
- 未来版本使用 `/v2`、`/v3` 等
- 保持向后兼容性

## 注意事项

- 使用 `proto3` 语法
- 正确设置 `go_package` 选项
- 字段编号不能重复
- 遵循 protobuf 最佳实践

