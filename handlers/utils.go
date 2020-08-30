package handlers

import (
	"log"
	"net/http"

	"github.com/BarTar213/movies-service/models"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

const (
	invalidHeadersErr      = "invalid account headers"
	invalidRequestBodyErr  = "invalid request body"
	invalidMovieIdParamErr = "invalid param - movieId"
	invalidCommentIdParamErr = "invalid param - commentId"
)

func handlePostgresError(c *gin.Context, l *log.Logger, err error) {
	if err == pg.ErrNoRows {
		c.JSON(http.StatusBadRequest, models.Response{Error: "resource with given identification doesn't exist"})
		return
	}
	l.Println(err)
	c.JSON(http.StatusInternalServerError, models.Response{Error: "storage error"})
}
