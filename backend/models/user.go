package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleStudent UserRole = "student"
	RoleExpert  UserRole = "expert"
	RoleAdmin   UserRole = "admin"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Role      UserRole       `json:"role" gorm:"not null"`
	Name      string         `json:"name"`
	Avatar    string         `json:"avatar"`
	Bio       string         `json:"bio"`
	Links     string         `json:"links"` // JSON string for external links
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Files     []File   `json:"files,omitempty" gorm:"foreignKey:UserID"`
	Ratings   []Rating `json:"ratings,omitempty" gorm:"foreignKey:RaterID"`
	Received  []Rating `json:"received,omitempty" gorm:"foreignKey:RatedID"`
	Follows   []Follow `json:"follows,omitempty" gorm:"foreignKey:FollowerID"`
	Followers []Follow `json:"followers,omitempty" gorm:"foreignKey:FollowedID"`
}

type Follow struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	FollowerID uint      `json:"follower_id" gorm:"not null"`
	FollowedID uint      `json:"followed_id" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`

	Follower User `json:"follower,omitempty" gorm:"foreignKey:FollowerID"`
	Followed User `json:"followed,omitempty" gorm:"foreignKey:FollowedID"`
}
