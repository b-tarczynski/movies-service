package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BarTar213/movies-service/api"
	"github.com/BarTar213/movies-service/config"
	"github.com/BarTar213/movies-service/storage"
)

func main() {
	conf := config.NewConfig("movies-service.yml")
	logger := log.New(os.Stdout, "", log.LstdFlags)

	logger.Printf("%+v\n", conf)

	postgres, err := storage.NewPostgres(&conf.Postgres)
	if err != nil {
		logger.Fatalln(err)
	}

	a := api.NewApi(
		api.WithConfig(conf),
		api.WithLogger(logger),
		api.WithStorage(postgres),
	)

	go a.Run()
	logger.Print("started app")

	shutDownSignal := make(chan os.Signal)
	signal.Notify(shutDownSignal, syscall.SIGINT, syscall.SIGTERM)

	<-shutDownSignal
	logger.Print("exited from app")
}
