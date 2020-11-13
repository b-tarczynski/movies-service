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
	notificator "github.com/BarTar213/notificator/client"
	"github.com/gin-gonic/gin"
)

func main() {
	conf := config.NewConfig("movies-service.yml")
	logger := log.New(os.Stdout, "", log.LstdFlags)

	logger.Printf("%+v\n", conf)

	if conf.Api.Release {
		gin.SetMode(gin.ReleaseMode)
	}

	logger.Println("Connecting to postgresql")
	postgres, err := storage.NewPostgres(&conf.Postgres)
	if err != nil {
		logger.Fatalln(err)
	}

	logger.Println("Connecting to TMDB client")
	tmdbClient := tmdb.NewClient(5*time.Second, conf)

	logger.Println("Connecting to notificator")
	notificatorCli := notificator.New(conf.Notificator.Address, conf.Api.Timeout)

	a := api.NewApi(
		api.WithConfig(conf),
		api.WithLogger(logger),
		api.WithStorage(postgres),
		api.WithTmdbClient(tmdbClient),
		api.WithNotificator(notificatorCli),
	)

	go a.Run()
	logger.Print("started app")

	shutDownSignal := make(chan os.Signal)
	signal.Notify(shutDownSignal, syscall.SIGINT, syscall.SIGTERM)

	<-shutDownSignal
	logger.Print("exited from app")
}
