package controllers

import (
	"go-web/dtos"
	"go-web/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

// Register 注册新用户
func (ac *AuthController) Register(c *gin.Context) {
	var req dtos.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := ac.AuthService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		switch err.(type) {
		case *services.UserExistsError:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		}
		return
	}

	// 构造响应
	userResponse := dtos.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		RoleID:   user.RoleID,
		Role:     user.Role.Name,
	}

	response := dtos.AuthResponse{
		Token: token,
		User:  userResponse,
	}

	c.JSON(http.StatusCreated, response)
}

// Login 用户登录
func (ac *AuthController) Login(c *gin.Context) {
	var req dtos.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := ac.AuthService.Login(req.Username, req.Password)
	if err != nil {
		switch err.(type) {
		case *services.InvalidCredentialsError:
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		}
		return
	}

	// 构造响应
	userResponse := dtos.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		RoleID:   user.RoleID,
		Role:     user.Role.Name,
	}

	response := dtos.AuthResponse{
		Token: token,
		User:  userResponse,
	}

	c.JSON(http.StatusOK, response)
}
