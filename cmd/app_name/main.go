package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BarTar213/go-template/api"
	"github.com/BarTar213/go-template/config"
)

func main() {
	conf := config.NewConfig("app_name.yml")
	logger := log.New(os.Stdout, "", log.LstdFlags)

	a := api.NewApi(
		api.WithConfig(conf),
		api.WithLogger(logger),
	)

	go a.Run()
	logger.Print("started app")

	shutDownSignal := make(chan os.Signal)
	signal.Notify(shutDownSignal, syscall.SIGINT, syscall.SIGTERM)

	<-shutDownSignal
	logger.Print("exited from app")
}
