package service

import (
	"context"
	"fmt"
	"rating-system/internal/dto"
	"rating-system/internal/repo"
)

type RatingServiceIFace interface {
	GetRating(ctx context.Context, username string) (*dto.RatingResponse, error)
	UpdateRating(ctx context.Context, username string, delta int) error
}

type ratingService struct {
	repo repo.RatingRepository
}

func NewRatingService(repo repo.RatingRepository) RatingServiceIFace {
	return &ratingService{repo: repo}
}

func (r *ratingService) GetRating(ctx context.Context, username string) (*dto.RatingResponse, error) {
	rating, err := r.repo.GetRatingRepo(ctx, username)
	if err != nil {
		return nil, err
	}
	return &dto.RatingResponse{Stars: rating.Stars}, nil
}

func (r *ratingService) UpdateRating(ctx context.Context, username string, delta int) error {
	current, err := r.repo.GetRatingRepo(ctx, username)
	if err != nil {
		return fmt.Errorf("failed to get current rating: %w", err)
	}

	newStars := current.Stars + delta
	if newStars < 0 {
		newStars = 0
	}
	if newStars > 100 {
		newStars = 100
	}
	fmt.Println(newStars)

	if err := r.repo.UpdateRatingRepo(ctx, username, newStars); err != nil {
		return fmt.Errorf("failed to update rating: %w", err)
	}
	return nil
}
