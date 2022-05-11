package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"presentio-server-posts/src/v0/repo"
	"presentio-server-posts/src/v0/util"
	"strconv"
)

type PostsHandler struct {
	PostsRepo repo.PostsRepo
}

func CreatePostsHandler(group *gin.RouterGroup, postsRepo repo.PostsRepo) {
	handler := PostsHandler{
		PostsRepo: postsRepo,
	}

	group.GET("/:id", handler.getPost)
	//group.POST("/")
	//group.GET("/recommended")
	//group.DELETE("/:id")
}

func (h *PostsHandler) getPost(c *gin.Context) {
	_, err := util.ValidateAccessTokenHeader(c.GetHeader("Authorization"))

	if err != nil {
		util.HandleTokenError(err, c)

		return
	}

	postId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.Status(404)

		return
	}

	post, err := h.PostsRepo.FindById(postId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(404)
		} else {
			c.Status(500)
		}

		return
	}

	c.Header("Cache-Control", "public, max-age=18000")
	c.Header("Pragma", "")
	c.Header("Expires", "")

	c.JSON(200, post)
}
