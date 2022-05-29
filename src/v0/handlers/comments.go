package handlers

import (
	"github.com/gin-gonic/gin"
	"presentio-server-posts/src/v0/models"
	"presentio-server-posts/src/v0/repo"
	"presentio-server-posts/src/v0/util"
	"strconv"
)

type CommentsHandler struct {
	PostsRepo    repo.PostsRepo
	CommentsRepo repo.CommentsRepo
}

func SetupCommentsHandler(group *gin.RouterGroup, handler *CommentsHandler) {
	group.POST("/:id", handler.createComment)
	group.GET("/:id/:page", handler.getPostComments)
}

type CommentParams struct {
	Text string
}

func (h *CommentsHandler) createComment(c *gin.Context) {
	var params CommentParams

	err := c.ShouldBindJSON(&params)

	if err != nil {
		c.Status(422)
		return
	}

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

	rows, err := h.PostsRepo.IncrementComments(postId)

	if err != nil {
		c.Status(500)
		return
	}

	if rows == 0 {
		c.Status(404)
		return
	}

	comment := &models.Comment{
		Text:   params.Text,
		UserID: claims.ID,
		PostID: postId,
	}

	err = h.CommentsRepo.Create(comment)

	if err != nil {
		c.Status(500)
		return
	}

	c.JSON(201, comment)
}

func (h *CommentsHandler) getPostComments(c *gin.Context) {
	postId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.Status(404)
		return
	}

	page, err := strconv.ParseInt(c.Param("page"), 10, 32)

	if err != nil {
		c.Status(404)
		return
	}

	_, err = util.ValidateAccessTokenHeader(c.GetHeader("Authorization"))

	if err != nil {
		c.Status(util.HandleTokenError(err))
		return
	}

	comments, err := h.CommentsRepo.GetPostComments(postId, int(page))

	if err != nil {
		c.Status(500)
		return
	}

	c.JSON(200, comments)
}
