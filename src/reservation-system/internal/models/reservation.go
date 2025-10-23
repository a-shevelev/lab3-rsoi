package models

import (
	"time"

	"github.com/google/uuid"
)

type Reservation struct {
	ID             int64     `db:"id"`
	ReservationUID uuid.UUID `db:"reservation_uid"`
	Username       string    `db:"username"`
	BookUID        uuid.UUID `db:"book_uid"`
	LibraryUID     uuid.UUID `db:"library_uid"`
	Status         string    `db:"status"  validate:"oneof=RENTED RETURNED EXPIRED"`
	StartDate      time.Time `db:"start_date"`
	TillDate       time.Time `db:"till_date"`
}
