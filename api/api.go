package api

import (
	"log"

	"github.com/BarTar213/movies-service/config"
	"github.com/BarTar213/movies-service/middleware"
	"github.com/BarTar213/movies-service/storage"
	"github.com/BarTar213/movies-service/tmdb"
	notificator "github.com/BarTar213/notificator/client"
	"github.com/gin-gonic/gin"
)

type Api struct {
	Port        string
	Router      *gin.Engine
	Config      *config.Config
	Storage     storage.Storage
	TmdbClient  tmdb.Client
	Notificator notificator.Client
	Logger      *log.Logger
}

func WithConfig(conf *config.Config) func(a *Api) {
	return func(a *Api) {
		a.Config = conf
	}
}

func WithLogger(logger *log.Logger) func(a *Api) {
	return func(a *Api) {
		a.Logger = logger
	}
}

func WithStorage(storage storage.Storage) func(a *Api) {
	return func(a *Api) {
		a.Storage = storage
	}
}

func WithTmdbClient(tmdb tmdb.Client) func(a *Api) {
	return func(a *Api) {
		a.TmdbClient = tmdb
	}
}

func WithNotificator(notificator notificator.Client) func(a *Api) {
	return func(a *Api) {
		a.Notificator = notificator
	}
}

func NewApi(options ...func(api *Api)) *Api {
	a := &Api{
		Router: gin.Default(),
	}

	for _, option := range options {
		option(a)
	}

	moviesHndl := NewMovieHandlers(a.Storage, a.TmdbClient, a.Logger)
	commentsHndl := NewCommentHandlers(a.Storage, a.Logger)

	a.Router.Use(gin.Recovery())

	standard := a.Router.Group("")
	{
		movies := standard.Group("/movies")
		{
			movies.GET("", moviesHndl.ListMovies)
			movies.GET("/:movieId", moviesHndl.GetMovie)
			movies.GET("/:movieId/credits", moviesHndl.GetCredits)
		}

		comments := standard.Group("/comments")
		{
			comments.GET("", commentsHndl.ListComments)
		}
	}

	authorized := a.Router.Group("")
	authorized.Use(middleware.CheckAccount())
	{
		movies := authorized.Group("/movies")
		{
			movies.POST("/:movieId/like", moviesHndl.LikeMovie)

			movies.GET("/:movieId/rating", moviesHndl.GetRating)
			movies.POST("/:movieId/rating", moviesHndl.RateMovie)
			movies.DELETE("/:movieId/rating", moviesHndl.DeleteRating)
		}

		favourites := authorized.Group("/favourites")
		{
			favourites.GET("", moviesHndl.ListLikedMovies)
			favourites.GET("/:movieId", moviesHndl.CheckLiked)
			favourites.GET("/:movieId/comments", commentsHndl.ListLikedComments)
		}

		comments := authorized.Group("/comments")
		{
			comments.POST("", commentsHndl.AddComment)
			comments.POST("/:commId/like", commentsHndl.LikeComment)
			comments.PUT("/:commId", commentsHndl.UpdateComment)
			comments.DELETE("/:commId", commentsHndl.DeleteComment)
		}

		authorized.GET("/rating", moviesHndl.ListRatedMovies)
	}

	return a
}

func (a *Api) Run() error {
	return a.Router.Run(a.Config.Api.Port)
}
