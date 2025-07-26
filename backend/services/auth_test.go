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

// TestRegister_Success tests the successful registration of a new user.
func TestRegister_Success(t *testing.T) {
	// 1. Setup
	// Create instances of our mocks
	mockUserRepo := new(mocks.MockUserRepository)
	mockRoleRepo := new(mocks.MockRoleRepository)

	// Create a mock config
	cfg := &config.Config{
		App: config.AppConfig{
			DefaultRole: "user",
		},
		JWT: config.JWTConfig{
			Secret:     "test-secret", // Use a dedicated secret for testing
			Expiration: 3600,
		},
	}

	// Instantiate the service with the mocks
	authService := NewAuthService(cfg, mockUserRepo, mockRoleRepo)

	// Define the input for our test
	username := "testuser"
	email := "test@example.com"
	password := "password123"

	// 2. Define Mock Expectations
	// We expect FindByUsernameOrEmail to be called with our test data.
	// Since the user doesn't exist, it should return nil and an error.
	mockUserRepo.On("FindByUsernameOrEmail", username, email).Return(nil, errors.New("not found"))

	// We expect FindByName for the default role to be called.
	// We'll return a mock role.
	mockRole := &models.Role{
		Model: gorm.Model{ID: 1},
		Name:  "user",
	}
	mockRoleRepo.On("FindByName", "user").Return(mockRole, nil)

	// We expect Create to be called with a user object.
	// The `mock.Anything` matcher is used because the user object is created inside the service,
	// so we can't match the exact pointer. We just care that it's a *models.User.
	// It should return no error.
	mockUserRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

	// We expect LoadRole to be called, and it should succeed.
	mockUserRepo.On("LoadRole", mock.AnythingOfType("*models.User")).Run(func(args mock.Arguments) {
		// This function simulates the behavior of LoadRole: it populates the Role field.
		userArg := args.Get(0).(*models.User)
		userArg.Role = *mockRole
	}).Return(nil)

	// 3. Execution
	// Call the method we are testing
	user, token, err := authService.Register(username, email, password)

	// 4. Assertions
	// Use testify's assert package to check the results.
	assert.NoError(t, err)                   // We expect no error
	assert.NotNil(t, user)                   // We expect a user object to be returned
	assert.NotEmpty(t, token)                // We expect a token to be generated
	assert.Equal(t, username, user.Username) // Check if the username matches
	assert.Equal(t, "user", user.Role.Name)  // Check if the role is correct

	// We can also assert that our mocks were called as expected.
	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}

// TestRegister_UserExists tests registration when the user already exists.
func TestRegister_UserExists(t *testing.T) {
	// 1. Setup
	mockUserRepo := new(mocks.MockUserRepository)
	// Role repository will not be used in this case, so we can pass nil or an empty mock.
	mockRoleRepo := new(mocks.MockRoleRepository)
	cfg := &config.Config{} // Config is not used in this path

	authService := NewAuthService(cfg, mockUserRepo, mockRoleRepo)

	username := "existinguser"
	email := "exists@example.com"
	password := "password123"

	// 2. Define Mock Expectations
	// We expect FindByUsernameOrEmail to be called.
	// This time, it returns an existing user, indicating the user already exists.
	existingUser := &models.User{Username: username, Email: email}
	mockUserRepo.On("FindByUsernameOrEmail", username, email).Return(existingUser, nil)

	// 3. Execution
	user, token, err := authService.Register(username, email, password)

	// 4. Assertions
	assert.Error(t, err)                      // We expect an error
	assert.IsType(t, &UserExistsError{}, err) // The error should be of type UserExistsError
	assert.Nil(t, user)                       // No user object should be returned
	assert.Empty(t, token)                    // No token should be generated

	// Assert that Create was never called
	mockUserRepo.AssertNotCalled(t, "Create", mock.Anything)
}

// TestLogin_Success tests successful user login.
func TestLogin_Success(t *testing.T) {
	// 1. Setup
	mockUserRepo := new(mocks.MockUserRepository)
	mockRoleRepo := new(mocks.MockRoleRepository) // Not used in this path, but needed for constructor
	cfg := &config.Config{
		JWT: config.JWTConfig{Secret: "test-secret", Expiration: 3600},
	}
	authService := NewAuthService(cfg, mockUserRepo, mockRoleRepo)

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

	// 2. Define Mock Expectations
	mockUserRepo.On("FindByUsername", username).Return(existingUser, nil)
	mockUserRepo.On("LoadRole", existingUser).Run(func(args mock.Arguments) {
		userArg := args.Get(0).(*models.User)
		userArg.Role = *mockRole
	}).Return(nil)

	// 3. Execution
	user, token, err := authService.Login(username, password)

	// 4. Assertions
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
	assert.Equal(t, username, user.Username)
	mockUserRepo.AssertExpectations(t)
}

// TestLogin_InvalidCredentials tests login with a wrong password.
func TestLogin_InvalidCredentials(t *testing.T) {
	// 1. Setup
	mockUserRepo := new(mocks.MockUserRepository)
	mockRoleRepo := new(mocks.MockRoleRepository)
	cfg := &config.Config{}
	authService := NewAuthService(cfg, mockUserRepo, mockRoleRepo)

	username := "testuser"
	correctPassword := "password123"
	wrongPassword := "wrongpassword"
	hashedPassword, _ := utils.HashPassword(correctPassword)

	existingUser := &models.User{
		Model:    gorm.Model{ID: 1},
		Username: username,
		Password: hashedPassword,
	}

	// 2. Define Mock Expectations
	mockUserRepo.On("FindByUsername", username).Return(existingUser, nil)

	// 3. Execution
	user, token, err := authService.Login(username, wrongPassword)

	// 4. Assertions
	assert.Error(t, err)
	assert.IsType(t, &InvalidCredentialsError{}, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
	mockUserRepo.AssertExpectations(t)
}
