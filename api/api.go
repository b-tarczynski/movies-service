package api

import (
	"log"
	"net/http"

	"github.com/BarTar213/movies-service/config"
	"github.com/BarTar213/movies-service/handlers"
	"github.com/BarTar213/movies-service/storage"
	"github.com/gin-gonic/gin"
)

type Api struct {
	Port   string
	Router *gin.Engine
	Config *config.Config
	Logger *log.Logger
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

	h := handlers.NewMovieHandlers(a.Storage, a.Logger)

	a.Router.GET("/", a.health)
	a.Router.GET("/movies/:id", h.GetMovie)

	return a
}

func (a *Api) Run() error {
	return a.Router.Run(a.Config.Api.Port)
}

func (a *Api) health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "healthy")
}
