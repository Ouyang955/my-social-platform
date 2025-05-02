package handler

import (
	"my-social-platform/internal/middleware"
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

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := service.Register(input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

// LoginHandler - 处理用户登录请求
func LoginHandler(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := service.Login(input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := middleware.GenerateJWT(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// ProfileHandler - 获取当前登录用户信息（JWT解析后）
func ProfileHandler(c *gin.Context) {
	// 模拟获取 JWT 中的用户名（你也可以解析 token 并做真实用户查找）
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	parsedToken, _ := middleware.ParseJWT(token)
	claims := parsedToken.Claims.(jwt.MapClaims)

	username := claims["username"].(string)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Welcome to your profile",
		"username": username,
	})
}
