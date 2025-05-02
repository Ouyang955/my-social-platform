package main

import (
	"my-social-platform/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 路由设置
	r.GET("/", handler.HelloHandler)

	// 启动服务器
	r.Run(":8080")
}
