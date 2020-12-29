package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/BarTar213/movies-service/models"
	"github.com/BarTar213/movies-service/storage"
	"github.com/BarTar213/movies-service/tmdb"
	"github.com/BarTar213/movies-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
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
		handlePostgresError(c, h.logger, err, movieCommentResource)
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

func (h *MovieHandlers) ListLikedMovies(c *gin.Context) {
	params := &models.PaginationParams{OrderBy: "revenue DESC"}
	err := c.ShouldBindQuery(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: invalidPaginationQueryParams})
		return
	}

	account := utils.GetAccount(c)

	movies, err := h.storage.ListLikedMovies(account.ID, params)
	if err != nil {
		handlePostgresError(c, h.logger, err, movieCommentResource)
		return
	}

	c.JSON(http.StatusOK, models.Response{Data: movies})
}

func (h *MovieHandlers) CheckLiked(c *gin.Context) {
	movieId, err := strconv.Atoi(c.Param(movieIdKey))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}

	account := utils.GetAccount(c)

	likedMovie := &models.LikedMovie{
		MovieId: movieId,
		UserId:  account.ID,
	}
	liked, err := h.storage.CheckLiked(likedMovie)
	if err != nil {
		handlePostgresError(c, h.logger, err, movieCommentResource)
		return
	}

	c.JSON(http.StatusOK, map[string]bool{"liked": liked})
}

func (h *MovieHandlers) RateMovie(c *gin.Context) {
	movieId, err := strconv.Atoi(c.Param(movieIdKey))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}

	account := utils.GetAccount(c)

	rating := &models.Rating{}
	err = c.ShouldBindJSON(rating)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: invalidRequestBodyErr})
		return
	}
	rating.UserId = account.ID
	rating.MovieId = movieId
	rating.CreateDate = time.Now()

	err = h.storage.AddRating(rating)
	if err != nil {
		handlePostgresError(c, h.logger, err, ratingResource)
		return
	}

	c.JSON(http.StatusCreated, rating)
}

func (h MovieHandlers) DeleteRating(c *gin.Context) {
	movieId, err := strconv.Atoi(c.Param(movieIdKey))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}

	account := utils.GetAccount(c)

	rating := &models.Rating{
		UserId:  account.ID,
		MovieId: movieId,
	}
	err = h.storage.DeleteRating(rating)
	if err != nil {
		handlePostgresError(c, h.logger, err, ratingResource)
		return
	}

	c.JSON(http.StatusOK, &models.Response{})
}

func (h *MovieHandlers) ListRatedMovies(c *gin.Context) {
	params := &models.PaginationParams{OrderBy: "create_date DESC"}
	err := c.ShouldBindQuery(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: invalidPaginationQueryParams})
		return
	}

	account := utils.GetAccount(c)

	ratings, err := h.storage.ListRatedMovies(account.ID)
	if err != nil {
		handlePostgresError(c, h.logger, err, ratingResource)
		return
	}

	c.JSON(http.StatusOK, ratings)
}

func (h *MovieHandlers) GetRating(c *gin.Context) {
	movieId, err := strconv.Atoi(c.Param(movieIdKey))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}

	account := utils.GetAccount(c)

	rating := &models.Rating{
		UserId:  account.ID,
		MovieId: movieId,
		Rating: intPointer(0),
	}
	err = h.storage.GetRating(rating)
	if err != nil && err != pg.ErrNoRows {
		handlePostgresError(c, h.logger, err, ratingResource)
		return
	}

	c.JSON(http.StatusOK, rating)
}
