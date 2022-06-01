package repo

import (
	"database/sql"
	"gorm.io/gorm"
	"presentio-server-posts/src/v0/models"
)

type FavoritesRepo struct {
	db *gorm.DB
}

func CreateFavoritesRepo(db *gorm.DB) FavoritesRepo {
	return FavoritesRepo{
		db: db,
	}
}

func (r *FavoritesRepo) GetUserFavorites(userId int64, page int) ([]int64, error) {
	var results []int64

	tx := r.db.
		Where("user_id = ?", userId).
		Offset(page*20).
		Limit(20).
		Order("id DESC").
		Model(&models.Favorite{}).
		Pluck("post_id", &results)

	return results, tx.Error
}

func (r *FavoritesRepo) Create(favorite *models.Favorite) error {
	return r.db.Create(favorite).Error
}

func (r *FavoritesRepo) FindByIds(userId int64, postId int64) (*models.Favorite, error) {
	var favorite models.Favorite

	tx := r.db.
		Where("user_id = ?", userId).
		Where("post_id = ?", postId).
		First(&favorite)

	return &favorite, tx.Error
}

func (r *FavoritesRepo) Delete(userId int64, postId int64) (int64, error) {
	tx := r.db.
		Where("user_id = ?", userId).
		Where("post_id = ?", postId).
		Delete(&models.Favorite{})

	return tx.RowsAffected, tx.Error
}

func (r *FavoritesRepo) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return r.db.Transaction(fc, opts...)
}
