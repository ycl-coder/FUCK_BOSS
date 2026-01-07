# 安装 gRPC Web 支持

## 问题

后端需要 `github.com/improbable-eng/grpc-web/go/grpcweb` 包来支持 gRPC Web，但该包尚未安装。

## 安装步骤

在 `backend` 目录下运行：

```bash
cd backend
go get github.com/improbable-eng/grpc-web/go/grpcweb
go mod tidy
```

## 验证

安装完成后，运行：

```bash
go run cmd/server/main.go
```

如果看到日志显示 "gRPC server started (with gRPC Web support)"，说明安装成功。

## 如果安装失败

如果遇到权限问题，可以尝试：

```bash
# 修复 npm 权限（如果相关）
sudo chown -R $(whoami) ~/.npm

# 或者使用 sudo（不推荐，但可以临时解决）
sudo go get github.com/improbable-eng/grpc-web/go/grpcweb
```

