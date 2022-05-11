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

	result := r.db.Where("id = ?", postId).First(&post)

	return &post, result.Error
}
