# Fuck Boss Frontend

Fuck Boss 平台的前端应用，使用 React + TypeScript + Vite 构建。

## 技术栈

- **React 19**: UI 框架
- **TypeScript 5**: 类型安全
- **Vite 7**: 构建工具
- **React Router 7**: 前端路由
- **Zustand**: 状态管理
- **Ant Design**: UI 组件库
- **TanStack Query**: 数据获取和缓存
- **gRPC Web**: 与后端 gRPC 服务通信
- **Day.js**: 日期处理

## 快速开始

### 安装依赖

```bash
npm install
```

### 开发模式

```bash
npm run dev
```

应用将在 `http://localhost:3000` 启动。

### 构建生产版本

```bash
npm run build
```

### 预览生产构建

```bash
npm run preview
```

## 项目结构

```
frontend/
├── src/
│   ├── api/              # API 客户端
│   │   └── grpc/         # gRPC Web 客户端
│   ├── features/         # 功能模块
│   │   ├── post/         # 帖子功能
│   │   └── search/       # 搜索功能
│   ├── shared/           # 共享组件和工具
│   │   ├── components/   # 共享组件
│   │   ├── hooks/        # 自定义 Hooks
│   │   └── types/        # TypeScript 类型
│   └── app/              # 应用配置
│       └── routes.tsx    # 路由配置
├── public/               # 静态资源
└── package.json
```

## 环境变量

创建 `.env` 文件配置环境变量：

```env
VITE_GRPC_URL=http://localhost:50051
VITE_API_BASE_URL=http://localhost:50051
```

## 开发指南

详见 `docs/development/` 目录下的文档。
