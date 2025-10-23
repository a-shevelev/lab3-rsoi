package server

import (
	"fmt"
	"net/http"

	"rating-system/internal/handlers/http/v1"
	"rating-system/internal/repo"
	"rating-system/internal/service"
	"rating-system/pkg/postgres"

	"github.com/gin-gonic/gin"
	//log "github.com/sirupsen/logrus"
	//log "github.com/sirupsen/logrus"
)

type Server struct {
	Host      string `envconfig:"HOST"`
	Port      int    `envconfig:"PORT" required:"true"`
	DB        postgres.Client
	GinRouter *gin.Engine
}

func New(dbc postgres.Client, host string, port int) (*Server, error) {
	s := &Server{
		Host:      host,
		Port:      port,
		DB:        dbc,
		GinRouter: gin.Default(),
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

	rateRepo := repo.NewRatingRepo(s.DB)
	rateService := service.NewRatingService(rateRepo)
	rateHandler := handlers.New(rateService)
	rateHandler.RegisterRoutes(v1)

	return nil
}

func (s *Server) Run() error {
	return s.GinRouter.Run(fmt.Sprintf("%s:%d", s.Host, s.Port))
}
