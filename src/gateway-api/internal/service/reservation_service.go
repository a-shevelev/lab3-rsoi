package service

import (
	"errors"
	"fmt"
	"gateway-api/internal/client"
	"gateway-api/internal/dto"
	"gateway-api/pkg/ext"
)

type ReservationService struct {
	ClientRes  *client.Reservation
	ClientLib  *client.Library
	ClientRate *client.Rating
}

func NewReservationService(clRes *client.Reservation, clLib *client.Library, clRate *client.Rating) *ReservationService {
	return &ReservationService{ClientRes: clRes, ClientLib: clLib, ClientRate: clRate}
}

func (s *ReservationService) Get(username string) ([]dto.ReservationFullResponse, error) {
	raw, err := s.ClientRes.Get(username)
	if err != nil {
		return nil, err
	}

	result := make([]dto.ReservationFullResponse, 0, len(raw))
	for _, r := range raw {
		book, err := s.ClientLib.GetBookByUID(r.BookUID)
		if err != nil {
			return nil, err
		}

		lib, err := s.ClientLib.GetLibraryByUID(r.LibraryUID)
		if err != nil {
			return nil, err
		}

		fullRes := dto.ReservationToFull(r, dto.BookToRaw(*book), *lib)
		result = append(result, fullRes)
	}
	return result, nil
}

func (s *ReservationService) CreateReservation(username string, req dto.CreateReservationRequest) (*dto.ReservationFullResponse, error) {

	resCount, err := s.ClientRes.GetCurrentAmount(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get current amount: %s", err)
	}

	starsCount, err := s.ClientRate.Get(username)
	if err != nil {
		if errors.Is(err, ext.ServiceUnavailableError) {
			return nil, ext.RatingServiceUnavailableError
		}
		return nil, fmt.Errorf("failed to get rating: %s", err)
	}
	books, err := s.ClientLib.GetBookByUID(req.BookUID)
	if err != nil {
		if errors.Is(err, ext.ServiceUnavailableError) {
			return nil, ext.LibraryServiceUnavailableError
		}
		return nil, fmt.Errorf("failed to get book by uid: %s", err)
	}
	if books.AvailableCount < 0 {
		return nil, ext.BookNotAvailableError
	}
	if resCount >= starsCount.Stars {
		return nil, fmt.Errorf("You rented maximum amount of books", resCount)
	}
	result, err := s.ClientRes.Create(username, req)
	if err != nil {
		if errors.Is(err, ext.ServiceUnavailableError) {
			return nil, ext.ReservationServiceUnavailableError
		}
		return nil, fmt.Errorf("failed to create reservation: %s", err)
	}
	err = s.ClientLib.UpdateBookCount(result.LibraryUID, result.BookUID, -1)
	if err != nil {

		return nil, fmt.Errorf("failed to update book count: %s", err)
	}

	book, err := s.ClientLib.GetBookByUID(result.BookUID)
	if err != nil {
		err := s.ClientRes.DeleteReservation(result.ReservationUID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete book: %s", err)
		}
		return nil, err
	}
	lib, err := s.ClientLib.GetLibraryByUID(result.LibraryUID)
	if err != nil {
		return nil, err
	}
	fullRes := dto.ReservationToFull(*result, dto.BookToRaw(*book), *lib)

	return &fullRes, nil
}

func (s *ReservationService) ReturnBook(
	username string,
	req dto.ReturnReservationRequest,
	reservationUID string) error {
	rate := 1
	err := s.ClientRes.UpdateStatus(reservationUID, req.Date)
	if err != nil {
		return fmt.Errorf("failed to update status: %s", err)
	}
	res, err := s.ClientRes.GetByUID(reservationUID)
	if err != nil {
		return fmt.Errorf("failed to get reservation by uid: %s", err)
	}
	if res.Status == "EXPIRED" {
		rate = -10
	}
	book, err := s.ClientLib.GetBookByUID(res.BookUID)
	if err != nil {
		return fmt.Errorf("failed to get book by uid: %s", err)
	}
	if book.Condition != req.Condition {
		rate = -10
		err = s.ClientLib.UpdateBookCondition(res.BookUID, req.Condition)
		if err != nil {
			return fmt.Errorf("failed to update book condition: %s", err)
		}
	}

	err = s.ClientLib.UpdateBookCount(res.LibraryUID, res.BookUID, 1)
	if err != nil {
		return fmt.Errorf("failed to update book count: %s", err)
	}
	err = s.ClientRate.Update(username, rate)
	if err != nil {
		return fmt.Errorf("failed to update rate: %s", err)
	}
	return nil
}
