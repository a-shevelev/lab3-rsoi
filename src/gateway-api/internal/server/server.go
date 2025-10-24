package server

import (
	"fmt"
	"gateway-api/internal/client"
	"gateway-api/internal/dto"
	handlers "gateway-api/internal/handlers/http/v1"
	"gateway-api/internal/rabbitmq"
	"gateway-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Server struct {
	Host              string `envconfig:"HOST"`
	Port              int    `envconfig:"PORT" required:"true"`
	LibraryClient     *client.Library
	RatingClient      *client.Rating
	ReservationClient *client.Reservation
	GinRouter         *gin.Engine
	RmqChannel        *amqp.Channel
	libQueue          string
	ratingQueue       string
	reservationQueue  string
}

func New(
	host string,
	port int,
	libSys *client.Library,
	rateSys *client.Rating,
	resSys *client.Reservation,
	rmqChannel *amqp.Channel,
) (*Server, error) {
	s := &Server{
		Host:              host,
		Port:              port,
		GinRouter:         gin.Default(),
		LibraryClient:     libSys,
		RatingClient:      rateSys,
		ReservationClient: resSys,
		RmqChannel:        rmqChannel,
		libQueue:          "lib-status-queue",
		ratingQueue:       "rate-status-queue",
		reservationQueue:  "reservation-status-queue",
	}

	if err := s.initRoutes(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) initRoutes() error {
	s.GinRouter.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "pong"})
	})

	s.GinRouter.GET("/manage/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	v1 := s.GinRouter.Group("/api/v1")

	libService := service.NewLibraryService(s.LibraryClient)
	libHandler := handlers.NewLibraryHandler(libService)
	libHandler.RegisterRoutes(v1)

	rateService := service.NewRatingService(s.RatingClient)
	rateHandler := handlers.NewRatingHandler(rateService)
	rateHandler.RegisterRoutes(v1)

	reservationService := service.NewReservationService(
		s.ReservationClient,
		s.LibraryClient,
		s.RatingClient,
		s.RmqChannel,
		s.libQueue,
		s.ratingQueue,
		s.reservationQueue,
	)
	reservationHandler := handlers.NewReservationHandler(reservationService)
	reservationHandler.RegisterRoutes(v1)

	queues := []string{
		s.reservationQueue,
		s.libQueue,
		s.ratingQueue,
	}

	for _, q := range queues {
		_, err := s.RmqChannel.QueueDeclare(
			q,
			true,  // durable
			false, // autoDelete
			false, // exclusive
			false, // noWait
			nil,   // args
		)
		if err != nil {
			log.Fatalf("failed to declare queue %s: %v", q, err)
		}
	}

	err := rabbitmq.RunRetryWorker(s.RmqChannel, s.reservationQueue, func(evt dto.ReturnRetryEvent) error {
		return s.ReservationClient.UpdateStatus(evt.ReservationUID, evt.Date)
	})
	if err != nil {
		return err
	}

	err = rabbitmq.RunRetryWorker(s.RmqChannel, s.libQueue, func(evt dto.ReturnRetryEvent) error {
		if err := s.LibraryClient.UpdateBookCount(evt.LibraryUID, evt.BookUID, +1); err != nil {
			return err
		}
		if evt.Condition != "" {
			return s.LibraryClient.UpdateBookCondition(evt.BookUID, evt.Condition)
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = rabbitmq.RunRetryWorker(s.RmqChannel, s.ratingQueue, func(evt dto.ReturnRetryEvent) error {
		return s.RatingClient.Update(evt.Username, evt.RateDelta)
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Run() error {
	return s.GinRouter.Run(fmt.Sprintf("%s:%d", s.Host, s.Port))
}
