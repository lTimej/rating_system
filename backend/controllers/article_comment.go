package controllers

import (
	"net/http"
	"rating_system/config"
	"rating_system/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ArticleCommentController struct{}

// CreateArticleComment 创建文章评论
func (acc *ArticleCommentController) CreateArticleComment(c *gin.Context) {
	userID, _ := c.Get("user_id")
	articleID := c.Param("id")

	var req struct {
		Content  string `json:"content" binding:"required"`
		ParentID *uint  `json:"parent_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证文章是否存在且为公开已发布文章
	var article models.Article
	if err := config.DB.First(&article, articleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	if !article.IsPublic || article.Status != models.ArticleStatusPublished {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot comment on this article"})
		return
	}

	// 如果是回复评论，验证父评论是否存在
	if req.ParentID != nil {
		var parentComment models.ArticleComment
		if err := config.DB.Where("id = ? AND article_id = ?", *req.ParentID, articleID).First(&parentComment).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parent comment not found"})
			return
		}
	}

	// 创建评论
	comment := models.ArticleComment{
		ArticleID: article.ID,
		UserID:    userID.(uint),
		Content:   req.Content,
		ParentID:  req.ParentID,
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// 预加载用户信息
	if err := config.DB.Preload("User").First(&comment, comment.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment created successfully",
		"comment": comment,
	})
}

// GetArticleComments 获取文章评论列表
func (acc *ArticleCommentController) GetArticleComments(c *gin.Context) {
	articleID := c.Param("id")

	// 验证文章是否存在且为公开已发布文章
	var article models.Article
	if err := config.DB.First(&article, articleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	if !article.IsPublic || article.Status != models.ArticleStatusPublished {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot view comments of this article"})
		return
	}

	var comments []models.ArticleComment
	// 只获取顶级评论（没有父评论的评论）
	if err := config.DB.Where("article_id = ? AND parent_id IS NULL", articleID).
		Preload("User").
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Preload("User").Order("created_at ASC")
		}).
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
	})
}

// DeleteArticleComment 删除文章评论
func (acc *ArticleCommentController) DeleteArticleComment(c *gin.Context) {
	commentID := c.Param("comment_id")
	userID, _ := c.Get("user_id")

	var comment models.ArticleComment
	if err := config.DB.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// 检查权限：只有评论作者才能删除
	if comment.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// 软删除评论
	if err := config.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment deleted successfully",
	})
}
