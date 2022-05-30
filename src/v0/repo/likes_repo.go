package repo

import (
	"database/sql"
	"gorm.io/gorm"
	"presentio-server-posts/src/v0/models"
)

type LikesRepo struct {
	db *gorm.DB
}

func CreateLikesRepo(db *gorm.DB) LikesRepo {
	return LikesRepo{
		db,
	}
}

func (r *LikesRepo) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return r.db.Transaction(fc, opts...)
}

func (r *LikesRepo) FindByIds(userId int64, postId int64) (*models.Like, error) {
	var like models.Like

	result := r.db.
		Where("user_id = ?", userId).
		Where("post_id = ?", postId).
		First(&like)

	return &like, result.Error
}

func (r *LikesRepo) Create(like *models.Like) error {
	return r.db.Create(like).Error
}

func (r *LikesRepo) Delete(userId int64, postId int64) (int64, error) {
	tx := r.db.
		Where("user_id = ? and post_id = ?", userId, postId).
		Delete(&models.Like{})

	return tx.RowsAffected, tx.Error
}
