package repo

import (
	"database/sql"
	"gorm.io/gorm"
	"presentio-server-posts/src/v0/models"
)

type CommentsRepo struct {
	db *gorm.DB
}

func CreateCommentsRepo(db *gorm.DB) CommentsRepo {
	return CommentsRepo{
		db,
	}
}

func (r *CommentsRepo) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *CommentsRepo) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return r.db.Transaction(fc, opts...)
}

func (r *CommentsRepo) GetPostComments(postId int64, page int) ([]models.Comment, error) {
	var comments []models.Comment

	result := r.db.
		Where("post_id = ?", postId).
		Joins("User").
		Offset(page * 20).
		Limit(20).
		Order("id DESC").
		Find(&comments).
		Error

	return comments, result
}
