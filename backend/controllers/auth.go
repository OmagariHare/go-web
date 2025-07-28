package controllers

import (
	"go-web/dtos"
	"go-web/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthService services.AuthServiceInterface
}

func NewAuthController(authService services.AuthServiceInterface) *AuthController {
	return &AuthController{AuthService: authService}
}

// Register 注册新用户
func (ac *AuthController) Register(c *gin.Context) {
	var req dtos.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	user, token, err := ac.AuthService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		_ = c.Error(err)
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
		_ = c.Error(err)
		return
	}

	user, token, err := ac.AuthService.Login(req.Username, req.Password)
	if err != nil {
		_ = c.Error(err)
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
