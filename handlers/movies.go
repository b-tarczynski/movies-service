package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/BarTar213/movies-service/models"
	"github.com/BarTar213/movies-service/storage"
	"github.com/gin-gonic/gin"
)

const (
	movieIdKey = "movieId"
)

type MovieHandlers struct {
	postgres storage.Storage
	logger   *log.Logger
}

func NewMovieHandlers(postgres storage.Storage, logger *log.Logger) *MovieHandlers {
	return &MovieHandlers{
		postgres: postgres,
		logger:   logger,
	}
}

func (h *MovieHandlers) GetMovie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param(movieIdKey))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}

	movie := &models.Movie{Id: id}
	err = h.postgres.GetMovie(movie)
	if err != nil {
		handlePostgresError(c, h.logger, err)
		return
	}

	c.JSON(http.StatusOK, movie)
}
