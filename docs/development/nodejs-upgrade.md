# Node.js 升级指南

## 问题说明

Vite 7 需要 Node.js 20.19+ 或 22.12+，但当前使用的是 Node.js 18.17.0，导致开发服务器启动失败。

错误信息：
```
TypeError: crypto.hash is not a function
```

这是因为 `crypto.hash` 是 Node.js 20+ 才有的 API。

## 解决方案

### 使用 nvm 升级 Node.js（推荐）

如果你已经安装了 nvm，可以按以下步骤升级：

#### 1. 查看可用的 Node.js 版本

```bash
# 查看所有可用的 LTS 版本
nvm ls-remote --lts

# 或查看最新的 LTS 版本
nvm ls-remote --lts | tail -5
```

#### 2. 安装 Node.js 20 LTS（推荐）

```bash
# 安装最新的 Node.js 20 LTS 版本
nvm install 20

# 或安装特定版本
nvm install 20.19.0
```

#### 3. 使用新版本

```bash
# 切换到 Node.js 20
nvm use 20

# 验证版本
node --version
# 应该显示 v20.x.x
```

#### 4. 设置为默认版本（可选）

```bash
# 将 Node.js 20 设置为默认版本
nvm alias default 20
```

#### 5. 重新安装前端依赖

```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
```

#### 6. 验证前端可以运行

```bash
# 验证构建
npm run build

# 启动开发服务器
npm run dev
```

### 使用 Node.js 22（可选）

如果你想要使用最新的 Node.js 22：

```bash
# 安装 Node.js 22 LTS
nvm install 22

# 切换到 Node.js 22
nvm use 22

# 设置为默认版本
nvm alias default 22
```

### 不使用 nvm 的情况

如果你没有安装 nvm，可以：

1. **安装 nvm**（推荐）：
   ```bash
   # macOS/Linux
   curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
   
   # 重新加载 shell 配置
   source ~/.zshrc  # 或 ~/.bashrc
   ```

2. **直接从官网下载**：
   - 访问 [Node.js 官网](https://nodejs.org/)
   - 下载并安装 Node.js 20 LTS 或 22 LTS
   - 注意：这可能会覆盖系统默认的 Node.js 版本

## 验证

升级后，验证 Node.js 版本：

```bash
node --version
# 应该显示 v20.x.x 或 v22.x.x

npm --version
# 应该显示对应的 npm 版本
```

然后验证前端项目：

```bash
cd frontend
npm run build
npm run dev
```

## 项目配置

为了确保团队成员使用正确的 Node.js 版本，建议在项目根目录创建 `.nvmrc` 文件：

```bash
echo "20" > .nvmrc
```

这样，团队成员可以使用 `nvm use` 自动切换到正确的版本。

## 常见问题

### Q: 升级后 npm 包需要重新安装吗？

A: 是的，建议删除 `node_modules` 和 `package-lock.json` 后重新安装：

```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
```

### Q: 如何在不同项目间切换 Node.js 版本？

A: 使用 nvm：

```bash
# 切换到 Node.js 20
nvm use 20

# 切换到 Node.js 18（如果需要）
nvm use 18

# 查看已安装的版本
nvm ls
```

### Q: 升级后其他项目会受影响吗？

A: 如果使用 nvm，不会影响系统默认的 Node.js 版本。每个项目可以独立使用不同的 Node.js 版本。

## 相关文档

- [前端验证指南](./frontend-verification.md)
- [开发环境设置指南](./setup-guide.md)

