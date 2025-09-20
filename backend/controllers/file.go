package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"rating_system/config"
	"rating_system/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type FileController struct{}

func (fc *FileController) UploadFile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// 创建上传目录
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// 生成唯一文件名
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d_%d%s", userID, time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// 保存文件
	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// 确定文件类型
	fileType := getFileType(header.Header.Get("Content-Type"))

	// 保存文件信息到数据库
	fileModel := models.File{
		UserID:      userID.(uint),
		FileName:    header.Filename,
		FilePath:    filePath,
		FileType:    fileType,
		FileSize:    header.Size,
		MimeType:    header.Header.Get("Content-Type"),
		Title:       c.PostForm("title"),
		Description: c.PostForm("description"),
		IsPublic:    c.PostForm("is_public") != "false",
	}

	if err := config.DB.Create(&fileModel).Error; err != nil {
		// 删除已上传的文件
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"file":    fileModel,
	})
}

func (fc *FileController) GetUserFiles(c *gin.Context) {
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

	var files []models.File
	query := config.DB.Where("user_id = ?", targetUserID)

	// 如果不是查看自己的文件，只显示公开的文件
	if targetUserID != currentUserID.(uint) {
		query = query.Where("is_public = ?", true)
	}

	if err := query.Preload("User").Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"files": files,
	})
}

func (fc *FileController) GetFile(c *gin.Context) {
	fileID := c.Param("id")
	currentUserID, _ := c.Get("user_id")

	var file models.File
	if err := config.DB.Preload("User").First(&file, fileID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// 检查权限：只有文件所有者或公开文件才能访问
	if file.UserID != currentUserID.(uint) && !file.IsPublic {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"file": file,
	})
}

func (fc *FileController) DownloadFile(c *gin.Context) {
	fileID := c.Param("id")
	currentUserID, _ := c.Get("user_id")

	var file models.File
	if err := config.DB.First(&file, fileID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// 检查权限
	if file.UserID != currentUserID.(uint) && !file.IsPublic {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found on disk"})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.FileName))
	c.Header("Content-Type", file.MimeType)
	c.File(file.FilePath)
}

func (fc *FileController) DeleteFile(c *gin.Context) {
	fileID := c.Param("id")
	currentUserID, _ := c.Get("user_id")

	var file models.File
	if err := config.DB.First(&file, fileID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// 检查权限：只有文件所有者才能删除
	if file.UserID != currentUserID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// 删除数据库记录
	if err := config.DB.Delete(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file record"})
		return
	}

	// 删除物理文件
	if err := os.Remove(file.FilePath); err != nil {
		// 记录错误但不返回失败，因为数据库记录已删除
		fmt.Printf("Failed to delete physical file: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File deleted successfully",
	})
}

func getFileType(mimeType string) models.FileType {
	if strings.HasPrefix(mimeType, "image/") {
		return models.FileTypeImage
	} else if strings.HasPrefix(mimeType, "video/") {
		return models.FileTypeVideo
	} else if strings.Contains(mimeType, "document") ||
		strings.Contains(mimeType, "pdf") ||
		strings.Contains(mimeType, "text") ||
		strings.Contains(mimeType, "application/msword") ||
		strings.Contains(mimeType, "application/vnd.openxmlformats") {
		return models.FileTypeDocument
	}
	return models.FileTypeOther
}
