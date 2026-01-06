# Backend - 全国公司曝光平台

## 项目结构

采用 DDD 分层架构：

- `cmd/server/` - 应用程序入口
- `internal/domain/` - 领域层（实体、值对象、Repository 接口）
- `internal/application/` - 应用层（Use Cases）
- `internal/infrastructure/` - 基础设施层（PostgreSQL、Redis 实现）
- `internal/presentation/` - 表现层（gRPC Handlers）
- `pkg/` - 可复用的公共包
- `api/proto/` - Protocol Buffers 定义
- `test/` - 测试文件

## 开发环境要求

- Go 1.21+
- PostgreSQL 14+
- Redis 7.0+
- Docker & Docker Compose（推荐）

## 快速开始

```bash
# 安装依赖
go mod tidy

# 运行测试
go test ./...

# 启动服务
go run cmd/server/main.go
```

## 更多信息

详见 `docs/development/` 目录下的开发文档。

