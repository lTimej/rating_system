package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"rating_system/config"
	"rating_system/routes"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func serveSPA(c *gin.Context) {
	// 检查是否是API请求
	if strings.HasPrefix(c.Request.URL.Path, "/api") {
		c.JSON(404, gin.H{"error": "API endpoint not found"})
		return
	}

	// 检查请求是否接受HTML
	accept := c.Request.Header.Get("Accept")
	if !strings.Contains(accept, "text/html") {
		c.Status(404)
		return
	}

	// 尝试读取index.html文件
	content, err := ioutil.ReadFile("dist/index.html")
	if err != nil {
		fmt.Printf("Error reading dist/index.html: %v\n", err)
		// 如果找不到dist/index.html，重定向到根路径
		c.Redirect(302, "/")
		return
	}

	// 返回SPA应用
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, string(content))
}

func main() {
	// 初始化数据库
	config.InitDatabase()

	// 创建Gin实例
	r := gin.Default()

	// 设置文件上传大小限制 (50MB)
	r.MaxMultipartMemory = 50 << 20 // 50 MB

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

	// 设置API路由
	routes.SetupRoutes(r)

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	// 前端静态文件服务 - 使用static中间件服务静态资源
	r.Use(static.Serve("/", static.LocalFile("./static/dist", false)))

	// 处理前端路由 (Vue Router history mode)
	// 这个必须放在最后，作为fallback处理所有未匹配的路由
	r.NoRoute(serveSPA)

	// 启动服务器
	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
