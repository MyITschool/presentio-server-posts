package models

type Like struct {
	ID     int64 `gorm:"primaryKey"`
	UserID int64 `binding:"required"`
	PostID int64 `binding:"required"`
}
