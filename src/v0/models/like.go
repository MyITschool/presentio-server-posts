package models

type Like struct {
	ID     int64 `gorm:"primaryKey" json:"id"`
	UserID int64 `json:"userId" binding:"required"`
	PostID int64 `json:"postId" binding:"required"`
}
