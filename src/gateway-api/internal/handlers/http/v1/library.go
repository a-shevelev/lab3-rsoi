package handlers

import (
	"gateway-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LibraryHandler struct {
	Service *service.LibraryService
}

func NewLibraryHandler(s *service.LibraryService) *LibraryHandler {
	return &LibraryHandler{Service: s}
}

func (h *LibraryHandler) RegisterRoutes(rg *gin.RouterGroup) {
	routes := rg.Group("/libraries")
	routes.GET("/", h.GetLibraries)
	routes.GET("/:uid/books", h.GetLibraryBooks)

}

func (h *LibraryHandler) GetLibraries(c *gin.Context) {
	city := c.Query("city")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	res, err := h.Service.GetLibraries(city, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *LibraryHandler) GetLibraryBooks(c *gin.Context) {
	libraryUid := c.Param("uid")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	showAll := c.DefaultQuery("showAll", "false") == "true"

	res, err := h.Service.GetLibraryBooks(libraryUid, page, size, showAll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
