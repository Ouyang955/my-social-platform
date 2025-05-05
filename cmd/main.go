package main

import (
	"log"
	"my-social-platform/internal/handler"
	"my-social-platform/internal/middleware"
	"my-social-platform/internal/pkg/logger"
	"my-social-platform/internal/repository"

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
// 5. 将CORS中间件应用到Gin路由器

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
	// 这段代码配置了跨域资源共享(CORS)策略，允许前端应用与后端API进行安全通信
	// 1. 创建默认CORS配置
	config := cors.DefaultConfig()
	// 2. 允许的前端源(这里是本地开发服务器地址)
	config.AllowOrigins = []string{"http://localhost:3000"}
	// 3. 允许的HTTP请求头(包括认证所需的Authorization头)
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	// 4. 允许的HTTP方法(GET查询、POST创建、PUT更新、DELETE删除、OPTIONS预检请求)
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	// 5. 将CORS中间件应用到Gin路由器
	r.Use(cors.New(config))

	// 注册和登录接口
	r.POST("/register", handler.RegisterHandler)
	r.POST("/login", handler.LoginHandler)

	// 需要认证的路由组
	// 1. 创建一个名为"/api"的路由组,所有需要认证的API都将放在这个组下
	// 2. 使用middleware.JWTAuthMiddleware()中间件进行JWT令牌验证
	//    - 该中间件会检查请求头中的Authorization字段
	//    - 验证JWT令牌的有效性和完整性
	//    - 如果令牌无效或过期,会阻止请求并返回401未授权错误
	//    - 如果令牌有效,会将用户信息存储在上下文中并允许请求继续
	// 3. 大括号内定义了所有需要认证才能访问的API端点
	authorized := r.Group("/api")
	authorized.Use(middleware.JWTAuthMiddleware())
	// 在这个路由组下的接口:
	// - GET /api/profile: 获取当前登录用户的个人资料
	// 这些接口都需要有效的JWT令牌才能访问
	// JWT令牌保存了用户的登录信息，通过验证令牌可以确认用户的身份和登录状态
	// 所有需要用户登录后才能访问的功能都应该放在这个路由组下
	{
		// 保护的API接口
		// GET /api/profile - 获取当前登录用户的个人资料
		// 这个端点只有在请求中包含有效的JWT令牌时才能访问
		// 处理函数ProfileHandler会从JWT令牌中提取用户信息并返回
		authorized.GET("/profile", handler.ProfileHandler)
	}

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
