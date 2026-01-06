#!/bin/bash

# generate.sh - 生成 gRPC Go 代码
# 使用 protoc 从 .proto 文件生成 Go 代码

set -e

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# 项目根目录（backend）
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Proto 文件目录
PROTO_DIR="$PROJECT_ROOT/api/proto"
# 输出目录
OUTPUT_DIR="$PROJECT_ROOT/api/proto"

# 检查 protoc 是否安装
if ! command -v protoc &> /dev/null; then
    echo "错误: protoc 未安装"
    echo "请安装 Protocol Buffers 编译器:"
    echo "  macOS: brew install protobuf"
    echo "  Linux: apt-get install protobuf-compiler"
    exit 1
fi

# 检查 protoc-gen-go 是否安装
if ! command -v protoc-gen-go &> /dev/null; then
    echo "错误: protoc-gen-go 未安装"
    echo "请安装: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
    exit 1
fi

# 检查 protoc-gen-go-grpc 是否安装
if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "错误: protoc-gen-go-grpc 未安装"
    echo "请安装: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
    exit 1
fi

echo "开始生成 gRPC Go 代码..."

# 生成 content/v1 的代码
CONTENT_PROTO_DIR="$PROTO_DIR/content/v1"
if [ -d "$CONTENT_PROTO_DIR" ]; then
    echo "生成 content/v1 的代码..."
    protoc \
        --go_out="$OUTPUT_DIR" \
        --go_opt=paths=source_relative \
        --go-grpc_out="$OUTPUT_DIR" \
        --go-grpc_opt=paths=source_relative \
        --proto_path="$PROTO_DIR" \
        "$CONTENT_PROTO_DIR"/*.proto
    
    # 如果文件生成到了错误的位置，移动到正确位置
    if [ -f "$PROJECT_ROOT/api/proto/fuck_boss/backend/api/proto/content/v1/content.pb.go" ]; then
        echo "移动文件到正确位置..."
        cp "$PROJECT_ROOT/api/proto/fuck_boss/backend/api/proto/content/v1/content.pb.go" \
           "$CONTENT_PROTO_DIR/content.pb.go"
        rm -rf "$PROJECT_ROOT/api/proto/fuck_boss"
    fi
    
    echo "✓ content/v1 代码生成完成"
else
    echo "警告: $CONTENT_PROTO_DIR 目录不存在"
fi

echo "所有 gRPC Go 代码生成完成！"

