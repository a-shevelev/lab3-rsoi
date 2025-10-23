package server

import (
	"fmt"
	"gateway-api/internal/client"
	handlers "gateway-api/internal/handlers/http/v1"
	"gateway-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Host              string `envconfig:"HOST"`
	Port              int    `envconfig:"PORT" required:"true"`
	LibraryClient     *client.Library
	RatingClient      *client.Rating
	ReservationClient *client.Reservation
	GinRouter         *gin.Engine
}

func New(host string, port int, libSys *client.Library, rateSys *client.Rating, resSys *client.Reservation) (*Server, error) {
	s := &Server{
		Host:              host,
		Port:              port,
		GinRouter:         gin.Default(),
		LibraryClient:     libSys,
		RatingClient:      rateSys,
		ReservationClient: resSys,
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

	reservationService := service.NewReservationService(s.ReservationClient, s.LibraryClient, s.RatingClient)
	reservationHandler := handlers.NewReservationHandler(reservationService)
	reservationHandler.RegisterRoutes(v1)

	return nil
}

func (s *Server) Run() error {
	return s.GinRouter.Run(fmt.Sprintf("%s:%d", s.Host, s.Port))
}
