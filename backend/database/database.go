package database

import (
	"go-web/config"
	"go-web/models"
	"log"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(cfg *config.Config) {
	var err error
	DB, err = gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")

	// 自动迁移数据模型
	err = DB.AutoMigrate(&models.User{}, &models.Role{}, &gormadapter.CasbinRule{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed")
}