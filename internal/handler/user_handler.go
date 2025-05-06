package handler

import (
	"my-social-platform/internal/pkg/logger"
	"my-social-platform/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UpdateProfileHandler 处理用户资料更新请求
func UpdateProfileHandler(c *gin.Context) {
	// 获取当前登录用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// 获取客户端IP
	clientIP := c.ClientIP()
	username, _ := c.Get("username")

	// 解析请求体
	var input struct {
		Nickname string `json:"nickname"`
		Bio      string `json:"bio"`
		Avatar   string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Log(logger.ERROR, "UPDATE_PROFILE", username.(string), clientIP, "无效的请求格式")
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求格式"})
		return
	}

	// 更新用户资料
	err := service.UpdateUserProfile(userID.(uint), input.Nickname, input.Bio)
	if err != nil {
		logger.Log(logger.ERROR, "UPDATE_PROFILE", username.(string), clientIP, "更新资料失败: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新资料失败"})
		return
	}

	// 如果提供了头像URL，更新头像
	if input.Avatar != "" {
		err = service.UpdateAvatar(userID.(uint), input.Avatar)
		if err != nil {
			logger.Log(logger.ERROR, "UPDATE_PROFILE", username.(string), clientIP, "更新头像失败: "+err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新头像失败"})
			return
		}
	}

	// 获取更新后的用户资料
	updatedProfile, err := service.GetUserProfileByID(userID.(uint))
	if err != nil {
		logger.Log(logger.ERROR, "UPDATE_PROFILE", username.(string), clientIP, "获取更新后的资料失败: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取更新后的资料失败"})
		return
	}

	logger.Log(logger.INFO, "UPDATE_PROFILE", username.(string), clientIP, "用户资料更新成功")
	c.JSON(http.StatusOK, gin.H{"message": "资料更新成功", "user": updatedProfile})
}
