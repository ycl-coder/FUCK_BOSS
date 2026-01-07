# Frontend E2E Tests

前端端到端测试，使用 Playwright 框架。

## 安装

### 1. 安装依赖

```bash
npm install
```

### 2. 安装 Playwright 浏览器

```bash
npx playwright install chromium
```

或者安装所有浏览器：

```bash
npx playwright install
```

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

## 前置条件

### 1. Node.js 版本

确保使用 Node.js 20+：

```bash
# 使用 nvm
nvm use

# 或手动切换
nvm use 20
```

### 2. 后端服务

确保后端 gRPC 服务正在运行：

```bash
cd ../backend
go run cmd/server/main.go
```

或者使用 Docker Compose：

```bash
make docker-up
```

### 3. 数据库和 Redis

确保 PostgreSQL 和 Redis 服务可用（通过 Docker Compose 或本地安装）。

## 注意事项

1. **测试数据**: 测试会创建真实的测试数据，可能需要定期清理
2. **网络超时**: 如果网络较慢，可能需要调整 `timeout` 配置
3. **并发测试**: 默认并行运行测试，如果遇到资源问题，可以在 `playwright.config.ts` 中调整 `workers` 数量

## 故障排除

### 错误：`playwright: command not found`

**原因**: Playwright 依赖未安装或浏览器未安装

**解决方案**:
```bash
# 1. 安装依赖
npm install

# 2. 安装浏览器
npx playwright install chromium
```

### 错误：无法连接到服务器

确保前端开发服务器正在运行：
```bash
npm run dev
```

### 错误：gRPC 错误

确保后端服务正在运行：
```bash
cd ../backend
go run cmd/server/main.go
```

### 错误：Node.js 版本不兼容

升级到 Node.js 20+：
```bash
nvm use 20
# 或
nvm install 20 && nvm use 20
```

### 浏览器未安装

运行：
```bash
npx playwright install chromium
```
