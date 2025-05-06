package handler

import (
	"fmt"
	"my-social-platform/internal/middleware"
	"my-social-platform/internal/model"
	"my-social-platform/internal/pkg/logger"
	"my-social-platform/internal/service"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// RegisterHandler - 处理用户注册请求
func RegisterHandler(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// 获取客户端IP
	clientIP := c.ClientIP()

	// 1. 将请求体绑定到input结构体
	// ShouldBindJSON类似于SpringBoot中的@RequestBody注解
	// 它会自动将JSON请求体解析并映射到input结构体中
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Log(logger.ERROR, "REGISTER", input.Username, clientIP, "Invalid input format")
		// 如果解析失败,返回400错误
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 2. 调用service层的Register方法处理注册逻辑
	// 类似于SpringBoot中注入Service并调用其方法
	user, err := service.Register(input.Username, input.Password)
	if err != nil {
		logger.Log(logger.ERROR, "REGISTER", input.Username, clientIP, "Failed to register user: "+err.Error())
		// 如果注册失败,返回500错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	// 3. 注册成功，记录日志
	logger.Log(logger.INFO, "REGISTER", input.Username, clientIP, "User registered successfully")

	// 3. 注册成功,返回201状态码和用户信息
	// gin.H相当于Java中的Map或ResponseEntity
	// StatusCreated(201)表示资源创建成功
	// 使用gin.Context的JSON方法返回响应
	// - http.StatusCreated (201) 表示资源创建成功的状态码
	// - gin.H 是一个map[string]interface{}的简写,用于构建JSON响应
	// - {"user": user} 创建一个包含user对象的JSON响应体
	//
	// 最终作用:
	// 1. 向客户端返回201状态码,表示用户注册成功
	// 2. 返回一个JSON格式的响应体,包含新注册用户的信息
	// 3. 客户端(如前端页面)可以通过解析这个JSON响应来获取用户信息,
	//    进行后续操作(如自动登录、跳转到个人主页等)
	c.JSON(http.StatusCreated, gin.H{"user": user})
}

// LoginHandler - 处理用户登录请求

// 这个函数处理用户的登录请求,主要做以下几件事:
//
// 1. 解析请求体
// - 定义一个匿名结构体来接收用户名和密码
// - 使用ShouldBindJSON将JSON请求体解析到结构体中
// - 如果解析失败返回400错误
//
// 2. 验证用户身份
// - 调用service.Login验证用户名和密码
// - 如果验证失败返回401未授权错误
//
// 3. 生成JWT令牌
// - 使用middleware.GenerateJWT为用户生成JWT令牌
// - 如果生成失败返回500服务器错误
//
// 4. 返回令牌
// - 登录成功时返回200状态码和JWT令牌
// - 前端可以保存这个令牌用于后续的认证请求
func LoginHandler(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// 获取客户端IP
	clientIP := c.ClientIP()

	// 使用gin框架的ShouldBindJSON方法将请求体中的JSON数据解析到input结构体中
	// 如果解析失败(比如JSON格式错误或缺少必要字段),err会包含具体错误信息
	// 这一步相当于SpringBoot中使用@RequestBody注解自动将JSON请求体绑定到对象
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Log(logger.ERROR, "LOGIN", input.Username, clientIP, "Invalid input format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var userModel model.User
	if err := service.GetUserByUsername(input.Username, &userModel); err != nil {
		logger.Log(logger.WARNING, "LOGIN", input.Username, clientIP, "Login failed: "+err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	userDTO, err := service.Login(input.Username, input.Password)
	if err != nil {
		logger.Log(logger.WARNING, "LOGIN", input.Username, clientIP, "Login failed: "+err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := middleware.GenerateJWT(userModel)
	if err != nil {
		logger.Log(logger.ERROR, "LOGIN", input.Username, clientIP, "Failed to generate token: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	logger.Log(logger.INFO, "LOGIN", input.Username, clientIP, "User logged in successfully")
	c.JSON(http.StatusOK, gin.H{"token": token, "user": userDTO})
}

// ProfileHandler - 获取当前登录用户信息（JWT解析后）
func ProfileHandler(c *gin.Context) {
	// 获取客户端IP
	clientIP := c.ClientIP()

	// 从JWT中获取用户信息
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	parsedToken, err := middleware.ParseJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
		return
	}

	claims := parsedToken.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	userID := uint(claims["user_id"].(float64))

	// 获取用户完整信息
	userProfile, err := service.GetUserProfileByID(userID)
	if err != nil {
		logger.Log(logger.ERROR, "PROFILE", username, clientIP, "获取用户资料失败: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户资料失败"})
		return
	}

	// 获取用户的帖子
	posts, err := service.GetPostsByUserIDService(userID)
	if err != nil {
		logger.Log(logger.ERROR, "PROFILE", username, clientIP, "获取用户帖子失败: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户帖子失败"})
		return
	}

	// 记录访问日志
	logger.Log(logger.INFO, "PROFILE", username, clientIP, "用户访问个人资料")

	// 返回完整的用户资料和帖子
	c.JSON(http.StatusOK, gin.H{
		"user":  userProfile,
		"posts": posts,
	})
}

// GetAllPostsHandler - 获取所有帖子
func GetAllPostsHandler(c *gin.Context) {
	// 获取客户端IP
	clientIP := c.ClientIP()

	// 获取用户信息（如果已登录）
	_, exists := c.Get("user_id")
	username, _ := c.Get("username")

	// 记录访问日志
	if exists {
		logger.Log(logger.INFO, "POSTS", username.(string), clientIP, "User accessed all posts")
	} else {
		logger.Log(logger.INFO, "POSTS", "guest", clientIP, "Guest accessed all posts")
	}

	// 调用服务层获取所有帖子
	posts, err := service.GetAllPostsService()
	if err != nil {
		logger.Log(logger.ERROR, "POSTS", "system", clientIP, "获取帖子失败: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取帖子失败"})
		return
	}

	// 处理图片URL格式，确保前端可以正确显示
	for i := range posts {
		// 如果图片路径是相对路径，转换为绝对URL
		if posts[i].Images != "" && !strings.HasPrefix(posts[i].Images, "http") {
			// 如果是上传文件路径，添加服务器域名
			if strings.HasPrefix(posts[i].Images, "/uploads/") {
				// 在生产环境中应该使用配置的服务器域名
				// posts[i].Images = "https://your-domain.com" + posts[i].Images
				// 开发环境使用本地地址
				posts[i].Images = "http://localhost:8080" + posts[i].Images
			}
		}
	}

	// 记录找到的帖子数量
	logger.Log(logger.INFO, "POSTS", "system", clientIP, fmt.Sprintf("找到 %d 条帖子", len(posts)))

	// 返回帖子列表
	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

// GetUserPostsHandler - 获取指定用户的所有帖子
func GetUserPostsHandler(c *gin.Context) {
	// 从JWT中获取当前登录的用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// 获取客户端IP
	clientIP := c.ClientIP()
	username, _ := c.Get("username")

	// 记录访问日志
	logger.Log(logger.INFO, "USER_POSTS", username.(string), clientIP, "User accessed their posts")

	// 调用服务层获取用户的所有帖子
	posts, err := service.GetPostsByUserIDService(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取帖子失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

// We've moved this functionality to FileUploadImageHandler in file_handler.go
