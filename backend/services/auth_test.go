// package services_test 包含了对services包的单元测试。
package services

import (
	"errors"
	"go-web/config"
	"go-web/mocks"
	"go-web/models"
	"go-web/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// TestRegister_Success 测试新用户成功注册的场景。
func TestRegister_Success(t *testing.T) {
	// 1. 准备阶段 (Setup)
	// 创建仓库（repository）的模拟（mock）实例
	mockUserRepo := new(mocks.MockUserRepository)
	mockRoleRepo := new(mocks.MockRoleRepository)

	// 创建一个用于测试的模拟配置
	cfg := &config.Config{
		App: config.AppConfig{
			DefaultRole: "user",
		},
		JWT: config.JWTConfig{
			Secret:     "test-secret", // 在测试中使用独立的密钥
			Expiration: 3600,
		},
	}

	// 使用模拟实例来初始化被测试的服务
	authService := NewAuthService(cfg, mockUserRepo, mockRoleRepo)

	// 定义测试的输入数据
	username := "testuser"
	email := "test@example.com"
	password := "password123"

	// 2. 定义模拟期望 (Mock Expectations)
	// 期望 FindByUsernameOrEmail 方法被调用，并且因为用户不存在，所以返回错误。
	mockUserRepo.On("FindByUsernameOrEmail", username, email).Return(nil, errors.New("not found"))

	// 期望 FindByName 方法被调用以查找默认角色。
	mockRole := &models.Role{
		Model: gorm.Model{ID: 1},
		Name:  "user",
	}
	mockRoleRepo.On("FindByName", "user").Return(mockRole, nil)

	// 期望 Create 方法被一个用户对象调用。
	// mock.AnythingOfType 用于匹配在服务内部创建的对象指针。
	mockUserRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

	// 期望 LoadRole 方法被调用，并模拟其行为：填充用户的Role字段。
	mockUserRepo.On("LoadRole", mock.AnythingOfType("*models.User")).Run(func(args mock.Arguments) {
		userArg := args.Get(0).(*models.User)
		userArg.Role = *mockRole
	}).Return(nil)

	// 3. 执行阶段 (Execution)
	// 调用我们正在测试的方法
	user, token, err := authService.Register(username, email, password)

	// 4. 断言阶段 (Assertions)
	// 使用 testify/assert 包来验证结果
	assert.NoError(t, err)                   // 期望没有错误发生
	assert.NotNil(t, user)                   // 期望返回一个用户对象
	assert.NotEmpty(t, token)                // 期望生成一个token
	assert.Equal(t, username, user.Username) // 验证用户名是否正确
	assert.Equal(t, "user", user.Role.Name)  // 验证用户角色是否正确

	// 验证所有的模拟期望都已满足
	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}

// TestRegister_UserExists 测试当用户已存在时注册失败的场景。
func TestRegister_UserExists(t *testing.T) {
	// 1. 准备阶段
	mockUserRepo := new(mocks.MockUserRepository)
	// 在这个测试用例中，RoleRepository 和 Config 不会被使用，因此我们可以传入一个空的mock或者nil。
	authService := NewAuthService(&config.Config{}, mockUserRepo, new(mocks.MockRoleRepository))

	username := "existinguser"
	email := "exists@example.com"
	password := "password123"

	// 2. 定义模拟期望
	// 期望 FindByUsernameOrEmail 被调用，并返回一个已存在的用户，表示用户已存在。
	existingUser := &models.User{Username: username, Email: email}
	mockUserRepo.On("FindByUsernameOrEmail", username, email).Return(existingUser, nil)

	// 3. 执行阶段
	user, token, err := authService.Register(username, email, password)

	// 4. 断言阶段
	assert.Error(t, err)                      // 期望有错误发生
	assert.IsType(t, &UserExistsError{}, err) // 错误类型应为 UserExistsError
	assert.Nil(t, user)                       // 不应返回用户对象
	assert.Empty(t, token)                    // 不应生成token

	// 验证 Create 方法从未被调用
	mockUserRepo.AssertNotCalled(t, "Create", mock.Anything)
}

// TestLogin_Success 测试用户成功登录的场景。
func TestLogin_Success(t *testing.T) {
	// 1. 准备阶段
	mockUserRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{
		JWT: config.JWTConfig{Secret: "test-secret", Expiration: 3600},
	}
	authService := NewAuthService(cfg, mockUserRepo, new(mocks.MockRoleRepository))

	username := "testuser"
	password := "password123"
	hashedPassword, _ := utils.HashPassword(password)

	mockRole := &models.Role{Model: gorm.Model{ID: 1}, Name: "user"}
	existingUser := &models.User{
		Model:    gorm.Model{ID: 1},
		Username: username,
		Password: hashedPassword,
		RoleID:   1,
	}

	// 2. 定义模拟期望
	mockUserRepo.On("FindByUsername", username).Return(existingUser, nil)
	mockUserRepo.On("LoadRole", existingUser).Run(func(args mock.Arguments) {
		userArg := args.Get(0).(*models.User)
		userArg.Role = *mockRole
	}).Return(nil)

	// 3. 执行阶段
	user, token, err := authService.Login(username, password)

	// 4. 断言阶段
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
	assert.Equal(t, username, user.Username)
	mockUserRepo.AssertExpectations(t)
}

// TestLogin_InvalidCredentials 测试使用错误密码登录失败的场景。
func TestLogin_InvalidCredentials(t *testing.T) {
	// 1. 准备阶段
	mockUserRepo := new(mocks.MockUserRepository)
	authService := NewAuthService(&config.Config{}, mockUserRepo, new(mocks.MockRoleRepository))

	username := "testuser"
	correctPassword := "password123"
	wrongPassword := "wrongpassword"
	hashedPassword, _ := utils.HashPassword(correctPassword)

	existingUser := &models.User{
		Model:    gorm.Model{ID: 1},
		Username: username,
		Password: hashedPassword,
	}

	// 2. 定义模拟期望
	mockUserRepo.On("FindByUsername", username).Return(existingUser, nil)

	// 3. 执行阶段
	user, token, err := authService.Login(username, wrongPassword)

	// 4. 断言阶段
	assert.Error(t, err)
	assert.IsType(t, &InvalidCredentialsError{}, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
	mockUserRepo.AssertExpectations(t)
}
