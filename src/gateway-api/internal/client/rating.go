package client

import (
	"encoding/json"
	"fmt"
	"gateway-api/internal/dto"
	"gateway-api/pkg/circuit"
	"gateway-api/pkg/ext"
	"net/http"
	"time"
)

type Rating struct {
	BaseURL    string `envconfig:"BASE_URL"`
	HTTPClient *http.Client
	GetBreaker *circuit.Breaker
}

func NewRating(baseURL string) *Rating {
	return &Rating{
		BaseURL:    baseURL,
		HTTPClient: http.DefaultClient,
		GetBreaker: circuit.NewBreaker(3, 5*time.Second, 60*time.Second, 3),
	}
}

func (c *Rating) isHealthy() bool {
	resp, err := c.HTTPClient.Get(fmt.Sprintf("%s/manage/health", c.BaseURL))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (c *Rating) Get(username string) (*dto.UserRatingResponse, error) {
	action := func() (*dto.UserRatingResponse, error) {
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

	fallback := func() *dto.UserRatingResponse {
		return &dto.UserRatingResponse{
			Stars: 0,
		}
	}

	return circuit.WithCircuitBreaker(c.GetBreaker, action, fallback, c.isHealthy)
}

func (c *Rating) Update(username string, stars int) error {

	if !c.isHealthy() {
		return ext.ServiceUnavailableError
	}
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
