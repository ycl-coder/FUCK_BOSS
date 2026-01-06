# scripts - 脚本文件

开发和部署相关的脚本。

## 脚本

- **generate.sh** - gRPC 代码生成脚本（protoc）

## 使用

### gRPC 代码生成

使用 Makefile 命令（推荐）：

```bash
make generate-proto
```

或直接运行脚本：

```bash
cd backend && ./scripts/generate.sh
```

### 前置要求

1. **安装 protoc**:
   ```bash
   # macOS
   brew install protobuf
   
   # Linux
   apt-get install protobuf-compiler
   ```

2. **安装 Go 插件**:
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

### 生成的文件

脚本会从 `api/proto/content/v1/*.proto` 生成以下文件：

- `api/proto/content/v1/*.pb.go` - Protocol Buffers 消息类型
- `api/proto/content/v1/*_grpc.pb.go` - gRPC 服务接口

生成的文件会放在与 `.proto` 文件相同的目录中。

