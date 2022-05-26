package models

type Like struct {
	ID     int64 `gorm:"primaryKey" json:"id"`
	UserID int64 `binding:"required" json:"userId"`
	PostID int64 `binding:"required" json:"postId"`
}
