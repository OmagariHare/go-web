package routers

import (
	"go-web/config"
	"go-web/controllers"
	"go-web/database"
	"go-web/middleware"
	"go-web/repositories"
	"go-web/services"
	"go-web/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Add the error handling middleware as the first middleware.
	r.Use(middleware.ErrorHandler())

	// 添加CORS中间件
	r.Use(middleware.CORSMiddleware(cfg))

	// 添加日志中间件
	r.Use(middleware.LoggingMiddleware(utils.Logger))

	// 初始化Casbin
	err := middleware.InitCasbin(cfg)
	if err != nil {
		panic("Failed to initialize Casbin: " + err.Error())
	}

	// 创建数据库连接
	db := database.DB

	// 创建仓储实例
	userRepository := repositories.NewGormUserRepository(db)
	roleRepository := repositories.NewGormRoleRepository(db)

	// 创建服务实例
	authService := services.NewAuthService(cfg, userRepository, roleRepository)
	userService := services.NewUserService(userRepository)

	// 创建控制器实例
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)

	// Public routes (no authentication required)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	auth := r.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	// 受保护的路由（需要认证和授权）
	users := r.Group("/users")
	users.Use(middleware.AuthMiddleware(cfg))
	users.Use(middleware.CasbinMiddleware())
	{
		users.GET("/", userController.GetUsers)
		users.GET("/:id", userController.GetUser)
		users.PUT("/:id", userController.UpdateUser)
		users.DELETE("/:id", userController.DeleteUser)
	}

	return r
}
