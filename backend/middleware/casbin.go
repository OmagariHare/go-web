package middleware

import (
	"go-web/config"
	"go-web/database"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Enforcer *casbin.Enforcer

func InitCasbin(cfg *config.Config) error {
	// 从配置中读取模型定义
	modelConf := cfg.Casbin.Model

	// 处理YAML中的多行字符串
	modelConf = strings.ReplaceAll(modelConf, "\\n", "\n")

	// 定义Casbin模型
	m, err := model.NewModelFromString(modelConf)
	if err != nil {
		return err
	}

	// 创建GORM适配器
	a, err := gormadapter.NewAdapterByDB(database.DB)
	if err != nil {
		return err
	}

	// 创建Casbin执行器
	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		return err
	}

	// 加载策略
	if err := e.LoadPolicy(); err != nil {
		return err
	}

	// 添加默认策略（如果策略为空）
	policies, err := e.GetPolicy()
	if err != nil {
		return err
	}
	if len(policies) == 0 {
		if _, err := e.AddPolicy("admin", "*", "*"); err != nil {
			return err
		}
		if _, err := e.AddPolicy("user", "/users/:id", "GET"); err != nil {
			return err
		}
		if _, err := e.AddPolicy("user", "/users/:id", "PUT"); err != nil {
			return err
		}
		if _, err := e.AddPolicy("anonymous", "/auth/register", "POST"); err != nil {
			return err
		}
		if _, err := e.AddPolicy("anonymous", "/auth/login", "POST"); err != nil {
			return err
		}
	}

	Enforcer = e
	return nil
}

func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use the global enforcer
		middleware := CasbinMiddlewareWithEnforcer(Enforcer)
		middleware(c)
	}
}

// CasbinMiddlewareWithEnforcer creates a middleware with a specific enforcer instance.
func CasbinMiddlewareWithEnforcer(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户角色，如果不存在则默认为anonymous
		role, exists := c.Get("role")
		if !exists {
			role = "anonymous"
		}

		// Check permission
		ok, err := e.Enforce(role.(string), c.Request.URL.Path, c.Request.Method)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred when authorizing user"})
			c.Abort()
			return
		}

		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this resource"})
			c.Abort()
			return
		}

		c.Next()
	}
}
