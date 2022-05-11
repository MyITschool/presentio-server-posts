package models

import "time"

type Post struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	User      User      `gorm:"foreignKey:user_id" json:"user" binding:"required"`
	UserId    int64     `binding:"required"`
	Text      string    `json:"text" binding:"required"`
	CreatedAt time.Time `json:"createdAt" binding:"required" gorm:"type:timestamp"`
}
