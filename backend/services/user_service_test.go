// package services_test 包含了对services包的单元测试。
package services

import (
	"go-web/mocks"
	"go-web/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestGetUsers_Success 测试成功获取用户列表的场景。
func TestGetUsers_Success(t *testing.T) {
	// 1. 准备阶段
	mockUserRepo := new(mocks.MockUserRepository)
	userService := NewUserService(mockUserRepo)

	// 2. 定义模拟期望
	mockedUsers := []models.User{
		{Model: gorm.Model{ID: 1}, Username: "user1"},
		{Model: gorm.Model{ID: 2}, Username: "user2"},
	}
	mockUserRepo.On("FindAll").Return(mockedUsers, nil)

	// 3. 执行阶段
	users, err := userService.GetUsers()

	// 4. 断言阶段
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 2)
	assert.Equal(t, "user1", users[0].Username)
	mockUserRepo.AssertExpectations(t)
}

// TestGetUser_Success 测试成功获取单个用户信息的场景。
func TestGetUser_Success(t *testing.T) {
	// 1. 准备阶段
	mockUserRepo := new(mocks.MockUserRepository)
	userService := NewUserService(mockUserRepo)

	userID := uint(1)

	// 2. 定义模拟期望
	mockedUser := &models.User{Model: gorm.Model{ID: userID}, Username: "testuser"}
	mockUserRepo.On("FindByID", userID).Return(mockedUser, nil)

	// 3. 执行阶段
	user, err := userService.GetUser(userID)

	// 4. 断言阶段
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	mockUserRepo.AssertExpectations(t)
}

// TestUpdateUser_Success_AsAdmin 测试管理员成功更新用户信息的场景。
func TestUpdateUser_Success_AsAdmin(t *testing.T) {
	// 1. 准备阶段
	mockUserRepo := new(mocks.MockUserRepository)
	userService := NewUserService(mockUserRepo)

	targetUserID := uint(2)
	adminUserID := uint(1)
	adminUserRole := "admin"

	updateData := &models.User{
		Username: "updated_username",
		Email:    "updated_email@example.com",
		RoleID:   3,
	}

	originalUser := &models.User{
		Model:    gorm.Model{ID: targetUserID},
		Username: "original_username",
		Email:    "original@example.com",
		RoleID:   2,
	}

	// 2. 定义模拟期望
	mockUserRepo.On("FindByID", targetUserID).Return(originalUser, nil)
	mockUserRepo.On("Update", originalUser).Return(nil)
	mockUserRepo.On("LoadRole", originalUser).Return(nil)

	// 3. 执行阶段
	updatedUser, err := userService.UpdateUser(targetUserID, adminUserID, adminUserRole, updateData)

	// 4. 断言阶段
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.Equal(t, "updated_username", updatedUser.Username)
	assert.Equal(t, "updated_email@example.com", updatedUser.Email)
	assert.Equal(t, uint(3), updatedUser.RoleID)
	mockUserRepo.AssertExpectations(t)
}

// TestUpdateUser_PermissionDenied 测试因权限不足而更新失败的场景。
func TestUpdateUser_PermissionDenied(t *testing.T) {
	// 1. 准备阶段
	mockUserRepo := new(mocks.MockUserRepository)
	userService := NewUserService(mockUserRepo)

	targetUserID := uint(2)
	requestingUserID := uint(3)  // 一个不同的用户ID
	requestingUserRole := "user" // 非管理员角色

	updateData := &models.User{Username: "new_name"}

	// 2. 执行阶段
	updatedUser, err := userService.UpdateUser(targetUserID, requestingUserID, requestingUserRole, updateData)

	// 3. 断言阶段
	assert.Error(t, err)                      // 期望有错误发生
	assert.Equal(t, ErrPermissionDenied, err) // 错误应为权限不足
	assert.Nil(t, updatedUser)                // 不应返回用户对象

	// 确保没有调用任何仓库方法
	mockUserRepo.AssertNotCalled(t, "FindByID")
	mockUserRepo.AssertNotCalled(t, "Update")
}

// TestDeleteUser_Success 测试成功删除用户的场景。
func TestDeleteUser_Success(t *testing.T) {
	// 1. 准备阶段
	mockUserRepo := new(mocks.MockUserRepository)
	userService := NewUserService(mockUserRepo)

	userID := uint(1)
	userToDelete := &models.User{Model: gorm.Model{ID: userID}}

	// 2. 定义模拟期望
	mockUserRepo.On("FindByID", userID).Return(userToDelete, nil)
	mockUserRepo.On("Delete", userToDelete).Return(nil)

	// 3. 执行阶段
	err := userService.DeleteUser(userID)

	// 4. 断言阶段
	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}
