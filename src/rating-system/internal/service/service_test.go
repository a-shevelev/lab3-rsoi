package service_test

import (
	"context"
	"errors"
	"rating-system/internal/models"
	"rating-system/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRatingRepo struct {
	mock.Mock
}

func (m *MockRatingRepo) GetRatingRepo(ctx context.Context, username string) (*models.Rating, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Rating), args.Error(1)
}

func (m *MockRatingRepo) UpdateRatingRepo(ctx context.Context, username string, stars int) error {
	args := m.Called(ctx, username, stars)
	return args.Error(0)
}

// --- Тесты ---

func TestGetRating_Success(t *testing.T) {
	mockRepo := new(MockRatingRepo)
	svc := service.NewRatingService(mockRepo)

	username := "user1"
	mockRepo.On("GetRatingRepo", mock.Anything, username).Return(&models.Rating{Stars: 42}, nil)

	resp, err := svc.GetRating(context.Background(), username)
	assert.NoError(t, err)
	assert.Equal(t, 42, resp.Stars)
	mockRepo.AssertExpectations(t)
}

func TestGetRating_Error(t *testing.T) {
	mockRepo := new(MockRatingRepo)
	svc := service.NewRatingService(mockRepo)

	username := "user1"
	mockRepo.On("GetRatingRepo", mock.Anything, username).Return(nil, errors.New("db error"))

	resp, err := svc.GetRating(context.Background(), username)
	assert.Nil(t, resp)
	assert.ErrorContains(t, err, "db error")
	mockRepo.AssertExpectations(t)
}

func TestUpdateRating_Success(t *testing.T) {
	mockRepo := new(MockRatingRepo)
	svc := service.NewRatingService(mockRepo)

	username := "user1"
	currentStars := 50
	delta := 10

	mockRepo.On("GetRatingRepo", mock.Anything, username).Return(&models.Rating{Stars: currentStars}, nil)
	mockRepo.On("UpdateRatingRepo", mock.Anything, username, currentStars+delta).Return(nil)

	err := svc.UpdateRating(context.Background(), username, delta)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateRating_ClampToZero(t *testing.T) {
	mockRepo := new(MockRatingRepo)
	svc := service.NewRatingService(mockRepo)

	username := "user1"
	currentStars := 5
	delta := -10

	mockRepo.On("GetRatingRepo", mock.Anything, username).Return(&models.Rating{Stars: currentStars}, nil)
	mockRepo.On("UpdateRatingRepo", mock.Anything, username, 0).Return(nil)

	err := svc.UpdateRating(context.Background(), username, delta)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateRating_ClampToMax(t *testing.T) {
	mockRepo := new(MockRatingRepo)
	svc := service.NewRatingService(mockRepo)

	username := "user1"
	currentStars := 95
	delta := 10

	mockRepo.On("GetRatingRepo", mock.Anything, username).Return(&models.Rating{Stars: currentStars}, nil)
	mockRepo.On("UpdateRatingRepo", mock.Anything, username, 100).Return(nil)

	err := svc.UpdateRating(context.Background(), username, delta)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateRating_GetError(t *testing.T) {
	mockRepo := new(MockRatingRepo)
	svc := service.NewRatingService(mockRepo)

	username := "user1"
	mockRepo.On("GetRatingRepo", mock.Anything, username).Return(nil, errors.New("db error"))

	err := svc.UpdateRating(context.Background(), username, 10)
	assert.ErrorContains(t, err, "failed to get current rating")
	mockRepo.AssertExpectations(t)
}

func TestUpdateRating_UpdateError(t *testing.T) {
	mockRepo := new(MockRatingRepo)
	svc := service.NewRatingService(mockRepo)

	username := "user1"
	mockRepo.On("GetRatingRepo", mock.Anything, username).Return(&models.Rating{Stars: 50}, nil)
	mockRepo.On("UpdateRatingRepo", mock.Anything, username, 60).Return(errors.New("update failed"))

	err := svc.UpdateRating(context.Background(), username, 10)
	assert.ErrorContains(t, err, "failed to update rating")
	mockRepo.AssertExpectations(t)
}
