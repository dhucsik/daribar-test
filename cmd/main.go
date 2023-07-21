package main

import (
	"log"
	"net/http"

	"github.com/dhucsik/daribar-test/config"
	"github.com/dhucsik/daribar-test/handler"
	"github.com/dhucsik/daribar-test/middleware"
	"github.com/dhucsik/daribar-test/storage"
)

func main() {
	log.Fatalln(run())
}

func run() error {
	conf := config.NewConfig()

	log.Println("Config - ", conf)

	fanout, err := storage.ListenForDataUpdates(conf)
	if err != nil {
		return err
	}

	hdler := handler.NewHandler(fanout, new(storage.Storage), new(middleware.JWTMiddleware))

	http.Handle("/event", hdler)
	return http.ListenAndServe(conf.Port, nil)
}
