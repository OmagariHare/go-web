// package models 定义了应用程序中使用的数据结构，这些结构将映射到数据库中的表。
package models

import (
	"gorm.io/gorm"
)

// User 代表系统中的一个用户。
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null" json:"username"` // 用户的唯一名称
	Email    string `gorm:"uniqueIndex;not null" json:"email"`    // 用户的唯一电子邮件地址
	Password string `gorm:"not null" json:"-"`                    // 用户的哈希密码，json:"-" 表示在JSON序列化时忽略此字段
	RoleID   uint   `gorm:"not null" json:"role_id"`              // 关联的角色ID
	Role     Role   `json:"role"`                                 // 用户所属的角色（通过RoleID进行关联）
}
