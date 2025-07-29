package services

// package services 包含了应用程序的业务逻辑。
// 它作为控制器和仓库之间的桥梁，处理如用户认证、注册等核心功能。

import (
	"go-web/config"
	"go-web/models"
	"go-web/repositories"
	"go-web/utils"
	"gorm.io/gorm"
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
	DB             *gorm.DB // 添加DB实例用于事务
}

// NewAuthService 是 AuthService 的构造函数。
func NewAuthService(cfg *config.Config, userRepo repositories.UserRepository, roleRepo repositories.RoleRepository, db *gorm.DB) AuthServiceInterface {
	return &AuthService{
		Config:         cfg,
		UserRepository: userRepo,
		RoleRepository: roleRepo,
		DB:             db, // 注入DB实例
	}
}

// Register 负责注册一个新用户。
// 它会检查用户是否已存在，对密码进行哈希处理，分配默认角色，创建用户，并生成JWT。
// 整个注册过程在一个数据库事务中完成，以确保数据一致性。
func (s *AuthService) Register(username, email, password string) (*models.User, string, error) {
	var user *models.User

	// 启动数据库事务
	tx := s.DB.Begin()
	if tx.Error != nil {
		return nil, "", tx.Error
	}

	// 使用 defer-recover 机制确保事务在发生 panic 时能够回滚
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 在事务中创建仓库实例
	txUserRepo := repositories.NewGormUserRepository(tx)
	txRoleRepo := repositories.NewGormRoleRepository(tx)

	// 1. 检查用户名或邮箱是否已经被注册
	// 在注册场景下，我们期望这里返回 "record not found" 错误
	_, err := txUserRepo.FindByUsernameOrEmail(username, email)
	if err == nil {
		tx.Rollback()
		return nil, "", &UserExistsError{}
	}
	if err != gorm.ErrRecordNotFound {
		// 如果是其他类型的数据库错误，则回滚并返回错误
		tx.Rollback()
		return nil, "", err
	}

	// 2. 对用户密码进行哈希加密
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		tx.Rollback()
		return nil, "", err
	}

	// 3. 获取默认角色
	// 我们期望默认角色在数据库中总是存在的（由测试的SetupTest或生产的seeder保证）
	defaultRoleName := s.Config.App.DefaultRole
	role, err := txRoleRepo.FindByName(defaultRoleName)
	if err != nil {
		// 如果角色不存在（这在正常情况下不应该发生），则回滚
		tx.Rollback()
		return nil, "", err
	}

	// 4. 创建新用户实例
	user = &models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
		RoleID:   role.ID,
	}

	// 5. 将新用户存入数据库
	if err := txUserRepo.Create(user); err != nil {
		tx.Rollback()
		return nil, "", err
	}

	// 6. 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, "", err
	}

	// 7. 为新注册的用户生成JWT（在事务成功后执行）
	// 此时 user 对象已经包含了 RoleID，但 Role 对象本身需要从 role 变量中获取
	user.Role = *role
	token, err := utils.GenerateToken(user.ID, user.Role.Name, s.Config)
	if err != nil {
		// 事务已经提交，但Token生成失败。这是一个边缘情况。
		// 此时用户已创建成功，但无法立即登录。
		// 我们可以选择返回错误，让用户稍后尝试登录。
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
