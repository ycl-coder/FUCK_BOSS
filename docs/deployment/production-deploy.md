# 生产环境部署指南

本文档介绍如何在生产环境中部署 Fuck Boss 平台。

## 前置要求

- Linux 服务器（推荐 Ubuntu 22.04 LTS 或 CentOS 8+）
- Docker 20.10+ 和 Docker Compose 2.0+
- 至少 4GB RAM，8GB 推荐
- 至少 50GB 可用磁盘空间
- 域名和 SSL 证书（如需要 HTTPS）

## 架构建议

### 推荐架构

```
                    ┌─────────────┐
                    │   Nginx     │  (反向代理/负载均衡)
                    │  (可选)     │
                    └──────┬──────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
   ┌────▼────┐       ┌────▼────┐       ┌────▼────┐
   │ Backend │       │ Backend │       │ Backend │
   │  (gRPC) │       │  (gRPC) │       │  (gRPC) │
   └────┬────┘       └────┬────┘       └────┬────┘
        │                  │                  │
        └──────────────────┼──────────────────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
   ┌────▼────┐       ┌────▼────┐       ┌────▼────┐
   │PostgreSQL│      │  Redis  │       │  Redis  │
   │ (主从)  │      │ (集群)  │       │ (集群)  │
   └─────────┘      └─────────┘       └─────────┘
```

### 单机部署（小型应用）

适合：日访问量 < 10,000，数据量 < 100GB

- 1 个 Backend 实例
- 1 个 PostgreSQL 实例
- 1 个 Redis 实例

### 高可用部署（中大型应用）

适合：日访问量 > 10,000，需要高可用性

- 多个 Backend 实例（负载均衡）
- PostgreSQL 主从复制
- Redis 集群或哨兵模式

## 部署步骤

### 1. 服务器准备

```bash
# 更新系统
sudo apt-get update && sudo apt-get upgrade -y

# 安装 Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 安装 Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 配置 Docker 用户组
sudo usermod -aG docker $USER
```

### 2. 项目部署

```bash
# 创建项目目录
sudo mkdir -p /opt/fuck_boss
sudo chown $USER:$USER /opt/fuck_boss
cd /opt/fuck_boss

# 克隆项目（或上传代码）
git clone <repository-url> .

# 或使用 scp 上传
# scp -r /path/to/fuck_boss user@server:/opt/fuck_boss
```

### 3. 配置生产环境

#### 创建生产环境配置文件

```bash
# 创建生产环境配置目录
mkdir -p /opt/fuck_boss/config/production

# 创建 .env 文件
cat > /opt/fuck_boss/.env <<EOF
# 数据库配置
POSTGRES_USER=fuck_boss_user
POSTGRES_PASSWORD=$(openssl rand -base64 32)
POSTGRES_DB=fuck_boss_prod

# Redis 配置
REDIS_PASSWORD=$(openssl rand -base64 32)

# gRPC 配置
GRPC_PORT=50051

# 日志配置
LOG_LEVEL=info
LOG_FORMAT=json
EOF

# 设置文件权限
chmod 600 /opt/fuck_boss/.env
```

#### 修改 docker-compose.yml

创建生产环境专用的 `docker-compose.prod.yml`：

```yaml
version: '3.8'

services:
  postgres:
    environment:
      POSTGRES_PASSWORD_FILE: /run/secrets/postgres_password
    secrets:
      - postgres_password
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backups:/backups
    restart: always
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

  redis:
    command: redis-server --requirepass ${REDIS_PASSWORD} --appendonly yes
    volumes:
      - redis_data:/data
    restart: always
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

  backend:
    environment:
      FUCK_BOSS_DATABASE_PASSWORD_FILE: /run/secrets/postgres_password
      FUCK_BOSS_REDIS_PASSWORD: ${REDIS_PASSWORD}
    secrets:
      - postgres_password
    restart: always
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '1'
          memory: 1G
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

secrets:
  postgres_password:
    file: ./secrets/postgres_password.txt

volumes:
  postgres_data:
  redis_data:
```

### 4. 安全配置

#### 创建密钥文件

```bash
# 创建密钥目录
mkdir -p /opt/fuck_boss/secrets
chmod 700 /opt/fuck_boss/secrets

# 生成数据库密码
openssl rand -base64 32 > /opt/fuck_boss/secrets/postgres_password.txt
chmod 600 /opt/fuck_boss/secrets/postgres_password.txt
```

#### 配置防火墙

```bash
# 安装 UFW（Ubuntu）
sudo apt-get install ufw

# 允许 SSH
sudo ufw allow 22/tcp

# 允许 gRPC 端口（仅内网，或通过 Nginx 代理）
sudo ufw allow from 10.0.0.0/8 to any port 50051

# 启用防火墙
sudo ufw enable
```

### 5. 启动服务

```bash
cd /opt/fuck_boss

# 构建镜像
docker-compose -f docker-compose.yml -f docker-compose.prod.yml build

# 启动服务
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# 查看状态
docker-compose -f docker-compose.yml -f docker-compose.prod.yml ps
```

### 6. 配置 Nginx 反向代理（可选）

如果需要通过 HTTP/HTTPS 访问 gRPC 服务：

```nginx
# /etc/nginx/sites-available/fuck_boss
upstream grpc_backend {
    server localhost:50051;
}

server {
    listen 80;
    server_name api.fuck_boss.com;

    # gRPC 代理
    location / {
        grpc_pass grpc://grpc_backend;
        grpc_set_header Host $host;
        grpc_set_header X-Real-IP $remote_addr;
    }
}

# HTTPS 配置
server {
    listen 443 ssl http2;
    server_name api.fuck_boss.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        grpc_pass grpc://grpc_backend;
        grpc_set_header Host $host;
        grpc_set_header X-Real-IP $remote_addr;
    }
}
```

## 监控和日志

### 配置日志轮转

```bash
# 创建日志轮转配置
sudo tee /etc/logrotate.d/docker-containers <<EOF
/var/lib/docker/containers/*/*.log {
    rotate 7
    daily
    compress
    size=10M
    missingok
    delaycompress
    copytruncate
}
EOF
```

### 设置监控

推荐使用以下监控工具：

1. **Prometheus + Grafana**
   - 监控服务指标
   - 可视化仪表板

2. **ELK Stack**
   - 日志收集和分析
   - 集中式日志管理

3. **健康检查脚本**

```bash
#!/bin/bash
# /opt/fuck_boss/scripts/health-check.sh

HEALTH_URL="http://localhost:50051"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" $HEALTH_URL)

if [ $STATUS -ne 200 ]; then
    echo "Service unhealthy, status: $STATUS"
    # 发送告警
    exit 1
fi
```

## 备份策略

### 数据库备份

```bash
#!/bin/bash
# /opt/fuck_boss/scripts/backup-db.sh

BACKUP_DIR="/opt/fuck_boss/backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/fuck_boss_$DATE.sql"

# 创建备份
docker-compose exec -T postgres pg_dump -U postgres fuck_boss > $BACKUP_FILE

# 压缩备份
gzip $BACKUP_FILE

# 删除 7 天前的备份
find $BACKUP_DIR -name "*.sql.gz" -mtime +7 -delete

# 上传到云存储（可选）
# aws s3 cp $BACKUP_FILE.gz s3://backup-bucket/
```

### 设置定时备份

```bash
# 添加到 crontab
crontab -e

# 每天凌晨 2 点备份
0 2 * * * /opt/fuck_boss/scripts/backup-db.sh
```

## 性能优化

### PostgreSQL 优化

编辑 `docker-compose.prod.yml`：

```yaml
postgres:
  environment:
    POSTGRES_INITDB_ARGS: "-E UTF8 --locale=C"
  command: >
    postgres
    -c shared_buffers=256MB
    -c effective_cache_size=1GB
    -c maintenance_work_mem=64MB
    -c checkpoint_completion_target=0.9
    -c wal_buffers=16MB
    -c default_statistics_target=100
    -c random_page_cost=1.1
    -c effective_io_concurrency=200
    -c work_mem=4MB
    -c min_wal_size=1GB
    -c max_wal_size=4GB
```

### Redis 优化

```yaml
redis:
  command: >
    redis-server
    --maxmemory 2gb
    --maxmemory-policy allkeys-lru
    --appendonly yes
    --appendfsync everysec
```

### Backend 优化

```yaml
backend:
  environment:
    FUCK_BOSS_DATABASE_MAX_OPEN_CONNS: "200"
    FUCK_BOSS_DATABASE_MAX_IDLE_CONNS: "50"
    FUCK_BOSS_REDIS_POOL_SIZE: "100"
```

## 高可用部署

### PostgreSQL 主从复制

参考 PostgreSQL 官方文档配置主从复制。

### Redis 哨兵模式

```yaml
redis-sentinel:
  image: docker.m.daocloud.io/redis:7-alpine
  command: redis-sentinel /etc/redis/sentinel.conf
  volumes:
    - ./redis/sentinel.conf:/etc/redis/sentinel.conf
```

### 负载均衡

使用 Nginx 或 HAProxy 进行负载均衡：

```nginx
upstream grpc_backend {
    least_conn;
    server backend1:50051;
    server backend2:50051;
    server backend3:50051;
}
```

## 故障恢复

### 服务恢复

```bash
# 重启服务
docker-compose restart

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs --tail=100 backend
```

### 数据恢复

```bash
# 恢复数据库备份
docker-compose exec -T postgres psql -U postgres fuck_boss < backup.sql
```

## 安全建议

1. **定期更新**：保持 Docker 和系统更新
2. **最小权限**：使用非 root 用户运行容器
3. **网络安全**：使用防火墙限制访问
4. **密钥管理**：使用 Docker Secrets 或外部密钥管理服务
5. **审计日志**：启用审计日志记录
6. **定期备份**：自动化备份流程

## 相关文档

- [Docker 部署指南](./docker-deploy.md)
- [开发环境设置](../development/setup-guide.md)
- [监控和告警](../monitoring/README.md)

