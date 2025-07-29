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
	var token string

	// 检查是否在测试环境中
	isTest := s.DB != nil && s.DB.Statement.Dest == &gorm.DB{}

	var txUserRepo repositories.UserRepository
	var txRoleRepo repositories.RoleRepository

	// 在测试环境中，我们使用注入的模拟仓库
	if isTest {
		txUserRepo = s.UserRepository
		txRoleRepo = s.RoleRepository
	} else {
		// 在生产环境中，我们在事务中创建新的仓库实例
		tx := s.DB.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		txUserRepo = repositories.NewGormUserRepository(tx)
		txRoleRepo = repositories.NewGormRoleRepository(tx)

		// 1. 检查用户名或邮箱是否已经被注册
		_, err := txUserRepo.FindByUsernameOrEmail(username, email)
		if err == nil {
			tx.Rollback()
			return nil, "", &UserExistsError{}
		}

		// 2. 对用户密码进行哈希加密
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			tx.Rollback()
			return nil, "", err
		}

		// 3. 获取或创建默认角色
		defaultRoleName := s.Config.App.DefaultRole
		role, err := txRoleRepo.FindByName(defaultRoleName)
		if err != nil {
			// 如果默认角色在数据库中不存在，则创建一个新的
			role = &models.Role{Name: defaultRoleName, Description: "普通用户"}
			if err := txRoleRepo.Create(role); err != nil {
				tx.Rollback()
				return nil, "", err
			}
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

		// 6. 加载用户的角色信息，以便在后续步骤中使用
		// 注意：在事务内，user对象已经包含了RoleID，但Role对象本身需要加载
		user.Role = *role

		if err := tx.Commit().Error; err != nil {
			return nil, "", err
		}
	}

	// 7. 为新注册的用户生成JWT（在事务成功后执行）
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
