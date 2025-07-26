package services

import (
	"go-web/config"
	"go-web/models"
	"go-web/repositories"
	"go-web/utils"
)

// UserExistsError 是用户已存在时返回的错误
type UserExistsError struct{}

func (e *UserExistsError) Error() string {
	return "user already exists"
}

// InvalidCredentialsError 是凭证无效时返回的错误
type InvalidCredentialsError struct{}

func (e *InvalidCredentialsError) Error() string {
	return "invalid username or password"
}

// AuthServiceInterface defines the contract for authentication services.
type AuthServiceInterface interface {
	Register(username, email, password string) (*models.User, string, error)
	Login(username, password string) (*models.User, string, error)
}

// AuthService 提供认证相关的业务逻辑
type AuthService struct {
	Config         *config.Config
	UserRepository repositories.UserRepository
	RoleRepository repositories.RoleRepository
}

func NewAuthService(cfg *config.Config, userRepo repositories.UserRepository, roleRepo repositories.RoleRepository) AuthServiceInterface {
	return &AuthService{
		Config:         cfg,
		UserRepository: userRepo,
		RoleRepository: roleRepo,
	}
}

// Register 注册新用户
func (s *AuthService) Register(username, email, password string) (*models.User, string, error) {
	// 检查用户是否已存在
	_, err := s.UserRepository.FindByUsernameOrEmail(username, email)
	if err == nil {
		return nil, "", &UserExistsError{}
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	// 获取默认角色
	defaultRoleName := s.Config.App.DefaultRole
	role, err := s.RoleRepository.FindByName(defaultRoleName)
	if err != nil {
		// 如果默认角色不存在，则创建
		role = &models.Role{Name: defaultRoleName, Description: "Regular user"}
		if err := s.RoleRepository.Create(role); err != nil {
			return nil, "", err
		}
	}

	// 创建新用户
	user := &models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
		RoleID:   role.ID,
	}

	if err := s.UserRepository.Create(user); err != nil {
		return nil, "", err
	}

	// 加载用户角色
	if err := s.UserRepository.LoadRole(user); err != nil {
		return nil, "", err
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Role.Name, s.Config)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Login 用户登录
func (s *AuthService) Login(username, password string) (*models.User, string, error) {
	// 查找用户
	user, err := s.UserRepository.FindByUsername(username)
	if err != nil {
		return nil, "", &InvalidCredentialsError{}
	}

	// 验证密码
	if err := utils.CheckPasswordHash(password, user.Password); err != nil {
		return nil, "", &InvalidCredentialsError{}
	}

	// 加载用户角色
	if err := s.UserRepository.LoadRole(user); err != nil {
		return nil, "", err
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Role.Name, s.Config)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
