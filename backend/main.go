package main

import (
	"io/ioutil"
	"log"
	"rating_system/config"
	"rating_system/routes"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func buildFront(engine *gin.Engine) {
	engine.NoRoute(func(c *gin.Context) {
		accept := c.Request.Header.Get("Accept")
		flag := strings.Contains(accept, "text/html")
		if flag {
			content, err := ioutil.ReadFile("dist/index.html")
			if (err) != nil {
				c.Writer.WriteHeader(404)
				c.Writer.WriteString("Not Found")
				return
			}
			c.Writer.WriteHeader(200)
			c.Writer.Header().Add("Accept", "text/html")
			c.Writer.Write((content))
			c.Writer.Flush()
		}
	})
	engine.Use(static.Serve("/", static.LocalFile("static/dist", true)))
}

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
	buildFront(r)

	// 设置API路由
	routes.SetupRoutes(r)

	// // 静态文件服务
	// r.Static("/uploads", "./uploads")

	// // 前端静态文件服务
	// r.Static("/static", "./static")
	// r.StaticFile("/favicon.ico", "./static/favicon.ico")

	// // 根路径处理
	// r.GET("/", func(c *gin.Context) {
	// 	c.File("./static/index.html")
	// })

	// // 处理前端路由 (Vue Router history mode)
	// r.NoRoute(func(c *gin.Context) {
	// 	// 如果是API请求，返回404
	// 	if strings.HasPrefix(c.Request.URL.Path, "/api") {
	// 		c.JSON(404, gin.H{"error": "API endpoint not found"})
	// 		return
	// 	}
	// 	// 否则返回前端应用
	// 	c.File("./static/index.html")
	// })

	// 启动服务器
	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
