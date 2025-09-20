package models

import (
	"time"

	"gorm.io/gorm"
)

type FileType string

const (
	FileTypeDocument FileType = "document"
	FileTypeVideo    FileType = "video"
	FileTypeImage    FileType = "image"
	FileTypeOther    FileType = "other"
)

type File struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	FileName    string         `json:"file_name" gorm:"not null"`
	FilePath    string         `json:"file_path" gorm:"not null"`
	FileType    FileType       `json:"file_type" gorm:"not null"`
	FileSize    int64          `json:"file_size"`
	MimeType    string         `json:"mime_type"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	IsPublic    bool           `json:"is_public" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User    User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Ratings []Rating `json:"ratings,omitempty" gorm:"foreignKey:FileID"`
}
