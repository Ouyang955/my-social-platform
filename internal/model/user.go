package model

import "github.com/jinzhu/gorm"

// 用户模型
type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"password"`
}

// TableName 自定义表名
func (User) TableName() string {
	return "users"
}
