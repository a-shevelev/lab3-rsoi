package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gateway-api/internal/dto"
	"net/http"
	"net/url"
)

type Library struct {
	BaseURL    string `envconfig:"BASE_URL"`
	HTTPClient *http.Client
}

func NewLibrary(baseURL string) *Library {
	return &Library{
		BaseURL:    baseURL,
		HTTPClient: http.DefaultClient,
	}
}

func (c *Library) GetLibraries(city string, page, size int) (*dto.LibraryPaginationResponse, error) {
	u, _ := url.Parse(fmt.Sprintf("%s/api/v1/libraries", c.BaseURL))
	q := u.Query()
	q.Set("city", city)
	if page > 0 {
		q.Set("page", fmt.Sprintf("%d", page))
	}
	if size > 0 {
		q.Set("size", fmt.Sprintf("%d", size))
	}
	u.RawQuery = q.Encode()

	resp, err := c.HTTPClient.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result dto.LibraryPaginationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Library) GetLibraryBooks(libraryUid string, page, size int, showAll bool) (*dto.LibraryBookPaginationResponse, error) {
	u, _ := url.Parse(fmt.Sprintf("%s/api/v1/libraries/%s/books", c.BaseURL, libraryUid))
	q := u.Query()
	if page > 0 {
		q.Set("page", fmt.Sprintf("%d", page))
	}
	if size > 0 {
		q.Set("size", fmt.Sprintf("%d", size))
	}
	if showAll {
		q.Set("showAll", "true")
	}
	u.RawQuery = q.Encode()

	resp, err := c.HTTPClient.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result dto.LibraryBookPaginationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Library) GetLibraryByUID(libraryUid string) (*dto.LibraryResponse, error) {
	resp, err := c.HTTPClient.Get(fmt.Sprintf("%s/api/v1/libraries/%s/", c.BaseURL, libraryUid))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result dto.LibraryResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Library) GetBookByUID(bookUid string) (*dto.BookResponse, error) {
	resp, err := c.HTTPClient.Get(fmt.Sprintf("%s/api/v1/books/%s/", c.BaseURL, bookUid))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result dto.BookResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Library) UpdateBookCondition(bookUid string, condition string) error {
	reqBody, _ := json.Marshal(map[string]string{"condition": condition})
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/v1/books/%s/condition", c.BaseURL, bookUid), bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update book condition, status: %d", resp.StatusCode)
	}
	return nil
}

func (c *Library) UpdateBookCount(libraryUid, bookUid string, delta int) error {
	req, _ := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/api/v1/library/%s/books/%s/count/%d/", c.BaseURL, libraryUid, bookUid, delta),
		nil,
	)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update book count, status: %d", resp.StatusCode)
	}
	return nil
}
