# Frontend E2E Tests

前端端到端测试，使用 Playwright 框架。

## 测试覆盖

### 1. 创建和查看流程 (`create-and-view.spec.ts`)
- 创建帖子并验证出现在列表中
- 导航到创建页面并验证表单
- 验证表单验证错误

### 2. 搜索流程 (`search.spec.ts`)
- 导航到搜索页面
- 执行搜索
- 按城市筛选搜索结果
- 从搜索结果导航到详情页
- 显示空状态

### 3. 导航流程 (`navigation.spec.ts`)
- 在主要页面之间导航
- 验证所有页面都显示 header
- 验证活动菜单项高亮

## 运行测试

### 安装依赖

```bash
npm install
```

### 安装 Playwright 浏览器

```bash
npx playwright install
```

### 运行所有 E2E 测试

```bash
npm run test:e2e
```

### 以 UI 模式运行（推荐用于调试）

```bash
npm run test:e2e:ui
```

### 以有头模式运行（可以看到浏览器）

```bash
npm run test:e2e:headed
```

### 调试模式

```bash
npm run test:e2e:debug
```

### 运行特定测试文件

```bash
npx playwright test test/e2e/flows/create-and-view.spec.ts
```

## 配置

测试配置在 `playwright.config.ts` 中：

- **baseURL**: 默认 `http://localhost:8000`（可通过 `PLAYWRIGHT_BASE_URL` 环境变量覆盖）
- **webServer**: 自动启动开发服务器（如果未运行）
- **浏览器**: 默认使用 Chromium（可配置 Firefox、WebKit）

## 环境变量

- `PLAYWRIGHT_BASE_URL`: 覆盖测试的基础 URL
- `CI`: 在 CI 环境中自动启用重试和并行限制

## 测试报告

测试运行后会生成 HTML 报告：

```bash
npx playwright show-report
```

## 注意事项

1. **后端服务**: 确保后端 gRPC 服务正在运行（`http://localhost:50051`）
2. **数据库**: 测试使用真实的后端服务，需要确保数据库和 Redis 可用
3. **测试数据**: 测试会创建真实的测试数据，可能需要定期清理
4. **网络超时**: 如果网络较慢，可能需要调整 `timeout` 配置

## 故障排除

### 测试失败：无法连接到服务器

确保前端开发服务器正在运行：
```bash
npm run dev
```

### 测试失败：gRPC 错误

确保后端服务正在运行：
```bash
cd ../backend
go run cmd/server/main.go
```

### 浏览器未安装

运行：
```bash
npx playwright install
```

