package main

import (
	"context"
	"gateway-api/internal/client"
	"gateway-api/internal/server"

	log "github.com/sirupsen/logrus"

	"github.com/kelseyhightower/envconfig"
)

func main() {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		panic(err)
	}

	clientLibrary := client.NewLibrary(cfg.LibrarySystem.BaseURL)
	clientRating := client.NewRating(cfg.RatingSystem.BaseURL)
	clientReservation := client.NewReservation(cfg.ReservationSystem.BaseURL)

	srv, err := server.New(cfg.Server.Host, cfg.Server.Port,
		clientLibrary, clientRating, clientReservation)
	if err != nil {
		log.WithError(err).Error("failed to initialize server")
	}

	err = srv.Run()
	if err != nil {
		log.WithError(err).Error("server shutdown")
		return
	}

}
