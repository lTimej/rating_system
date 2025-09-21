package controllers

import (
	"net/http"
	"rating_system/config"
	"rating_system/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentRatingController struct{}

type CreateCommentRatingRequest struct {
	CommentID uint   `json:"comment_id" binding:"required"`
	Score     int    `json:"score" binding:"required,min=1,max=5"`
	Content   string `json:"content"`
}

// CreateCommentRating 创建评论评价
func (crc *CommentRatingController) CreateCommentRating(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req CreateCommentRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查评论是否存在
	var comment models.ArticleComment
	if err := config.DB.Preload("User").First(&comment, req.CommentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论未找到"})
		return
	}

	// 不能评价自己的评论
	if comment.UserID == userID.(uint) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能评价自己的评论"})
		return
	}

	// 检查是否已经评价过这个评论
	var existingRating models.ArticleCommentRating
	if err := config.DB.Where("comment_id = ? AND rater_id = ?", req.CommentID, userID).First(&existingRating).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "您已经评价过这条评论"})
		return
	}

	// 创建评价
	rating := models.ArticleCommentRating{
		CommentID: req.CommentID,
		RaterID:   userID.(uint),
		RatedID:   comment.UserID,
		Score:     req.Score,
		Content:   req.Content,
	}

	if err := config.DB.Create(&rating).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建评价失败"})
		return
	}

	// 预加载关联数据
	if err := config.DB.Preload("Comment").Preload("Rater").Preload("Rated").First(&rating, rating.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "加载评价失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "评论评价创建成功",
		"rating":  rating,
	})
}

// GetCommentRatings 获取评论的评价列表
func (crc *CommentRatingController) GetCommentRatings(c *gin.Context) {
	commentID := c.Param("comment_id")

	var ratings []models.ArticleCommentRating
	if err := config.DB.Where("comment_id = ?", commentID).
		Preload("Rater").
		Preload("Rated").
		Order("created_at DESC").
		Find(&ratings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评价失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ratings": ratings,
	})
}

// GetUserCommentRatings 获取用户的评论评价统计
func (crc *CommentRatingController) GetUserCommentRatings(c *gin.Context) {
	userIDParam := c.Param("user_id")
	currentUserID, _ := c.Get("user_id")

	var targetUserID uint
	if userIDParam == "me" {
		targetUserID = currentUserID.(uint)
	} else {
		id, err := strconv.ParseUint(userIDParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		targetUserID = uint(id)
	}

	// 获取收到的评价
	var receivedRatings []models.ArticleCommentRating
	if err := config.DB.Where("rated_id = ?", targetUserID).
		Preload("Rater").
		Preload("Comment").
		Order("created_at DESC").
		Find(&receivedRatings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取收到的评价失败"})
		return
	}

	// 获取给出的评价
	var givenRatings []models.ArticleCommentRating
	if err := config.DB.Where("rater_id = ?", targetUserID).
		Preload("Rated").
		Preload("Comment").
		Order("created_at DESC").
		Find(&givenRatings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取给出的评价失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"received_ratings": receivedRatings,
		"given_ratings":    givenRatings,
	})
}

// DeleteCommentRating 删除评论评价
func (crc *CommentRatingController) DeleteCommentRating(c *gin.Context) {
	ratingID := c.Param("rating_id")
	userID, _ := c.Get("user_id")

	var rating models.ArticleCommentRating
	if err := config.DB.First(&rating, ratingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评价未找到"})
		return
	}

	// 检查权限：只有评价者才能删除
	if rating.RaterID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "访问被拒绝"})
		return
	}

	if err := config.DB.Delete(&rating).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除评价失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "评价删除成功",
	})
}
