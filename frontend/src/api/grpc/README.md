# gRPC Web Client

gRPC Web 客户端配置和实现。

## 当前状态

⚠️ **注意**: 当前使用的是 Mock 实现。需要从 Protocol Buffers 文件生成 TypeScript 代码。

## 生成 TypeScript 代码

### 方法 1: 使用 protoc-gen-grpc-web (推荐)

```bash
# 安装 protoc-gen-grpc-web
npm install -D protoc-gen-grpc-web

# 生成 TypeScript 代码
protoc \
  --plugin=protoc-gen-grpc-web=./node_modules/.bin/protoc-gen-grpc-web \
  --grpc-web_out=import_style=typescript,mode=grpcwebtext:./src/api/grpc \
  --proto_path=../backend/api/proto \
  ../backend/api/proto/content/v1/content.proto
```

### 方法 2: 使用 @grpc/grpc-web (当前使用)

需要手动创建客户端包装器，或使用工具生成。

## 使用示例

```typescript
import { contentServiceClient } from '@/api/grpc/contentClient'

// 创建帖子
const result = await contentServiceClient.createPost({
  company: '测试公司',
  cityCode: 'beijing',
  cityName: '北京',
  content: '这是一个测试内容',
})

// 获取帖子列表
const posts = await contentServiceClient.listPosts('beijing', 1, 20)

// 获取帖子详情
const post = await contentServiceClient.getPost('post-id')

// 搜索帖子
const results = await contentServiceClient.searchPosts({
  keyword: '测试',
  cityCode: 'beijing',
  page: 1,
  pageSize: 20,
})
```

## 错误处理

所有 gRPC 错误都会被转换为标准的 Error 对象。

```typescript
try {
  const post = await contentServiceClient.getPost('invalid-id')
} catch (error) {
  console.error('Failed to get post:', error)
  // Handle error
}
```

## 配置

gRPC 服务器地址通过环境变量配置：

```env
VITE_GRPC_URL=http://localhost:50051
```

或在代码中指定：

```typescript
const client = createContentServiceClient('http://localhost:50051')
```

