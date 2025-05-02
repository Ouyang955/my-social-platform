package model

import "github.com/jinzhu/gorm"

// 用户模型
type User struct {
	gorm.Model
	UserName string `json:"username" gorm:"unique;not null"`
	PassWord string `json:"password"`
}

// TableName 自定义表名
func (User) TableName() string {
	return "users"
}
