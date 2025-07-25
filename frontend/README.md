# Go-Web 前端项目

本项目是 Go-Web 应用的前端部分，经过精心重构，采用 Vue 3、Vite、Element Plus 和 Tailwind CSS 构建。它提供了一个美观、现代且响应式的用户界面，包括一个全新设计的登录页和功能完善的仪表盘。

---

### ✨ **视觉预览**

*（此处为新版登录页面的效果图占位符。想象一个拥有动态渐变背景、精致的毛玻璃质感卡片和优雅 UI 组件的登录界面。）*

![登录界面预览](placeholder.png "新版登录界面预览图")

---

## 🚀 核心功能

- **极致的 UI/UX**: 专注于美学和用户体验的全面重新设计。
- **动态登录页**: 沉浸式的登录界面，具有流畅的渐变背景和优雅的毛玻璃效果。
- **功能仪表盘**: 为已认证用户设计的经典后台布局，包含侧边栏、顶部用户信息和主内容区。
- **响应式设计**: 基于 Tailwind CSS 构建，确保在桌面和移动设备上都有一致的、出色的视觉体验。
- **状态管理**: 使用 Pinia 进行集中式状态管理，高效处理认证和用户信息。
- **类型安全**: 项目完全使用 TypeScript 编写，保证了代码的健壮性和可维护性。
- **丰富的图标**: 集成 `@element-plus/icons-vue`，提供丰富的 UI 图标选择。

---

## 🛠️ 技术栈

- **核心框架**: [Vue 3](https://cn.vuejs.org/) (使用组合式 API)
- **构建工具**: [Vite](https://cn.vitejs.dev/)
- **UI 组件库**: [Element Plus](https://element-plus.org/zh-CN/)
- **CSS 方案**: [Tailwind CSS](https://www.tailwindcss.cn/)
- **状态管理**: [Pinia](https://pinia.vuejs.org/zh/)
- **路由**: [Vue Router](https://router.vuejs.org/zh/)
- **开发语言**: [TypeScript](https://www.typescriptlang.org/)
- **图标库**: [@element-plus/icons-vue](https://element-plus.org/zh-CN/component/icon.html)

---

## ⚙️ 本地开发指南

请遵循以下步骤，在您的本地机器上启动项目用于开发和测试。

### 环境要求

- [Node.js](https://nodejs.org/) (建议版本 16.x 或更高)
- [npm](https://www.npmjs.com/) 或其他包管理器，如 [pnpm](https://pnpm.io/zh/)、[yarn](https://yarnpkg.com/zh-Hans/)

### 安装与配置

1.  **克隆代码库** (如果尚未克隆):
    ```sh
    git clone <repository-url>
    ```

2.  **进入前端项目目录**:
    ```sh
    cd go-web/frontend
    ```

3.  **安装依赖**:
    ```sh
    npm install
    ```

4.  **API 代理配置**:
    为了解决开发环境下的跨域问题，Vite 已配置了 API 代理。默认会将 `/api` 前缀的请求转发至 `http://localhost:8080`。如果您的后端服务运行在不同地址，请修改 `vite.config.ts` 文件中的 `proxy` 配置。

### 启动项目

- **开发模式**:
  执行以下命令以启动 Vite 开发服务器，该模式支持热更新。
  ```sh
  npm run dev
  ```
  应用将在 `http://localhost:5173` (或下一个可用端口) 上运行。

- **生产环境构建**:
  如需为生产环境打包，请运行：
  ```sh
  npm run build
  ```
  优化后的静态文件将生成在 `dist/` 目录中。

---

## 📁 目录结构详解

```
frontend/
├── public/              # 无需 Vite 处理的公共静态资源
├── src/
│   ├── assets/          # 样式、字体等静态资源
│   ├── components/      # 可复用的 Vue 组件
│   ├── router/          # Vue Router 路由配置 (index.ts)
│   ├── services/        # API 服务调用 (authService.ts 等)
│   ├── store/           # Pinia 状态管理模块 (authStore.ts)
│   ├── types/           # TypeScript 类型定义
│   ├── utils/           # 工具函数 (如 http 客户端封装)
│   ├── views/           # 页面级组件 (LoginView.vue, DashboardView.vue)
│   ├── App.vue          # 应用根组件
│   └── main.ts          # 应用主入口文件
├── .gitignore           # Git 忽略配置文件
├── index.html           # Vite 的主 HTML 入口文件
├── package.json         # 项目依赖和脚本定义
├── postcss.config.js    # PostCSS 配置文件
├── tailwind.config.js   # Tailwind CSS 配置文件
├── tsconfig.json        # TypeScript 编译器配置
└── vite.config.ts       # Vite 配置文件
```
