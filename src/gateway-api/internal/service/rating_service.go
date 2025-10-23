package service

import (
	"gateway-api/internal/client"
	"gateway-api/internal/dto"
)

type RatingService struct {
	Client *client.Rating
}

func NewRatingService(client *client.Rating) *RatingService {
	return &RatingService{Client: client}
}

func (s *RatingService) GetRating(username string) (*dto.UserRatingResponse, error) {
	return s.Client.Get(username)
}
