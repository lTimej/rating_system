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
		c.JSON(http.StatusNotFound, gin.H{"error": "文章未找到"})
		return
	}

	if !article.IsPublic || article.Status != models.ArticleStatusPublished {
		c.JSON(http.StatusForbidden, gin.H{"error": "无法评论此文章"})
		return
	}

	// 如果是回复评论，验证父评论是否存在
	if req.ParentID != nil {
		var parentComment models.ArticleComment
		if err := config.DB.Where("id = ? AND article_id = ?", *req.ParentID, articleID).First(&parentComment).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "父评论未找到"})
			return
		}
	}

	// 创建评论
	comment := models.ArticleComment{
		ArticleID:           article.ID,
		UserID:              userID.(uint),
		Content:             req.Content,
		ParentID:            req.ParentID,
		IsAnonymousToOthers: true, // 默认对其他人匿名，只有被回复者能看到姓名
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建评论失败"})
		return
	}

	// 预加载用户信息
	if err := config.DB.Preload("User").First(&comment, comment.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "加载评论失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "评论创建成功",
		"comment": comment,
	})
}

// GetArticleComments 获取文章评论列表
func (acc *ArticleCommentController) GetArticleComments(c *gin.Context) {
	articleID := c.Param("id")
	userID, _ := c.Get("user_id")

	// 验证文章是否存在且为公开已发布文章
	var article models.Article
	if err := config.DB.First(&article, articleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章未找到"})
		return
	}

	if !article.IsPublic || article.Status != models.ArticleStatusPublished {
		c.JSON(http.StatusForbidden, gin.H{"error": "无法查看此文章的评论"})
		return
	}

	var comments []models.ArticleComment

	// 所有用户都可以看到所有评论，但姓名可见性由processCommentVisibility处理
	if err := config.DB.Where("article_id = ? AND parent_id IS NULL", articleID).
		Preload("User").
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Preload("User").Order("created_at ASC")
		}).
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论失败"})
		return
	}

	// 处理姓名可见性：回复者对被回复者匿名，但被回复者能看到回复者姓名
	processedComments := acc.processCommentVisibility(comments, userID.(uint))

	c.JSON(http.StatusOK, gin.H{
		"comments": processedComments,
	})
}

// DeleteArticleComment 删除文章评论
func (acc *ArticleCommentController) DeleteArticleComment(c *gin.Context) {
	commentID := c.Param("comment_id")
	userID, _ := c.Get("user_id")

	var comment models.ArticleComment
	if err := config.DB.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论未找到"})
		return
	}

	// 检查权限：只有评论作者才能删除
	if comment.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "访问被拒绝"})
		return
	}

	// 软删除评论
	if err := config.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除评论失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "评论删除成功",
	})
}

// processCommentVisibility 处理评论的姓名可见性
// 规则：
// 1. 回复者对其他人匿名，但被回复者能看到回复者姓名
// 2. 父评论作者对回复者也匿名显示
func (acc *ArticleCommentController) processCommentVisibility(comments []models.ArticleComment, currentUserID uint) []models.ArticleComment {
	processedComments := make([]models.ArticleComment, len(comments))
	
	for i, comment := range comments {
		processedComments[i] = comment
		
		// 处理父评论作者的姓名可见性
		// 如果当前用户不是父评论作者本人，则隐藏父评论作者姓名
		if comment.UserID != currentUserID {
			processedComments[i].User = models.User{
				ID:       comment.User.ID,
				Username: "匿名用户",
				Name:     "匿名用户",
				Role:     comment.User.Role,
			}
		}
		
		// 处理回复的可见性
		if len(comment.Replies) > 0 {
			processedReplies := make([]models.ArticleComment, len(comment.Replies))
			
			for j, reply := range comment.Replies {
				processedReplies[j] = reply
				
				// 回复者姓名可见性规则：
				// 1. 如果当前用户是被回复的人（父评论作者），可以看到回复者真实姓名
				// 2. 如果当前用户是回复者本人，可以看到自己的真实姓名
				// 3. 其他情况下，回复者姓名显示为匿名
				if reply.IsAnonymousToOthers && 
				   comment.UserID != currentUserID && 
				   reply.UserID != currentUserID {
					// 对其他人匿名显示
					processedReplies[j].User = models.User{
						ID:       reply.User.ID,
						Username: "匿名用户",
						Name:     "匿名用户",
						Role:     reply.User.Role,
					}
				}
				// 否则显示真实姓名（被回复者或回复者本人可见）
			}
			
			processedComments[i].Replies = processedReplies
		}
	}
	
	return processedComments
}

// GetLatestComments 获取最新的文章评论（用于首页显示）
func (acc *ArticleCommentController) GetLatestComments(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	userRole, _ := c.Get("user_role")
	role := userRole.(models.UserRole)

	var comments []models.ArticleComment
	
	if role == models.RoleStudent {
		// 学生只能看到推送文章的评论
		if err := config.DB.
			Joins("JOIN articles ON articles.id = article_comments.article_id").
			Joins("JOIN article_pushes ON articles.id = article_pushes.article_id").
			Where("article_pushes.user_id = ? AND articles.status = ?", userID, models.ArticleStatusPublished).
			Preload("User").
			Preload("Article").
			Order("article_comments.created_at DESC").
			Limit(10).
			Find(&comments).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取最新评论失败"})
			return
		}
	} else {
		// 专家和管理员可以看到所有公开文章的评论
		if err := config.DB.
			Joins("JOIN articles ON articles.id = article_comments.article_id").
			Where("articles.is_public = ? AND articles.status = ?", true, models.ArticleStatusPublished).
			Preload("User").
			Preload("Article").
			Order("article_comments.created_at DESC").
			Limit(10).
			Find(&comments).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取最新评论失败"})
			return
		}
	}

	// 处理姓名可见性
	processedComments := acc.processLatestCommentsVisibility(comments, userID.(uint))

	c.JSON(http.StatusOK, gin.H{
		"comments": processedComments,
	})
}

// processLatestCommentsVisibility 处理最新评论的姓名可见性
// 规则：只有评论者本人能看到自己的真实姓名，其他人都看到匿名用户
func (acc *ArticleCommentController) processLatestCommentsVisibility(comments []models.ArticleComment, currentUserID uint) []models.ArticleComment {
	processedComments := make([]models.ArticleComment, len(comments))
	
	for i, comment := range comments {
		processedComments[i] = comment
		
		// 简化的匿名逻辑：只有评论者本人能看到自己的真实姓名
		if comment.UserID != currentUserID {
			processedComments[i].User = models.User{
				ID:       comment.User.ID,
				Username: "匿名用户",
				Name:     "匿名用户",
				Role:     comment.User.Role,
			}
		}
	}
	
	return processedComments
}
