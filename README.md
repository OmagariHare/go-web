# Go-Web 全栈开发项目

欢迎来到 Go-Web 项目！这是一个现代化的全栈 Web 应用，后端采用 Go 语言构建，前端则由 Vue.js 驱动。项目结构经过精心设计，将前后端代码清晰地分离到各自的目录中，以实现更好的模块化和可维护性。

---

## 🚀 项目概览

- **后端 (Backend)**: 一个强大的 RESTful API 服务，基于 Go、Gin、GORM 和 Casbin 构建，为整个应用提供安全、高效的数据支持和权限管理。

- **前端 (Frontend)**: 一个美观且响应迅速的用户界面，采用 Vue 3、Vite、Element Plus 和 Tailwind CSS 重构，旨在提供卓越的用户体验。

---

## 🛠️ 技术栈概览

| 技术领域 | 主要技术与框架                                           |
| :------- | :------------------------------------------------------- |
| **后端**   | Go, Gin, PostgreSQL, GORM, JWT, Casbin, Viper, Zap, Redis  |
| **前端**   | Vue 3, Vite, TypeScript, Element Plus, Tailwind CSS, Pinia |

---

## 📁 项目结构

项目顶层目录结构清晰地划分了后端和前端的开发环境：

```
go-web/
├── backend/         # 包含所有 Go 后端应用代码
│   └── README.md    # 后端项目的详细中文文档
├── frontend/        # 包含所有 Vue.js 前端应用代码
│   └── README.md    # 前端项目的详细中文文档
└── README.md        # 当前的项目总览 README
```

---

## ⚡ 快速启动指南

要完整地运行此应用，您需要分别启动后端和前端两个服务。

### 1. 启动后端服务

首先，进入后端目录，并遵循其文档中的详细步骤进行设置和启动。

```sh
cd backend
# 请参考 backend/README.md 中的详细说明
go run main.go
```

### 2. 启动前端服务

在 **另一个新的终端** 中，进入前端目录，并启动 Vite 开发服务器。

```sh
cd frontend
# 请参考 frontend/README.md 中的详细说明
npm run dev
```

当两个服务都成功运行后，您就可以在浏览器中通过访问 `http://localhost:5173` (或终端提示的其他端口) 来查看和使用该应用了。

---

## 📚 详细文档

关于项目各部分的详细信息（包括环境配置、API 端点、目录结构详解等），请参阅其对应目录下的 README 文件：

- **[后端文档](./backend/README.md)**: 指导您如何配置数据库、设置环境变量以及运行 Go 服务。
- **[前端文档](./frontend/README.md)**: 指导您如何安装前端依赖以及启动本地开发服务器。

---

## 🤝 贡献代码

欢迎任何形式的贡献！如果您有任何改进建议或发现了问题，请随时提交 Pull Request 或创建 Issue 进行讨论。

## 📄 开源许可证

本项目基于 MIT 许可证开源。
