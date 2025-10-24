package handlers

import (
	"errors"
	"gateway-api/internal/dto"
	"gateway-api/internal/service"
	"gateway-api/pkg/ext"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReservationHandler struct {
	Service *service.ReservationService
}

func NewReservationHandler(service *service.ReservationService) *ReservationHandler {
	return &ReservationHandler{Service: service}
}

func (h *ReservationHandler) RegisterRoutes(rg *gin.RouterGroup) {
	routes := rg.Group("/reservations")
	routes.GET("/", h.GetReservations)
	routes.POST("/", h.CreateReservation)
	routes.POST("/:uid/return/", h.ReturnBook)
}

func (h *ReservationHandler) GetReservations(c *gin.Context) {
	username := c.GetHeader("X-User-Name")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Name header required"})
		return
	}

	reservations, err := h.Service.Get(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reservations)
}

func (h *ReservationHandler) CreateReservation(c *gin.Context) {
	username := c.GetHeader("X-User-Name")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Name header required"})
		return
	}
	var req dto.CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reservation, err := h.Service.CreateReservation(username, req)
	if err != nil {
		if errors.Is(err, ext.LibraryServiceUnavailableError) {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": ext.LibraryServiceUnavailableError.Error()})
			return
		}
		if errors.Is(err, ext.RatingServiceUnavailableError) {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": ext.RatingServiceUnavailableError.Error()})
			return
		}
		if errors.Is(err, ext.ReservationServiceUnavailableError) {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": ext.ReservationServiceUnavailableError.Error()})
			return
		}
		if errors.Is(err, ext.BookNotAvailableError) {
			c.JSON(http.StatusBadRequest, gin.H{"message": ext.BookNotAvailableError.Error()})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reservation)
}

func (h *ReservationHandler) ReturnBook(c *gin.Context) {
	username := c.GetHeader("X-User-Name")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Name header required"})
		return
	}
	var req dto.ReturnReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var reqURI struct {
		ReservationUID string `uri:"uid" binding:"required"`
	}

	if err := c.ShouldBindUri(&reqURI); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.ReturnBook(username, req, reqURI.ReservationUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
