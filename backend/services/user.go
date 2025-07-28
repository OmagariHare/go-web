package services

// package services 包含了应用程序的业务逻辑。
// 这个文件特别关注与用户管理相关的功能。

import (
	"errors"
	"go-web/models"
	"go-web/repositories"
)

// ErrPermissionDenied 在用户尝试执行未授权的操作时返回。
var ErrPermissionDenied = errors.New("权限不足")

// UserServiceInterface 定义了用户服务应实现的功能契约。
type UserServiceInterface interface {
	// GetUsers 获取所有用户的列表。
	GetUsers() ([]models.User, error)
	// GetUser 根据ID获取单个用户的详细信息。
	GetUser(id uint) (*models.User, error)
	// UpdateUser 更新指定ID的用户信息。
	UpdateUser(targetUserID, currentUserID uint, currentUserRole string, updateUser *models.User) (*models.User, error)
	// DeleteUser 删除指定ID的用户。
	DeleteUser(id uint) error
}

// UserService 提供了用户管理相关的业务逻辑实现。
// 它依赖于用户仓库。
type UserService struct {
	UserRepository repositories.UserRepository
}

// NewUserService 是 UserService 的构造函数。
func NewUserService(userRepo repositories.UserRepository) UserServiceInterface {
	return &UserService{UserRepository: userRepo}
}

// GetUsers 获取所有用户的列表。
func (s *UserService) GetUsers() ([]models.User, error) {
	return s.UserRepository.FindAll()
}

// GetUser 获取单个用户的详细信息。
func (s *UserService) GetUser(id uint) (*models.User, error) {
	return s.UserRepository.FindByID(id)
}

// UpdateUser 更新用户信息。
// 它包含了权限检查逻辑：
// - 用户可以更新自己的信息。
// - 管理员（admin）可以更新任何人的信息。
// - 只有管理员可以更改用户的角色。
func (s *UserService) UpdateUser(targetUserID, currentUserID uint, currentUserRole string, updateUser *models.User) (*models.User, error) {
	// 权限检查：如果目标用户不是当前用户，并且当前用户不是管理员，则拒绝访问。
	if targetUserID != currentUserID && currentUserRole != "admin" {
		return nil, ErrPermissionDenied
	}

	// 从数据库获取最新的用户信息
	user, err := s.UserRepository.FindByID(targetUserID)
	if err != nil {
		return nil, err
	}

	// 更新允许修改的字段
	if updateUser.Username != "" {
		user.Username = updateUser.Username
	}
	if updateUser.Email != "" {
		user.Email = updateUser.Email
	}
	if updateUser.RoleID != 0 {
		// 只有管理员可以修改角色ID
		if currentUserRole == "admin" {
			user.RoleID = updateUser.RoleID
		}
	}

	// 将更新后的用户信息保存到数据库
	if err := s.UserRepository.Update(user); err != nil {
		return nil, err
	}

	// 重新加载用户的角色信息，以确保返回的数据是完整的
	if err := s.UserRepository.LoadRole(user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser 删除一个用户。
func (s *UserService) DeleteUser(id uint) error {
	// 首先需要根据ID找到对应的用户实体
	user, err := s.UserRepository.FindByID(id)
	if err != nil {
		return err
	}

	// 然后删除该用户
	return s.UserRepository.Delete(user)
}
