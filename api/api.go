package api

import (
	"log"

	"github.com/BarTar213/movies-service/config"
	"github.com/BarTar213/movies-service/middleware"
	"github.com/BarTar213/movies-service/storage"
	"github.com/BarTar213/movies-service/tmdb"
	"github.com/gin-gonic/gin"
)

type Api struct {
	Port       string
	Router     *gin.Engine
	Config     *config.Config
	Logger     *log.Logger
	Storage    storage.Storage
	TmdbClient tmdb.Client
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

func NewApi(options ...func(api *Api)) *Api {
	a := &Api{
		Router: gin.Default(),
	}

	for _, option := range options {
		option(a)
	}

	mh := NewMovieHandlers(a.Storage, a.TmdbClient, a.Logger)
	ch := NewCommentHandlers(a.Storage, a.Logger)

	a.Router.Use(gin.Recovery())

	movies := a.Router.Group("/movies")
	{
		movies.GET("", mh.ListMovies)
		movies.GET("/:movieId", mh.GetMovie)
		movies.GET("/:movieId/credits", mh.GetCredits)

		authorized := movies.Group("")
		authorized.Use(middleware.CheckAccount())
		{
			authorized.POST("/:movieId/like", mh.LikeMovie)
		}
	}

	comments := a.Router.Group("/comments")
	{
		comments.GET("", ch.GetComments)

		authorized := comments.Group("")
		authorized.Use(middleware.CheckAccount())
		{
			authorized.POST("", ch.AddComment)
			authorized.POST("/:commId/like", ch.LikeComment)
			authorized.PUT("/:commId", ch.UpdateComment)
			authorized.DELETE("/:commId", ch.DeleteComment)
		}
	}

	return a
}

func (a *Api) Run() error {
	return a.Router.Run(a.Config.Api.Port)
}
