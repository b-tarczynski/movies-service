package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/BarTar213/movies-service/models"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

const (
	invalidRequestBodyErr        = "invalid request body"
	invalidMovieIdParamErr       = "invalid param - commentId"
	invalidCommentIdParamErr     = "invalid param - commentId"
	invalidLikedParamErr         = "invalid param - liked"
	invalidPaginationQueryParams = "invalid pagination query params"

	likedParam      = "liked"
	movieIdQuery    = "movie_id"
	movieResource   = "movie"
	commentResource = "comment"
)

func handlePostgresError(c *gin.Context, l *log.Logger, err error, resource string) {
	if err == pg.ErrNoRows {
		c.JSON(http.StatusBadRequest, models.Response{Error: fmt.Sprintf("%s with given identification doesn't exist", resource)})
		return
	}
	l.Println(err)

	msg := ""
	pgErr, ok := err.(pg.Error)
	if ok {
		switch pgErr.Field('C') {
		case "23503":
			msg = fmt.Sprintf("%s with given information doesn't exists", resource)
		case "23505":
			msg = fmt.Sprintf("%s with given information already exists", resource)
		}
		if len(msg) > 0 {
			c.JSON(http.StatusBadRequest, models.Response{Error: msg})
			return
		}
	}

	c.JSON(http.StatusInternalServerError, models.Response{Error: "storage error"})
}
