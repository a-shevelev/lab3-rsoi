package ext

import "errors"

var (
	ServiceUnavailableError            = errors.New("service unavailable")
	RatingServiceUnavailableError      = errors.New("Bonus Service unavailable")
	LibraryServiceUnavailableError     = errors.New("Library Service unavailable")
	ReservationServiceUnavailableError = errors.New("Reservation Service unavailable")
	BookNotAvailableError              = errors.New("Book not available")
)
