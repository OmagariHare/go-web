package models

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex;not null" json:"name"`
	Description string `json:"description"`
}