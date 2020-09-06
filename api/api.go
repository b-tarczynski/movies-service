package api

import (
	"log"

	"github.com/BarTar213/movies-service/config"
	"github.com/BarTar213/movies-service/middleware"
	"github.com/BarTar213/movies-service/storage"
	"github.com/gin-gonic/gin"
)

type Api struct {
	Port    string
	Router  *gin.Engine
	Config  *config.Config
	Logger  *log.Logger
	Storage storage.Storage
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

func NewApi(options ...func(api *Api)) *Api {
	a := &Api{
		Router: gin.Default(),
	}

	for _, option := range options {
		option(a)
	}

	mh := NewMovieHandlers(a.Storage, a.Logger)
	ch := NewCommentHandlers(a.Storage, a.Logger)

	a.Router.Use(gin.Recovery())

	movies := a.Router.Group("/movies")
	{
		movies.GET("", mh.ListMovies)
		movies.GET("/:commentId", mh.GetMovie)

		authorized := movies.Group("")
		authorized.Use(middleware.CheckAccount())
		{
			authorized.POST("/:commentId/like", mh.LikeMovie)
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
