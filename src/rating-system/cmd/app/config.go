package main

import (
	"rating-system/internal/server"
	"rating-system/pkg/postgres"
)

type Config struct {
	Server server.Server   `envconfig:"SERVER"`
	DB     postgres.Config `envconfig:"DB"`
}
