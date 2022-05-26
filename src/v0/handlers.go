package v0

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"presentio-server-posts/src/v0/handlers"
	"presentio-server-posts/src/v0/repo"
)

type Config struct {
	Db *gorm.DB
}

func SetupRouter(group *gin.RouterGroup, config *Config) {
	postsRepo := repo.CreatePostsRepo(config.Db)
	likesRepo := repo.CreateLikesRepo(config.Db)

	handlers.SetupPostsHandler(group.Group("/posts"), &handlers.PostsHandler{
		PostsRepo: postsRepo,
	})

	handlers.SetupLikesHandler(group.Group("/likes"), &handlers.LikesHandler{
		PostsRepo: postsRepo,
		LikesRepo: likesRepo,
	})
}
