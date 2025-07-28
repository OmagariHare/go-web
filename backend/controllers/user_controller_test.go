package controllers

import (
	"encoding/json"
	"fmt"
	"go-web/config"
	"go-web/dtos"
	"go-web/middleware"
	"go-web/models"
	"go-web/services"
	"go-web/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Make sure MockUserService implements the interface
var _ services.UserServiceInterface = (*MockUserService)(nil)

// MockUserService is a mock for the UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUsers() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserService) GetUser(id uint) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(targetUserID, currentUserID uint, currentUserRole string, updateUser *models.User) (*models.User, error) {
	args := m.Called(targetUserID, currentUserID, currentUserRole, updateUser)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// setupCorrectIsolatedUserRouter sets up a fully isolated router for testing user endpoints.
func setupCorrectIsolatedUserRouter(e *casbin.Enforcer) (*gin.Engine, *MockUserService, *config.Config) {
	gin.SetMode(gin.TestMode)

	mockUserService := new(MockUserService)
	userController := NewUserController(mockUserService)

	// THIS IS THE FIX: Added a non-zero Expiration to the JWT config.
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "a-truly-isolated-secret",
			Expiration: 60, // 60 seconds
		},
	}

	router := gin.New()
	router.Use(middleware.AuthMiddleware(cfg))
	router.Use(middleware.CasbinMiddlewareWithEnforcer(e))

	router.GET("/users", userController.GetUsers)

	return router, mockUserService, cfg
}

func TestGetUsers_Endpoint_Success_Correct(t *testing.T) {
	casbinModel, _ := model.NewModelFromString(`[request_definition]
r = sub, obj, act
[policy_definition]
p = sub, obj, act
[role_definition]
g = _, _
[policy_effect]
e = some(where (p.eft == allow))
[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act`)
	casbinEnforcer, _ := casbin.NewEnforcer(casbinModel)

	_, err := casbinEnforcer.AddPolicy("admin", "/users", "GET")
	assert.NoError(t, err)

	router, mockUserService, cfg := setupCorrectIsolatedUserRouter(casbinEnforcer)

	mockedUsers := []models.User{{Model: gorm.Model{ID: 1}, Username: "user1"}}
	mockUserService.On("GetUsers").Return(mockedUsers, nil)

	token, err := utils.GenerateToken(1, "admin", cfg)
	assert.NoError(t, err)

	req, _ := http.NewRequest(http.MethodGet, "/users", http.NoBody)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseUsers []dtos.UserResponse
	err = json.Unmarshal(w.Body.Bytes(), &responseUsers)
	assert.NoError(t, err)
	assert.Len(t, responseUsers, 1)
	mockUserService.AssertExpectations(t)
}

func TestGetUsers_Endpoint_Forbidden_Correct(t *testing.T) {
	casbinModel, _ := model.NewModelFromString(`[request_definition]
r = sub, obj, act
[policy_definition]
p = sub, obj, act
[role_definition]
g = _, _
[policy_effect]
e = some(where (p.eft == allow))
[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act`)
	casbinEnforcer, _ := casbin.NewEnforcer(casbinModel)

	// Add policy for admin, but not for user
	_, err := casbinEnforcer.AddPolicy("admin", "/users", "GET")
	assert.NoError(t, err)

	router, _, cfg := setupCorrectIsolatedUserRouter(casbinEnforcer)

	// Token is for a 'user', who does not have permission
	token, err := utils.GenerateToken(2, "user", cfg)
	assert.NoError(t, err)

	req, _ := http.NewRequest(http.MethodGet, "/users", http.NoBody)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
