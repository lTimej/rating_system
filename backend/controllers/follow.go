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
	followedID := c.Param("user_id")

	// 检查被关注用户是否存在
	var followedUser models.User
	if err := config.DB.First(&followedUser, followedID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户未找到"})
		return
	}

	// 不能关注自己
	if followerID.(uint) == followedUser.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能关注自己"})
		return
	}

	// 检查是否已经关注
	var existingFollow models.Follow
	if err := config.DB.Where("follower_id = ? AND followed_id = ?", followerID, followedID).First(&existingFollow).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "已经关注了该用户"})
		return
	}

	// 创建关注关系
	follow := models.Follow{
		FollowerID: followerID.(uint),
		FollowedID: followedUser.ID,
	}

	if err := config.DB.Create(&follow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "关注用户失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "关注成功",
		"follow":  follow,
	})
}

func (fc *FollowController) UnfollowUser(c *gin.Context) {
	followerID, _ := c.Get("user_id")
	followedID := c.Param("user_id")

	var follow models.Follow
	if err := config.DB.Where("follower_id = ? AND followed_id = ?", followerID, followedID).First(&follow).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "关注关系未找到"})
		return
	}

	if err := config.DB.Delete(&follow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取消关注失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "取消关注成功",
	})
}

func (fc *FollowController) GetFollowing(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var follows []models.Follow
	if err := config.DB.Where("follower_id = ?", userID).
		Preload("Followed").Find(&follows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取关注列表失败"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取粉丝列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"followers": follows,
	})
}
