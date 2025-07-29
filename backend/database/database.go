package database

// 负责处理数据库连接和初始化。

import (
	"go-web/config"
	"go-web/models"
	"log"
	"time"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB 是一个全局的*gorm.DB对象，用于在整个应用程序中访问数据库。
var DB *gorm.DB

// NewPostgresDialector 创建一个新的 gorm.Dialector 用于 PostgreSQL。
// 封装了 postgres.Open，方便在测试中模拟。
func NewPostgresDialector(dsn string) gorm.Dialector {
	return postgres.Open(dsn)
}

// ConnectDB 使用提供的配置初始化数据库连接。
// 它会执行以下操作：
// 1. 使用配置中的DSN（数据源名称）打开一个到PostgreSQL数据库的连接。
// 2. 如果连接失败，将记录致命错误并终止程序。
// 3. 连接成功后，使用GORM的AutoMigrate功能自动迁移数据模型。
//   - 这将确保数据库中的表结构与代码中定义的模型（User, Role, CasbinRule）保持一致。
//
// 4. 如果迁移失败，同样会记录致命错误并终止程序。
func ConnectDB(cfg *config.Config) {
	var err error

	// 使用GORM和PostgreSQL驱动打开数据库连接
	DB, err = gorm.Open(NewPostgresDialector(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		// 如果连接失败，记录致命错误并退出程序
		log.Fatal("无法连接到数据库: ", err)
	}

	log.Println("数据库连接成功")

	// 配置数据库连接池
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("无法获取底层的sql.DB: ", err)
	}
	// 设置连接池中的最大空闲连接数
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	// 设置数据库的最大打开连接数
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	// 设置连接可被重用的最大时间
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Minute)

	// 自动迁移数据模型，确保表结构与模型定义一致
	// AutoMigrate 会创建或更新表以匹配 User, Role 和 CasbinRule 结构体
	err = DB.AutoMigrate(&models.User{}, &models.Role{}, &gormadapter.CasbinRule{})
	if err != nil {
		// 如果迁移失败，记录致命错误并退出程序
		log.Fatal("数据库迁移失败: ", err)
	}

	log.Println("数据库迁移完成")
}
