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

type HealthData struct {
	OrganizationName string `json:"organization_name"`
	OrganizationID   string `json:"organization_id"`
}

// Check performs a health check.
func (s *HealthService) Check() (*BaseResponse, error) {

	url := fmt.Sprintf("%s/ping", s.base.BasePath)

	log.Printf("Checking health on %s", url)

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

	log.Printf("Health check response: %s", string(body))

	var response BaseResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil

}
