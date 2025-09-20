package config

import (
	"log"
	"rating_system/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase() {
	var err error

	// 连接SQLite数据库
	DB, err = gorm.Open(sqlite.Open("rating_system.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移数据库表
	err = DB.AutoMigrate(
		&models.User{},
		&models.File{},
		&models.Rating{},
		&models.Feedback{},
		&models.Follow{},
		&models.FileComment{},
		&models.Article{},
		&models.ArticleComment{},
		&models.ArticleLike{},
		&models.ArticleCommentRating{},
		&models.ArticlePush{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 创建默认管理员账户
	createDefaultAdmin()

	log.Println("Database initialized successfully")
}

func createDefaultAdmin() {
	var count int64
	DB.Model(&models.User{}).Where("role = ?", models.RoleAdmin).Count(&count)

	if count == 0 {
		admin := models.User{
			Username: "admin",
			Email:    "admin@example.com",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Role:     models.RoleAdmin,
			Name:     "系统管理员",
			IsActive: true,
		}

		if err := DB.Create(&admin).Error; err != nil {
			log.Printf("Failed to create default admin: %v", err)
		} else {
			log.Println("Default admin created: username=admin, password=password")
		}
	}
}
