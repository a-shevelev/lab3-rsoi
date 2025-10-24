package service_test

import (
	"context"
	"reservation-system/internal/dto"
	"reservation-system/internal/models"
	"reservation-system/internal/service"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockReservationRepo struct {
	mock.Mock
}

func (m *MockReservationRepo) CreateReservation(ctx context.Context, r models.Reservation) (*models.Reservation, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(*models.Reservation), args.Error(1)
}

func (m *MockReservationRepo) GetReservationByUID(ctx context.Context, uid string) (*models.Reservation, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(*models.Reservation), args.Error(1)
}

func (m *MockReservationRepo) GetReservations(ctx context.Context, username string) ([]models.Reservation, error) {
	args := m.Called(ctx, username)
	return args.Get(0).([]models.Reservation), args.Error(1)
}

func (m *MockReservationRepo) GetCurrentReservationsAmount(ctx context.Context, username string) (uint64, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockReservationRepo) UpdateReservationStatus(ctx context.Context, uid uuid.UUID, status string) error {
	args := m.Called(ctx, uid, status)
	return args.Error(0)
}

func (m *MockReservationRepo) Delete(ctx context.Context, reservationUID string) error {
	args := m.Called(ctx, reservationUID)
	return args.Error(0)
}

func TestCreateReservation_Success(t *testing.T) {
	mockRepo := new(MockReservationRepo)
	svc := service.NewReservationService(mockRepo)

	req := dto.CreateReservationRequest{
		BookUID:    uuid.New().String(),
		LibraryUID: uuid.New().String(),
		TillDate:   time.Now().Add(24 * time.Hour).Format("2006-01-02"),
	}
	username := "user"

	expectedRes := &models.Reservation{
		ReservationUID: uuid.New(),
		Username:       username,
		BookUID:        uuid.MustParse(req.BookUID),
		LibraryUID:     uuid.MustParse(req.LibraryUID),
		Status:         "RENTED",
		StartDate:      time.Now().UTC().Truncate(24 * time.Hour),
		TillDate:       time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour),
	}

	mockRepo.On("CreateReservation", mock.Anything, mock.AnythingOfType("models.Reservation")).
		Return(expectedRes, nil)

	res, err := svc.CreateReservation(context.Background(), req, username)

	assert.NoError(t, err)
	assert.Equal(t, expectedRes.Username, res.Username)
	assert.Equal(t, "RENTED", res.Status)
	mockRepo.AssertExpectations(t)
}

func TestCreateReservation_InvalidTillDate(t *testing.T) {
	mockRepo := new(MockReservationRepo)
	svc := service.NewReservationService(mockRepo)

	req := dto.CreateReservationRequest{
		BookUID:    uuid.New().String(),
		LibraryUID: uuid.New().String(),
		TillDate:   time.Now().Add(-24 * time.Hour).Format("2006-01-02"), // раньше текущей даты
	}
	username := "user"

	res, err := svc.CreateReservation(context.Background(), req, username)

	assert.Nil(t, res)
	assert.ErrorContains(t, err, "tillDate must be greater than startDate")
}

func TestUpdateStatus_ReturnedAndExpired(t *testing.T) {
	mockRepo := new(MockReservationRepo)
	svc := service.NewReservationService(mockRepo)

	reservationUID := uuid.New()
	//startDate := time.Now().Add(-3 * 24 * time.Hour).UTC().Truncate(24 * time.Hour)
	tillDate := time.Now().Add(-1 * 24 * time.Hour).UTC().Truncate(24 * time.Hour)

	mockRepo.On("GetReservationByUID", mock.Anything, reservationUID.String()).
		Return(&models.Reservation{
			ReservationUID: reservationUID,
			Status:         "RENTED",
			TillDate:       tillDate,
		}, nil)

	mockRepo.On("UpdateReservationStatus", mock.Anything, reservationUID, "EXPIRED").Return(nil)

	err := svc.UpdateStatus(context.Background(), reservationUID, time.Now().Format("2006-01-02"))
	assert.NoError(t, err)

	mockRepo.AssertCalled(t, "UpdateReservationStatus", mock.Anything, reservationUID, "EXPIRED")
}

func TestUpdateStatus_AlreadyReturned(t *testing.T) {
	mockRepo := new(MockReservationRepo)
	svc := service.NewReservationService(mockRepo)

	reservationUID := uuid.New()
	mockRepo.On("GetReservationByUID", mock.Anything, reservationUID.String()).
		Return(&models.Reservation{
			ReservationUID: reservationUID,
			Status:         "RETURNED",
		}, nil)

	err := svc.UpdateStatus(context.Background(), reservationUID, time.Now().Format("2006-01-02"))
	assert.ErrorContains(t, err, "book has already been returned")
}

func TestGetReservations(t *testing.T) {
	mockRepo := new(MockReservationRepo)
	svc := service.NewReservationService(mockRepo)

	username := "user"
	mockRepo.On("GetReservations", mock.Anything, username).Return([]models.Reservation{
		{Username: username},
	}, nil)

	res, err := svc.GetReservations(context.Background(), username)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, username, res[0].Username)
}
