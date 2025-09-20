package controllers

import (
	"net/http"
	"rating_system/config"
	"rating_system/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RatingController struct{}

type CreateRatingRequest struct {
	RatedID     uint   `json:"rated_id" binding:"required"`
	FileID      *uint  `json:"file_id"`
	Content     string `json:"content" binding:"required"`
	Score       int    `json:"score" binding:"min=1,max=5"`
	IsAnonymous bool   `json:"is_anonymous"`
}

func (rc *RatingController) CreateRating(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")

	var req CreateRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查被评价用户是否存在
	var ratedUser models.User
	if err := config.DB.First(&ratedUser, req.RatedID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rated user not found"})
		return
	}

	// 不能评价自己
	if userID.(uint) == req.RatedID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot rate yourself"})
		return
	}

	// 检查权限：学生只能评价学生，专家可以评价所有人
	if userRole.(models.UserRole) == models.RoleStudent && ratedUser.Role != models.RoleStudent {
		c.JSON(http.StatusForbidden, gin.H{"error": "Students can only rate other students"})
		return
	}

	// 如果指定了文件ID，检查文件是否存在且属于被评价用户
	if req.FileID != nil {
		var file models.File
		if err := config.DB.First(&file, *req.FileID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
		if file.UserID != req.RatedID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File does not belong to the rated user"})
			return
		}
	}

	// 创建评价
	rating := models.Rating{
		RaterID:     userID.(uint),
		RatedID:     req.RatedID,
		FileID:      req.FileID,
		Content:     req.Content,
		Score:       req.Score,
		IsAnonymous: req.IsAnonymous,
	}

	if err := config.DB.Create(&rating).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create rating"})
		return
	}

	// 加载关联数据
	config.DB.Preload("Rater").Preload("Rated").Preload("File").First(&rating, rating.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Rating created successfully",
		"rating":  rating,
	})
}

func (rc *RatingController) GetUserRatings(c *gin.Context) {
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

	// 获取用户收到的评价
	var receivedRatings []models.Rating
	query := config.DB.Where("rated_id = ?", targetUserID)

	// 根据需求：评价者身份公开，被评价者身份匿名
	// 只有被评价者本人可以看到针对自己的评价
	if targetUserID != currentUserID.(uint) {
		// 其他人看不到这个用户收到的评价
		receivedRatings = []models.Rating{}
	} else {
		query.Preload("Rater").Preload("File").Preload("Feedbacks").Find(&receivedRatings)
	}

	// 获取用户给出的评价（只有本人可见）
	var givenRatings []models.Rating
	if targetUserID == currentUserID.(uint) {
		config.DB.Where("rater_id = ?", targetUserID).
			Preload("Rated").Preload("File").Preload("Feedbacks").
			Find(&givenRatings)
	}

	c.JSON(http.StatusOK, gin.H{
		"received_ratings": receivedRatings,
		"given_ratings":    givenRatings,
	})
}

func (rc *RatingController) GetRating(c *gin.Context) {
	ratingID := c.Param("id")
	currentUserID, _ := c.Get("user_id")

	var rating models.Rating
	if err := config.DB.Preload("Rater").Preload("Rated").Preload("File").Preload("Feedbacks").First(&rating, ratingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rating not found"})
		return
	}

	// 检查权限：只有评价者和被评价者可以查看
	if rating.RaterID != currentUserID.(uint) && rating.RatedID != currentUserID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rating": rating,
	})
}

func (rc *RatingController) CreateFeedback(c *gin.Context) {
	ratingID := c.Param("rating_id")
	userID, _ := c.Get("user_id")

	var req struct {
		IsHelpful bool `json:"is_helpful" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查评价是否存在
	var rating models.Rating
	if err := config.DB.First(&rating, ratingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rating not found"})
		return
	}

	// 检查是否已经给过反馈
	var existingFeedback models.Feedback
	if err := config.DB.Where("rating_id = ? AND user_id = ?", ratingID, userID).First(&existingFeedback).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Feedback already exists"})
		return
	}

	// 创建反馈
	feedback := models.Feedback{
		RatingID:  rating.ID,
		UserID:    userID.(uint),
		IsHelpful: req.IsHelpful,
	}

	if err := config.DB.Create(&feedback).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create feedback"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Feedback created successfully",
		"feedback": feedback,
	})
}

func (rc *RatingController) GetPublicRatings(c *gin.Context) {
	// 获取公开的评价（用于管理员推送给学生查看）
	var ratings []models.Rating

	// 只显示针对学生的评价，且文件是公开的
	query := config.DB.Joins("JOIN users ON users.id = ratings.rated_id").
		Where("users.role = ?", models.RoleStudent)

	// 如果有文件关联，确保文件是公开的
	query = query.Joins("LEFT JOIN files ON files.id = ratings.file_id").
		Where("ratings.file_id IS NULL OR files.is_public = ?", true)

	if err := query.Preload("Rater").Preload("File").Preload("Feedbacks").Find(&ratings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get ratings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ratings": ratings,
	})
}
