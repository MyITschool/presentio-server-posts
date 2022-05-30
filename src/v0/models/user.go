package models

type User struct {
	ID        int64  `gorm:"primaryKey" json:"id"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required"`
	PFPUrl    string `json:"pfpUrl" binding:"required"`
	Bio       string `json:"bio" binding:"required"`
	Followers int64  `json:"followers" binding:"required"`
	Following int64  `json:"following" binding:"required"`
}
