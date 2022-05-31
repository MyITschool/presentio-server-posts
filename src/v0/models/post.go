package models

import (
	"github.com/lib/pq"
	"time"
)

type Post struct {
	ID           int64          `gorm:"primaryKey" json:"id"`
	UserID       int64          `json:"userId" binding:"required"`
	User         *User          `json:"user" binding:"required"`
	Text         string         `json:"text" binding:"required"`
	CreatedAt    time.Time      `json:"createdAt" binding:"required" gorm:"type:timestamp"`
	Reposts      int64          `json:"reposts" binding:"required"`
	Likes        int64          `json:"likes" binding:"required"`
	Comments     int64          `json:"comments" binding:"required"`
	Tags         []*Tag         `json:"tags" gorm:"many2many:post_tags"`
	SourceID     *int64         `json:"sourceId" binding:"required"`
	Source       *Post          `json:"source" binding:"required"`
	SourceUserID *int64         `json:"sourceUserId" binding:"required"`
	SourceUser   *User          `json:"sourceUser" binding:"required"`
	PhotoRatio   float64        `json:"photoRatio" binding:"required"`
	Attachments  pq.StringArray `json:"attachments" binding:"required" gorm:"type:varchar[]"`
	Lang         string         `json:"lang" binding:"required"`
	Liked        Like           `json:"liked" binding:"required"`
	Favorite     Favorite       `json:"favorite" binding:"required"`
	Deleted      bool           `json:"deleted" binding:"required"`
	Own          bool           `json:"own" gorm:"-"`
}
