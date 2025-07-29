package controllers

import (
	"go-web/dtos"
	"go-web/models"
	"go-web/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService services.UserServiceInterface
}

func NewUserController(userService services.UserServiceInterface) *UserController {
	return &UserController{UserService: userService}
}

// GetUsers 获取用户列表
func (uc *UserController) GetUsers(c *gin.Context) {
	users, err := uc.UserService.GetUsers()
	if err != nil {
		_ = c.Error(err)
		return
	}

	var userResponses []dtos.UserResponse
	for i := range users {
		user := &users[i]
		userResponses = append(userResponses, dtos.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			RoleID:   user.RoleID,
			Role:     user.Role.Name,
		})
	}

	c.JSON(http.StatusOK, userResponses)
}

// GetUser 获取单个用户信息
func (uc *UserController) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	user, err := uc.UserService.GetUser(uint(id))
	if err != nil {
		_ = c.Error(err)
		return
	}

	userResponse := dtos.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		RoleID:   user.RoleID,
		Role:     user.Role.Name,
	}

	c.JSON(http.StatusOK, userResponse)
}

// UpdateUser 更新用户信息
func (uc *UserController) UpdateUser(c *gin.Context) {
	targetUserID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	var req dtos.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	// 从DTO创建更新模型，避免不安全的字段绑定
	updateData := &models.User{
		Username: req.Username,
		Email:    req.Email,
		RoleID:   req.RoleID,
	}

	currentUserID := c.GetUint("user_id")
	currentUserRole, _ := c.Get("role")

	user, err := uc.UserService.UpdateUser(uint(targetUserID), currentUserID, currentUserRole.(string), updateData)
	if err != nil {
		_ = c.Error(err)
		return
	}

	userResponse := dtos.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		RoleID:   user.RoleID,
		Role:     user.Role.Name,
	}

	c.JSON(http.StatusOK, userResponse)
}

// DeleteUser 删除用户
func (uc *UserController) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = uc.UserService.DeleteUser(uint(id))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
