package main

import (
	"fmt"
	"laravel-go/framework/core"
	"laravel-go/framework/database"
	"laravel-go/framework/http"
	"laravel-go/framework/queue"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("=== Laravel-Go 博客系统示例 ===")

	// 初始化应用
	app := core.NewApplication()

	// 加载配置
	app.LoadConfig("config")

	// 初始化数据库
	db, err := database.NewConnection(app.Config.Get("database.default").(string))
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	app.Container.Singleton("db", db)

	// 初始化队列
	queue.Init()
	memoryQueue := queue.NewMemoryQueue()
	queue.QueueManager.Extend("memory", memoryQueue)
	queue.QueueManager.SetDefaultQueue("memory")

	// 启动队列工作进程
	worker := queue.NewWorker(memoryQueue, "default")
	worker.SetOnCompleted(func(job queue.Job) {
		fmt.Printf("任务完成: %s\n", string(job.GetPayload()))
	})
	worker.SetOnFailed(func(job queue.Job, err error) {
		fmt.Printf("任务失败: %s - %v\n", string(job.GetPayload()), err)
	})

	if err := worker.Start(); err != nil {
		log.Fatal("启动队列工作进程失败:", err)
	}

	// 创建HTTP服务器
	server := http.NewServer(app)

	// 注册路由
	registerRoutes(server)

	// 启动服务器
	go func() {
		port := app.Config.Get("app.port", "8080").(string)
		fmt.Printf("博客系统启动在 http://localhost:%s\n", port)
		if err := server.Start(":" + port); err != nil {
			log.Fatal("启动服务器失败:", err)
		}
	}()

	// 等待中断信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("正在关闭博客系统...")

	// 停止队列工作进程
	if err := worker.Stop(); err != nil {
		log.Printf("停止队列工作进程失败: %v", err)
	}

	// 关闭数据库连接
	if err := db.Close(); err != nil {
		log.Printf("关闭数据库连接失败: %v", err)
	}

	fmt.Println("博客系统已关闭")
}

func registerRoutes(server *http.Server) {
	// 注册控制器
	authController := &AuthController{}
	postController := &PostController{}
	userController := &UserController{}

	// Web 路由
	server.Get("/", func(ctx *http.Context) {
		ctx.JSON(200, map[string]interface{}{
			"message": "欢迎使用 Laravel-Go 博客系统",
			"version": "1.0.0",
		})
	})

	// 认证路由
	server.Post("/auth/register", authController.Register)
	server.Post("/auth/login", authController.Login)
	server.Post("/auth/logout", authController.Logout)

	// 文章路由
	server.Get("/posts", postController.Index)
	server.Get("/posts/:id", postController.Show)
	server.Post("/posts", postController.Store)
	server.Put("/posts/:id", postController.Update)
	server.Delete("/posts/:id", postController.Destroy)

	// 用户路由
	server.Get("/users/profile", userController.Profile)
	server.Put("/users/profile", userController.UpdateProfile)
} 