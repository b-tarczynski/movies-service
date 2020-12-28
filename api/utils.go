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
	invalidMovieIdParamErr       = "invalid param - movieId"
	invalidCommentIdParamErr     = "invalid param - commentId"
	invalidLikedParamErr         = "invalid param - liked"
	invalidPaginationQueryParams = "invalid pagination query params"

	likedParam           = "liked"
	movieIdQuery         = "movie_id"
	movieResource        = "movie"
	movieCommentResource = "movie comment"
	commentResource      = "comment"
	creditsResource      = "credits"
	ratingResource       = "rating"
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

func handleTMDBError(c *gin.Context, l *log.Logger, status int, err error, resource string) {
	if err != nil {
		l.Println(err)
		c.JSON(http.StatusInternalServerError, models.Response{Error: "api error"})
		return
	}

	switch status {
	case http.StatusNotFound:
		c.JSON(http.StatusNotFound, models.Response{Error: fmt.Sprintf("%s with given information already exists", resource)})
		return
	}

	c.JSON(http.StatusInternalServerError, models.Response{Error: "api error"})
}
