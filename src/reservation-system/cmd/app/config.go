package main

import (
	"reservation-system/internal/server"
	"reservation-system/pkg/postgres"
)

type Config struct {
	Server server.Server   `envconfig:"SERVER"`
	DB     postgres.Config `envconfig:"DB"`
}
