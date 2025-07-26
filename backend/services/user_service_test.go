package services

import (
	"go-web/mocks"
	"go-web/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetUsers_Success(t *testing.T) {
	// 1. Setup
	mockUserRepo := new(mocks.MockUserRepository)
	userService := NewUserService(mockUserRepo)

	// 2. Define Mock Expectations
	mockedUsers := []models.User{
		{Model: gorm.Model{ID: 1}, Username: "user1"},
		{Model: gorm.Model{ID: 2}, Username: "user2"},
	}
	mockUserRepo.On("FindAll").Return(mockedUsers, nil)

	// 3. Execution
	users, err := userService.GetUsers()

	// 4. Assertions
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 2)
	assert.Equal(t, "user1", users[0].Username)
	mockUserRepo.AssertExpectations(t)
}

func TestGetUser_Success(t *testing.T) {
	// 1. Setup
	mockUserRepo := new(mocks.MockUserRepository)
	userService := NewUserService(mockUserRepo)

	userID := uint(1)

	// 2. Define Mock Expectations
	mockedUser := &models.User{Model: gorm.Model{ID: userID}, Username: "testuser"}
	mockUserRepo.On("FindByID", userID).Return(mockedUser, nil)

	// 3. Execution
	user, err := userService.GetUser(userID)

	// 4. Assertions
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	mockUserRepo.AssertExpectations(t)
}
