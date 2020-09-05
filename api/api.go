package api

import (
	"log"

	"github.com/BarTar213/movies-service/config"
	"github.com/BarTar213/movies-service/handlers"
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

	mh := handlers.NewMovieHandlers(a.Storage, a.Logger)
	ch := handlers.NewCommentHandlers(a.Storage, a.Logger)

	movies := a.Router.Group("/movies")
	{
		movies.GET("", mh.ListMovies)
		movies.GET("/:movieId", mh.GetMovie)
		movies.POST("/:movieId/like", mh.LikeMovie)
	}

	comments := a.Router.Group("/comments")
	{
		comments.GET("", ch.GetComments)
		comments.POST("", ch.AddComment)
		comments.PUT("/:commId", ch.UpdateComment)
		comments.DELETE("/:commId", ch.DeleteComment)
		comments.POST("/:commId/like", ch.LikeComment)
	}

	return a
}

func (a *Api) Run() error {
	return a.Router.Run(a.Config.Api.Port)
}
