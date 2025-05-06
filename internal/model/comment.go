package model

import (
	"time"
)

// Comment 评论模型
type Comment struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
	PostID    uint       `json:"post_id" gorm:"index"`        // 关联的帖子ID
	UserID    uint       `json:"user_id" gorm:"index"`        // 评论用户ID
	Content   string     `json:"content"`                     // 评论内容
	LikeCount int        `json:"like_count" gorm:"default:0"` // 评论点赞数
}

// TableName 自定义表名
func (Comment) TableName() string {
	return "comment"
}
