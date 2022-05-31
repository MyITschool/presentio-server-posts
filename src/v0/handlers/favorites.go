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

type FavoritesHandler struct {
	PostsRepo     repo.PostsRepo
	FavoritesRepo repo.FavoritesRepo
}

func SetupFavoritesHandler(group *gin.RouterGroup, handler *FavoritesHandler) {
	group.GET("/:page", handler.getFavorites)
	group.POST("/:id", handler.addFavorite)
	group.DELETE("/:id", handler.deleteFavorite)
}

func (h *FavoritesHandler) getFavorites(c *gin.Context) {
	token, err := util.ValidateAccessTokenHeader(c.GetHeader("Authorization"))

	if err != nil {
		c.Status(util.HandleTokenError(err))
		return
	}

	page, err := strconv.Atoi(c.Param("page"))

	if err != nil {
		c.Status(404)

		return
	}

	claims, ok := token.Claims.(*util.UserClaims)

	if !ok {
		c.Status(403)
		return
	}

	favoriteIds, err := h.FavoritesRepo.GetUserFavorites(claims.ID, page)

	if err != nil {
		c.Status(500)
		return
	}

	if len(favoriteIds) == 0 {
		c.JSON(200, make([]interface{}, 0))
		return
	}

	favoritePosts, err := h.PostsRepo.FindIdIn(favoriteIds, claims.ID)

	if err != nil {
		c.Status(500)
		return
	}

	c.JSON(200, favoritePosts)
}

func (h *FavoritesHandler) addFavorite(c *gin.Context) {
	token, err := util.ValidateAccessTokenHeader(c.GetHeader("Authorization"))

	if err != nil {
		c.Status(util.HandleTokenError(err))
		return
	}

	postId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.Status(404)

		return
	}

	claims, ok := token.Claims.(*util.UserClaims)

	if !ok {
		c.Status(403)
		return
	}

	err = h.FavoritesRepo.Transaction(func(tx *gorm.DB) error {
		postsRepo := repo.CreatePostsRepo(tx)
		favoritesRepo := repo.CreateFavoritesRepo(tx)

		_, err := favoritesRepo.FindByIds(claims.ID, postId)

		if err == nil {
			c.Status(409)
			return nil
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		_, err = postsRepo.FindMinimal(postId)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(404)
			return nil
		}

		if err != nil {
			return err
		}

		err = favoritesRepo.Create(&models.Favorite{
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
	}
}

func (h *FavoritesHandler) deleteFavorite(c *gin.Context) {
	token, err := util.ValidateAccessTokenHeader(c.GetHeader("Authorization"))

	if err != nil {
		c.Status(util.HandleTokenError(err))
		return
	}

	postId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.Status(404)

		return
	}

	claims, ok := token.Claims.(*util.UserClaims)

	if !ok {
		c.Status(403)
		return
	}

	err = h.FavoritesRepo.Transaction(func(tx *gorm.DB) error {
		favoritesRepo := repo.CreateFavoritesRepo(tx)

		_, err := favoritesRepo.FindByIds(claims.ID, postId)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(409)
			return nil
		}

		if err != nil {
			return err
		}

		_, err = favoritesRepo.Delete(claims.ID, postId)

		if err != nil {
			return err
		}

		c.Status(204)
		return nil
	})

	if err != nil {
		c.Status(500)
	}
}
