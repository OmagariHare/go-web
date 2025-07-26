// package repositories 提供了数据访问的抽象层。
// 它定义了与数据库交互的接口，以及使用GORM实现的具体逻辑。
package repositories

import (
	"go-web/models"

	"gorm.io/gorm"
)

// UserRepository 定义了与用户数据相关的操作接口。
// 这种接口定义方式有利于实现依赖倒置，方便进行单元测试。
type UserRepository interface {
	// FindByUsername 根据用户名查找用户。
	FindByUsername(username string) (*models.User, error)
	// FindByUsernameOrEmail 根据用户名或邮箱查找用户。
	FindByUsernameOrEmail(username, email string) (*models.User, error)
	// Create 创建一个新用户。
	Create(user *models.User) error
	// FindAll 获取所有用户列表。
	FindAll() ([]models.User, error)
	// FindByID 根据用户ID查找用户。
	FindByID(id uint) (*models.User, error)
	// Update 更新一个已存在的用户信息。
	Update(user *models.User) error
	// Delete 删除一个用户。
	Delete(user *models.User) error
	// LoadRole 加载用户的角色信息。
	LoadRole(user *models.User) error
}

// GormUserRepository 是 UserRepository 的GORM实现。
// 它持有一个 *gorm.DB 实例，用于执行数据库操作。
type GormUserRepository struct {
	DB *gorm.DB
}

// NewGormUserRepository 是一个构造函数，用于创建一个新的 GormUserRepository 实例。
func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{DB: db}
}

// FindByUsername 实现了 UserRepository 接口的 FindByUsername 方法。
func (r *GormUserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsernameOrEmail 实现了 UserRepository 接口的 FindByUsernameOrEmail 方法。
func (r *GormUserRepository) FindByUsernameOrEmail(username, email string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("username = ? OR email = ?", username, email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 实现了 UserRepository 接口的 Create 方法。
func (r *GormUserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

// FindAll 实现了 UserRepository 接口的 FindAll 方法。
// 使用 Preload("Role") 来预加载关联的角色信息，避免N+1查询问题。
func (r *GormUserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	if err := r.DB.Preload("Role").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// FindByID 实现了 UserRepository 接口的 FindByID 方法。
// 同样使用 Preload("Role") 来预加载角色信息。
func (r *GormUserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.DB.Preload("Role").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 实现了 UserRepository 接口的 Update 方法。
// GORM的Save方法会自动更新记录的所有字段。
func (r *GormUserRepository) Update(user *models.User) error {
	return r.DB.Save(user).Error
}

// Delete 实现了 UserRepository 接口的 Delete 方法。
func (r *GormUserRepository) Delete(user *models.User) error {
	return r.DB.Delete(user).Error
}

// LoadRole 实现了 UserRepository 接口的 LoadRole 方法。
// 它使用GORM的Association方法来显式加载用户关联的角色信息。
func (r *GormUserRepository) LoadRole(user *models.User) error {
	return r.DB.Model(user).Association("Role").Find(&user.Role)
}
