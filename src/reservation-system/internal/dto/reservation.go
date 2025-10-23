package dto

import (
	"reservation-system/internal/models"
)

type CreateReservationRequest struct {
	BookUID    string `json:"bookUid" binding:"required"`
	LibraryUID string `json:"libraryUid" binding:"required"`
	TillDate   string `json:"tillDate" binding:"required,datetime=2006-01-02"`
}

type ReservationResponse struct {
	ReservationUID string `json:"reservationUid"`
	Username       string `json:"username"`
	BookUID        string `json:"bookUid"`
	LibraryUID     string `json:"libraryUid"`
	Status         string `json:"status"`
	StartDate      string `json:"startDate"`
	TillDate       string `json:"tillDate"`
}

type ReservationsListResponse struct {
	Items []ReservationResponse
}

func ToReservationDTO(m *models.Reservation) ReservationResponse {
	return ReservationResponse{
		ReservationUID: m.ReservationUID.String(),
		Username:       m.Username,
		BookUID:        m.BookUID.String(),
		LibraryUID:     m.LibraryUID.String(),
		Status:         m.Status,
		StartDate:      m.StartDate.Format("2006-01-02"),
		TillDate:       m.TillDate.Format("2006-01-02"),
	}
}

func ToReservationsDTO(list []models.Reservation) []ReservationResponse {
	out := make([]ReservationResponse, 0, len(list))
	for _, r := range list {
		out = append(out, ToReservationDTO(&r))
	}
	return out
}
