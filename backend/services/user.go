package services

import (
	"errors"
	"go-web/models"
	"go-web/repositories"
)

var ErrPermissionDenied = errors.New("permission denied")

// UserServiceInterface defines the contract for user-related services.
type UserServiceInterface interface {
	GetUsers() ([]models.User, error)
	GetUser(id uint) (*models.User, error)
	UpdateUser(targetUserID, currentUserID uint, currentUserRole string, updateUser *models.User) (*models.User, error)
	DeleteUser(id uint) error
}

// UserService 提供用户相关的业务逻辑
type UserService struct {
	UserRepository repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserServiceInterface {
	return &UserService{UserRepository: userRepo}
}

// GetUsers 获取用户列表
func (s *UserService) GetUsers() ([]models.User, error) {
	return s.UserRepository.FindAll()
}

// GetUser 获取单个用户信息
func (s *UserService) GetUser(id uint) (*models.User, error) {
	return s.UserRepository.FindByID(id)
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(targetUserID, currentUserID uint, currentUserRole string, updateUser *models.User) (*models.User, error) {
	// 检查权限
	if targetUserID != currentUserID && currentUserRole != "admin" {
		return nil, ErrPermissionDenied
	}

	user, err := s.UserRepository.FindByID(targetUserID)
	if err != nil {
		return nil, err
	}

	// 更新用户信息
	if updateUser.Username != "" {
		user.Username = updateUser.Username
	}
	if updateUser.Email != "" {
		user.Email = updateUser.Email
	}
	if updateUser.RoleID != 0 {
		// 只有管理员可以更改角色
		if currentUserRole == "admin" {
			user.RoleID = updateUser.RoleID
		}
	}

	if err := s.UserRepository.Update(user); err != nil {
		return nil, err
	}

	// 加载用户角色
	if err := s.UserRepository.LoadRole(user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id uint) error {
	user, err := s.UserRepository.FindByID(id)
	if err != nil {
		return err
	}

	return s.UserRepository.Delete(user)
}
