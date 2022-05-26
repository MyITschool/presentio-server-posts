package models

type Comment struct {
	ID     int64  `gorm:"primaryKey" json:"id"`
	UserID int64  `json:"userId" binding:"required"`
	User   *User  `json:"user" binding:"required"`
	PostID int64  `json:"postId" binding:"required"`
	Text   string `json:"text" binding:"required"`
}
