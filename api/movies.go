package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/BarTar213/movies-service/models"
	"github.com/BarTar213/movies-service/storage"
	"github.com/BarTar213/movies-service/tmdb"
	"github.com/BarTar213/movies-service/utils"
	"github.com/gin-gonic/gin"
)

const (
	movieIdKey = "movieId"
)

type MovieHandlers struct {
	storage storage.Storage
	tmdb    tmdb.Client
	logger  *log.Logger
}

func NewMovieHandlers(postgres storage.Storage, tmdb tmdb.Client, logger *log.Logger) *MovieHandlers {
	return &MovieHandlers{
		storage: postgres,
		tmdb:    tmdb,
		logger:  logger,
	}
}

func (h *MovieHandlers) GetMovie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param(movieIdKey))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}

	movie := &models.Movie{Id: id}
	err = h.storage.GetMovie(movie)
	if err != nil {
		handlePostgresError(c, h.logger, err, movieResource)
		return
	}
	go h.AddRecentViewedMovie(c.Copy(), id)

	c.JSON(http.StatusOK, movie)
}

func (h *MovieHandlers) ListMovies(c *gin.Context) {
	params := models.PaginationParams{OrderBy: "revenue DESC"}
	err := c.ShouldBindQuery(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: invalidPaginationQueryParams})
		return
	}

	title := fmt.Sprintf("%%%s%%", c.Query("title"))

	movies, err := h.storage.ListMovies(title, &params)
	if err != nil {
		handlePostgresError(c, h.logger, err, movieResource)
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

	account := utils.GetAccount(c)

	liked, err := strconv.ParseBool(c.Query(likedParam))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: invalidLikedParamErr})
		return
	}

	if liked {
		err = h.storage.DeleteMovieLike(account.ID, movieId)
	} else {
		err = h.storage.LikeMovie(account.ID, movieId)
	}
	if err != nil {
		handlePostgresError(c, h.logger, err, movieResource)
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

	err = h.storage.AddRecentViewedMovie(account.ID, movieId)
	if err != nil {
		h.logger.Printf("addRecentViewedMovie: %s", err)
	}
}
