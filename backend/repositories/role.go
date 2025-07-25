package repositories

import (
	"go-web/models"

	"gorm.io/gorm"
)

type RoleRepository interface {
	FindByName(name string) (*models.Role, error)
	Create(role *models.Role) error
}

type GormRoleRepository struct {
	DB *gorm.DB
}

func NewGormRoleRepository(db *gorm.DB) *GormRoleRepository {
	return &GormRoleRepository{DB: db}
}

func (r *GormRoleRepository) FindByName(name string) (*models.Role, error) {
	var role models.Role
	if err := r.DB.Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *GormRoleRepository) Create(role *models.Role) error {
	return r.DB.Create(role).Error
}
