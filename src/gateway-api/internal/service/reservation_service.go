package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"gateway-api/internal/client"
	"gateway-api/internal/dto"
	"gateway-api/pkg/ext"

	"github.com/streadway/amqp"
)

type ReservationService struct {
	ClientRes        *client.Reservation
	ClientLib        *client.Library
	ClientRate       *client.Rating
	rmqChannel       *amqp.Channel
	libQueue         string
	ratingQueue      string
	reservationQueue string
}

func mapUnavailable(err error, target error) error {
	if errors.Is(err, ext.ServiceUnavailableError) {
		return target
	}
	return err
}

func NewReservationService(
	clRes *client.Reservation,
	clLib *client.Library,
	clRate *client.Rating,
	rmqChannel *amqp.Channel,
	libQ string,
	ratingQ string,
	reservationQ string,
) *ReservationService {
	return &ReservationService{
		ClientRes:        clRes,
		ClientLib:        clLib,
		ClientRate:       clRate,
		rmqChannel:       rmqChannel,
		libQueue:         libQ,
		ratingQueue:      ratingQ,
		reservationQueue: reservationQ,
	}
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

func (s *ReservationService) ReturnBook_(
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

func (s *ReservationService) ReturnBook(username string, req dto.ReturnReservationRequest, reservationUID string) error {
	rate := 1

	if err := s.ClientRes.UpdateStatus(reservationUID, req.Date); err != nil {
		if errors.Is(err, ext.ServiceUnavailableError) {
			evt := dto.ReturnRetryEvent{
				Username:       username,
				ReservationUID: reservationUID,
				Date:           req.Date,
				RateDelta:      rate,
			}
			s.enqueueReturn(evt, s.reservationQueue)
			return nil // пользователю success
		}
		return fmt.Errorf("failed to update reservation status: %w", err)
	}

	res, err := s.ClientRes.GetByUID(reservationUID)
	if err != nil {
		if errors.Is(err, ext.ServiceUnavailableError) {
			evt := dto.ReturnRetryEvent{
				Username:       username,
				ReservationUID: reservationUID,
				Date:           req.Date,
				RateDelta:      rate,
			}
			s.enqueueReturn(evt, s.reservationQueue)
			return nil
		}
		return fmt.Errorf("failed to get reservation by uid: %w", err)
	}

	// Если истек срок — штраф
	if res.Status == "EXPIRED" {
		rate = -10
	}

	if err := s.ClientLib.UpdateBookCount(res.LibraryUID, res.BookUID, +1); err != nil {
		if errors.Is(err, ext.ServiceUnavailableError) {
			evt := dto.ReturnRetryEvent{
				Username:       username,
				ReservationUID: reservationUID,
				BookUID:        res.BookUID,
				LibraryUID:     res.LibraryUID,
				RateDelta:      rate,
				Condition:      req.Condition,
			}
			s.enqueueReturn(evt, s.libQueue)
			return nil
		}
		return fmt.Errorf("failed to update book count: %w", err)
	}

	book, err := s.ClientLib.GetBookByUID(res.BookUID)
	if err != nil {
		return fmt.Errorf("failed to get book by uid: %s", err)
	}

	if req.Condition != "" && req.Condition != book.Condition {
		if err := s.ClientLib.UpdateBookCondition(res.BookUID, req.Condition); err != nil {
			if errors.Is(err, ext.ServiceUnavailableError) {
				evt := dto.ReturnRetryEvent{
					Username:       username,
					ReservationUID: reservationUID,
					BookUID:        res.BookUID,
					LibraryUID:     res.LibraryUID,
					RateDelta:      rate,
					Condition:      req.Condition,
				}
				s.enqueueReturn(evt, s.libQueue)
				return nil
			}
			return fmt.Errorf("failed to update book condition: %w", err)
		}
	}

	if err := s.ClientRate.Update(username, rate); err != nil {
		if errors.Is(err, ext.ServiceUnavailableError) {
			evt := dto.ReturnRetryEvent{
				Username:       username,
				ReservationUID: reservationUID,
				BookUID:        res.BookUID,
				LibraryUID:     res.LibraryUID,
				RateDelta:      rate,
			}
			s.enqueueReturn(evt, s.ratingQueue)
			return nil
		}
		return fmt.Errorf("failed to update user rating: %w", err)
	}

	return nil
}

func (s *ReservationService) enqueueReturn(evt dto.ReturnRetryEvent, queue string) {
	body, _ := json.Marshal(evt)
	_ = s.rmqChannel.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
