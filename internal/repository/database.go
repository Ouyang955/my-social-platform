package repository

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
		
// InitDB - 初始化MySQL数据库连接
func InitDB() {
	// 连接字符串
	dsn := "root:123456@tcp(localhost:3306)/social_platform?charset=utf8mb4&parseTime=True&loc=Local"

	// 尝试连接数据库
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connection established.")

	// 自动迁移数据库结构（可选）
	DB.AutoMigrate()
}

// CloseDB - 关闭数据库连接
func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get SQL DB instance:", err)
	}
	sqlDB.Close()
}
