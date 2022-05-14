package repo

import (
	"gorm.io/gorm"
	"presentio-server-posts/src/v0/models"
)

type PostsRepo struct {
	db *gorm.DB
}

func CreatePostsRepo(db *gorm.DB) PostsRepo {
	return PostsRepo{
		db,
	}
}

func (r *PostsRepo) FindById(postId int64) (*models.Post, error) {
	var post models.Post

	result := r.db.
		Joins("User").
		Joins("Source").
		Joins("SourceUser").
		Where("posts.id = ?", postId).
		First(&post)

	return &post, result.Error
}

func (r *PostsRepo) Create(post *models.Post) error {
	return r.db.Create(&post).Error
}

func (r *PostsRepo) DeleteWithGuard(postId int64, userID int64) error {
	return r.db.
		Where("id = ?", postId).
		Where("user_id = ?", userID).
		Model((*models.Post)(nil)).
		Updates(map[string]interface{}{"text": "", "deleted": true, "reposts": 0, "likes": 0}).
		Error
}

func (r *PostsRepo) GetUserPosts(userId int64, page int) ([]models.Post, error) {
	var posts []models.Post

	err := r.db.
		Where("posts.user_id = ?", userId).
		Where("posts.deleted = false").
		Joins("User").
		Joins("Source").
		Joins("SourceUser").
		Limit(20).
		Offset(20 * page).
		Find(&posts).
		Error

	return posts, err
}
