package handlers

import (
	"gateway-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rating)
}
