# Docker 部署指南

本文档介绍如何使用 Docker 和 Docker Compose 部署 Fuck Boss 平台。

## 前置要求

- Docker 20.10+ 
- Docker Compose 2.0+
- 至少 2GB 可用内存
- 至少 10GB 可用磁盘空间

## 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd fuck_boss
```

### 2. 配置环境变量（可选）

默认配置已包含在 `docker-compose.yml` 中，如需自定义，可以：

- 修改 `docker-compose.yml` 中的环境变量
- 或创建 `.env` 文件（Docker Compose 会自动读取）

示例 `.env` 文件：

```env
# 数据库配置
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your_secure_password
POSTGRES_DB=fuck_boss

# Redis 配置（可选）
REDIS_PASSWORD=

# gRPC 端口
GRPC_PORT=50051
```

### 3. 构建和启动服务

```bash
# 构建 Docker 镜像
make docker-build

# 启动所有服务（PostgreSQL, Redis, Backend）
make docker-up

# 查看服务状态
make docker-ps

# 查看日志
make docker-logs
```

### 4. 验证服务

```bash
# 检查服务健康状态
docker-compose ps

# 测试 gRPC 服务（需要 grpcurl）
grpcurl -plaintext localhost:50051 list

# 查看后端日志
docker-compose logs -f backend
```

## 服务说明

### PostgreSQL 数据库

- **容器名**: `fuck_boss_postgres`
- **端口**: `5432`
- **数据持久化**: `postgres_data` volume
- **默认用户**: `postgres`
- **默认密码**: `postgres_password`（生产环境请修改）

### Redis 缓存

- **容器名**: `fuck_boss_redis`
- **端口**: `6379`
- **数据持久化**: `redis_data` volume
- **持久化**: AOF (Append Only File) 已启用

### Backend gRPC 服务

- **容器名**: `fuck_boss_backend`
- **端口**: `50051`
- **健康检查**: 每 30 秒检查一次
- **依赖**: 等待 PostgreSQL 和 Redis 健康后启动

## 常用命令

### 启动和停止

```bash
# 启动所有服务
make docker-up
# 或
docker-compose up -d

# 停止所有服务
make docker-down
# 或
docker-compose down

# 重启服务
make docker-restart
# 或
docker-compose restart
```

### 查看日志

```bash
# 查看所有服务日志
make docker-logs
# 或
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f backend
docker-compose logs -f postgres
docker-compose logs -f redis
```

### 服务管理

```bash
# 查看服务状态
make docker-ps
# 或
docker-compose ps

# 进入容器
docker-compose exec backend sh
docker-compose exec postgres psql -U postgres -d fuck_boss

# 执行数据库迁移（如果需要）
docker-compose exec backend ./server
```

### 数据管理

```bash
# 备份 PostgreSQL 数据
docker-compose exec postgres pg_dump -U postgres fuck_boss > backup.sql

# 恢复 PostgreSQL 数据
docker-compose exec -T postgres psql -U postgres fuck_boss < backup.sql

# 清理所有数据（危险操作）
make docker-clean
# 或
docker-compose down -v
```

## 配置说明

### 环境变量

Backend 服务支持以下环境变量（格式：`FUCK_BOSS_<SECTION>_<FIELD>`）：

#### 数据库配置

- `FUCK_BOSS_DATABASE_HOST`: 数据库主机（默认: postgres）
- `FUCK_BOSS_DATABASE_PORT`: 数据库端口（默认: 5432）
- `FUCK_BOSS_DATABASE_USER`: 数据库用户（默认: postgres）
- `FUCK_BOSS_DATABASE_PASSWORD`: 数据库密码
- `FUCK_BOSS_DATABASE_DBNAME`: 数据库名称（默认: fuck_boss）
- `FUCK_BOSS_DATABASE_SSLMODE`: SSL 模式（默认: disable）

#### Redis 配置

- `FUCK_BOSS_REDIS_HOST`: Redis 主机（默认: redis）
- `FUCK_BOSS_REDIS_PORT`: Redis 端口（默认: 6379）
- `FUCK_BOSS_REDIS_PASSWORD`: Redis 密码（默认: 空）
- `FUCK_BOSS_REDIS_DB`: Redis 数据库编号（默认: 0）

#### gRPC 配置

- `FUCK_BOSS_GRPC_PORT`: gRPC 服务端口（默认: 50051）

#### 日志配置

- `FUCK_BOSS_LOG_LEVEL`: 日志级别（debug/info/warn/error，默认: info）
- `FUCK_BOSS_LOG_FORMAT`: 日志格式（json/text/console，默认: json）

### 网络配置

所有服务运行在 `fuck_boss_network` 网络中，服务之间可以通过服务名互相访问：

- `postgres` - PostgreSQL 服务
- `redis` - Redis 服务
- `backend` - Backend gRPC 服务

## 故障排查

### 服务无法启动

1. **检查端口占用**：
   ```bash
   # 检查端口是否被占用
   lsof -i :5432  # PostgreSQL
   lsof -i :6379  # Redis
   lsof -i :50051 # gRPC
   ```

2. **查看服务日志**：
   ```bash
   docker-compose logs backend
   docker-compose logs postgres
   docker-compose logs redis
   ```

3. **检查健康状态**：
   ```bash
   docker-compose ps
   # 查看健康检查状态
   ```

### 数据库连接失败

1. **检查数据库是否就绪**：
   ```bash
   docker-compose exec postgres pg_isready -U postgres
   ```

2. **检查环境变量**：
   ```bash
   docker-compose exec backend env | grep FUCK_BOSS_DATABASE
   ```

3. **测试数据库连接**：
   ```bash
   docker-compose exec backend sh
   # 在容器内测试连接
   ```

### Redis 连接失败

1. **检查 Redis 是否就绪**：
   ```bash
   docker-compose exec redis redis-cli ping
   ```

2. **检查 Redis 配置**：
   ```bash
   docker-compose exec backend env | grep FUCK_BOSS_REDIS
   ```

### 构建失败

1. **清理构建缓存**：
   ```bash
   docker-compose build --no-cache backend
   ```

2. **检查 Dockerfile**：
   ```bash
   docker build -t test-backend ./backend
   ```

## 生产环境建议

### 安全配置

1. **修改默认密码**：
   - 修改 `docker-compose.yml` 中的 `POSTGRES_PASSWORD`
   - 使用强密码（至少 16 个字符，包含大小写字母、数字和特殊字符）

2. **启用 SSL**：
   - 配置 PostgreSQL SSL 证书
   - 设置 `FUCK_BOSS_DATABASE_SSLMODE=require`

3. **限制网络访问**：
   - 使用防火墙限制端口访问
   - 仅允许必要的 IP 访问数据库和 Redis

### 性能优化

1. **资源限制**：
   ```yaml
   services:
     backend:
       deploy:
         resources:
           limits:
             cpus: '2'
             memory: 2G
           reservations:
             cpus: '1'
             memory: 1G
   ```

2. **数据库优化**：
   - 调整 PostgreSQL 连接池大小
   - 配置适当的缓存策略

3. **Redis 优化**：
   - 配置 Redis 内存限制
   - 启用持久化策略

### 监控和日志

1. **日志收集**：
   - 配置日志驱动（如 `json-file`）
   - 使用日志聚合工具（如 ELK Stack）

2. **健康检查**：
   - 配置监控系统（如 Prometheus）
   - 设置告警规则

3. **备份策略**：
   - 定期备份数据库
   - 备份 Redis 数据（如需要）

## 升级和维护

### 更新服务

```bash
# 拉取最新代码
git pull

# 重新构建镜像
make docker-build

# 重启服务
make docker-restart
```

### 数据库迁移

数据库迁移会在服务启动时自动执行。如果需要手动迁移：

```bash
# 进入后端容器
docker-compose exec backend sh

# 运行迁移（如果实现了迁移工具）
```

## 卸载

```bash
# 停止并删除所有容器
make docker-down

# 删除所有数据卷（危险操作）
make docker-clean
```

## 相关文档

- [生产环境部署指南](./production-deploy.md)
- [开发环境设置](../development/setup-guide.md)
- [故障排查指南](../development/troubleshooting.md)

