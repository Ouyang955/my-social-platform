package main

import (
	"log"
	"my-social-platform/internal/handler"
	"my-social-platform/internal/middleware"
	"my-social-platform/internal/pkg/logger"
	"my-social-platform/internal/repository"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS 跨域资源共享 https://blog.csdn.net/leah126/article/details/141624726
// 跨域资源共享(CORS)是一种安全策略，允许浏览器从不同的源(域)请求资源
// 浏览器打开网页时 会找前端服务器（3000），得到静态资源、html等信息
// 然后浏览器会运行前端服务器返回的js代码，js代码请求后端api时，会进行跨域请求
// 跨域请求会先进行OPTIONS请求（预检请求），询问后端是否允许跨域请求
// 后端允许跨域请求后，浏览器会进行真正的请求

// 1. 创建默认CORS配置
// 2. 允许的前端源(这里是本地开发服务器地址)
// 3. 允许的HTTP请求头(包括认证所需的Authorization头)
// 4. 允许的HTTP方法(GET查询、POST创建、PUT更新、DELETE删除、OPTIONS预检请求)
// 5. 允许携带Cookie等身份凭证
// 6. 缓存预检请求结果
// 7. 将CORS中间件应用到Gin路由器

func main() {
	// 初始化日志系统
	if err := logger.InitLogger(); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Close()

	// 初始化数据库连接
	repository.InitDB()
	defer repository.CloseDB()

	// 创建gin引擎
	r := gin.Default()

	// 配置CORS(跨域资源共享)
	// 开发环境下允许所有源访问，解决跨域问题
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // 允许前端开发服务器访问
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))

	// 静态文件服务
	// 创建uploads目录（如果不存在）
	os.MkdirAll("uploads/images", 0755)
	r.Static("/uploads", "./uploads")

	// 前端静态文件服务
	// 确保编译后的前端文件存在于build目录
	frontendDir := "./frontend/build"
	if _, err := os.Stat(frontendDir); !os.IsNotExist(err) {
		r.Static("/static", filepath.Join(frontendDir, "static"))
		r.StaticFile("/favicon.ico", filepath.Join(frontendDir, "favicon.ico"))
	} else {
		// 开发环境下使用前端开发服务器
		log.Println("Frontend build directory not found, using API mode only")
	}

	// 注册和登录接口
	r.POST("/register", handler.RegisterHandler)
	r.POST("/login", handler.LoginHandler)

	// 公开的帖子API - 不需要登录也能获取帖子列表
	r.GET("/api/posts", handler.GetAllPostsHandler)

	// 需要认证的路由组
	authorized := r.Group("/api")
	authorized.Use(middleware.JWTAuthMiddleware())
	{
		// 用户资料
		authorized.GET("/profile", handler.ProfileHandler)
		authorized.PUT("/profile", handler.UpdateProfileHandler)

		// 帖子相关API
		authorized.POST("/posts", handler.CreatePostHandler)
		authorized.GET("/posts/:id", handler.GetPostDetailHandler)
		authorized.GET("/user/posts", handler.GetUserPostsHandler)

		// 图片上传接口 - 需要登录才能上传图片
		authorized.POST("/upload/image", handler.FileUploadImageHandler)
	}

	// 公开的图片获取接口 - 不需要登录也能查看图片
	r.GET("/api/images/:filename", handler.GetImageHandler)

	// 前端路由处理 - 将所有未匹配的路由重定向到前端应用
	if _, err := os.Stat(frontendDir); !os.IsNotExist(err) {
		r.NoRoute(func(c *gin.Context) {
			c.File(filepath.Join(frontendDir, "index.html"))
		})
	}

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
