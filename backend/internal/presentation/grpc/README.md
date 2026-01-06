# grpc - gRPC 处理器

gRPC 服务的处理器实现，调用 Application Layer 的 Use Cases。

## 结构

- **content_handler.go** - ContentService gRPC 实现

## ContentService

实现 Protocol Buffers 定义的 ContentService 接口。

```go
service ContentService {
  rpc CreatePost(CreatePostRequest) returns (CreatePostResponse);
  rpc ListPosts(ListPostsRequest) returns (ListPostsResponse);
  rpc GetPost(GetPostRequest) returns (GetPostResponse);
  rpc SearchPosts(SearchPostsRequest) returns (SearchPostsResponse);
}
```

## 实现

```go
type ContentService struct {
    createUseCase *application.CreatePostUseCase
    listUseCase   *application.ListPostsUseCase
    getUseCase    *application.GetPostUseCase
    searchUseCase *application.SearchPostsUseCase
}

func (s *ContentService) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {
    // 转换请求为 Command
    // 调用 UseCase
    // 转换响应
}
```

## 错误处理

将应用错误转换为 gRPC 状态码：

- 验证错误 → `InvalidArgument`
- 未找到 → `NotFound`
- 限流错误 → `ResourceExhausted`
- 内部错误 → `Internal`

## 注意事项

- 使用 context 支持取消和超时
- 所有错误必须转换为 gRPC 状态码
- 不直接暴露 Domain Entity

