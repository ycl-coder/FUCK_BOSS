# gRPC Web 设置指南

## 问题

前端无法直接调用 gRPC 服务（浏览器不支持 gRPC），需要使用 gRPC Web。

## 解决方案

### 方案 1: 使用 grpcweb 包（推荐）

在后端添加 gRPC Web 支持，使浏览器可以通过 HTTP 调用 gRPC 服务。

#### 1. 安装依赖

```bash
cd backend
go get github.com/improbable-eng/grpc-web/go/grpcweb
go mod tidy
```

#### 2. 后端已配置

后端代码已经配置了 gRPC Web 支持（`backend/cmd/server/main.go`），只需要安装依赖即可。

#### 3. 前端客户端

前端已经实现了 gRPC Web 客户端（`frontend/src/api/grpc/contentClient.ts`），使用 fetch API 调用后端。

### 方案 2: 使用 Envoy 代理（生产环境推荐）

在生产环境中，可以使用 Envoy 作为 gRPC Web 代理。

## 当前状态

- ✅ 后端代码已配置 gRPC Web 支持
- ✅ 前端客户端已实现
- ⚠️ 需要安装后端依赖：`go get github.com/improbable-eng/grpc-web/go/grpcweb`

## 安装步骤

1. **安装后端依赖**：
   ```bash
   cd backend
   go get github.com/improbable-eng/grpc-web/go/grpcweb
   go mod tidy
   ```

2. **重启后端服务**：
   ```bash
   go run cmd/server/main.go
   ```

3. **验证前端连接**：
   前端应该能够正常调用后端 gRPC 服务。

## 故障排除

### 错误：`could not import github.com/improbable-eng/grpc-web/go/grpcweb`

**原因**: 依赖未安装

**解决方案**:
```bash
cd backend
go get github.com/improbable-eng/grpc-web/go/grpcweb
go mod tidy
```

### 错误：CORS 错误

**原因**: 后端未允许前端域名

**解决方案**: 检查后端代码中的 `WithOriginFunc` 配置，确保允许前端域名。

### 错误：404 Not Found

**原因**: gRPC Web URL 格式不正确

**解决方案**: 确保前端使用正确的 URL 格式：`/package.service/method`

## 参考

- [grpc-web 文档](https://github.com/improbable-eng/grpc-web)
- [gRPC Web 规范](https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-WEB.md)

