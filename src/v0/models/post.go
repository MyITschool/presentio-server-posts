package models

import "time"

type Post struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    int64     `binding:"required"`
	User      User      `gorm:"foreignKey:UserID;references:id;" json:"user" binding:"required"`
	Text      string    `json:"text" binding:"required"`
	CreatedAt time.Time `json:"createdAt" binding:"required" gorm:"type:timestamp"`
}
