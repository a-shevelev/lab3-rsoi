package service

import (
	"context"
	"errors"
	"fmt"
	"reservation-system/internal/dto"
	"reservation-system/internal/models"
	"reservation-system/internal/repo"
	"time"

	"github.com/google/uuid"
)

const (
	tillTimeError = "tillDate must be greater than startDate"
)

type ReservationServiceIFace interface {
	CreateReservation(ctx context.Context, req dto.CreateReservationRequest, username string) (*models.Reservation, error)
	GetReservation(ctx context.Context, uid string) (*models.Reservation, error)
	GetReservations(ctx context.Context, username string) ([]models.Reservation, error)
	GetCurrentAmount(ctx context.Context, username string) (uint64, error)
	UpdateStatus(ctx context.Context, reservationUID uuid.UUID, status string) error
	DeleteReservation(ctx context.Context, reservationUID string) error
}

type reservationService struct {
	repo repo.ReservationRepo
}

func NewReservationService(r repo.ReservationRepo) ReservationServiceIFace {
	return &reservationService{repo: r}
}

func (r *reservationService) CreateReservation(ctx context.Context, req dto.CreateReservationRequest, username string) (*models.Reservation, error) {
	bookUID, err := uuid.Parse(req.BookUID)
	if err != nil {
		return nil, err
	}
	libraryUID, err := uuid.Parse(req.LibraryUID)
	if err != nil {
		return nil, err
	}
	tillDate, err := time.Parse("2006-01-02", req.TillDate)
	if err != nil {
		return nil, err
	}
	startDate := time.Now().UTC().Truncate(24 * time.Hour)

	if tillDate.Equal(startDate) || tillDate.Before(startDate) {
		return nil, fmt.Errorf(tillTimeError)
	}
	res := models.Reservation{
		ReservationUID: uuid.New(),
		Username:       username,
		BookUID:        bookUID,
		LibraryUID:     libraryUID,
		Status:         "RENTED",
		StartDate:      startDate,
		TillDate:       tillDate,
	}
	return r.repo.CreateReservation(ctx, res)
}

func (r *reservationService) GetReservation(ctx context.Context, uid string) (*models.Reservation, error) {
	return r.repo.GetReservationByUID(ctx, uid)
}

func (r *reservationService) GetReservations(ctx context.Context, username string) ([]models.Reservation, error) {
	return r.repo.GetReservations(ctx, username)
}

func (r *reservationService) GetCurrentAmount(ctx context.Context, username string) (uint64, error) {
	return r.repo.GetCurrentReservationsAmount(ctx, username)
}

func (r *reservationService) UpdateStatus(ctx context.Context, reservationUID uuid.UUID, date string) error {
	res, err := r.repo.GetReservationByUID(ctx, reservationUID.String())
	if err != nil {
		return err
	}
	if res.Status != "RENTED" {
		return errors.New("book has already been returned")
	}
	returnDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return err
	}
	status := "RETURNED"
	fmt.Println(returnDate)
	if returnDate.After(res.TillDate) {
		status = "EXPIRED"
	}
	return r.repo.UpdateReservationStatus(ctx, reservationUID, status)
}

func (r *reservationService) DeleteReservation(ctx context.Context, reservationUID string) error {
	return r.repo.Delete(ctx, reservationUID)
}
