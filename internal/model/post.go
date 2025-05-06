package model

import (
	"time"
)

type Post struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at" gorm:"index"`
	UserID       uint       `json:"user_id" gorm:"index"`           // 发帖用户ID
	Content      string     `json:"content"`                        // 帖子文字内容
	Images       string     `json:"images"`                         // 帖子图片URL数组
	Tag          string     `json:"tag"`                            // 帖子标签
	LikeCount    int        `json:"like_count" gorm:"default:0"`    // 帖子点赞数
	FavCount     int        `json:"fav_count" gorm:"default:0"`     // 帖子收藏数
	CommentCount int        `json:"comment_count" gorm:"default:0"` // 评论数
}

// 表名：post
func (Post) TableName() string {
	return "post"
}
