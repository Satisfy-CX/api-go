package scxapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// HealthService provides health-check endpoints.
type HealthService struct {
	base *SCXService
}

// HealthResponse is the response from a health check.
type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// Check performs a health check.
func (s *HealthService) Check() (*HealthResponse, error) {
	log.Printf("Checking health on %s", s.base.BasePath)

	url := fmt.Sprintf("%s/ping", s.base.BasePath)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.base.APIKey)

	resp, err := s.base.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("health check failed: %s", string(body))
	}

	var out HealthResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
