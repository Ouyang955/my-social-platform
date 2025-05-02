package repository

import (
	"log"
	"my-social-platform/internal/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB *gorm.DB

// InitDB 初始化数据库的链接
func InitDB() {
	var err error
	DB, err = gorm.Open("sqlite3", "social_platform.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移数据库
	DB.AutoMigrate(&model.User{})
}

// 关闭数据库
func CloseDB() {
	DB.Close()
}
