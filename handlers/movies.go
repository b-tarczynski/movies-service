package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/BarTar213/movies-service/models"
	"github.com/BarTar213/movies-service/storage"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

const (
	idKey = "id"
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
	id, err := strconv.Atoi(c.Param(idKey))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}

	movie := &models.Movie{Id: id}
	err = h.postgres.GetMovie(movie)
	if err != nil {
		h.handlePostgresError(c, err)
		return
	}

	c.JSON(http.StatusOK, movie)
}

func (h *MovieHandlers) handlePostgresError(c *gin.Context, err error) {
	if err == pg.ErrNoRows {
		c.JSON(http.StatusBadRequest, models.Response{Error: "resource with given identification doesn't exist"})
		return
	}
	h.logger.Println(err)
	c.JSON(http.StatusInternalServerError, models.Response{Error: "storage error"})
}
