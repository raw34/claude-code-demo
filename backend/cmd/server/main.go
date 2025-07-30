package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/user/user-management/internal/config"
	"github.com/user/user-management/internal/database"
	"github.com/user/user-management/internal/handlers"
	"github.com/user/user-management/internal/middleware"
	"github.com/user/user-management/internal/repository"
	"github.com/user/user-management/internal/service"
)

func main() {
	// 加载环境变量，优先级：
	// 1. 当前目录的 .env
	// 2. 项目根目录的 .env
	if err := godotenv.Load(); err != nil {
		// 尝试加载项目根目录的 .env
		if err := godotenv.Load("../.env"); err != nil {
			log.Println("No .env file found in current or parent directory")
		}
	}

	// 加载配置
	cfg := config.Load()

	// 连接数据库
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 连接Redis
	redisClient, err := database.ConnectRedis(cfg.Redis)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redisClient.Close()

	// 初始化仓库
	userRepo := repository.NewUserRepository(db)

	// 初始化服务
	sessionService := service.NewSessionService(redisClient)
	authService := service.NewAuthService(userRepo, sessionService, cfg.JWT.Secret, cfg.JWT.AccessTokenExpiry)
	userService := service.NewUserService(userRepo)

	// 初始化处理器
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	// 设置Gin模式
	gin.SetMode(os.Getenv("GIN_MODE"))

	// 创建路由
	router := gin.Default()

	// 中间件
	router.Use(middleware.CORS())
	router.Use(middleware.ErrorHandler())

	// API路由组
	api := router.Group("/api/v1")
	{
		// 认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", middleware.Auth(authService), authHandler.Logout)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// 用户路由（需要认证）
		users := api.Group("/users")
		users.Use(middleware.Auth(authService))
		{
			users.GET("", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
			users.GET("/profile", userHandler.GetProfile)
			users.PUT("/profile", userHandler.UpdateProfile)
		}
	}

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		// 检查Redis连接
		_, err := redisClient.Ping(c).Result()
		if err != nil {
			c.JSON(500, gin.H{"status": "unhealthy", "redis": "down"})
			return
		}
		
		// 检查数据库连接
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			c.JSON(500, gin.H{"status": "unhealthy", "database": "down"})
			return
		}
		
		c.JSON(200, gin.H{"status": "healthy", "database": "up", "redis": "up"})
	})

	// 启动服务器
	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}