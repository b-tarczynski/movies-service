package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/BarTar213/movies-service/models"
	"github.com/BarTar213/movies-service/storage"
	"github.com/BarTar213/movies-service/utils"
	notificator "github.com/BarTar213/notificator/client"
	"github.com/BarTar213/notificator/senders"
	"github.com/gin-gonic/gin"
)

const (
	commentId = "commId"
)

type CommentHandlers struct {
	storage     storage.Storage
	notificator notificator.Client
	logger      *log.Logger
}

func NewCommentHandlers(storage storage.Storage, notificator notificator.Client, logger *log.Logger) *CommentHandlers {
	return &CommentHandlers{
		storage:     storage,
		notificator: notificator,
		logger:      logger,
	}
}

func (h *CommentHandlers) ListComments(c *gin.Context) {
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

	comments, err := h.storage.ListMovieComments(id, &params)
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
		comment := &models.Comment{}
		err = h.storage.LikeComment(account.ID, commentId, comment)
		if err == nil {
			go h.sendNotification(comment, account)
		}
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

func (h *CommentHandlers) ListLikedComments(c *gin.Context) {
	account := utils.GetAccount(c)

	id, err := strconv.Atoi(c.Param(movieIdKey))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}

	commentIds, err := h.storage.ListLikedCommentsForMovie(id, account.ID)
	if err != nil {
		handlePostgresError(c, h.logger, err, commentResource)
		return
	}

	c.JSON(http.StatusOK, commentIds)
}


func (h *CommentHandlers) sendNotification(comment *models.Comment, account *models.AccountInfo) {
	internal := &senders.Internal{
		ResourceID: comment.Id,
		Resource:   "comment",
		Tag:        fmt.Sprintf("movie/%d", comment.MovieId),
		Recipients: []int{comment.UserId},
		Data: map[string]string{
			"user": account.Login,
		},
	}
	_, _, err := h.notificator.SendInternal(context.Background(), "commentLike", internal)
	if err != nil{
		h.logger.Printf("send internal notification: %s", err)
	}
}
