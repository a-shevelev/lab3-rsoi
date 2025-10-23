package handlers

import (
	"errors"
	"fmt"
	"lab2-rsoi/library-system/internal/dto"
	"lab2-rsoi/library-system/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LibraryHandler struct {
	service service.LibraryServiceIface
}

type GetBookLibRequest struct {
	UID string `uri:"uid" binding:"required"`
}

func New(service service.LibraryServiceIface) *LibraryHandler {
	return &LibraryHandler{service: service}
}
func (h *LibraryHandler) RegisterRoutes(rg *gin.RouterGroup) {
	libraryRoutes := rg.Group("libraries")
	{
		libraryRoutes.GET("", h.GetLibraries)
		libraryRoutes.GET("/:uid/books/", h.GetBooks)
		libraryRoutes.GET("/:uid/", h.GetLibraryByUid)
	}
	rg.GET("/books/:uid/", h.GetBookInfoByUid)
	rg.PUT("/books/:uid/condition", h.UpdateBookCondition)
	rg.PUT("/library/:libraryUid/books/:bookUid/count/:delta/", h.UpdateBookCount)
}

func (h *LibraryHandler) GetLibraries(c *gin.Context) {
	var req dto.GetLibrariesRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.City == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city is required"})
		return
	}

	page := req.Page
	size := req.Size
	fmt.Println(req.Page, req.Size)
	if page == 0 {
		page = 0
	}
	if size == 0 {
		size = 0
	}

	resp, err := h.service.ListLibraries(c, req.City, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *LibraryHandler) GetBooks(c *gin.Context) {
	var req dto.GetBooksRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	libraryUID, err := uuid.Parse(req.LibraryUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid libraryUid"})
		return
	}

	page := req.Page
	size := req.Size
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 10
	}

	resp, err := h.service.ListBooks(c, libraryUID, req.ShowAll, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *LibraryHandler) GetLibraryByUid(c *gin.Context) {
	var req GetBookLibRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	libraryUID, err := uuid.Parse(req.UID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid libraryUid"})
		return
	}

	resp, err := h.service.GetLibraryByUID(c, libraryUID)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *LibraryHandler) GetBookInfoByUid(c *gin.Context) {
	var req GetBookLibRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookUID, err := uuid.Parse(req.UID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid libraryUid"})
		return
	}

	resp, err := h.service.GetBookByUID(c, bookUID)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *LibraryHandler) UpdateBookCondition(c *gin.Context) {
	var req struct {
		Condition string `json:"condition" binding:"required"`
	}
	var uriReq struct {
		UID string `uri:"uid" binding:"required"`
	}

	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookUID, err := uuid.Parse(uriReq.UID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bookUid"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateBookCondition(c, bookUID, req.Condition); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "condition updated"})
}

func (h *LibraryHandler) UpdateBookCount(c *gin.Context) {
	var uriReq struct {
		BookUID    string `uri:"bookUid" binding:"required"`
		LibraryUID string `uri:"libraryUid" binding:"required"`
		Delta      int    `uri:"delta" binding:"required"`
	}

	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookUID, err := uuid.Parse(uriReq.BookUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bookUid"})
		return
	}
	libraryUID, err := uuid.Parse(uriReq.LibraryUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bookUid"})
		return
	}

	if err := h.service.UpdateBookCount(c, bookUID, libraryUID, uriReq.Delta); err != nil {
		if errors.Is(err, service.CountOfBooksIsZero) {
			c.JSON(http.StatusBadRequest, gin.H{"error": service.CountOfBooksIsZero})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "count updated"})
}
