package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gateway-api/internal/dto"
	"net/http"
)

type Reservation struct {
	BaseURL    string `envconfig:"BASE_URL"`
	HTTPClient *http.Client
}

// NewReservationClient создаёт новый клиент для ReservationService
func NewReservation(baseURL string) *Reservation {
	return &Reservation{
		BaseURL:    baseURL,
		HTTPClient: http.DefaultClient,
	}
}

func (c *Reservation) Get(username string) ([]dto.ReservationResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/reservation", c.BaseURL), nil)
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
	fmt.Println(resp.Body)

	var result []dto.ReservationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Reservation) GetByUID(uid string) (*dto.ReservationResponse, error) {
	req, err := http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/reservation/%s", c.BaseURL, uid), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result dto.ReservationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Reservation) GetCurrentAmount(username string) (int, error) {
	req, err := http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/reservation/amount", c.BaseURL), nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("X-User-Name", username)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var payload struct {
		Amount int `json:"amount"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return 0, err
	}

	return payload.Amount, nil
}

func (c *Reservation) Create(username string, req dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
	body, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/api/v1/reservation", c.BaseURL),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("X-User-Name", username)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result dto.ReservationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Reservation) UpdateStatus(uid string, date string) error {
	body, err := json.Marshal(map[string]string{"date": date})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/api/v1/reservation/%s", c.BaseURL, uid),
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}
