package service

import (
	"errors"
	"my-social-platform/internal/model"
	"my-social-platform/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword - 哈希密码
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword - 校验密码
func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// Register - 注册用户
func Register(username, password string) (*model.User, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: username,
		Password: hashedPassword,
	}

	if err := repository.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Login - 用户登录
func Login(username, password string) (*model.User, error) {
	var user model.User
	if err := repository.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	if !VerifyPassword(user.Password, password) {
		return nil, errors.New("invalid password")
	}

	return &user, nil
}
