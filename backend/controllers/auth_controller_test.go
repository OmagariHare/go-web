package controllers

import (
	"bytes"
	"encoding/json"
	"go-web/dtos"
	"go-web/middleware"
	"go-web/models"
	"go-web/services"
	"go-web/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Make sure MockAuthService implements the interface
var _ services.AuthServiceInterface = (*MockAuthService)(nil)

// MockAuthService is a mock implementation of AuthService for testing controllers.
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(username, email, password string) (*models.User, string, error) {
	args := m.Called(username, email, password)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).(*models.User), args.String(1), args.Error(2)
}

func (m *MockAuthService) Login(username, password string) (*models.User, string, error) {
	args := m.Called(username, password)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).(*models.User), args.String(1), args.Error(2)
}

func setupAuthTestRouter() (*gin.Engine, *MockAuthService) {
	gin.SetMode(gin.TestMode)
	utils.InitLogger("debug", "", 100, 3, 7, false) // Initialize logger for tests
	mockAuthService := new(MockAuthService)
	authController := NewAuthController(mockAuthService)

	router := gin.Default()
	router.Use(middleware.ErrorHandler())
	router.POST("/auth/register", authController.Register)
	router.POST("/auth/login", authController.Login)

	return router, mockAuthService
}

func TestRegister_Endpoint_Success(t *testing.T) {
	// 1. Setup
	router, mockAuthService := setupAuthTestRouter()

	// 2. Define Mock Expectations
	registerReq := dtos.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	mockedUser := &models.User{
		Username: registerReq.Username,
		Email:    registerReq.Email,
		Role:     models.Role{Name: "user"},
	}
	mockedToken := "mocked-jwt-token"

	mockAuthService.On("Register", registerReq.Username, registerReq.Email, registerReq.Password).Return(mockedUser, mockedToken, nil)

	// 3. Execution
	jsonValue, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 4. Assertions
	assert.Equal(t, http.StatusCreated, w.Code, "HTTP status code should be 201 Created")

	var response dtos.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Should be able to unmarshal the response body")
	assert.Equal(t, mockedToken, response.Token, "Response token should match the mocked token")
	assert.Equal(t, registerReq.Username, response.User.Username, "Response username should match the request username")

	mockAuthService.AssertExpectations(t)
}

func TestRegister_Endpoint_UserExists(t *testing.T) {
	// 1. Setup
	router, mockAuthService := setupAuthTestRouter()

	// 2. Define Mock Expectations
	registerReq := dtos.RegisterRequest{
		Username: "existinguser",
		Email:    "exists@example.com",
		Password: "password123",
	}

	// Simulate the service returning a UserExistsError
	mockAuthService.On("Register", registerReq.Username, registerReq.Email, registerReq.Password).Return(nil, "", &services.UserExistsError{})

	// 3. Execution
	jsonValue, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 4. Assertions
	assert.Equal(t, http.StatusConflict, w.Code, "HTTP status code should be 409 Conflict")

	var errorResponse middleware.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, errorResponse.Code)
	assert.Contains(t, errorResponse.Message, "用户已存在")

	mockAuthService.AssertExpectations(t)
}

func TestLogin_Endpoint_Success(t *testing.T) {
	// 1. Setup
	router, mockAuthService := setupAuthTestRouter()

	// 2. Define Mock Expectations
	loginReq := dtos.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	mockedUser := &models.User{
		Username: loginReq.Username,
		Role:     models.Role{Name: "user"},
	}
	mockedToken := "mocked-jwt-token-for-login"

	mockAuthService.On("Login", loginReq.Username, loginReq.Password).Return(mockedUser, mockedToken, nil)

	// 3. Execution
	jsonValue, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 4. Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response dtos.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, mockedToken, response.Token)
	assert.Equal(t, loginReq.Username, response.User.Username)

	mockAuthService.AssertExpectations(t)
}

func TestLogin_Endpoint_InvalidCredentials(t *testing.T) {
	// 1. Setup
	router, mockAuthService := setupAuthTestRouter()

	// 2. Define Mock Expectations
	loginReq := dtos.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	mockAuthService.On("Login", loginReq.Username, loginReq.Password).Return(nil, "", &services.InvalidCredentialsError{})

	// 3. Execution
	jsonValue, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 4. Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var errorResponse middleware.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, errorResponse.Code)
	assert.Contains(t, errorResponse.Message, "无效的用户名或密码")

	mockAuthService.AssertExpectations(t)
}
