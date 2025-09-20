package routes

import (
	"rating_system/controllers"
	"rating_system/middleware"
	"rating_system/models"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// 初始化控制器
	authController := &controllers.AuthController{}
	fileController := &controllers.FileController{}
	ratingController := &controllers.RatingController{}
	adminController := &controllers.AdminController{}
	followController := &controllers.FollowController{}

	// 公开路由
	public := r.Group("/api")
	{
		public.POST("/login", authController.Login)
		public.POST("/register", authController.Register)
		public.GET("/ratings/public", ratingController.GetPublicRatings)
	}

	// 需要认证的路由
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// 用户相关
		protected.GET("/profile", authController.GetProfile)
		protected.PUT("/profile", authController.UpdateProfile)

		// 文件相关
		protected.POST("/files", fileController.UploadFile)
		protected.GET("/files/user/:user_id", fileController.GetUserFiles)
		protected.GET("/files/:id", fileController.GetFile)
		protected.GET("/files/:id/download", fileController.DownloadFile)
		protected.DELETE("/files/:id", fileController.DeleteFile)

		// 评价相关
		protected.POST("/ratings", ratingController.CreateRating)
		protected.GET("/ratings/user/:user_id", ratingController.GetUserRatings)
		protected.GET("/ratings/:id", ratingController.GetRating)
		protected.POST("/ratings/:rating_id/feedback", ratingController.CreateFeedback)

		// 关注相关
		protected.POST("/follow/:user_id", followController.FollowUser)
		protected.DELETE("/follow/:user_id", followController.UnfollowUser)
		protected.GET("/following", followController.GetFollowing)
		protected.GET("/followers", followController.GetFollowers)
	}

	// 管理员路由
	admin := r.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.RequireRole(models.RoleAdmin))
	{
		// 用户管理
		admin.GET("/users", adminController.GetUsers)
		admin.GET("/users/:id", adminController.GetUser)
		admin.POST("/users", adminController.CreateUser)
		admin.PUT("/users/:id", adminController.UpdateUser)
		admin.DELETE("/users/:id", adminController.DeleteUser)

		// 评价管理
		admin.GET("/ratings", adminController.GetRatings)
		admin.DELETE("/ratings/:id", adminController.DeleteRating)

		// 统计信息
		admin.GET("/stats", adminController.GetStats)
	}
}
