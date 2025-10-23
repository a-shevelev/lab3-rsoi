package client

import (
	"encoding/json"
	"fmt"
	"gateway-api/internal/dto"
	"net/http"
)

type Rating struct {
	BaseURL    string `envconfig:"BASE_URL"`
	HTTPClient *http.Client
}

// NewRatingClient создаёт новый клиент для RatingService
func NewRating(baseURL string) *Rating {
	return &Rating{
		BaseURL:    baseURL,
		HTTPClient: http.DefaultClient,
	}
}

func (c *Rating) Get(username string) (*dto.UserRatingResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/rating", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-User-Name", username)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result dto.UserRatingResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Rating) Update(username string, stars int) error {
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/rating/stars/%d/", c.BaseURL, stars), nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-User-Name", username)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	return nil
}
