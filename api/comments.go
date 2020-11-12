package api

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/BarTar213/movies-service/models"
	"github.com/BarTar213/movies-service/storage"
	"github.com/BarTar213/movies-service/utils"
	"github.com/gin-gonic/gin"
)

const (
	commentId = "commId"
)

type CommentHandlers struct {
	storage storage.Storage
	logger  *log.Logger
}

func NewCommentHandlers(storage storage.Storage, logger *log.Logger) *CommentHandlers {
	return &CommentHandlers{
		storage: storage,
		logger:  logger,
	}
}

func (h *CommentHandlers) GetComments(c *gin.Context) {
	id, err := strconv.Atoi(c.Query(movieIdQuery))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: invalidMovieIdParamErr})
		return
	}

	params := models.PaginationParams{OrderBy: "create_date"}
	err = c.ShouldBindQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: invalidPaginationQueryParams})
		return
	}

	comments, err := h.storage.GetMovieComments(id, &params)
	if err != nil {
		handlePostgresError(c, h.logger, err, commentResource)
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (h *CommentHandlers) LikeComment(c *gin.Context) {
	account := utils.GetAccount(c)

	commentId, err := strconv.Atoi(c.Param(commentId))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: invalidCommentIdParamErr})
		return
	}

	liked, err := strconv.ParseBool(c.Query(likedParam))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: invalidLikedParamErr})
		return
	}

	if liked {
		err = h.storage.DeleteCommentLike(account.ID, commentId)
	} else {
		err = h.storage.LikeComment(account.ID, commentId)
	}
	if err != nil {
		handlePostgresError(c, h.logger, err, commentResource)
		return
	}

	c.JSON(http.StatusOK, models.Response{})
}

func (h *CommentHandlers) AddComment(c *gin.Context) {
	account := utils.GetAccount(c)

	movieId, err := strconv.Atoi(c.Query(movieIdQuery))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: invalidMovieIdParamErr})
		return
	}

	comment := models.Comment{}
	err = c.ShouldBindJSON(&comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: invalidRequestBodyErr})
		return
	}

	now := time.Now()
	comment.UpdateDate = now
	comment.CreateDate = now
	comment.MovieId = movieId
	comment.UserId = account.ID

	err = h.storage.AddMovieComment(&comment)
	if err != nil {
		handlePostgresError(c, h.logger, err, commentResource)
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (h *CommentHandlers) UpdateComment(c *gin.Context) {
	account := utils.GetAccount(c)

	commentId, err := strconv.Atoi(c.Param(commentId))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: invalidCommentIdParamErr})
		return
	}

	comment := models.Comment{}
	err = c.ShouldBindJSON(&comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: invalidRequestBodyErr})
		return
	}

	comment.Id = commentId
	comment.UserId = account.ID
	comment.UpdateDate = time.Now()

	err = h.storage.UpdateComment(&comment)
	if err != nil {
		handlePostgresError(c, h.logger, err, commentResource)
		return
	}

	c.JSON(http.StatusOK, comment)
}

func (h *CommentHandlers) DeleteComment(c *gin.Context) {
	account := utils.GetAccount(c)

	commentId, err := strconv.Atoi(c.Param(commentId))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: invalidCommentIdParamErr})
		return
	}

	comment := models.Comment{}
	comment.Id = commentId
	comment.UserId = account.ID

	err = h.storage.DeleteComment(&comment)
	if err != nil {
		handlePostgresError(c, h.logger, err, commentResource)
		return
	}

	c.JSON(http.StatusOK, models.Response{})
}
