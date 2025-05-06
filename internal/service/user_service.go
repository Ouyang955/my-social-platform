package service

import (
	"my-social-platform/internal/dto"
	"my-social-platform/internal/repository"
)

// GetUserProfileByID 根据用户ID获取完整个人资料
func GetUserProfileByID(id uint) (*dto.UserDTO, error) {
	user, err := repository.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	return &dto.UserDTO{
		ID:          user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Bio:         user.Bio,
		FollowCount: user.FollowCount,
		FansCount:   user.FansCount,
		LikeCount:   user.LikeCount,
	}, nil
}

// GetUserProfileByUsername 根据用户名获取完整个人资料
func GetUserProfileByUsername(username string) (*dto.UserDTO, error) {
	user, err := repository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	return &dto.UserDTO{
		ID:          user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Bio:         user.Bio,
		FollowCount: user.FollowCount,
		FansCount:   user.FansCount,
		LikeCount:   user.LikeCount,
	}, nil
}

// UpdateUserProfile 更新用户资料
func UpdateUserProfile(userID uint, nickname string, bio string) error {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return err
	}

	// 更新字段
	user.Nickname = nickname
	user.Bio = bio

	// 保存更改
	return repository.UpdateUserProfile(user)
}

// UpdateAvatar 更新用户头像
func UpdateAvatar(userID uint, avatarURL string) error {
	return repository.UpdateAvatar(userID, avatarURL)
}
