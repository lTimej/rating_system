package controllers

import (
	"net/http"
	"rating_system/config"
	"rating_system/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ArticleController struct{}

// CreateArticle 创建文章
func (ac *ArticleController) CreateArticle(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Title      string `json:"title" binding:"required"`
		Content    string `json:"content" binding:"required"`
		Summary    string `json:"summary"`
		CoverImage string `json:"cover_image"`
		Tags       string `json:"tags"`
		Category   string `json:"category"`
		IsPublic   bool   `json:"is_public"`
		Status     string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认状态
	status := models.ArticleStatusDraft
	if req.Status == "published" {
		status = models.ArticleStatusPublished
	}

	article := models.Article{
		AuthorID:   userID.(uint),
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		CoverImage: req.CoverImage,
		Tags:       req.Tags,
		Category:   req.Category,
		Status:     status,
		IsPublic:   req.IsPublic,
	}

	// 如果是发布状态，设置发布时间
	if status == models.ArticleStatusPublished {
		now := time.Now()
		article.PublishedAt = &now
	}

	if err := config.DB.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article"})
		return
	}

	// 预加载作者信息
	if err := config.DB.Preload("Author").First(&article, article.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Article created successfully",
		"article": article,
	})
}

// GetArticles 获取文章列表
func (ac *ArticleController) GetArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	category := c.Query("category")
	authorID := c.Query("author_id")
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	currentUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userRole, _ := c.Get("user_role")
	role := userRole.(models.UserRole)

	var articles []models.Article
	var total int64

	// 根据用户角色决定可见的文章
	if role == models.RoleStudent {
		// 学生只能看到推送给他们的文章
		query := config.DB.Table("articles").
			Joins("INNER JOIN article_pushes ON articles.id = article_pushes.article_id").
			Where("article_pushes.user_id = ? AND articles.status = ?", currentUserID, models.ArticleStatusPublished)

		// 按分类筛选
		if category != "" {
			query = query.Where("articles.category = ?", category)
		}

		// 按作者筛选
		if authorID != "" {
			query = query.Where("articles.author_id = ?", authorID)
		}

		query.Count(&total)

		if err := query.Preload("Author").
			Order("articles.created_at DESC").
			Limit(pageSize).
			Offset(offset).
			Find(&articles).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get articles"})
			return
		}
	} else {
		// 专家和管理员可以看到所有公开且已发布的文章（或自己的文章）
		query := config.DB.Model(&models.Article{})

		// 只显示公开且已发布的文章（除非是查看自己的文章）
		if authorID != strconv.Itoa(int(currentUserID.(uint))) {
			query = query.Where("is_public = ? AND status = ?", true, models.ArticleStatusPublished)
		}

		// 按分类筛选
		if category != "" {
			query = query.Where("category = ?", category)
		}

		// 按作者筛选
		if authorID != "" {
			query = query.Where("author_id = ?", authorID)
		}

		// 按状态筛选（仅作者本人可见）
		if status != "" && authorID == strconv.Itoa(int(currentUserID.(uint))) {
			query = query.Where("status = ?", status)
		}

		query.Count(&total)

		if err := query.Preload("Author").
			Order("created_at DESC").
			Limit(pageSize).
			Offset(offset).
			Find(&articles).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get articles"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"articles": articles,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetArticle 获取单篇文章
func (ac *ArticleController) GetArticle(c *gin.Context) {
	articleID := c.Param("id")
	currentUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userRole, _ := c.Get("user_role")
	role := userRole.(models.UserRole)

	var article models.Article
	if err := config.DB.Preload("Author").First(&article, articleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	// 根据用户角色检查访问权限
	if role == models.RoleStudent {
		// 学生只能访问推送给他们的文章
		var push models.ArticlePush
		if err := config.DB.Where("article_id = ? AND user_id = ?", articleID, currentUserID).First(&push).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		
		// 文章必须是已发布状态
		if article.Status != models.ArticleStatusPublished {
			c.JSON(http.StatusForbidden, gin.H{"error": "Article not available"})
			return
		}
	} else {
		// 专家和管理员的权限检查：只有公开且已发布的文章或作者本人可以查看
		if (!article.IsPublic || article.Status != models.ArticleStatusPublished) &&
			article.AuthorID != currentUserID.(uint) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	// 增加浏览次数（不是作者本人访问时）
	if article.AuthorID != currentUserID.(uint) {
		config.DB.Model(&article).Update("view_count", gorm.Expr("view_count + 1"))
		article.ViewCount++
	}

	c.JSON(http.StatusOK, gin.H{
		"article": article,
	})
}

// UpdateArticle 更新文章
func (ac *ArticleController) UpdateArticle(c *gin.Context) {
	articleID := c.Param("id")
	userID, _ := c.Get("user_id")

	var article models.Article
	if err := config.DB.First(&article, articleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	// 检查权限：只有作者可以编辑
	if article.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var req struct {
		Title      string `json:"title"`
		Content    string `json:"content"`
		Summary    string `json:"summary"`
		CoverImage string `json:"cover_image"`
		Tags       string `json:"tags"`
		Category   string `json:"category"`
		IsPublic   *bool  `json:"is_public"`
		Status     string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Summary != "" {
		updates["summary"] = req.Summary
	}
	if req.CoverImage != "" {
		updates["cover_image"] = req.CoverImage
	}
	if req.Tags != "" {
		updates["tags"] = req.Tags
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}
	if req.Status != "" {
		updates["status"] = req.Status
		// 如果从草稿改为发布，设置发布时间
		if req.Status == "published" && article.Status == models.ArticleStatusDraft {
			now := time.Now()
			updates["published_at"] = &now
		}
	}

	if err := config.DB.Model(&article).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article"})
		return
	}

	// 重新加载文章数据
	if err := config.DB.Preload("Author").First(&article, article.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Article updated successfully",
		"article": article,
	})
}

// DeleteArticle 删除文章
func (ac *ArticleController) DeleteArticle(c *gin.Context) {
	articleID := c.Param("id")
	userID, _ := c.Get("user_id")

	var article models.Article
	if err := config.DB.First(&article, articleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	// 检查权限：只有作者可以删除
	if article.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// 软删除文章
	if err := config.DB.Delete(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Article deleted successfully",
	})
}

// LikeArticle 点赞/取消点赞文章
func (ac *ArticleController) LikeArticle(c *gin.Context) {
	articleID := c.Param("id")
	userID, _ := c.Get("user_id")

	// 检查文章是否存在
	var article models.Article
	if err := config.DB.First(&article, articleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	// 检查是否已经点赞
	var existingLike models.ArticleLike
	err := config.DB.Where("article_id = ? AND user_id = ?", articleID, userID).First(&existingLike).Error

	if err == gorm.ErrRecordNotFound {
		// 没有点赞，创建点赞记录
		like := models.ArticleLike{
			ArticleID: article.ID,
			UserID:    userID.(uint),
		}
		if err := config.DB.Create(&like).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like article"})
			return
		}
		// 增加点赞数
		config.DB.Model(&article).Update("like_count", gorm.Expr("like_count + 1"))
		c.JSON(http.StatusOK, gin.H{"message": "Article liked", "liked": true})
	} else if err == nil {
		// 已经点赞，取消点赞
		if err := config.DB.Delete(&existingLike).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlike article"})
			return
		}
		// 减少点赞数
		config.DB.Model(&article).Update("like_count", gorm.Expr("like_count - 1"))
		c.JSON(http.StatusOK, gin.H{"message": "Article unliked", "liked": false})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
	}
}
