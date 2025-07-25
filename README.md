# Go Web 开发教学项目

本项目是一个使用 Go 语言进行 Web 开发的完整示例，整合了现代 Web 开发中常用的技术栈：

- Web 框架：[Gin](https://gin-gonic.com/)
- ORM 框架：[GORM](https://gorm.io/)
- 数据库：[PostgreSQL](https://www.postgresql.org/)
- 认证授权：[JWT](https://jwt.io/) + [Casbin](https://casbin.org/)
- 配置管理：[Viper](https://github.com/spf13/viper)
- 日志系统：[Zap](https://github.com/uber-go/zap)
- 缓存系统：[Redis](https://redis.io/)
- API文档：[Swagger](https://swagger.io/)
- 前端框架：[Vue 3](https://v3.vuejs.org/) + [TypeScript](https://www.typescriptlang.org/)

## 目录结构

```
go-web/
├── README.md
├── go.mod
├── go.sum
├── main.go
├── config/
│   ├── config.go
│   └── config.yaml
├── models/
│   ├── user.go
│   └── role.go
├── database/
│   └── database.go
├── cache/
│   └── redis.go
├── middleware/
│   ├── jwt.go
│   ├── casbin.go
│   ├── logging.go
│   └── cors.go
├── utils/
│   ├── jwt.go
│   └── logger.go
├── controllers/
│   ├── auth.go
│   └── user.go
├── services/
│   ├── user.go
│   └── auth.go
├── repositories/
│   ├── user.go
│   └── role.go
├── dtos/
│   ├── user.go
│   └── auth.go
├── validators/
│   └── user.go
├── docs/
│   └── swagger.yaml
├── routers/
│   └── router.go
├── logs/
│   └── app.log
└── frontend/
    ├── public/
    ├── src/
    │   ├── assets/
    │   ├── components/
    │   ├── views/
    │   │   ├── LoginView.vue
    │   │   └── DashboardView.vue
    │   ├── router/
    │   │   └── index.ts
    │   ├── store/
    │   │   ├── index.ts
    │   │   └── authStore.ts
    │   ├── services/
    │   │   ├── authService.ts
    │   │   └── userService.ts
    │   ├── utils/
    │   │   └── http.ts
    │   ├── types/
    │   │   └── user.ts
    │   ├── App.vue
    │   └── main.ts
    ├── index.html
    ├── package.json
    ├── vite.config.ts
    ├── tsconfig.json
    ├── tsconfig.node.json
    ├── .env
    └── .gitignore
```

## 技术栈详解

### Gin Web 框架
Gin 是一个用 Go 编写的 HTTP web 框架。它具有高性能、生产力高的特点，是目前最流行的 Go web 框架之一。

### GORM ORM 框架
GORM 是一个强大的 Go 语言 ORM 库，支持多种数据库，包括 PostgreSQL、MySQL、SQLite 等。它提供了链式 API，使数据库操作变得简单直观。

### PostgreSQL 数据库
PostgreSQL 是一个功能强大的开源关系型数据库系统，以其可靠性、功能 robustness 和性能而闻名。

### JWT 身份验证
JWT (JSON Web Token) 是一种开放标准(RFC 7519)，用于在各方之间安全地传输信息。在 Web 应用中，它常用于身份验证和信息交换。

### Casbin 访问控制
Casbin 是一个强大的、高效的开源访问控制库，支持各种访问控制模型，包括 ACL、RBAC、ABAC 等。

### Viper 配置管理
Viper 是一个完整的 Go 配置解决方案，支持多种配置格式（JSON、TOML、YAML、HCL、envfile 和 Java properties），可以处理所有类型的配置需求和格式。

### Zap 日志系统
Zap 是 Uber 开发的快速、结构化的日志库，提供了高性能的日志记录功能。

### Redis 缓存系统
Redis 是一个开源的内存数据结构存储系统，常用于缓存、会话存储等场景。

### Swagger API 文档
Swagger 是一套开放源代码的软件工具，用于设计、构建、记录和使用 RESTful Web 服务。

### Vue 3 + TypeScript 前端
Vue 3 是流行的前端框架，结合 TypeScript 提供了类型安全和更好的开发体验。

## 项目功能

1. 用户注册和登录
2. JWT Token 生成和验证
3. 基于角色的访问控制 (RBAC)
4. 用户管理接口
5. 权限管理
6. 灵活的配置管理
7. 结构化日志记录
8. Redis缓存支持
9. API文档自动生成
10. 数据验证
11. 错误处理和统一响应格式
12. 跨域支持
13. 前端登录和仪表板界面

## 安装和运行

### 后端环境要求

- Go 1.24 或更高版本
- PostgreSQL 数据库
- Redis (可选，用于缓存功能)

### 前端环境要求

- Node.js 16+ 或更高版本
- npm 或 yarn 包管理器

### 后端安装步骤

1. 克隆项目到本地：
   ```
   git clone <repository-url>
   ```

2. 安装依赖：
   ```
   go mod tidy
   ```

3. 配置数据库：
   在 `config/config.yaml` 中配置数据库连接信息

4. 运行后端项目：
   ```
   go run main.go
   ```

### 前端安装步骤

1. 进入前端目录：
   ```
   cd frontend
   ```

2. 安装依赖：
   ```
   npm install
   ```

3. 运行前端开发服务器：
   ```
   npm run dev
   ```

4. 访问应用：
   打开浏览器访问 http://localhost:3000

## API 接口说明

### 认证相关接口

- `POST /auth/register` - 用户注册
- `POST /auth/login` - 用户登录
- `POST /auth/logout` - 用户登出

### 用户相关接口

- `GET /users` - 获取用户列表 (需要管理员权限)
- `GET /users/:id` - 获取用户详情 (需要相应权限)
- `PUT /users/:id` - 更新用户信息 (需要相应权限)
- `DELETE /users/:id` - 删除用户 (需要管理员权限)

## 完整后端架构组件说明

### 1. 配置管理层 (Config)
- 配置文件管理
- 环境变量支持
- 配置热加载

### 2. 数据访问层 (Database/Repository)
- 数据库连接管理
- ORM封装
- 数据访问对象(DAO)

### 3. 缓存层 (Cache)
- Redis/Memcached集成
- 缓存策略实现
- 分布式缓存支持

### 4. 业务逻辑层 (Service)
- 核心业务逻辑实现
- 事务管理
- 业务规则封装

### 5. 控制器层 (Controller)
- HTTP请求处理
- 参数验证
- 响应格式化

### 6. 中间件层 (Middleware)
- 认证授权
- 日志记录
- 请求追踪
- 限流控制
- 跨域处理

### 7. 工具类层 (Utils)
- 加密解密
- 字符串处理
- 时间处理
- 文件操作

### 8. 数据传输层 (DTO)
- 请求数据封装
- 响应数据封装
- 数据转换

### 9. 验证层 (Validator)
- 输入数据验证
- 业务规则验证

### 10. 文档层 (Docs)
- API文档
- 技术文档
- 使用说明

### 11. 日志层 (Logs)
- 系统日志记录
- 错误日志追踪
- 性能日志分析

## 前端项目说明

### 技术栈
- Vue 3 Composition API
- TypeScript
- Vue Router
- Pinia (状态管理)
- Axios (HTTP客户端)
- Vite (构建工具)

### 功能特性
1. 用户登录界面
2. 仪表板界面
3. 路由权限控制
4. JWT Token管理
5. 响应式设计

### 目录结构
```
frontend/
├── public/              # 静态资源
├── src/
│   ├── assets/          # 样式和静态资源
│   ├── components/      # 可复用组件
│   ├── views/           # 页面视图
│   ├── router/          # 路由配置
│   ├── store/           # 状态管理 (Pinia)
│   ├── services/        # 业务服务层
│   ├── utils/           # 工具函数
│   ├── types/           # TypeScript类型定义
│   ├── App.vue          # 根组件
│   └── main.ts          # 入口文件
├── index.html           # HTML模板
├── package.json         # 项目依赖
├── vite.config.ts       # Vite配置
├── tsconfig.json        # TypeScript配置
├── .env                 # 环境变量配置
└── .gitignore           # Git忽略文件配置
```

### 核心组件说明

#### 1. 路由 (Router)
- 使用 Vue Router 实现前端路由
- 包含登录和仪表板两个主要页面
- 实现路由守卫进行权限控制

#### 2. 状态管理 (Store)
- 使用 Pinia 进行状态管理
- 实现认证状态管理 (authStore)
- 统一管理用户信息和认证状态

#### 3. 服务层 (Services)
- authService: 处理认证相关业务逻辑
- userService: 处理用户相关业务逻辑
- 封装了与后端API的交互

#### 4. 工具函数 (Utils)
- http.ts: 封装 Axios 实例，统一处理请求和响应
- 包含请求拦截器和响应拦截器
- 自动添加认证Token

#### 5. 类型定义 (Types)
- 定义了用户、认证等相关 TypeScript 类型
- 提供类型安全保证

### 开发环境代理配置
前端开发服务器配置了代理，将 `/api` 请求转发到后端服务器：
```typescript
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

### 前端功能实现

#### 登录页面 (LoginView.vue)
- 用户名和密码输入表单
- 表单验证和错误提示
- 调用认证服务进行登录
- 登录成功后跳转到仪表板

#### 仪表板页面 (DashboardView.vue)
- 显示当前用户信息
- 提供登出功能
- 路由权限验证

#### 状态管理 (authStore.ts)
- 管理用户认证状态
- 封装登录和登出逻辑
- 持久化存储Token和用户信息

## 学习路径

1. 理解项目结构和各组件作用
2. 学习如何使用 Gin 构建 RESTful API
3. 掌握 GORM 进行数据库操作
4. 实现 JWT 身份验证机制
5. 学习使用 Casbin 进行权限控制
6. 使用 Viper 进行配置管理
7. 集成 Zap 日志系统
8. 实现 Redis 缓存功能
9. 使用 Swagger 生成 API 文档
10. 了解中间件的使用方式
11. 掌握错误处理和统一响应格式
12. 实现数据验证
13. 学习 Vue 3 + TypeScript 前端开发
14. 实现前后端交互

## 代码实现步骤

### 第一步：项目初始化和依赖安装

初始化项目并安装所需依赖：

```bash
# 后端依赖
go mod init go-web
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go get -u github.com/golang-jwt/jwt/v4
go get -u github.com/casbin/casbin/v2
go get -u github.com/spf13/viper
go get -u go.uber.org/zap
go get -u github.com/go-redis/redis/v8
go get -u github.com/swaggo/swag/cmd/swag

# 前端依赖
npm install vue@next vue-router@next pinia axios
npm install -D typescript @vitejs/plugin-vue @vitejs/plugin-vue-jsx vite vue-tsc
```

### 第二步：配置数据库连接

创建数据库配置文件，使用 Viper 管理配置，使用 GORM 连接 PostgreSQL。

### 第三步：定义数据模型

定义用户、角色等数据模型，并使用 GORM 进行数据库迁移。

### 第四步：实现 JWT 认证

创建 JWT 工具函数，实现登录验证和 Token 生成。

### 第五步：集成 Casbin 权限控制

配置 Casbin，定义权限模型和策略（使用 YAML 格式），实现基于角色的访问控制。

### 第六步：构建 API 接口

使用 Gin 构建 RESTful API，实现用户注册、登录、管理等功能。

### 第七步：添加中间件

创建认证和权限中间件，保护 API 接口。

### 第八步：集成日志系统

使用 Zap 实现结构化日志记录。

### 第九步：添加缓存支持

集成 Redis 实现缓存功能。

### 第十步：生成 API 文档

使用 Swagger 生成 API 文档。

### 第十一步：开发前端界面

使用 Vue 3 + TypeScript 开发登录和仪表板界面。

## 扩展学习

1. 添加单元测试
2. 实现分布式追踪 (如 Jaeger)
3. 添加消息队列 (如 RabbitMQ, Kafka)
4. 实现微服务架构
5. 添加容器化部署 (Docker, Kubernetes)
6. 实现负载均衡
7. 添加监控和告警 (Prometheus, Grafana)
8. 实现 CI/CD 流程
9. 添加健康检查和探针
10. 实现数据备份和恢复
11. 添加国际化支持
12. 实现服务端渲染(SSR)

通过本项目，你将掌握使用 Go 进行现代 Web 开发的核心技能，以及 Vue 3 + TypeScript 前端开发技能，为进一步开发更复杂的应用打下坚实基础。