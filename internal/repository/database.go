package repository

import (
	"log"
	"my-social-platform/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB - 初始化MySQL数据库连接
func InitDB() {
	// 连接字符串
	dsn := "root:123456@tcp(127.0.0.1:3306)/social_platform?charset=utf8mb4&parseTime=True&loc=Local"

	// 尝试连接数据库
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 迁移时禁用外键约束
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connection established.")

	// 检查数据库中是否已存在所需的表
	if !DB.Migrator().HasTable(&model.User{}) ||
		!DB.Migrator().HasTable(&model.Post{}) ||
		!DB.Migrator().HasTable(&model.Comment{}) {
		log.Println("Tables do not exist, creating schema...")
		// 自动迁移表结构 - 只更新表结构，不删除数据
		err = DB.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{})
		if err != nil {
			log.Fatal("Failed to migrate database:", err)
		}
		log.Println("Database migration completed successfully.")
	} else {
		log.Println("Tables already exist, skipping migration.")
	}
}

// CloseDB - 关闭数据库连接
func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get SQL DB instance:", err)
	}
	sqlDB.Close()
}
