package repo

import (
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

func (r *PostsRepo) FindById(postId int64) (*models.Post, error) {
	var post models.Post

	result := r.db.
		Joins("User").
		Joins("Source").
		Joins("SourceUser").
		Preload("Tags").
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
		Preload("Tags").
		Limit(20).
		Offset(20 * page).
		Find(&posts).
		Error

	return posts, err
}

func (r *PostsRepo) FindByQuery(tags []string, keywords []string, page int) ([]models.Post, error) {
	var posts []models.Post

	tx := r.db.
		Where("posts.deleted = false").
		Joins("User").
		Joins("Source").
		Joins("SourceUser").
		Joins("JOIN post_tags pt ON posts.id = pt.post_id").
		Joins("JOIN tags t ON t.id = pt.tag_id").
		Preload("Tags").
		Limit(20).
		Offset(20 * page)

	if len(tags) > 0 {
		tx = tx.Where("t.name IN ?", tags)
	}

	if len(keywords) > 0 {
		tx = tx.Where("posts.ts @@ to_tsquery(getlang(posts.id), ?)", strings.Join(keywords, "|"))
	}

	err := tx.Find(&posts).Error

	return posts, err
}
