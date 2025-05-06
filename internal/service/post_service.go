package service

import (
	"my-social-platform/internal/model"
	"my-social-platform/internal/repository"
)

// 发帖的业务逻辑
func CreatePostService(post *model.Post) error {
	// 可加参数校验，内容审核等
	return repository.CreatePost(post)
}

// 根据id查找帖子
func GetPostByIDService(id uint) (*model.Post, error) {
	return repository.GetPostByID(id)
}

// 获取所有帖子
func GetAllPostsService() ([]*model.Post, error) {
	return repository.GetAllPosts()
}

// 根据用户ID获取帖子
func GetPostsByUserIDService(userID uint) ([]*model.Post, error) {
	return repository.GetPostsByUserID(userID)
}
