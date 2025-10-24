package main

import (
	"gateway-api/internal/client"
	"gateway-api/internal/server"
)

type Config struct {
	Server            server.Server      `envconfig:"SERVER"`
	ReservationSystem client.Reservation `envconfig:"RESERVATION_SYSTEM"`
	LibrarySystem     client.Library     `envconfig:"LIBRARY_SYSTEM"`
	RatingSystem      client.Rating      `envconfig:"RATING_SYSTEM"`
	RabbitMQ          string             `envconfig:"RABBITMQ"`
}
