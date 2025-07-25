package repositories

import (
	"go-web/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByUsername(username string) (*models.User, error)
	FindByUsernameOrEmail(username, email string) (*models.User, error)
	Create(user *models.User) error
	FindAll() ([]models.User, error)
	FindByID(id uint) (*models.User, error)
	Update(user *models.User) error
	Delete(user *models.User) error
	LoadRole(user *models.User) error
}

type GormUserRepository struct {
	DB *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{DB: db}
}

func (r *GormUserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) FindByUsernameOrEmail(username, email string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("username = ? OR email = ?", username, email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *GormUserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	if err := r.DB.Preload("Role").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *GormUserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.DB.Preload("Role").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) Update(user *models.User) error {
	return r.DB.Save(user).Error
}

func (r *GormUserRepository) Delete(user *models.User) error {
	return r.DB.Delete(user).Error
}

func (r *GormUserRepository) LoadRole(user *models.User) error {
	return r.DB.Model(user).Association("Role").Find(&user.Role)
}