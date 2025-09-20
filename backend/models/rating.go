package models

import (
	"time"

	"gorm.io/gorm"
)

type Rating struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	RaterID     uint           `json:"rater_id" gorm:"not null"`          // 评价者ID
	RatedID     uint           `json:"rated_id" gorm:"not null"`          // 被评价者ID
	FileID      *uint          `json:"file_id"`                           // 可选：针对特定文件的评价
	Content     string         `json:"content" gorm:"not null"`           // 评价内容
	Score       int            `json:"score" gorm:"default:0"`            // 评分 (1-5)
	IsAnonymous bool           `json:"is_anonymous" gorm:"default:false"` // 是否匿名评价
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Rater     User       `json:"rater,omitempty" gorm:"foreignKey:RaterID"`
	Rated     User       `json:"rated,omitempty" gorm:"foreignKey:RatedID"`
	File      *File      `json:"file,omitempty" gorm:"foreignKey:FileID"`
	Feedbacks []Feedback `json:"feedbacks,omitempty" gorm:"foreignKey:RatingID"`
}

type Feedback struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	RatingID  uint      `json:"rating_id" gorm:"not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	IsHelpful bool      `json:"is_helpful" gorm:"not null"` // true: 有帮助, false: 没有帮助
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	Rating Rating `json:"rating,omitempty" gorm:"foreignKey:RatingID"`
	User   User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
