package handlers

import (
	"net/http"
	"reservation-system/internal/dto"
	"reservation-system/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReservationHandler struct {
	service service.ReservationServiceIFace
}

func New(service service.ReservationServiceIFace) *ReservationHandler {
	return &ReservationHandler{service: service}
}

func (h *ReservationHandler) RegisterRoutes(rg *gin.RouterGroup) {
	resRoutes := rg.Group("/reservation")
	{
		resRoutes.POST("/", h.CreateReservation)
		resRoutes.GET("/:uid", h.GetReservation)
		resRoutes.GET("/", h.GetReservations)
		resRoutes.PUT("/:uid", h.UpdateStatus)
		resRoutes.GET("/amount", h.GetCurrentAmount)
		resRoutes.DELETE("/:uid", h.DeleteReservation)

	}
}

type GetUIDRequest struct {
	UID string `uri:"uid" binding:"required"`
}

func (h *ReservationHandler) GetReservation(c *gin.Context) {
	var GetUIDRequest GetUIDRequest
	if err := c.ShouldBindUri(&GetUIDRequest); err != nil {
	}

	res, err := h.service.GetReservation(c, GetUIDRequest.UID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ToReservationDTO(res))
}

func (h *ReservationHandler) CreateReservation(c *gin.Context) {
	username := c.GetHeader("X-User-Name")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Name header is required"})
		return
	}
	var req dto.CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.CreateReservation(c, req, username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.ToReservationDTO(res))
}

func (h *ReservationHandler) GetReservations(c *gin.Context) {
	username := c.GetHeader("X-User-Name")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Name header is required"})
		return
	}

	res, err := h.service.GetReservations(c, username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ToReservationsDTO(res))
}

func (h *ReservationHandler) GetCurrentAmount(c *gin.Context) {
	username := c.GetHeader("X-User-Name")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Name header is required"})
		return
	}

	amount, err := h.service.GetCurrentAmount(c, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"amount": amount})
}

func (h *ReservationHandler) UpdateStatus(c *gin.Context) {
	var uri GetUIDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		Date string `json:"date" binding:"required,datetime=2006-01-02"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid, err := uuid.Parse(uri.UID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UID"})
		return
	}

	if err := h.service.UpdateStatus(c, uid, req.Date); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ReservationHandler) DeleteReservation(c *gin.Context) {
	var GetUIDRequest GetUIDRequest
	if err := c.ShouldBindUri(&GetUIDRequest); err != nil {
	}

	err := h.service.DeleteReservation(c, GetUIDRequest.UID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
