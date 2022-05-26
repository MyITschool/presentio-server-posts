package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"presentio-server-posts/src/v0/models"
	"presentio-server-posts/src/v0/repo"
	"presentio-server-posts/src/v0/util"
	"strconv"
)

// likes, comments, reposts

type LikesHandler struct {
	PostsRepo repo.PostsRepo
	LikesRepo repo.LikesRepo
}

func SetupLikesHandler(group *gin.RouterGroup, handler *LikesHandler) {
	group.POST("/:id", handler.likePost)
	group.DELETE("/:id", handler.removeLike)
}

func (h *LikesHandler) likePost(c *gin.Context) {
	postId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.Status(404)
		return
	}

	token, err := util.ValidateAccessTokenHeader(c.GetHeader("Authorization"))

	if err != nil {
		c.Status(util.HandleTokenError(err))
		return
	}

	claims, ok := token.Claims.(*util.UserClaims)

	if !ok {
		c.Status(403)
		return
	}

	err = h.LikesRepo.Transaction(func(tx *gorm.DB) error {
		_, err = h.LikesRepo.FindByIds(claims.ID, postId)

		if err == nil {
			c.Status(409)
			return nil
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		rows, err := h.PostsRepo.IncrementLikes(postId)

		if rows == 0 {
			c.Status(404)
			return nil
		}

		if err != nil {
			return err
		}

		err = h.LikesRepo.Create(&models.Like{
			UserID: claims.ID,
			PostID: postId,
		})

		if err != nil {
			return err
		}

		c.Status(201)
		return nil
	})

	if err != nil {
		c.Status(500)
		return
	}
}

func (h *LikesHandler) removeLike(c *gin.Context) {
	postId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.Status(404)
		return
	}

	token, err := util.ValidateAccessTokenHeader(c.GetHeader("Authorization"))

	if err != nil {
		c.Status(util.HandleTokenError(err))
		return
	}

	claims, ok := token.Claims.(*util.UserClaims)

	if !ok {
		c.Status(403)
		return
	}

	err = h.LikesRepo.Transaction(func(tx *gorm.DB) error {
		_, err = h.LikesRepo.FindByIds(claims.ID, postId)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(409)
			return nil
		}

		if err != nil {
			return err
		}

		rows, err := h.PostsRepo.DecrementLikes(postId)

		if rows == 0 {
			c.Status(404)
			return nil
		}

		if err != nil {
			return err
		}

		_, err = h.LikesRepo.Delete(&models.Like{
			UserID: claims.ID,
			PostID: postId,
		})

		if err != nil {
			return err
		}

		c.Status(204)
		return nil
	})

	if err != nil {
		c.Status(500)
		return
	}
}
