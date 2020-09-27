package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BarTar213/movies-service/api"
	"github.com/BarTar213/movies-service/config"
	"github.com/BarTar213/movies-service/storage"
	"github.com/BarTar213/movies-service/tmdb"
	"github.com/gin-gonic/gin"
)

func main() {
	conf := config.NewConfig("movies-service.yml")
	logger := log.New(os.Stdout, "", log.LstdFlags)

	logger.Printf("%+v\n", conf)

	if conf.Api.Release {
		gin.SetMode(gin.ReleaseMode)
	}

	postgres, err := storage.NewPostgres(&conf.Postgres)
	if err != nil {
		logger.Fatalln(err)
	}
	tmdbClient := tmdb.NewClient(5*time.Second, conf)

	a := api.NewApi(
		api.WithConfig(conf),
		api.WithLogger(logger),
		api.WithStorage(postgres),
		api.WithTmdbClient(tmdbClient),
	)

	go a.Run()
	logger.Print("started app")

	shutDownSignal := make(chan os.Signal)
	signal.Notify(shutDownSignal, syscall.SIGINT, syscall.SIGTERM)

	<-shutDownSignal
	logger.Print("exited from app")
}
