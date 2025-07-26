// package repositories 提供了数据访问的抽象层。
// 它定义了与数据库交互的接口，以及使用GORM实现的具体逻辑。
package repositories

import (
	"go-web/models"

	"gorm.io/gorm"
)

// RoleRepository 定义了与角色数据相关的操作接口。
// 这种接口定义方式有利于实现依赖倒置，方便进行单元测试（可以使用mock实现）。
type RoleRepository interface {
	// FindByName 根据角色名称查找角色。
	FindByName(name string) (*models.Role, error)
	// Create 创建一个新的角色。
	Create(role *models.Role) error
}

// GormRoleRepository 是 RoleRepository 的GORM实现。
// 它持有一个 *gorm.DB 实例，用于执行数据库操作。
type GormRoleRepository struct {
	DB *gorm.DB
}

// NewGormRoleRepository 是一个构造函数，用于创建一个新的 GormRoleRepository 实例。
func NewGormRoleRepository(db *gorm.DB) *GormRoleRepository {
	return &GormRoleRepository{DB: db}
}

// FindByName 实现了 RoleRepository 接口的 FindByName 方法。
// 它通过GORM的Where子句查询指定名称的角色。
func (r *GormRoleRepository) FindByName(name string) (*models.Role, error) {
	var role models.Role
	if err := r.DB.Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// Create 实现了 RoleRepository 接口的 Create 方法。
// 它使用GORM的Create方法将一个新的角色记录插入到数据库中。
func (r *GormRoleRepository) Create(role *models.Role) error {
	return r.DB.Create(role).Error
}
