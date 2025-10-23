package dto

import "github.com/google/uuid"

type GetLibrariesRequest struct {
	City string `form:"city" binding:"required"`
	Page int    `form:"page"`
	Size int    `form:"size"`
}

type GetBooksRequest struct {
	LibraryUID string `uri:"uid" binding:"required"`
	Page       int    `form:"page"`
	Size       int    `form:"size"`
	ShowAll    bool   `form:"showAll"`
}

type LibraryResponse struct {
	ID         uint64    `json:"id,omitempty"`
	LibraryUID uuid.UUID `json:"libraryUid"`
	Name       string    `json:"name"`
	City       string    `json:"city"`
	Address    string    `json:"address"`
}

type LibraryPaginationResponse struct {
	Page          int               `json:"page"`
	PageSize      int               `json:"pageSize"`
	TotalElements int               `json:"totalElements"`
	Items         []LibraryResponse `json:"items"`
}

type BookPaginationResponse struct {
	Page          int            `json:"page"`
	PageSize      int            `json:"pageSize"`
	TotalElements int            `json:"totalElements"`
	Items         []BookResponse `json:"items"`
}

type BookResponse struct {
	ID             uint64    `json:"id,omitempty"`
	BookUID        uuid.UUID `json:"bookUid"`
	Name           string    `json:"name"`
	Author         *string   `json:"author"`
	Genre          *string   `json:"genre"`
	Condition      string    `json:"condition"`
	AvailableCount int       `json:"availableCount"`
}
