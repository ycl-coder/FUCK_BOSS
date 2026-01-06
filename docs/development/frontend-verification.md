# 前端验证指南

本文档说明如何验证前端项目可以正常运行。

## 验证步骤

### 1. 检查依赖安装

确保所有依赖已正确安装：

```bash
cd frontend
npm install
```

### 2. 检查 TypeScript 编译

运行 TypeScript 编译检查：

```bash
npm run build
```

**预期结果**：
- ✅ TypeScript 编译通过，无错误
- ✅ Vite 构建成功
- ✅ 生成 `dist/` 目录

**常见问题**：
- 如果出现 TypeScript 错误，检查 `tsconfig.json` 配置
- 如果出现模块找不到错误，运行 `npm install` 重新安装依赖

### 3. 启动开发服务器

启动开发服务器：

```bash
npm run dev
```

**预期结果**：
- ✅ 服务器启动成功
- ✅ 显示本地访问地址（通常是 `http://localhost:8000`）
- ✅ 浏览器可以访问页面

**验证点**：
1. 打开浏览器访问 `http://localhost:8000`
2. 页面应该正常加载，无控制台错误
3. 检查浏览器开发者工具（F12）：
   - Console 标签：无红色错误
   - Network 标签：资源加载成功（200 状态码）

### 4. 检查页面路由

验证各个路由页面是否可以访问：

- `/` - 首页（帖子列表）
- `/create` - 创建帖子页面
- `/post/:id` - 帖子详情页面
- `/search` - 搜索页面

**验证方法**：
1. 在浏览器地址栏输入路由地址
2. 检查页面是否正常渲染
3. 检查是否有路由错误（404）

### 5. 检查组件渲染

验证主要组件是否正常渲染：

- ✅ `AppHeader` - 顶部导航栏
- ✅ `HomePage` - 首页组件
- ✅ `CreatePostPage` - 创建页面组件
- ✅ `PostDetailPage` - 详情页面组件
- ✅ `SearchPage` - 搜索页面组件

**验证方法**：
1. 检查页面元素是否显示
2. 检查 Ant Design 组件样式是否正常
3. 检查是否有组件渲染错误

### 6. 检查配置

验证前端配置是否正确：

**检查文件**：`frontend/src/shared/config/index.ts`

```typescript
export const config = {
  apiBaseUrl: import.meta.env.VITE_API_BASE_URL || 'http://localhost:50051',
  grpcUrl: import.meta.env.VITE_GRPC_URL || 'http://localhost:50051',
}
```

**验证点**：
- ✅ 配置文件存在
- ✅ API 地址配置正确
- ✅ 环境变量可以正确读取

### 7. 检查构建产物

验证生产构建是否成功：

```bash
npm run build
npm run preview
```

**预期结果**：
- ✅ 构建成功，生成 `dist/` 目录
- ✅ `preview` 命令可以启动预览服务器
- ✅ 预览服务器可以正常访问

## 快速验证命令

### 一键验证脚本

创建验证脚本 `frontend/scripts/verify.sh`：

```bash
#!/bin/bash

set -e

echo "🔍 验证前端项目..."

echo "1. 检查依赖..."
if [ ! -d "node_modules" ]; then
  echo "❌ node_modules 不存在，运行 npm install..."
  npm install
else
  echo "✅ node_modules 存在"
fi

echo "2. 检查 TypeScript 编译..."
npm run build

echo "3. 检查开发服务器..."
echo "✅ 所有检查通过！"
echo ""
echo "运行 'npm run dev' 启动开发服务器"
```

### Makefile 命令

在项目根目录的 `Makefile` 中添加：

```makefile
# Frontend commands
frontend-install:
	cd frontend && npm install

frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build

frontend-preview:
	cd frontend && npm run preview

frontend-verify:
	cd frontend && npm run build && echo "✅ 前端构建验证通过"
```

## 验证清单

完成以下检查项，确保前端可以正常运行：

- [ ] **依赖安装**：`npm install` 成功，无错误
- [ ] **TypeScript 编译**：`npm run build` 成功，无类型错误
- [ ] **开发服务器**：`npm run dev` 启动成功
- [ ] **页面访问**：浏览器可以访问 `http://localhost:8000`
- [ ] **路由功能**：所有路由页面可以正常访问
- [ ] **组件渲染**：主要组件正常显示
- [ ] **控制台错误**：浏览器控制台无红色错误
- [ ] **资源加载**：所有静态资源加载成功
- [ ] **构建产物**：`npm run build` 生成 `dist/` 目录
- [ ] **配置正确**：API 地址等配置正确

## 常见问题

### 1. Node.js 版本错误

**问题**：
- `You are using Node.js 18.17.0. Vite requires Node.js version 20.19+ or 22.12+.`
- `TypeError: crypto.hash is not a function`

**原因**：Vite 7 需要 Node.js 20.19+ 或 22.12+，`crypto.hash` 是 Node.js 20+ 才有的 API。

**解决**：
1. 使用 nvm 升级 Node.js（推荐）：
   ```bash
   # 安装 Node.js 20 LTS
   nvm install 20
   
   # 切换到 Node.js 20
   nvm use 20
   
   # 验证版本
   node --version
   ```

2. 重新安装前端依赖：
   ```bash
   cd frontend
   rm -rf node_modules package-lock.json
   npm install
   ```

3. 详细步骤请参考 [Node.js 升级指南](./nodejs-upgrade.md)

### 2. 模块找不到错误

**问题**：`Cannot find module 'xxx'`

**解决**：
```bash
cd frontend
npm install
```

### 3. TypeScript 类型错误

**问题**：TypeScript 编译报错

**解决**：
- 检查 `tsconfig.json` 配置
- 检查导入路径是否正确
- 检查类型定义文件是否存在

### 4. 页面空白

**问题**：浏览器页面显示空白

**解决**：
- 检查浏览器控制台错误
- 检查路由配置是否正确
- 检查组件导入路径是否正确

### 5. ngrok 访问被阻止

**问题**：`Blocked request. This host ("xxx.ngrok-free.dev") is not allowed.`

**原因**：Vite 默认只允许 localhost 访问，通过 ngrok 等反向代理访问时需要配置允许的主机。

**解决**：
- 已在 `vite.config.ts` 中配置允许 ngrok 域名（`.ngrok-free.dev`, `.ngrok.io`, `.ngrok.app`）
- 如果使用其他反向代理，需要将域名添加到 `server.allowedHosts` 配置中
- 开发环境也可以设置 `host: true` 允许所有主机访问（已配置）

**配置说明**：
```typescript
server: {
  host: true, // 允许外部访问
  allowedHosts: [
    'localhost',
    '.ngrok-free.dev', // 允许所有 ngrok 域名
    '.ngrok.io',
    '.ngrok.app',
  ],
}
```

### 6. 样式不显示

**问题**：Ant Design 组件样式不显示

**解决**：
- 检查 `antd` 是否正确安装
- 检查 `ConfigProvider` 是否正确配置
- 检查 CSS 文件是否正确导入

## 下一步

前端验证通过后，可以继续：

1. **实现页面组件**：完成各个页面的具体实现
2. **集成 gRPC Web**：连接后端 gRPC 服务
3. **实现业务逻辑**：完成创建、列表、搜索等功能
4. **添加测试**：编写前端单元测试和 E2E 测试

## 相关文档

- [开发环境设置指南](./setup-guide.md)
- [开发指南](./development-guide.md)
- [测试指南](./testing-guide.md)

