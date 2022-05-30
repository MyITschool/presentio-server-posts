package repo

import (
	"database/sql"
	"gorm.io/gorm"
	"presentio-server-posts/src/v0/models"
	"strings"
)

type PostsRepo struct {
	db *gorm.DB
}

func CreatePostsRepo(db *gorm.DB) PostsRepo {
	return PostsRepo{
		db,
	}
}

func (r *PostsRepo) FindById(postId int64, userId int64) (*models.Post, error) {
	var post models.Post

	result := r.db.
		Joins("User").
		Joins("Source").
		Joins("SourceUser").
		Joins("Liked", r.db.Where(&models.Like{UserID: userId, PostID: postId})).
		Preload("Tags").
		Preload("Source.Tags").
		Where("posts.deleted = false").
		Where("posts.id = ?", postId).
		First(&post)

	return &post, result.Error
}

func (r *PostsRepo) Create(post *models.Post) error {
	return r.db.Create(&post).Error
}

func (r *PostsRepo) DeleteWithGuard(postId int64, userId int64) (int64, error) {
	tx := r.db.
		Where("id = ?", postId).
		Where("user_id = ?", userId).
		Where("deleted = false").
		Model((*models.Post)(nil)).
		Updates(map[string]interface{}{"text": "", "deleted": true, "reposts": 0, "likes": 0})

	return tx.RowsAffected, tx.Error
}

func (r *PostsRepo) GetUserPosts(userId int64, page int, myUserId int64) ([]models.Post, error) {
	var posts []models.Post

	err := r.db.
		Where("posts.user_id = ?", userId).
		Where("posts.deleted = false").
		Joins("User").
		Joins("Source").
		Joins("SourceUser").
		Joins("Liked", r.db.Where(&models.Like{UserID: myUserId})).
		Preload("Tags").
		Preload("Source.Tags").
		Limit(20).
		Offset(20 * page).
		Order("posts.id DESC").
		Find(&posts).
		Error

	return posts, err
}

func (r *PostsRepo) FindByQuery(tags []string, keywords []string, page int, userId int64) ([]models.Post, error) {
	var posts []models.Post

	tx := r.db.
		Distinct().
		Where("posts.deleted = false").
		Joins("Liked", "? = likes.user_id and posts.id = likes.post_id", userId).
		Joins("User").
		Joins("Source").
		Joins("SourceUser").
		Joins("JOIN post_tags pt ON posts.id = pt.post_id").
		Joins("JOIN tags t ON t.id = pt.tag_id").
		Preload("Tags").
		Preload("Source.Tags").
		Limit(20).
		Offset(20 * page).
		Order("posts.id")

	if len(tags) > 0 {
		tx = tx.Where("t.name IN ?", tags)
	}

	if len(keywords) > 0 {
		tx = tx.Where("posts.ts @@ to_tsquery(getlang(posts.id), ?)", strings.Join(keywords, "|"))
	}

	err := tx.Find(&posts).Error

	return posts, err
}

func (r *PostsRepo) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return r.db.Transaction(fc, opts...)
}

func (r *PostsRepo) IncrementLikes(postId int64) (int64, error) {
	tx := r.db.
		Exec("UPDATE posts SET likes = likes + 1 WHERE id = ? AND deleted = false", postId)

	return tx.RowsAffected, tx.Error
}

func (r *PostsRepo) DecrementLikes(postId int64) (int64, error) {
	tx := r.db.
		Exec("UPDATE posts SET likes = likes - 1 WHERE id = ? AND deleted = false", postId)

	return tx.RowsAffected, tx.Error
}

func (r *PostsRepo) IncrementComments(postId int64) (int64, error) {
	tx := r.db.
		Exec("UPDATE posts SET comments = comments + 1 WHERE id = ? AND deleted = false", postId)

	return tx.RowsAffected, tx.Error
}
