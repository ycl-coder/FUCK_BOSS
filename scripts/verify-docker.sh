#!/bin/bash

# verify-docker.sh - 验证 Docker 配置
# 检查 Dockerfile 和 docker-compose.yml 配置是否正确

set -e

echo "=== Docker 配置验证 ==="
echo ""

# 检查 Docker 和 Docker Compose
echo "1. 检查 Docker 环境..."
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装"
    exit 1
fi
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose 未安装"
    exit 1
fi
echo "✓ Docker: $(docker --version)"
echo "✓ Docker Compose: $(docker-compose --version)"
echo ""

# 验证 docker-compose.yml 语法
echo "2. 验证 docker-compose.yml 语法..."
if docker-compose config --quiet > /dev/null 2>&1; then
    echo "✓ docker-compose.yml 语法正确"
else
    echo "❌ docker-compose.yml 语法错误"
    docker-compose config
    exit 1
fi
echo ""

# 检查 Dockerfile 是否存在
echo "3. 检查 Dockerfile..."
if [ ! -f "backend/Dockerfile" ]; then
    echo "❌ backend/Dockerfile 不存在"
    exit 1
fi
echo "✓ backend/Dockerfile 存在"
echo ""

# 检查必要的文件
echo "4. 检查必要的文件..."
REQUIRED_FILES=(
    "backend/go.mod"
    "backend/go.sum"
    "backend/cmd/server/main.go"
    "backend/api/proto/content/v1/content.proto"
    "backend/scripts/generate.sh"
)

for file in "${REQUIRED_FILES[@]}"; do
    if [ ! -f "$file" ]; then
        echo "❌ 缺少必要文件: $file"
        exit 1
    fi
done
echo "✓ 所有必要文件存在"
echo ""

# 检查 .dockerignore
echo "5. 检查 .dockerignore..."
if [ -f ".dockerignore" ]; then
    echo "✓ .dockerignore 存在"
else
    echo "⚠️  .dockerignore 不存在（可选）"
fi
echo ""

# 显示配置摘要
echo "6. 配置摘要..."
echo "   - PostgreSQL: docker.m.daocloud.io/postgres:16-alpine"
echo "   - Redis: docker.m.daocloud.io/redis:7-alpine"
echo "   - Backend: 多阶段构建 (golang:1.24-alpine -> alpine:latest)"
echo "   - 网络: fuck_boss_network (bridge)"
echo "   - 数据卷: postgres_data, redis_data"
echo ""

echo "=== 验证完成 ==="
echo ""
echo "下一步："
echo "  1. 构建镜像: make docker-build"
echo "  2. 启动服务: make docker-up"
echo "  3. 查看日志: make docker-logs"
echo "  4. 查看状态: make docker-ps"

