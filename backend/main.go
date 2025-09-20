package main

import (
	"log"
	"rating_system/config"
	"rating_system/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	config.InitDatabase()

	// 创建Gin实例
	r := gin.Default()

	// 添加CORS中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	// 设置路由
	routes.SetupRoutes(r)

	// 启动服务器
	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
