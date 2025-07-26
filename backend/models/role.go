// package models 定义了应用程序中使用的数据结构，这些结构将映射到数据库中的表。
package models

import (
	"gorm.io/gorm"
)

// Role 代表系统中的一个角色。
// 每个用户都将被分配一个角色，用于权限控制。
type Role struct {
	gorm.Model         // gorm.Model 包含了ID, CreatedAt, UpdatedAt, DeletedAt等标准字段
	Name        string `gorm:"uniqueIndex;not null" json:"name"` // 角色的唯一名称，例如 "admin", "user"
	Description string `json:"description"`                      // 对角色的简单描述
}
