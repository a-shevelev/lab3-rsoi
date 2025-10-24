package handlers

import (
	"errors"
	"gateway-api/internal/service"
	"gateway-api/pkg/ext"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	RatingServiceUnavailable = errors.New("Bonus Service unavailable")
)

type RatingHandler struct {
	Service *service.RatingService
}

func NewRatingHandler(svc *service.RatingService) *RatingHandler {
	return &RatingHandler{Service: svc}
}

func (h *RatingHandler) RegisterRoutes(rg *gin.RouterGroup) {
	routes := rg.Group("/rating")
	routes.GET("/", h.GetRating)
}

func (h *RatingHandler) GetRating(c *gin.Context) {
	username := c.GetHeader("X-User-Name")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Name header required"})
		return
	}

	rating, err := h.Service.GetRating(username)
	if err != nil {
		if errors.Is(err, ext.ServiceUnavailableError) {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": RatingServiceUnavailable.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rating)
}
