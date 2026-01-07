# Ngrok 配置指南

本文档说明如何配置 ngrok 只暴露前端端口，并通过代理转发 API 请求到本地后端。

## 场景说明

- **前端**: 通过 ngrok 暴露（例如：`https://xxx.ngrok-free.dev`）
- **后端**: 在本地运行（`localhost:50051`），不通过 ngrok 暴露
- **目标**: 前端通过代理将 API 请求转发到本地后端

## 配置步骤

### 1. 启动本地后端服务

```bash
cd backend
go run cmd/server/main.go
```

后端将在 `http://localhost:50051` 运行。

### 2. 启动前端开发服务器

```bash
cd frontend
npm run dev
```

前端将在 `http://localhost:8000` 运行。

### 3. 使用 ngrok 暴露前端端口

```bash
ngrok http 8000
```

ngrok 会提供一个公网 URL，例如：`https://unsenescent-didactically-madison.ngrok-free.dev`

### 4. 配置说明

#### 开发环境（Vite）

前端配置已经设置好：

- **`frontend/vite.config.ts`**: 
  - `proxy` 配置会将 `/api` 请求转发到 `http://localhost:50051`
  - 当用户通过 ngrok URL 访问前端时，API 请求会通过 Vite proxy 转发到本地后端

- **`frontend/src/shared/config/index.ts`**:
  - 默认使用空字符串（相对路径）
  - API 请求会使用当前域名（ngrok URL），然后通过 Vite proxy 转发

#### 生产环境（Docker + Nginx）

如果使用 Docker 运行前端：

- **`frontend/nginx.conf`**:
  - 配置了 `/api` 代理到 `host.docker.internal:50051`
  - 这允许 Docker 容器访问宿主机上的后端服务

- **`docker-compose.yml`**:
  - 添加了 `extra_hosts` 配置，使 `host.docker.internal` 指向宿主机
  - 前端构建时使用空字符串作为 API URL，使用相对路径

## 工作原理

### 开发环境流程

1. 用户访问：`https://xxx.ngrok-free.dev`
2. 前端页面加载（通过 ngrok）
3. 前端发起 API 请求：`https://xxx.ngrok-free.dev/api/posts`
4. Vite proxy 拦截 `/api` 请求
5. Vite proxy 转发到：`http://localhost:50051/api/posts`
6. 后端处理请求并返回响应
7. 响应通过 Vite proxy 返回给前端

### 生产环境流程（Docker）

1. 用户访问：`https://xxx.ngrok-free.dev`
2. Nginx 提供前端静态文件
3. 前端发起 API 请求：`https://xxx.ngrok-free.dev/api/posts`
4. Nginx 拦截 `/api` 请求
5. Nginx 转发到：`http://host.docker.internal:50051/api/posts`
6. 后端处理请求并返回响应
7. 响应通过 Nginx 返回给前端

## 验证配置

### 1. 检查 Vite proxy 是否工作

打开浏览器开发者工具，查看 Network 标签：
- API 请求应该显示为 `https://xxx.ngrok-free.dev/api/posts`
- 请求应该成功返回数据

### 2. 检查后端日志

后端应该收到来自 `localhost` 的请求（通过 Vite proxy）。

### 3. 常见问题

#### 问题：API 请求失败，CORS 错误

**解决方案**: 
- 确保后端 CORS 配置允许 ngrok 域名
- 检查 `backend/internal/presentation/middleware/cors.go`

#### 问题：Docker 中无法访问宿主机后端

**解决方案**:
- 确保 `docker-compose.yml` 中配置了 `extra_hosts`
- Linux 系统可能需要使用 `172.17.0.1` 而不是 `host.docker.internal`

#### 问题：ngrok 显示 "Blocked request"

**解决方案**:
- 在 `vite.config.ts` 中添加 ngrok 域名到 `allowedHosts`
- 已经配置了 `.ngrok-free.dev`、`.ngrok.io`、`.ngrok.app`

## 环境变量

如果需要覆盖默认配置，可以创建 `.env` 文件：

```env
# frontend/.env
# 使用相对路径（推荐，通过代理转发）
VITE_GRPC_URL=
VITE_API_BASE_URL=

# 或直接指定后端 URL（不推荐，会绕过代理）
# VITE_GRPC_URL=http://localhost:50051
# VITE_API_BASE_URL=http://localhost:50051
```

## 注意事项

1. **ngrok URL 会变化**: 每次重启 ngrok，URL 可能会变化，需要更新配置
2. **后端必须在本地运行**: 确保后端服务在 `localhost:50051` 运行
3. **防火墙**: 确保本地防火墙允许 ngrok 访问
4. **HTTPS**: ngrok 提供 HTTPS，但后端可能使用 HTTP，这是正常的（通过代理转发）

