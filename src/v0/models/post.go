package models

import "time"

type Post struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	User      User      `json:"user" binding:"required"`
	UserID    int64     `binding:"required"`
	Text      string    `json:"text" binding:"required"`
	CreatedAt time.Time `json:"createdAt" binding:"required" gorm:"type:timestamp"`
}
