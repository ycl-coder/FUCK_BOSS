# 全国公司曝光平台

一个允许用户分城市匿名发布公司不当行为的平台。

## 项目结构

```
fuck_boss/
├── backend/          # 后端（Go + gRPC）
├── frontend/         # 前端（React）
├── docs/             # 文档
└── .spec-workflow/   # 规范工作流文档
```

## 技术栈

- **后端**: Go 1.21+ + gRPC-Go
- **前端**: React + TypeScript
- **数据库**: PostgreSQL 14+
- **缓存**: Redis 7.0+
- **架构**: DDD（领域驱动设计）

## 快速开始

详见各子目录的 README：
- [后端开发指南](backend/README.md)
- [前端开发指南](frontend/README.md)（待创建）
- [部署文档](docs/deployment/)（待创建）

## 开发文档

- [产品指导](.spec-workflow/steering/product.md)
- [技术栈](.spec-workflow/steering/tech.md)
- [项目结构](.spec-workflow/steering/structure.md)
- [功能规范](.spec-workflow/specs/content-management-v1/)

## 许可证

（待定）

