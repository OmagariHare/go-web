# Go-Web 后端服务

本项目是 Go-Web 应用的后端服务，基于 Go 语言和 Gin 框架构建。它为前端应用提供了一个安全、高性能且可扩展的 RESTful API。

---

## 🚀 核心功能

- **高性能 API**: 基于 [Gin](https://gin-gonic.com/) 框架，提供卓越的 HTTP 处理性能。
- **JWT 认证**: 使用 JSON Web Tokens (JWT) 实现无状态、安全的用户认证机制。
- **Casbin 权限控制**: 集成 [Casbin](https://casbin.org/) 实现基于角色的访问控制 (RBAC)，支持精细化的权限管理。
- **ORM 操作**: 采用 [GORM](https://gorm.io/) 作为 ORM 框架，简化数据库操作，默认支持 PostgreSQL。
- **Viper 配置管理**: 通过 [Viper](https://github.com/spf13/viper) 管理应用配置，支持从 `config.yaml` 文件中读取，灵活易用。
- **结构化日志**: 使用 [Zap](https://github.com/uber-go/zap) 和 `lumberjack` 实现高性能的结构化日志记录与自动分割归档。
- **中间件架构**: 清晰的中间件设计，用于处理 CORS、JWT 验证、日志记录和权限检查。
- **Redis 缓存**: 集成 [Go-Redis](https://github.com/go-redis/redis) 提供高速缓存支持。
- **Swagger API 文档**: 内置 Swagger 支持，方便生成和浏览 API 文档。

---

## 🛠️ 技术栈

- **核心框架**: [Gin](https://github.com/gin-gonic/gin)
- **数据库**: [PostgreSQL](https://www.postgresql.org/)
- **ORM**: [GORM](https://gorm.io/)
- **认证**: [JWT-Go](https://github.com/golang-jwt/jwt)
- **授权**: [Casbin](https://github.com/casbin/casbin)
- **配置**: [Viper](https://github.com/spf13/viper)
- **日志**: [Zap](https://github.com/uber-go/zap)
- **缓存**: [Go-Redis](https://github.com/go-redis/redis)
- **API 文档**: [Swagger](https://swagger.io/)

---

## ⚙️ 本地开发指南

请遵循以下步骤在本地环境中设置并运行后端服务。

### 环境要求

- [Go](https://golang.org/dl/) (建议版本 1.18+)
- [PostgreSQL](https://www.postgresql.org/download/) 数据库服务
- [Redis](https://redis.io/download) 服务 (可选)

### 安装与配置

1.  **进入后端目录**:
    ```sh
    cd backend
    ```

2.  **安装依赖**:
    此命令将下载并安装 `go.mod` 文件中定义的所有依赖项。
    ```sh
    go mod tidy
    ```

3.  **配置应用**:
    - 复制或重命名 `config/config.yaml.example` 为 `config/config.yaml`。
    - 打开 `config/config.yaml` 文件并更新以下配置：
      - `database`: 你的 PostgreSQL 数据库连接信息 (主机、端口、用户名、密码、数据库名)。
      - `redis`: 你的 Redis 连接信息。
      - `jwt`: 用于签发 JWT 的密钥 (secret)。

4.  **初始化数据库**:
    应用在启动时会自动执行 GORM 的 `AutoMigrate`，以根据 `models` 目录下的模型创建或更新数据库表结构。

### 启动服务

- **运行应用**:
  ```sh
  go run main.go
  ```
  服务启动后，将监听在 `config.yaml` 中配置的端口（默认为 `8080`）。

---

## 📄 API 文档

项目使用 Swagger 生成 API 文档。

- **Swagger 定义文件**: API 的 OpenAPI 规范定义在 `docs/swagger.yaml`。
- **在线文档**: 服务运行后，可以通过访问 `http://localhost:8080/swagger/index.html` (需在代码中配置) 来查看和测试 API。

---

## 📁 目录结构详解

```
backend/
├── cache/         # Redis 缓存相关逻辑
├── config/        # Viper 配置初始化及 config.yaml 配置文件
├── controllers/   # Gin 的控制器层 (Handler)，处理 HTTP 请求
├── database/      # 数据库连接 (GORM) 初始化
├── docs/          # Swagger API 文档文件
├── dtos/          # 数据传输对象 (Data Transfer Objects)，用于 API 的请求和响应
├── logs/          # 由 Zap 生成的日志文件存放目录
├── middleware/    # Gin 中间件 (JWT, CORS, Casbin, 日志等)
├── models/        # GORM 的数据库模型定义
├── repositories/  # 数据仓库层，封装数据库的增删改查操作
├── routers/       # API 路由定义
├── services/      # 核心业务逻辑层
├── utils/         # 通用工具函数 (如密码处理, JWT 生成)
├── validators/    # 请求参数校验器
├── go.mod         # Go 模块依赖文件
├── go.sum         # 依赖项的校验和
└── main.go        # 应用主入口
```