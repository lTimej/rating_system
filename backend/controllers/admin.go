package controllers

import (
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
	}

	config.DB.Model(&models.User{}).Count(&stats.TotalUsers)
	config.DB.Model(&models.User{}).Where("role = ?", models.RoleStudent).Count(&stats.TotalStudents)
	config.DB.Model(&models.User{}).Where("role = ?", models.RoleExpert).Count(&stats.TotalExperts)
	config.DB.Model(&models.Rating{}).Count(&stats.TotalRatings)
	config.DB.Model(&models.File{}).Count(&stats.TotalFiles)

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}
