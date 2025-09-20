package models

import (
	"time"

	"gorm.io/gorm"
)

type ArticleStatus string

const (
	ArticleStatusDraft     ArticleStatus = "draft"     // 草稿
	ArticleStatusPublished ArticleStatus = "published" // 已发布
	ArticleStatusArchived  ArticleStatus = "archived"  // 已归档
)

// Article 文章模型
type Article struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	AuthorID    uint           `json:"author_id" gorm:"not null"`    // 作者ID
	Title       string         `json:"title" gorm:"not null"`        // 文章标题
	Content     string         `json:"content" gorm:"type:text"`     // 文章内容（Markdown格式）
	Summary     string         `json:"summary"`                      // 文章摘要
	CoverImage  string         `json:"cover_image"`                  // 封面图片URL
	Tags        string         `json:"tags"`                         // 标签（JSON字符串）
	Category    string         `json:"category"`                     // 分类
	Status      ArticleStatus  `json:"status" gorm:"default:draft"` // 文章状态
	ViewCount   int            `json:"view_count" gorm:"default:0"`  // 浏览次数
	LikeCount   int            `json:"like_count" gorm:"default:0"`  // 点赞次数
	IsPublic    bool           `json:"is_public" gorm:"default:true"` // 是否公开
	PublishedAt *time.Time     `json:"published_at"`                 // 发布时间
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Author   User             `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	Comments []ArticleComment `json:"comments,omitempty" gorm:"foreignKey:ArticleID"`
	Likes    []ArticleLike    `json:"likes,omitempty" gorm:"foreignKey:ArticleID"`
}

// ArticleComment 文章评论模型
type ArticleComment struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	ArticleID        uint           `json:"article_id" gorm:"not null"`        // 文章ID
	UserID           uint           `json:"user_id" gorm:"not null"`           // 评论者ID
	Content          string         `json:"content" gorm:"not null"`           // 评论内容
	ParentID         *uint          `json:"parent_id"`                         // 父评论ID，用于回复
	IsAnonymousToOthers bool        `json:"is_anonymous_to_others" gorm:"default:true"` // 对其他人是否匿名（被回复者除外）
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Article Article          `json:"article,omitempty" gorm:"foreignKey:ArticleID"`
	User    User             `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Parent  *ArticleComment  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Replies []ArticleComment `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
}

// ArticleLike 文章点赞模型
type ArticleLike struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ArticleID uint      `json:"article_id" gorm:"not null"` // 文章ID
	UserID    uint      `json:"user_id" gorm:"not null"`    // 点赞用户ID
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	Article Article `json:"article,omitempty" gorm:"foreignKey:ArticleID"`
	User    User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// ArticleCommentRating 文章评论评价模型
type ArticleCommentRating struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CommentID uint      `json:"comment_id" gorm:"not null"` // 评论ID
	RaterID   uint      `json:"rater_id" gorm:"not null"`   // 评价者ID
	RatedID   uint      `json:"rated_id" gorm:"not null"`   // 被评价者ID（评论作者）
	Score     int       `json:"score" gorm:"not null"`      // 评分 (1-5)
	Content   string    `json:"content"`                    // 评价内容（可选）
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Comment ArticleComment `json:"comment,omitempty" gorm:"foreignKey:CommentID"`
	Rater   User           `json:"rater,omitempty" gorm:"foreignKey:RaterID"`
	Rated   User           `json:"rated,omitempty" gorm:"foreignKey:RatedID"`
}
