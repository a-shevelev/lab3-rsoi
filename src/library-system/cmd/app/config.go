package main

import (
	"lab2-rsoi/library-system/internal/server"
	"lab2-rsoi/library-system/pkg/postgres"
)

type Config struct {
	Server server.Server   `envconfig:"SERVER"`
	DB     postgres.Config `envconfig:"DB"`
}
