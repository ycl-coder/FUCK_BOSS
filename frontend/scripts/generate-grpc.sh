#!/bin/bash

# generate-grpc.sh - 生成 gRPC Web TypeScript 代码
# 从 Protocol Buffers 文件生成 TypeScript 客户端代码

set -e

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# 前端项目根目录
FRONTEND_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
# 后端项目根目录
BACKEND_ROOT="$(cd "$FRONTEND_ROOT/.." && pwd)"

# Proto 文件目录
PROTO_DIR="$BACKEND_ROOT/backend/api/proto"
# 输出目录
OUTPUT_DIR="$FRONTEND_ROOT/src/api/grpc"

# 检查 protoc 是否安装
if ! command -v protoc &> /dev/null; then
    echo "错误: protoc 未安装"
    echo "请安装 Protocol Buffers 编译器:"
    echo "  macOS: brew install protobuf"
    echo "  Linux: apt-get install protobuf-compiler"
    exit 1
fi

# 检查 protoc-gen-grpc-web 是否安装
if ! command -v protoc-gen-grpc-web &> /dev/null; then
    echo "错误: protoc-gen-grpc-web 未安装"
    echo "请安装: npm install -D protoc-gen-grpc-web"
    exit 1
fi

echo "开始生成 gRPC Web TypeScript 代码..."

# 生成 content/v1 的代码
CONTENT_PROTO_DIR="$PROTO_DIR/content/v1"
if [ -d "$CONTENT_PROTO_DIR" ]; then
    echo "生成 content/v1 的代码..."
    protoc \
        --plugin=protoc-gen-grpc-web="$FRONTEND_ROOT/node_modules/.bin/protoc-gen-grpc-web" \
        --grpc-web_out=import_style=typescript,mode=grpcwebtext:$OUTPUT_DIR \
        --proto_path="$PROTO_DIR" \
        "$CONTENT_PROTO_DIR"/*.proto
    
    echo "✓ content/v1 代码生成完成"
else
    echo "警告: $CONTENT_PROTO_DIR 目录不存在"
fi

echo "所有 gRPC Web TypeScript 代码生成完成！"
echo "输出目录: $OUTPUT_DIR"

