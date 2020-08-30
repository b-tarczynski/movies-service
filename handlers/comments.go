package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/BarTar213/movies-service/config"
	"github.com/BarTar213/movies-service/models"
	"github.com/BarTar213/movies-service/storage"
	"github.com/gin-gonic/gin"
)

type CommentHandlers struct {
	storage storage.Storage
	headers *config.Headers
	logger  *log.Logger
}

func NewCommentHandlers(storage storage.Storage, logger *log.Logger, headers *config.Headers) *CommentHandlers {
	return &CommentHandlers{
		storage: storage,
		headers: headers,
		logger:  logger,
	}
}

func (h *CommentHandlers) GetComments(c *gin.Context) {
	id, err := strconv.Atoi(c.Param(movieIdKey))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}

	comments, err := h.storage.GetMovieComments(id)
	if err != nil {
		handlePostgresError(c, h.logger, err)
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (h *CommentHandlers) AddComment(c *gin.Context) {
	account := models.AccountInfo{}
	err := c.ShouldBindHeader(&account)
	if err != nil {
		c.JSON(http.StatusForbidden, models.Response{Error: invalidHeadersErr})
		return
	}

	movieId, err := strconv.Atoi(c.Param(movieIdKey))
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
	comment.UserId = account.AccountId

	err = h.storage.AddMovieComment(&comment)
	if err != nil {
		handlePostgresError(c, h.logger, err)
		return
	}

	c.JSON(http.StatusOK, comment)
}
