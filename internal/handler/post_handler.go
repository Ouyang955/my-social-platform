package handler

import (
	"my-social-platform/internal/model"
	"my-social-platform/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 处理发帖请求
func CreatePostHandler(c *gin.Context) {
	var post model.Post
	// 绑定JSON请求体到post中
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}
	// 登陆时已经通过JWT中间件在上下文中注入了用户信息
	userID, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	post.UserID = userID.(uint)

	// 调用服务层创建帖子
	err := service.CreatePostService(&post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建帖子失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "帖子创建成功", "post_id": post.ID})
}

// 根据用户id查找帖子
func GetPostDetailHandler(c *gin.Context) {
	// 拿到用户id
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	post, err := service.GetPostByIDService(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "帖子不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"post": post})
}
