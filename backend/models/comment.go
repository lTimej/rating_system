package models

import (
	"time"

	"gorm.io/gorm"
)

// FileComment 文件评论模型
type FileComment struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	FileID    uint           `json:"file_id" gorm:"not null"`    // 文件ID
	UserID    uint           `json:"user_id" gorm:"not null"`    // 评论者ID
	Content   string         `json:"content" gorm:"not null"`    // 评论内容
	ParentID  *uint          `json:"parent_id"`                  // 父评论ID，用于回复
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	File     File          `json:"file,omitempty" gorm:"foreignKey:FileID"`
	User     User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Parent   *FileComment  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Replies  []FileComment `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
}
