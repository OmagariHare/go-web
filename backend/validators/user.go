package validators

import (
	"go-web/database"
	"go-web/models"
)

// CheckUserExists 检查用户是否存在
func CheckUserExists(username, email string) bool {
	var count int64
	database.DB.Model(&models.User{}).Where("username = ? OR email = ?", username, email).Count(&count)
	return count > 0
}

// GetUserByID 根据ID获取用户
func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := database.DB.Preload("Role").First(&user, id).Error
	return &user, err
}