// package services_test 包含了对services包的单元测试和集成测试。
package services_test

import (
	"go-web/config"
	"go-web/models"
	"go-web/repositories"
	"go-web/services"
	"go-web/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// AuthServiceTestSuite 是一个测试套件，用于组织与AuthService相关的测试。
type AuthServiceTestSuite struct {
	suite.Suite
	db      *gorm.DB
	cfg     *config.Config
	service services.AuthServiceInterface
	// 用于集成测试的真实仓库
	userRepo repositories.UserRepository
	roleRepo repositories.RoleRepository
}

// SetupSuite 在测试套件开始时运行，用于初始化数据库和配置。
func (suite *AuthServiceTestSuite) SetupSuite() {
	// 初始化内存数据库
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		suite.T().Fatalf("无法连接到内存数据库: %v", err)
	}
	suite.db = db

	// 自动迁移数据库模式
	err = suite.db.AutoMigrate(&models.User{}, &models.Role{})
	if err != nil {
		suite.T().Fatalf("数据库迁移失败: %v", err)
	}

	// 初始化配置
	suite.cfg = &config.Config{
		App: config.AppConfig{
			DefaultRole: "user",
		},
		JWT: config.JWTConfig{
			Secret:     "test-secret",
			Expiration: 3600,
		},
	}

	// 初始化用于集成测试的真实仓库
	suite.userRepo = repositories.NewGormUserRepository(suite.db)
	suite.roleRepo = repositories.NewGormRoleRepository(suite.db)

	// 初始化服务
	suite.service = services.NewAuthService(suite.cfg, suite.userRepo, suite.roleRepo, suite.db)
}

// TearDownTest 在每个测试后运行，用于清理数据库。
func (suite *AuthServiceTestSuite) TearDownTest() {
	// 清理所有表以确保测试隔离
	suite.db.Exec("DELETE FROM users")
	suite.db.Exec("DELETE FROM roles")
}

// TestAuthServiceTestSuite 运行测试套件。
func TestAuthServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}

// TestRegister_Success 测试新用户成功注册的场景（集成测试）。
func (suite *AuthServiceTestSuite) TestRegister_Success() {
	// 执行
	user, token, err := suite.service.Register("testuser", "test@example.com", "password123")

	// 断言
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.NotEmpty(suite.T(), token)
	assert.Equal(suite.T(), "testuser", user.Username)
	assert.Equal(suite.T(), "user", user.Role.Name)

	// 验证数据库中确实创建了用户
	dbUser, dbErr := suite.userRepo.FindByUsername("testuser")
	assert.NoError(suite.T(), dbErr)
	assert.NotNil(suite.T(), dbUser)
	assert.Equal(suite.T(), "testuser", dbUser.Username)
}

// TestRegister_UserExists 测试当用户已存在时注册失败的场景（集成测试）。
func (suite *AuthServiceTestSuite) TestRegister_UserExists() {
	// 准备：先创建一个用户
	existingUser := &models.User{
		Username: "existinguser",
		Email:    "exists@example.com",
		Password: "password123",
		Role:     models.Role{Name: "user", Description: "普通用户"},
	}
	err := suite.userRepo.Create(existingUser)
	assert.NoError(suite.T(), err)

	// 执行：尝试用相同的用户名再次注册
	user, token, err := suite.service.Register("existinguser", "another@example.com", "password123")

	// 断言
	assert.Error(suite.T(), err)
	assert.IsType(suite.T(), &services.UserExistsError{}, err)
	assert.Nil(suite.T(), user)
	assert.Empty(suite.T(), token)
}

// TestLogin_Success 测试用户成功登录的场景。
func (suite *AuthServiceTestSuite) TestLogin_Success() {
	// 准备
	hashedPassword, _ := utils.HashPassword("password123")
	role := models.Role{Name: "admin", Description: "管理员"}
	suite.roleRepo.Create(&role)
	user := &models.User{
		Username: "loginuser",
		Email:    "login@example.com",
		Password: hashedPassword,
		RoleID:   role.ID,
	}
	suite.userRepo.Create(user)

	// 执行
	loggedInUser, token, err := suite.service.Login("loginuser", "password123")

	// 断言
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), loggedInUser)
	assert.NotEmpty(suite.T(), token)
	assert.Equal(suite.T(), "loginuser", loggedInUser.Username)
	assert.Equal(suite.T(), "admin", loggedInUser.Role.Name)
}

// TestLogin_InvalidCredentials 测试使用错误密码登录失败的场景。
func (suite *AuthServiceTestSuite) TestLogin_InvalidCredentials() {
	// 准备
	hashedPassword, _ := utils.HashPassword("correct-password")
	user := &models.User{Username: "loginuser", Password: hashedPassword}
	suite.userRepo.Create(user)

	// 执行
	loggedInUser, token, err := suite.service.Login("loginuser", "wrong-password")

	// 断言
	assert.Error(suite.T(), err)
	assert.IsType(suite.T(), &services.InvalidCredentialsError{}, err)
	assert.Nil(suite.T(), loggedInUser)
	assert.Empty(suite.T(), token)
}

// TestLogin_UserNotFound 测试用户不存在时登录失败的场景。
func (suite *AuthServiceTestSuite) TestLogin_UserNotFound() {
	// 执行
	user, token, err := suite.service.Login("nonexistentuser", "password")

	// 断言
	assert.Error(suite.T(), err)
	assert.IsType(suite.T(), &services.InvalidCredentialsError{}, err) // 服务返回相同的错误类型以避免泄露用户信息
	assert.Nil(suite.T(), user)
	assert.Empty(suite.T(), token)
}
