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
	for _, user := range users {
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

	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	currentUserID := c.GetUint("user_id")
	currentUserRole, _ := c.Get("role")

	user, err := uc.UserService.UpdateUser(uint(targetUserID), currentUserID, currentUserRole.(string), &req)
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
