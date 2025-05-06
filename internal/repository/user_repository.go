package repository

import "my-social-platform/internal/model"

// GetUserByID 根据用户ID获取用户信息
func GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := DB.First(&user, id).Error
	return &user, err
}

// GetUserByUsername 根据用户名获取用户信息
func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

// UpdateUserProfile 更新用户个人资料
func UpdateUserProfile(user *model.User) error {
	return DB.Save(user).Error
}

// UpdateAvatar 更新用户头像
func UpdateAvatar(userID uint, avatarURL string) error {
	return DB.Model(&model.User{}).Where("id = ?", userID).Update("avatar", avatarURL).Error
}

// UpdateUserBio 更新用户个性签名
func UpdateUserBio(userID uint, bio string) error {
	return DB.Model(&model.User{}).Where("id = ?", userID).Update("bio", bio).Error
}
