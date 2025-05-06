package repository

import "my-social-platform/internal/model"

// 新建帖子
func CreatePost(post *model.Post) error {
	return DB.Create(post).Error
}

// 根据帖子id查询帖子
func GetPostByID(id uint) (*model.Post, error) {
	var post model.Post
	err := DB.First(&post, id).Error
	return &post, err
}

// 获取所有帖子
func GetAllPosts() ([]*model.Post, error) {
	var posts []*model.Post
	err := DB.Order("created_at DESC").Find(&posts).Error
	return posts, err
}

// 根据用户ID获取帖子
func GetPostsByUserID(userID uint) ([]*model.Post, error) {
	var posts []*model.Post
	err := DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&posts).Error
	return posts, err
}
