package handlers

import (
	"errors"
	"github.com/abadojack/whatlanggo"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"presentio-server-posts/src/v0/models"
	"presentio-server-posts/src/v0/repo"
	"presentio-server-posts/src/v0/util"
	"strconv"
	"strings"
	"time"
)

type PostsHandler struct {
	PostsRepo repo.PostsRepo
}

func SetupPostsHandler(group *gin.RouterGroup, handler *PostsHandler) {
	group.GET("/:id", handler.getPost)
	group.POST("/", handler.createPost)
	//group.GET("/recommended/:page", handler.getRecommended)
	group.DELETE("/:id", handler.deletePost)
	group.GET("/user/:id/:page", handler.getUserPosts)
	group.GET("/user/self/:page", handler.getUserPostsSelf)
	group.GET("/search/:page", handler.search)
}

func (h *PostsHandler) getPost(c *gin.Context) {
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

	postId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.Status(404)

		return
	}

	post, err := h.PostsRepo.FindById(postId, claims.ID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(404)
		} else {
			c.Status(500)
		}

		return
	}

	post.Own = post.UserID == claims.ID

	c.Header("Cache-Control", "public, max-age=18000")
	c.Header("Pragma", "")
	c.Header("Expires", "")
	c.Header("Vary", "")

	c.JSON(200, post)
}

type PostParams struct {
	Text         string
	Tags         []string
	Attachments  []string
	PhotoRatio   *float64
	SourceID     *int64
	SourceUserId *int64
}

func validateParams(params *PostParams) bool {
	if params.SourceUserId == nil && params.SourceID != nil {
		return false
	}

	if params.SourceUserId != nil && params.SourceID == nil {
		return false
	}

	if len(params.Text) < 1 {
		return false
	}

	if params.SourceID != nil {
		return params.Tags != nil || params.Attachments != nil || params.PhotoRatio != nil
	}

	if len(params.Tags) < 1 || len(params.Tags) > 5 {
		return false
	}

	if len(params.Attachments) < 1 || len(params.Attachments) > 10 {
		return false
	}

	if *params.PhotoRatio <= 0 {
		return false
	}

	return true
}

func (h *PostsHandler) createPost(c *gin.Context) {
	var params PostParams

	err := c.ShouldBindJSON(&params)

	if err != nil || !validateParams(&params) {
		c.Status(422)
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

	lang := strings.ToLower(whatlanggo.DetectLang(params.Text).String())

	tags := make([]models.Tag, 0, len(params.Tags))

	for i := 0; i < len(params.Tags); i++ {
		tags = append(tags, models.Tag{
			Name: params.Tags[i],
		})
	}

	if params.PhotoRatio == nil {
		r := 0.0
		params.PhotoRatio = &r
	}

	post := models.Post{
		UserID:       claims.ID,
		Text:         params.Text,
		CreatedAt:    time.Now(),
		SourceID:     params.SourceID,
		SourceUserID: params.SourceUserId,
		Lang:         lang,
		Attachments:  params.Attachments,
		PhotoRatio:   *params.PhotoRatio,
	}

	err = h.PostsRepo.Transaction(func(tx *gorm.DB) error {
		postsRepo := repo.CreatePostsRepo(tx)
		tagsRepo := repo.CreateTagsRepo(tx)

		if post.SourceID != nil {
			rows, err := postsRepo.IncrementReposts(*post.SourceID)

			if err != nil {
				return err
			}

			if rows == 0 {
				c.Status(404)
				return nil
			}
		}

		err = postsRepo.Create(&post)

		if err != nil {
			return err
		}

		if len(tags) == 0 {
			return nil
		}

		err = tagsRepo.BulkInsert(tags)

		if err != nil {
			return err
		}

		err = tagsRepo.BulkInsertRelation(tags, post.ID)

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

func (h *PostsHandler) deletePost(c *gin.Context) {
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

	rows, err := h.PostsRepo.DeleteWithGuard(postId, claims.ID)

	if err != nil {
		c.Status(500)
		return
	}

	if rows == 0 {
		c.Status(404)
		return
	}

	c.Status(204)
}

func (h *PostsHandler) getUserPosts(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.Status(404)

		return
	}

	h.doGetUserPosts(userId, c)
}

func (h *PostsHandler) getUserPostsSelf(c *gin.Context) {
	h.doGetUserPosts(-1, c)
}

func (h *PostsHandler) doGetUserPosts(userId int64, c *gin.Context) {
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

	if userId == -1 {
		userId = claims.ID
	}

	posts, err := h.PostsRepo.GetUserPosts(userId, page, claims.ID)

	if err != nil {
		c.Status(500)
		return
	}

	for i := 0; i < len(posts); i++ {
		posts[i].Own = posts[i].UserID == claims.ID
	}

	cache := "18000"

	if userId == claims.ID {
		cache = "300"
	}

	c.Header("Cache-Control", "public, max-age="+cache)
	c.Header("Pragma", "")
	c.Header("Expires", "")

	c.JSON(200, posts)
}

func (h *PostsHandler) search(c *gin.Context) {
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

	tags := c.QueryArray("tag")
	keywords := c.QueryArray("keyword")

	posts, err := h.PostsRepo.FindByQuery(tags, keywords, page, claims.ID)

	if err != nil {
		c.Status(500)
		return
	}

	for i := 0; i < len(posts); i++ {
		posts[i].Own = posts[i].UserID == claims.ID
	}

	c.Header("Cache-Control", "public, max-age=300")
	c.Header("Pragma", "")
	c.Header("Expires", "")
	c.JSON(200, posts)
}
