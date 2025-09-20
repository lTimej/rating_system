package controllers

import (
	"net/http"
	"rating_system/config"
	"rating_system/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentController struct{}

// CreateFileComment 创建文件评论
func (cc *CommentController) CreateFileComment(c *gin.Context) {
	userID, _ := c.Get("user_id")
	fileID := c.Param("id")

	var req struct {
		Content  string `json:"content" binding:"required"`
		ParentID *uint  `json:"parent_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证文件是否存在且为公开文件
	var file models.File
	if err := config.DB.First(&file, fileID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	if !file.IsPublic {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot comment on private file"})
		return
	}

	// 如果是回复评论，验证父评论是否存在
	if req.ParentID != nil {
		var parentComment models.FileComment
		if err := config.DB.Where("id = ? AND file_id = ?", *req.ParentID, fileID).First(&parentComment).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parent comment not found"})
			return
		}
	}

	// 创建评论
	comment := models.FileComment{
		FileID:   file.ID,
		UserID:   userID.(uint),
		Content:  req.Content,
		ParentID: req.ParentID,
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

// GetFileComments 获取文件评论列表
func (cc *CommentController) GetFileComments(c *gin.Context) {
	fileID := c.Param("id")

	// 验证文件是否存在且为公开文件
	var file models.File
	if err := config.DB.First(&file, fileID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	if !file.IsPublic {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot view comments of private file"})
		return
	}

	var comments []models.FileComment
	// 只获取顶级评论（没有父评论的评论）
	if err := config.DB.Where("file_id = ? AND parent_id IS NULL", fileID).
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

// DeleteFileComment 删除文件评论
func (cc *CommentController) DeleteFileComment(c *gin.Context) {
	commentID := c.Param("comment_id")
	userID, _ := c.Get("user_id")

	var comment models.FileComment
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
