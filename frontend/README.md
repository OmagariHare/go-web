# Go Web 前端项目

这是一个基于 Vue 3 + TypeScript 的前端项目，用于与 Go 后端进行交互。

## 技术栈

- Vue 3 (Composition API)
- TypeScript
- Vite
- Pinia (状态管理)
- Vue Router
- Axios

## 安装依赖

```bash
npm install
```

## 开发环境

```bash
npm run dev
```

这将启动开发服务器，通常在 http://localhost:3000 上访问。

## 构建生产版本

```bash
npm run build
```

## 项目结构

```
src/
├── assets/          # 静态资源
├── components/      # 可复用组件
├── views/           # 页面视图
├── router/          # 路由配置
├── store/           # 状态管理
├── services/        # 业务服务
├── utils/           # 工具函数
├── types/           # TypeScript 类型定义
├── App.vue          # 根组件
└── main.ts          # 入口文件
```

## 功能说明

1. 登录页面 - 用户可以通过用户名和密码登录
2. 仪表板页面 - 显示用户信息和登出功能
3. 路由保护 - 未登录用户无法访问仪表板页面
4. 状态管理 - 使用 Pinia 管理用户认证状态

## API 代理配置

Vite 配置了代理，将 `/api` 开头的请求代理到后端服务器：

```javascript
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api/, '')
    }
  }
}
```

请确保后端服务在 `http://localhost:8080` 上运行。