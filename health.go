package scxapi

import (
	"context"
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
func (s *HealthService) Check(ctx context.Context) (*BaseResponse, error) {
	const operation = "health.check"

	url := fmt.Sprintf("%s/ping", s.base.BasePath)

	log.Printf("Checking health on %s", url)

	req, err := http.NewRequestWithContext(requestContext(ctx), http.MethodPost, url, nil)
	if err != nil {
		return nil, &RequestBuildError{Operation: operation, Err: err}
	}

	req.Header.Set("Authorization", "Bearer "+s.base.APIKey)

	resp, err := s.base.Client.Do(req)
	if err != nil {
		return nil, &RequestExecuteError{Operation: operation, Err: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &ResponseReadError{Operation: operation, Err: err}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{
			Operation:  operation,
			StatusCode: resp.StatusCode,
			Body:       string(body),
		}
	}

	var response BaseResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, &ResponseDecodeError{Operation: operation, Err: err}
	}

	return &response, nil

}
