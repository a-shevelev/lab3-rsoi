package server

import (
	"fmt"
	handlers "lab2-rsoi/library-system/internal/handlers/http/v1"
	"lab2-rsoi/library-system/internal/repo"
	"lab2-rsoi/library-system/internal/service"
	"lab2-rsoi/library-system/pkg/postgres"
	"net/http"

	"github.com/gin-gonic/gin"
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
		c.JSON(200, gin.H{"msg": "pong"})
	})

	s.GinRouter.GET("/manage/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	//if err := s.InitDocsRoutes(); err != nil {
	//	log.Info("Docs routes initialization failed")
	//}

	v1 := s.GinRouter.Group("/api/v1")

	library := repo.NewLibraryRepo(s.DB)
	libraryService := service.NewLibraryService(library)
	libraryHandler := handlers.New(libraryService)
	libraryHandler.RegisterRoutes(v1)

	return nil
}

func (s *Server) Run() error {
	return s.GinRouter.Run(fmt.Sprintf("%s:%d", s.Host, s.Port))
}
