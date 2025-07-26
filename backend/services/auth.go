// package services 包含了应用程序的业务逻辑。
// 它作为控制器和仓库之间的桥梁，处理如用户认证、注册等核心功能。
package services

import (
	"go-web/config"
	"go-web/models"
	"go-web/repositories"
	"go-web/utils"
)

// UserExistsError 在尝试创建已存在的用户时返回。
type UserExistsError struct{}

// Error 实现了 error 接口。
func (e *UserExistsError) Error() string {
	return "用户已存在"
}

// InvalidCredentialsError 在登录时提供的用户名或密码不正确时返回。
type InvalidCredentialsError struct{}

// Error 实现了 error 接口。
func (e *InvalidCredentialsError) Error() string {
	return "无效的用户名或密码"
}

// AuthServiceInterface 定义了认证服务应实现的功能契约。
// 使用接口可以方便地在测试中替换真实的服务实现。
type AuthServiceInterface interface {
	// Register 处理新用户的注册逻辑。
	Register(username, email, password string) (*models.User, string, error)
	// Login 处理用户的登录逻辑。
	Login(username, password string) (*models.User, string, error)
}

// AuthService 提供了认证相关的业务逻辑实现。
// 它依赖于配置、用户仓库和角色仓库。
type AuthService struct {
	Config         *config.Config
	UserRepository repositories.UserRepository
	RoleRepository repositories.RoleRepository
}

// NewAuthService 是 AuthService 的构造函数。
func NewAuthService(cfg *config.Config, userRepo repositories.UserRepository, roleRepo repositories.RoleRepository) AuthServiceInterface {
	return &AuthService{
		Config:         cfg,
		UserRepository: userRepo,
		RoleRepository: roleRepo,
	}
}

// Register 负责注册一个新用户。
// 它会检查用户是否已存在，对密码进行哈希处理，分配默认角色，创建用户，并生成JWT。
func (s *AuthService) Register(username, email, password string) (*models.User, string, error) {
	// 1. 检查用户名或邮箱是否已经被注册
	_, err := s.UserRepository.FindByUsernameOrEmail(username, email)
	if err == nil {
		return nil, "", &UserExistsError{}
	}

	// 2. 对用户密码进行哈希加密
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	// 3. 获取或创建默认角色
	defaultRoleName := s.Config.App.DefaultRole
	role, err := s.RoleRepository.FindByName(defaultRoleName)
	if err != nil {
		// 如果默认角色在数据库中不存在，则创建一个新的
		role = &models.Role{Name: defaultRoleName, Description: "普通用户"}
		if err := s.RoleRepository.Create(role); err != nil {
			return nil, "", err
		}
	}

	// 4. 创建新用户实例
	user := &models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
		RoleID:   role.ID,
	}

	// 5. 将新用户存入数据库
	if err := s.UserRepository.Create(user); err != nil {
		return nil, "", err
	}

	// 6. 加载用户的角色信息，以便在后续步骤中使用
	if err := s.UserRepository.LoadRole(user); err != nil {
		return nil, "", err
	}

	// 7. 为新注册的用户生成JWT
	token, err := utils.GenerateToken(user.ID, user.Role.Name, s.Config)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Login 负责处理用户登录。
// 它会验证用户名和密码，如果成功，则生成一个新的JWT。
func (s *AuthService) Login(username, password string) (*models.User, string, error) {
	// 1. 根据用户名查找用户
	user, err := s.UserRepository.FindByUsername(username)
	if err != nil {
		return nil, "", &InvalidCredentialsError{}
	}

	// 2. 验证提供的密码是否与存储的哈希密码匹配
	if err := utils.CheckPasswordHash(password, user.Password); err != nil {
		return nil, "", &InvalidCredentialsError{}
	}

	// 3. 加载用户的角色信息
	if err := s.UserRepository.LoadRole(user); err != nil {
		return nil, "", err
	}

	// 4. 生成JWT
	token, err := utils.GenerateToken(user.ID, user.Role.Name, s.Config)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
