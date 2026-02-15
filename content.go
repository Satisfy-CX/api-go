package scxapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// ContentService provides content CRUD endpoints.
type ContentService struct {
	base *SCXService
}

type ContentType string

const (
	ContentTypeArticle          = ContentType("Article")
	ContentTypeAuthorPersona    = ContentType("Author Persona")
	ContentTypeAudiencePersona  = ContentType("Audience Persona")
	ContentTypeKeywordList      = ContentType("Keywords list")
	ContentTypeBrandEntity      = ContentType("Brand Entity")
	ContentTypeCompetitorEntity = ContentType("Competitor Entity")
	ContentTypeThirdPartyEntity = ContentType("Third Party Entity")
)

type Content struct {
	ID              string      `json:"id,omitempty"`
	User            string      `json:"user,omitempty"`
	Name            string      `json:"name,omitempty"`
	Title           string      `json:"title,omitempty"`
	Body            string      `json:"body,omitempty"`
	SameAs          []string    `json:"same_as,omitempty"`
	DifferentFrom   []string    `json:"different_from,omitempty"`
	Context         string      `json:"context,omitempty"`
	Language        string      `json:"language,omitempty"`
	Type            ContentType `json:"type,omitempty"`
	ImportSourceURL string      `json:"import_source_url,omitempty"`
	CreatedAt       time.Time   `json:"created_at,omitempty"`
	UpdatedAt       time.Time   `json:"updated_at,omitempty"`
}

type ContentListRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// type ContentListResponse struct {
// 	ContentLibrary  `json:"content_library"`
// }

// List returns content for the organization (GET /content).
func (s *ContentService) List(ctx context.Context, request ContentListRequest) (*BaseResponse, error) {
	url := fmt.Sprintf("%s/content", s.base.BasePath)
	if request.Page > 0 || request.PageSize > 0 {
		url += "?"
		if request.Page > 0 {
			url += fmt.Sprintf("page=%d", request.Page)
		}
		if request.PageSize > 0 {
			if request.Page > 0 {
				url += "&"
			}
			url += fmt.Sprintf("page_size=%d", request.PageSize)
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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
		return nil, fmt.Errorf("content list failed: %s", string(body))
	}

	var response BaseResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

type ContentManageStatus int

const (
	_ ContentManageStatus = iota
	StatusCreated
	StatusUpdated
	StatusDeleted
	StatusError
)

type ContentManageRequest struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	Title           string      `json:"title"`
	Body            string      `json:"body"`
	SameAs          []string    `json:"same_as"`
	DifferentFrom   []string    `json:"different_from"`
	Context         string      `json:"context"`
	Language        string      `json:"language"`
	Type            ContentType `json:"type"`
	ImportSourceURL string      `json:"import_source_url"`
}

type ContentManageResponse struct {
	ID     string              `json:"id"`
	Status ContentManageStatus `json:"status"`
}

// Get fetches a single content item by ID (GET /content/{id}).
func (s *ContentService) Get(ctx context.Context, id string) (*Content, error) {
	url := fmt.Sprintf("%s/content/%s", s.base.BasePath, url.PathEscape(id))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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
		return nil, fmt.Errorf("content get failed: %s", string(body))
	}

	var out Content
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Manage is used to manage the Content library. It creates or updates a content item.
func (s *ContentService) Manage(ctx context.Context, request ContentManageRequest) (*ContentManageResponse, error) {
	url := fmt.Sprintf("%s/content/%s", s.base.BasePath, url.PathEscape(request.ID))

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
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
		return nil, fmt.Errorf("content manage failed: %s", string(body))
	}

	var out ContentManageResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
