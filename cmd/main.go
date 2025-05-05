package main

import (
	"log"
	"my-social-platform/internal/handler"
	"my-social-platform/internal/middleware"
	"my-social-platform/internal/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库连接
	repository.InitDB()
	defer repository.CloseDB()

	// 创建gin引擎
	r := gin.Default()

	// 配置CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	r.Use(cors.New(config))

	// 注册和登录接口
	r.POST("/register", handler.RegisterHandler)
	r.POST("/login", handler.LoginHandler)

	// 需要认证的路由
	authorized := r.Group("/api")
	authorized.Use(middleware.JWTAuthMiddleware())
	{
		// 保护的API接口
		authorized.GET("/profile", handler.ProfileHandler)
	}

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
