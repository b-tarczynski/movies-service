package handlers

import (
	"fmt"
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
	go h.AddRecentViewedMovie(c.Copy(), id)

	c.JSON(http.StatusOK, movie)
}

func (h *MovieHandlers) ListMovies(c *gin.Context) {
	params := models.PaginationParams{OrderBy: "vote_average"}
	err := c.ShouldBindQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: invalidPaginationQueryParams})
		return
	}

	title := fmt.Sprintf("%%%s%%", c.Query("title"))

	movies, err := h.postgres.ListMovies(title, &params)
	if err != nil {
		handlePostgresError(c, h.logger, err)
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (h *MovieHandlers) LikeMovie(c *gin.Context) {
	movieId, err := strconv.Atoi(c.Param(movieIdKey))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}

	account := models.AccountInfo{}
	err = c.ShouldBindHeader(&account)
	if err != nil {
		c.JSON(http.StatusForbidden, models.Response{Error: invalidHeadersErr})
		return
	}

	err = h.postgres.LikeMovie(account.AccountId, movieId)
	if err != nil {
		handlePostgresError(c, h.logger, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{})
}

func (h *MovieHandlers) AddRecentViewedMovie(c *gin.Context, movieId int) {
	account := models.AccountInfo{}
	err := c.ShouldBindHeader(&account)
	if err != nil {
		return
	}

	err = h.postgres.AddRecentViewedMovie(account.AccountId, movieId)
	if err != nil {
		h.logger.Printf("addRecentViewedMovie: %s", err)
	}
}
