package model

import "time"

// 用户模型
type User struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`
	Username    string     `json:"username" gorm:"unique;not null"`
	Password    string     `json:"password"`
	Nickname    string     `json:"nickname"`
	Avatar      string     `json:"avatar"`
	Bio         string     `json:"bio"`                           // 个性签名
	FollowCount int        `json:"follow_count" gorm:"default:0"` // 关注数
	FansCount   int        `json:"fans_count" gorm:"default:0"`   // 粉丝数
	LikeCount   int        `json:"like_count" gorm:"default:0"`   // 获赞数
}

// TableName 自定义表名
func (User) TableName() string {
	return "users"
}
