package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
	Todos    []Todo `json:"todos,omitempty"`
}

type Todo struct {
	gorm.Model
	UserID      uint   `json:"user_id"`
	Title       string `gorm:"not null" json:"title"`
	Description string `json:"description"`
	IsCompleted bool   `gorm:"default:false" json:"is_completed"`
}
