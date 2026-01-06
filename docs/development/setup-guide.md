# 开发环境设置指南

本文档介绍如何设置 Fuck Boss 平台的开发环境。

## 前置要求

### 必需软件

- **Go 1.24+**: [下载地址](https://golang.org/dl/)
- **Node.js 20.19+ 或 22.12+**: [下载地址](https://nodejs.org/)（前端开发必需）
  - 推荐使用 nvm 管理 Node.js 版本：`nvm install 20`
  - 详见 [Node.js 升级指南](./nodejs-upgrade.md)
- **Docker 20.10+**: [下载地址](https://www.docker.com/get-started)
- **Docker Compose 2.0+**: 通常随 Docker Desktop 一起安装
- **Protocol Buffers 编译器**: 
  - macOS: `brew install protobuf`
  - Linux: `apt-get install protobuf-compiler`
- **protoc-gen-go 插件**: `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
- **protoc-gen-go-grpc 插件**: `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`

### 可选工具

- **grpcurl**: gRPC 调试工具 - `brew install grpcurl` 或 `go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest`
- **grpcui**: gRPC Web UI - `go install github.com/fullstorydev/grpcui/cmd/grpcui@latest`
- **Make**: 构建工具（macOS/Linux 通常已安装）

## 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd fuck_boss
```

### 2. 启动测试环境

```bash
# 启动 PostgreSQL 和 Redis（测试环境）
make test-up

# 等待服务就绪（约 10-30 秒）
docker-compose -f docker-compose.test.yml ps
```

### 3. 配置开发环境

```bash
cd backend

# 复制配置文件
cp config/config.example.yaml config/config.yaml

# 编辑配置文件（使用测试环境配置）
# 修改 config/config.yaml:
#   database.port: 5433
#   database.user: test_user
#   database.password: test_password
#   database.dbname: test_db
#   redis.port: 6380
```

### 4. 安装依赖

```bash
cd backend

# 下载 Go 依赖
go mod download

# 生成 gRPC 代码
make generate-proto
# 或
./scripts/generate.sh
```

### 5. 运行服务器

```bash
cd backend

# 运行服务器
go run cmd/server/main.go

# 或使用 Makefile
make run
```

### 6. 验证服务

```bash
# 使用 grpcurl 测试
grpcurl -plaintext localhost:50051 list

# 或使用 grpcui（Web UI）
grpcui -plaintext localhost:50051
```

## 开发工作流

### 运行测试

```bash
cd backend

# 运行所有测试
go test ./...

# 运行单元测试
make test-unit-usecase

# 运行集成测试（需要测试环境运行）
make test-integration

# 运行 E2E 测试
make test-e2e-grpc

# 查看测试覆盖率
make test-unit-usecase-coverage-html
```

### 代码生成

```bash
cd backend

# 生成 gRPC 代码
make generate-proto

# 格式化代码
go fmt ./...

# 导入排序
goimports -w .
```

### 代码检查

```bash
cd backend

# 运行 go vet
go vet ./...

# 运行 golangci-lint（如果已安装）
golangci-lint run
```

## 项目结构

```
fuck_boss/
├── backend/                 # 后端代码
│   ├── api/                # API 定义（Protocol Buffers）
│   ├── cmd/                # 应用程序入口
│   ├── config/             # 配置文件
│   ├── internal/           # 内部代码
│   │   ├── domain/         # 领域层（DDD）
│   │   ├── application/    # 应用层（Use Cases）
│   │   ├── infrastructure/ # 基础设施层
│   │   └── presentation/   # 表现层（gRPC Handlers）
│   ├── pkg/                # 公共包
│   ├── scripts/            # 脚本
│   └── test/               # 测试代码
│       ├── unit/           # 单元测试
│       ├── integration/    # 集成测试
│       └── e2e/            # E2E 测试
├── frontend/               # 前端代码（待开发）
├── docs/                   # 文档
├── docker-compose.yml      # 生产环境 Docker Compose
├── docker-compose.test.yml # 测试环境 Docker Compose
└── Makefile               # 构建脚本
```

## 配置说明

### 开发环境配置

开发环境使用 `config/config.yaml` 文件：

```yaml
database:
  host: localhost
  port: 5433          # 测试环境端口
  user: test_user
  password: test_password
  dbname: test_db

redis:
  host: localhost
  port: 6380          # 测试环境端口

grpc:
  port: 50051

log:
  level: debug        # 开发环境使用 debug
  format: console     # 开发环境使用 console 格式
```

### 环境变量

所有配置都可以通过环境变量覆盖：

```bash
export FUCK_BOSS_DATABASE_HOST=localhost
export FUCK_BOSS_DATABASE_PORT=5433
export FUCK_BOSS_LOG_LEVEL=debug
export FUCK_BOSS_LOG_FORMAT=console
```

## 数据库管理

### 运行迁移

数据库迁移会在服务启动时自动执行。如果需要手动迁移：

```bash
# 连接到测试数据库
docker-compose -f docker-compose.test.yml exec postgres-test psql -U test_user -d test_db

# 查看表结构
\dt

# 查看表数据
SELECT * FROM posts LIMIT 10;
```

### 清理测试数据

```bash
# 停止测试环境（会清理所有数据，因为使用 tmpfs）
make test-down

# 重新启动
make test-up
```

## 调试技巧

### 使用 Delve 调试器

```bash
# 安装 Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 调试运行
dlv debug ./cmd/server
```

### 查看日志

```bash
# 查看服务器日志
go run cmd/server/main.go

# 查看 Docker 日志
docker-compose -f docker-compose.test.yml logs -f postgres-test
docker-compose -f docker-compose.test.yml logs -f redis-test
```

### 使用 grpcui 调试

```bash
# 启动 grpcui（Web UI）
grpcui -plaintext localhost:50051

# 浏览器打开 http://localhost:8080
```

## 常见问题

### 端口被占用

```bash
# 检查端口占用
lsof -i :50051
lsof -i :5433
lsof -i :6380

# 修改配置文件中的端口
```

### 数据库连接失败

1. 确保测试环境已启动：`make test-up`
2. 检查配置文件的端口和凭据
3. 等待数据库就绪：`docker-compose -f docker-compose.test.yml ps`

### gRPC 代码生成失败

1. 确保 protoc 已安装：`protoc --version`
2. 确保插件已安装：
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```
3. 检查 PATH：`echo $PATH | grep go/bin`

### 测试失败

1. 确保测试环境运行：`make test-up`
2. 检查数据库和 Redis 连接
3. 查看测试日志：`go test -v ./test/integration/...`

## 下一步

- [开发指南](./development-guide.md) - 了解代码规范和架构
- [测试指南](./testing-guide.md) - 了解如何编写和运行测试
- [API 文档](../api/README.md) - 了解 API 接口

