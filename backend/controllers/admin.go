package controllers

import (
	"fmt"
	"net/http"
	"rating_system/config"
	"rating_system/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AdminController struct{}

func (ac *AdminController) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	role := c.Query("role")

	offset := (page - 1) * limit

	var users []models.User
	var total int64

	query := config.DB.Model(&models.User{})
	if role != "" {
		query = query.Where("role = ?", role)
	}

	query.Count(&total)

	if err := query.Offset(offset).Limit(limit).
		Preload("Files").Preload("Follows").Preload("Followers").
		Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (ac *AdminController) GetUser(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := config.DB.Preload("Files").Preload("Ratings").Preload("Received").
		Preload("Follows.Followed").Preload("Followers.Follower").
		First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (ac *AdminController) CreateUser(c *gin.Context) {
	var req struct {
		Username string          `json:"username" binding:"required"`
		Email    string          `json:"email" binding:"required,email"`
		Password string          `json:"password" binding:"required,min=6"`
		Role     models.UserRole `json:"role" binding:"required"`
		Name     string          `json:"name" binding:"required"`
		IsActive bool            `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名和邮箱是否已存在
	var existingUser models.User
	if err := config.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
		Name:     req.Name,
		IsActive: req.IsActive,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}

func (ac *AdminController) UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var req struct {
		Username string          `json:"username"`
		Email    string          `json:"email"`
		Role     models.UserRole `json:"role"`
		Name     string          `json:"name"`
		IsActive *bool           `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 检查用户名和邮箱是否被其他用户使用
	if req.Username != "" && req.Username != user.Username {
		var existingUser models.User
		if err := config.DB.Where("username = ? AND id != ?", req.Username, userID).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}
		user.Username = req.Username
	}

	if req.Email != "" && req.Email != user.Email {
		var existingUser models.User
		if err := config.DB.Where("email = ? AND id != ?", req.Email, userID).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
		user.Email = req.Email
	}

	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    user,
	})
}

func (ac *AdminController) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 软删除用户
	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

func (ac *AdminController) GetRatings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	offset := (page - 1) * limit

	var ratings []models.Rating
	var total int64

	config.DB.Model(&models.Rating{}).Count(&total)

	if err := config.DB.Offset(offset).Limit(limit).
		Preload("Rater").Preload("Rated").Preload("File").Preload("Feedbacks").
		Find(&ratings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get ratings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ratings": ratings,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

func (ac *AdminController) DeleteRating(c *gin.Context) {
	ratingID := c.Param("id")

	var rating models.Rating
	if err := config.DB.First(&rating, ratingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rating not found"})
		return
	}

	// 删除相关的反馈
	config.DB.Where("rating_id = ?", ratingID).Delete(&models.Feedback{})

	// 删除评价
	if err := config.DB.Delete(&rating).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete rating"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Rating deleted successfully",
	})
}

func (ac *AdminController) GetStats(c *gin.Context) {
	var stats struct {
		TotalUsers    int64 `json:"total_users"`
		TotalStudents int64 `json:"total_students"`
		TotalExperts  int64 `json:"total_experts"`
		TotalRatings  int64 `json:"total_ratings"`
		TotalFiles    int64 `json:"total_files"`
		TotalArticles int64 `json:"total_articles"`
	}

	config.DB.Model(&models.User{}).Count(&stats.TotalUsers)
	config.DB.Model(&models.User{}).Where("role = ?", models.RoleStudent).Count(&stats.TotalStudents)
	config.DB.Model(&models.User{}).Where("role = ?", models.RoleExpert).Count(&stats.TotalExperts)
	config.DB.Model(&models.Rating{}).Count(&stats.TotalRatings)
	config.DB.Model(&models.File{}).Count(&stats.TotalFiles)
	config.DB.Model(&models.Article{}).Count(&stats.TotalArticles)

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}

// GetArticles 获取文章列表（管理员）
func (ac *AdminController) GetArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")

	offset := (page - 1) * limit

	var articles []models.Article
	var total int64

	query := config.DB.Model(&models.Article{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	if err := query.Offset(offset).Limit(limit).
		Preload("Author").
		Order("created_at DESC").
		Find(&articles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get articles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"articles": articles,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

// PushArticle 推送文章给用户
func (ac *AdminController) PushArticle(c *gin.Context) {
	adminID, _ := c.Get("user_id")

	var req struct {
		ArticleID uint   `json:"article_id" binding:"required"`
		UserIDs   []uint `json:"user_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查文章是否存在
	var article models.Article
	if err := config.DB.First(&article, req.ArticleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	// 检查用户是否存在
	var users []models.User
	if err := config.DB.Where("id IN ?", req.UserIDs).Find(&users).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some users not found"})
		return
	}

	if len(users) != len(req.UserIDs) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some users not found"})
		return
	}

	// 创建推送记录
	var pushes []models.ArticlePush
	var skippedUsers []uint
	var skippedCount int

	for _, userID := range req.UserIDs {
		// 检查是否已经推送过
		var existingPush models.ArticlePush
		if err := config.DB.Where("article_id = ? AND user_id = ?", req.ArticleID, userID).First(&existingPush).Error; err == nil {
			skippedUsers = append(skippedUsers, userID)
			skippedCount++
			continue // 已经推送过，跳过
		}

		pushes = append(pushes, models.ArticlePush{
			ArticleID: req.ArticleID,
			UserID:    userID,
			AdminID:   adminID.(uint),
		})
	}

	if len(pushes) > 0 {
		if err := config.DB.Create(&pushes).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to push articles"})
			return
		}
	}

	// 构建响应消息
	response := gin.H{
		"pushed_count":  len(pushes),
		"skipped_count": skippedCount,
		"total_count":   len(req.UserIDs),
	}

	if len(pushes) > 0 && skippedCount > 0 {
		response["message"] = fmt.Sprintf("成功推送给 %d 个用户，跳过 %d 个已推送的用户", len(pushes), skippedCount)
	} else if len(pushes) > 0 {
		response["message"] = fmt.Sprintf("成功推送给 %d 个用户", len(pushes))
	} else {
		response["message"] = "所有选中的用户都已经推送过此文章"
	}

	c.JSON(http.StatusOK, response)
}

// GetArticlePushes 获取文章推送记录
func (ac *AdminController) GetArticlePushes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	offset := (page - 1) * limit

	var pushes []models.ArticlePush
	var total int64

	config.DB.Model(&models.ArticlePush{}).Count(&total)

	if err := config.DB.Offset(offset).Limit(limit).
		Preload("Article").
		Preload("User").
		Preload("Admin").
		Order("created_at DESC").
		Find(&pushes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pushes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pushes": pushes,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

// GetUserPushedArticles 获取推送给用户的文章列表
func (ac *AdminController) GetUserPushedArticles(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	offset := (page - 1) * limit

	var pushes []models.ArticlePush
	var total int64

	config.DB.Model(&models.ArticlePush{}).Where("user_id = ?", userID).Count(&total)

	if err := config.DB.Where("user_id = ?", userID).
		Preload("Article").
		Preload("Article.Author").
		Preload("Admin").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&pushes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pushed articles"})
		return
	}

	// 提取文章信息
	articles := make([]models.Article, len(pushes))
	for i, push := range pushes {
		articles[i] = push.Article
	}

	c.JSON(http.StatusOK, gin.H{
		"articles": articles,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

// GetArticlePushedUsers 获取文章已推送的用户列表
func (ac *AdminController) GetArticlePushedUsers(c *gin.Context) {
	articleID := c.Param("article_id")

	var pushes []models.ArticlePush
	if err := config.DB.Where("article_id = ?", articleID).
		Preload("User").
		Find(&pushes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pushed users"})
		return
	}

	// 提取用户ID列表
	userIDs := make([]uint, len(pushes))
	for i, push := range pushes {
		userIDs[i] = push.UserID
	}

	c.JSON(http.StatusOK, gin.H{
		"pushed_user_ids": userIDs,
		"pushes":          pushes,
	})
}
