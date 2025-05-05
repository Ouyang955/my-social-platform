package handler

import (
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

	// 模拟获取 JWT 中的用户名（你也可以解析 token 并做真实用户查找）
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	parsedToken, _ := middleware.ParseJWT(token)
	claims := parsedToken.Claims.(jwt.MapClaims)

	username := claims["username"].(string)

	// 记录访问日志
	logger.Log(logger.INFO, "PROFILE", username, clientIP, "User accessed profile")

	c.JSON(http.StatusOK, gin.H{
		"message":  "Welcome to your profile",
		"username": username,
	})
}
