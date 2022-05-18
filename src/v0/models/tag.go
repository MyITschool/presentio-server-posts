package models

type Tag struct {
	ID     int64 `json:"id" gorm:"primaryKey"`
	TagId  int64 `json:"tagId" binding:"required"`
	PostId int64 `json:"postId" binding:"required"`
}
