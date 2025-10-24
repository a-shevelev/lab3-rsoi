package dto

type CreateReservationRequest struct {
	BookUID    string `json:"bookUid" binding:"required"`
	LibraryUID string `json:"libraryUid" binding:"required"`
	TillDate   string `json:"tillDate" binding:"required,datetime=2006-01-02"`
}

type ReturnReservationRequest struct {
	Date      string `json:"date" binding:"required,datetime=2006-01-02"`
	Condition string `json:"condition" binding:"required"`
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

type ReservationFullResponse struct {
	ReservationUID string          `json:"reservationUid"`
	Username       string          `json:"username"`
	Book           BookResponseRaw `json:"book"`
	Library        LibraryResponse `json:"library"`
	Status         string          `json:"status"`
	StartDate      string          `json:"startDate"`
	TillDate       string          `json:"tillDate"`
}

func ReservationToFull(
	r ReservationResponse,
	book BookResponseRaw,
	library LibraryResponse,
) ReservationFullResponse {
	return ReservationFullResponse{
		ReservationUID: r.ReservationUID,
		Username:       r.Username,
		Status:         r.Status,
		StartDate:      r.StartDate,
		TillDate:       r.TillDate,
		Book:           book,
		Library:        library,
	}
}

type ReturnRetryEvent struct {
	Username       string `json:"username,omitzero"`
	ReservationUID string `json:"reservation_uid,omitzero"`
	BookUID        string `json:"book_uid,omitzero"`
	LibraryUID     string `json:"library_uid,omitzero"`
	RateDelta      int    `json:"rate_delta,omitzero"`
	Condition      string `json:"condition,omitzero"`
	Date           string `json:"date,omitzero"`
}
