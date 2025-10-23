package main

import (
	"context"
	"fmt"
	"reservation-system/internal/server"
	"reservation-system/pkg/postgres"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		panic(err)
	}
	fmt.Println(cfg.DB)

	db, err := postgres.Connect(ctx, cfg.DB)
	if err != nil {
		log.WithError(err).Error("Failed to connect to database")
		panic(err)
	}
	if err := db.Ping(ctx); err != nil {
		log.WithError(err).Error("Failed to ping database")
	}
	log.Info("Successfully connected to database reservations")
	defer db.Close()

	srv, err := server.New(db, cfg.Server.Host, cfg.Server.Port)
	if err != nil {
		log.WithError(err).Error("failed to initialize server")
	}

	err = srv.Run()
	if err != nil {
		log.WithError(err).Error("server shutdown")
		return
	}

}
