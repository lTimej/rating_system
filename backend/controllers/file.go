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

	// 添加调试日志
	fmt.Printf("Upload request from user %v\n", userID)
	fmt.Printf("Content-Type: %s\n", c.GetHeader("Content-Type"))
	fmt.Printf("Content-Length: %s\n", c.GetHeader("Content-Length"))
	fmt.Printf("User-Agent: %s\n", c.GetHeader("User-Agent"))

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Printf("FormFile error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要上传的文件"})
		return
	}
	defer file.Close()

	fmt.Printf("File info: name=%s, size=%d, type=%s\n", header.Filename, header.Size, header.Header.Get("Content-Type"))

	// 检查文件大小 (1000MB限制)
	const maxFileSize = 1000 << 20 // 1000 MB
	if header.Size > maxFileSize {
		fmt.Printf("File too large: %d bytes\n", header.Size)
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小不能超过1GB"})
		return
	}

	// 检查文件名长度
	if len(header.Filename) > 255 {
		fmt.Printf("Filename too long: %d characters\n", len(header.Filename))
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件名过长"})
		return
	}

	// 创建上传目录
	uploadDir := "uploads"
	fmt.Printf("Creating upload directory: %s\n", uploadDir)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		fmt.Printf("Failed to create upload directory: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建上传目录失败"})
		return
	}

	// 检查目录权限
	if info, err := os.Stat(uploadDir); err != nil {
		fmt.Printf("Failed to stat upload directory: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传目录不可访问"})
		return
	} else {
		fmt.Printf("Upload directory permissions: %v\n", info.Mode())
	}

	// 生成唯一文件名
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d_%d%s", userID, time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)
	fmt.Printf("Generated file path: %s\n", filePath)

	// 保存文件
	fmt.Printf("Creating file: %s\n", filePath)
	dst, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文件失败"})
		return
	}
	defer dst.Close()

	fmt.Printf("Copying file data...\n")
	if _, err := io.Copy(dst, file); err != nil {
		fmt.Printf("Failed to copy file data: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}
	fmt.Printf("File saved successfully\n")

	// 确定文件类型
	fileType := getFileType(header.Header.Get("Content-Type"))
	fmt.Printf("File type determined: %s\n", fileType)

	// 获取表单数据
	title := c.PostForm("title")
	description := c.PostForm("description")
	isPublicStr := c.PostForm("is_public")

	// 处理布尔值转换
	var isPublic bool
	if isPublicStr == "true" {
		isPublic = true
	} else if isPublicStr == "false" {
		isPublic = false
	} else {
		// 默认为true
		isPublic = true
	}

	fmt.Printf("Form data - title: %s, description: %s, is_public: %s (%t)\n", title, description, isPublicStr, isPublic)

	// 保存文件信息到数据库
	fileModel := models.File{
		UserID:      userID.(uint),
		FileName:    header.Filename,
		FilePath:    filePath,
		FileType:    fileType,
		FileSize:    header.Size,
		MimeType:    header.Header.Get("Content-Type"),
		Title:       title,
		Description: description,
		IsPublic:    isPublic,
	}

	fmt.Printf("Saving file model to database...\n")
	if err := config.DB.Create(&fileModel).Error; err != nil {
		fmt.Printf("Database error: %v\n", err)
		// 删除已上传的文件
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件信息失败"})
		return
	}

	fmt.Printf("File upload completed successfully, ID: %d\n", fileModel.ID)
	c.JSON(http.StatusOK, gin.H{
		"message": "文件上传成功",
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文件列表失败"})
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
		c.JSON(http.StatusNotFound, gin.H{"error": "文件未找到"})
		return
	}

	// 检查权限：只有文件所有者或公开文件才能访问
	if file.UserID != currentUserID.(uint) && !file.IsPublic {
		c.JSON(http.StatusForbidden, gin.H{"error": "访问被拒绝"})
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
		c.JSON(http.StatusNotFound, gin.H{"error": "文件未找到"})
		return
	}

	// 检查权限
	if file.UserID != currentUserID.(uint) && !file.IsPublic {
		c.JSON(http.StatusForbidden, gin.H{"error": "访问被拒绝"})
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "磁盘上未找到文件"})
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
		c.JSON(http.StatusNotFound, gin.H{"error": "文件未找到"})
		return
	}

	// 检查权限：只有文件所有者才能删除
	if file.UserID != currentUserID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "访问被拒绝"})
		return
	}

	// 删除数据库记录
	if err := config.DB.Delete(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文件记录失败"})
		return
	}

	// 删除物理文件
	if err := os.Remove(file.FilePath); err != nil {
		// 记录错误但不返回失败，因为数据库记录已删除
		fmt.Printf("Failed to delete physical file: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "文件删除成功",
	})
}

func (fc *FileController) GetPublicFiles(c *gin.Context) {
	var files []models.File

	// 获取所有公开的文件，按创建时间倒序排列
	if err := config.DB.Where("is_public = ?", true).
		Preload("User").
		Order("created_at DESC").
		Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取公开文件失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"files": files,
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
