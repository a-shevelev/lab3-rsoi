package handlers

import (
	"net/http"
	"rating-system/internal/dto"
	"rating-system/internal/service"

	"github.com/gin-gonic/gin"
)

type RatingHandler struct {
	service service.RatingServiceIFace
}

type GetBookLibRequest struct {
	UID string `uri:"uid" binding:"required"`
}

func New(service service.RatingServiceIFace) *RatingHandler {
	return &RatingHandler{service: service}
}
func (h *RatingHandler) RegisterRoutes(rg *gin.RouterGroup) {
	ratingRoutes := rg.Group("/rating")
	{
		ratingRoutes.GET("/", h.GetRatingHandler)
		ratingRoutes.PUT("/stars/:stars_diff", h.UpdateRatingHandler)
	}
}

// GET /api/v1/rating/
// Header: X-User-Name: {{username}}
func (h *RatingHandler) GetRatingHandler(c *gin.Context) {
	username := c.GetHeader("X-User-Name")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Name header is required"})
		return
	}

	resp, err := h.service.GetRating(c, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *RatingHandler) UpdateRatingHandler(c *gin.Context) {
	username := c.GetHeader("X-User-Name")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Name header is required"})
		return
	}

	var uriReq dto.UpdateRatingRequest

	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateRating(c, username, uriReq.StarsDiff); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "rating updated successfully"})
}
