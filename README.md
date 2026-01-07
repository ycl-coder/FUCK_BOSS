# 全国公司曝光平台 (Fuck Boss)

一个允许用户分城市匿名发布公司不当行为的平台。该平台旨在为职场人士提供一个安全、匿名的渠道，曝光公司的不当行为，帮助其他求职者和员工了解潜在的工作环境问题，促进职场透明度和公平性。

## 核心特性

- 🔒 **匿名发布**：完全匿名，保护用户隐私和安全
- 🏙️ **城市分类**：按城市组织内容，便于本地化信息查找
- 🔍 **智能搜索**：支持按公司名称、城市、关键词搜索
- 📱 **响应式设计**：支持桌面和移动端访问
- 🚀 **高性能**：采用 DDD 架构，支持高并发访问
- 🐳 **容器化部署**：完整的 Docker 支持，一键部署

## 技术栈

### 后端
- **语言**: Go 1.23+
- **框架**: gRPC-Go, Protocol Buffers
- **数据库**: PostgreSQL 16+
- **缓存**: Redis 7+
- **架构**: DDD（领域驱动设计）

### 前端
- **框架**: React 19 + TypeScript 5
- **构建工具**: Vite 7
- **UI 库**: Ant Design 6
- **状态管理**: Zustand
- **数据获取**: TanStack Query
- **路由**: React Router 7

### 基础设施
- **容器化**: Docker + Docker Compose
- **反向代理**: Nginx（生产环境）
- **开发工具**: Make, Protocol Buffers Compiler

## 项目结构

```
fuck_boss/
├── backend/                 # 后端服务
│   ├── api/proto/           # Protocol Buffers 定义
│   ├── cmd/server/          # 应用程序入口
│   ├── config/              # 配置文件
│   ├── internal/            # 内部代码（DDD 分层）
│   │   ├── domain/          # 领域层（实体、值对象、Repository 接口）
│   │   ├── application/     # 应用层（Use Cases）
│   │   ├── infrastructure/  # 基础设施层（PostgreSQL、Redis 实现）
│   │   └── presentation/    # 表现层（gRPC、REST API Handlers）
│   ├── pkg/                 # 可复用的公共包
│   ├── scripts/             # 脚本文件
│   └── test/                # 测试代码
│       ├── unit/            # 单元测试
│       ├── integration/     # 集成测试
│       └── e2e/            # 端到端测试
├── frontend/                # 前端应用
│   ├── src/
│   │   ├── api/             # API 客户端
│   │   ├── features/        # 功能模块
│   │   ├── shared/          # 共享组件和工具
│   │   └── app/             # 应用配置
│   ├── public/              # 静态资源
│   └── test/                # 前端测试
├── docs/                    # 项目文档
│   ├── deployment/          # 部署文档
│   └── development/         # 开发文档
├── .spec-workflow/          # 规范工作流文档
├── docker-compose.yml       # Docker Compose 配置
├── docker-compose.test.yml  # 测试环境配置
└── Makefile                 # 构建脚本
```

## 快速开始

### 前置要求

- **Go 1.23+**: [下载地址](https://golang.org/dl/)
- **Node.js 20+**: [下载地址](https://nodejs.org/)（推荐使用 nvm）
- **Docker 20.10+**: [下载地址](https://www.docker.com/get-started)
- **Docker Compose 2.0+**: 通常随 Docker Desktop 一起安装
- **Protocol Buffers 编译器**: 
  - macOS: `brew install protobuf`
  - Linux: `apt-get install protobuf-compiler`

### 使用 Docker Compose（推荐）

```bash
# 1. 克隆项目
git clone <repository-url>
cd fuck_boss

# 2. 启动所有服务（PostgreSQL、Redis、后端、前端）
docker-compose up -d

# 3. 等待服务启动（约 30 秒）
docker-compose ps

# 4. 访问应用
# 前端: http://localhost:8000
# 后端 gRPC: localhost:50051
```

### 本地开发

#### 1. 启动测试环境

```bash
# 启动 PostgreSQL 和 Redis（测试环境）
make test-up

# 等待服务就绪
docker-compose -f docker-compose.test.yml ps
```

#### 2. 配置后端

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

#### 3. 启动后端服务

```bash
cd backend

# 安装依赖
go mod download

# 生成 gRPC 代码
make generate-proto

# 启动服务
go run cmd/server/main.go
# 或使用 Makefile
make run
```

#### 4. 启动前端服务

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 应用将在 http://localhost:8000 启动
```

### 使用 ngrok 暴露前端

如果需要通过公网访问前端（例如用于演示或测试）：

```bash
# 1. 启动前端开发服务器
cd frontend
npm run dev

# 2. 使用 ngrok 暴露前端端口
ngrok http 8000

# 3. 访问 ngrok 提供的 URL
# API 请求会自动通过 Vite proxy 转发到本地后端
```

详细配置请参考：[Ngrok 配置指南](docs/development/ngrok-setup.md)

## 开发命令

### 后端

```bash
# 运行测试
make backend-test

# 运行单元测试
make test-unit-usecase

# 运行集成测试（需要测试环境运行）
make test-integration

# 查看测试覆盖率
make backend-test-coverage-html

# 代码格式化
cd backend && go fmt ./...

# 代码检查
make backend-lint
```

### 前端

```bash
# 安装依赖
make frontend-install

# 开发模式
make frontend-dev

# 构建生产版本
make frontend-build

# 运行 E2E 测试
cd frontend && npm run test:e2e
```

### Docker

```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 停止所有服务
docker-compose down

# 重建并启动
docker-compose up -d --build
```

## 文档

### 开发文档
- [开发环境设置指南](docs/development/setup-guide.md)
- [开发指南](docs/development/development-guide.md)
- [测试指南](docs/development/testing-guide.md)
- [gRPC Web 设置指南](docs/development/grpc-web-setup.md)
- [Ngrok 配置指南](docs/development/ngrok-setup.md)
- [前端验证指南](docs/development/frontend-verification.md)

### 部署文档
- [Docker 部署指南](docs/deployment/docker-deploy.md)
- [生产环境部署指南](docs/deployment/production-deploy.md)

### 子项目文档
- [后端开发指南](backend/README.md)
- [前端开发指南](frontend/README.md)

### 规范文档
- [产品指导](.spec-workflow/steering/product.md)
- [技术栈](.spec-workflow/steering/tech.md)
- [项目结构](.spec-workflow/steering/structure.md)
- [功能规范](.spec-workflow/specs/content-management-v1/)

## 功能特性

### 第一版本核心功能

1. **匿名发布**
   - 用户可以匿名发布公司不当行为
   - 支持选择城市分类
   - 包含公司名称、问题描述、发生时间等关键信息
   - 无需注册即可发布

2. **内容查看**
   - 按城市浏览曝光内容
   - 支持列表和详情两种视图
   - 显示发布时间、城市、公司名称等关键信息
   - 按时间倒序排列
   - 支持分页浏览

3. **搜索功能**
   - 支持按公司名称搜索
   - 支持按城市筛选
   - 支持关键词搜索内容描述
   - 搜索结果高亮显示匹配项

### API 接口

- **gRPC API**: `localhost:50051`（gRPC 协议）
- **REST API**: `localhost:50051/api/*`（JSON over HTTP）
  - `POST /api/posts` - 创建帖子
  - `GET /api/posts` - 获取帖子列表
  - `GET /api/posts/:id` - 获取帖子详情
  - `POST /api/posts/search` - 搜索帖子

## 测试

项目包含完整的测试套件：

- **单元测试**: 使用 Mock 隔离依赖，快速执行
- **集成测试**: 使用真实数据库和 Redis，验证技术实现
- **E2E 测试**: 使用 Docker Compose 完整环境，验证系统功能

运行测试：

```bash
# 后端测试
make backend-test
make test-integration

# 前端测试
cd frontend && npm run test:e2e
```

## 贡献指南

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'feat: Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

### 代码规范

- 遵循 Go 代码规范（见 [tech.md](.spec-workflow/steering/tech.md)）
- 遵循 DDD 架构原则
- 所有代码必须通过测试
- 提交前运行 `gofmt` 和 `go vet`

## 许可证

（待定）

## 联系方式

- 项目地址: [GitHub](https://github.com/ycl-coder/FUCK_BOSS)
- 问题反馈: [Issues](https://github.com/ycl-coder/FUCK_BOSS/issues)

---

**注意**: 本项目旨在促进职场透明度，请确保发布的内容真实、客观。恶意发布虚假信息的行为是不被允许的。
