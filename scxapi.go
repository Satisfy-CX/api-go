package scxapi

import (
	"net/http"
	"time"
)

// SCXService is the main API client. Use NewSCXService to construct it.
// Sub-services (Health, Content) are created with a reference to this client
// so they share BasePath, APIKey, and Client.
type SCXService struct {
	BasePath string
	APIKey   string
	Client   *http.Client

	health  *HealthService
	content *ContentService
}

// Health returns the health-check sub-service.
func (s *SCXService) Health() *HealthService { return s.health }

// Content returns the content sub-service.
func (s *SCXService) Content() *ContentService { return s.content }

const DefaultBasePath = "https://api.satisfycx.ai/api/v1"

// NewSCXService builds an SCXService.
func NewSCXService(apiKey string, customBasePath string) *SCXService {

	basePath := DefaultBasePath
	if customBasePath != "" {
		basePath = customBasePath
	}

	svc := &SCXService{
		BasePath: basePath,
		APIKey:   apiKey,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	svc.health = &HealthService{base: svc}
	svc.content = &ContentService{base: svc}

	return svc
}
