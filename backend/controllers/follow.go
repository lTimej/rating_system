package controllers

import (
	"net/http"
	"rating_system/config"
	"rating_system/models"

	"github.com/gin-gonic/gin"
)

type FollowController struct{}

func (fc *FollowController) FollowUser(c *gin.Context) {
	followerID, _ := c.Get("user_id")
	followerRole, _ := c.Get("user_role")
	followedID := c.Param("user_id")

	// 专家才能关注学生
	if followerRole.(models.UserRole) != models.RoleExpert {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only experts can follow users"})
		return
	}

	// 检查被关注用户是否存在且是学生
	var followedUser models.User
	if err := config.DB.First(&followedUser, followedID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if followedUser.Role != models.RoleStudent {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can only follow students"})
		return
	}

	// 不能关注自己
	if followerID.(uint) == followedUser.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot follow yourself"})
		return
	}

	// 检查是否已经关注
	var existingFollow models.Follow
	if err := config.DB.Where("follower_id = ? AND followed_id = ?", followerID, followedID).First(&existingFollow).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Already following this user"})
		return
	}

	// 创建关注关系
	follow := models.Follow{
		FollowerID: followerID.(uint),
		FollowedID: followedUser.ID,
	}

	if err := config.DB.Create(&follow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User followed successfully",
		"follow":  follow,
	})
}

func (fc *FollowController) UnfollowUser(c *gin.Context) {
	followerID, _ := c.Get("user_id")
	followedID := c.Param("user_id")

	var follow models.Follow
	if err := config.DB.Where("follower_id = ? AND followed_id = ?", followerID, followedID).First(&follow).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Follow relationship not found"})
		return
	}

	if err := config.DB.Delete(&follow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfollow user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User unfollowed successfully",
	})
}

func (fc *FollowController) GetFollowing(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var follows []models.Follow
	if err := config.DB.Where("follower_id = ?", userID).
		Preload("Followed").Find(&follows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get following list"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"following": follows,
	})
}

func (fc *FollowController) GetFollowers(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var follows []models.Follow
	if err := config.DB.Where("followed_id = ?", userID).
		Preload("Follower").Find(&follows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get followers list"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"followers": follows,
	})
}
