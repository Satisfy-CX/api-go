package scxapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestNewSCXService_DefaultBasePathAndSubservices(t *testing.T) {
	apiKey := "test-api-key"
	svc := NewSCXService(apiKey, "")

	if svc.BasePath != PublicAPIPath {
		t.Fatalf("expected default base path %q, got %q", PublicAPIPath, svc.BasePath)
	}
	if svc.APIKey != apiKey {
		t.Fatalf("expected API key %q, got %q", apiKey, svc.APIKey)
	}
	if svc.Client == nil {
		t.Fatal("expected HTTP client to be initialized")
	}
	if svc.Client.Timeout != 30*time.Second {
		t.Fatalf("expected timeout to be 30s, got %s", svc.Client.Timeout)
	}
	if svc.Health() == nil {
		t.Fatal("expected health service to be initialized")
	}
	if svc.Content() == nil {
		t.Fatal("expected content service to be initialized")
	}
}

func TestNewSCXService_CustomBasePath(t *testing.T) {
	customBasePath := "http://localhost:9999"
	svc := NewSCXService("test-key", customBasePath)

	if svc.BasePath != customBasePath {
		t.Fatalf("expected custom base path %q, got %q", customBasePath, svc.BasePath)
	}
}

func TestHealthCheck_Success(t *testing.T) {
	apiKey := "health-key"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/ping" {
			http.Error(w, "bad path", http.StatusNotFound)
			return
		}
		if r.Header.Get("Authorization") != "Bearer "+apiKey {
			http.Error(w, "missing auth", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"api_status":"ok",
			"api_message":"healthy",
			"health":{"organization_name":"Acme","organization_id":"org-1"}
		}`))
	}))
	defer server.Close()

	svc := NewSCXService(apiKey, server.URL)
	resp, err := svc.Health().Check()
	if err != nil {
		t.Fatalf("health check returned unexpected error: %v", err)
	}
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %q", resp.Status)
	}
	if resp.Message != "healthy" {
		t.Fatalf("expected message healthy, got %q", resp.Message)
	}
	if resp.Health == nil {
		t.Fatal("expected health payload to be present")
	}
	if resp.Health.OrganizationName != "Acme" || resp.Health.OrganizationID != "org-1" {
		t.Fatalf("unexpected health payload: %+v", *resp.Health)
	}
}

func TestHealthCheck_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer server.Close()

	svc := NewSCXService("health-key", server.URL)
	_, err := svc.Health().Check()
	if err == nil {
		t.Fatal("expected error for non-200 health response")
	}
	if !strings.Contains(err.Error(), "health check failed") {
		t.Fatalf("expected health check failure error, got %v", err)
	}
}

func TestContentLibrary_Success(t *testing.T) {
	ctx := context.Background()
	apiKey := "content-key"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/content/library" {
			http.Error(w, "bad path", http.StatusNotFound)
			return
		}
		if r.Header.Get("Authorization") != "Bearer "+apiKey {
			http.Error(w, "missing auth", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"api_status":"ok",
			"api_message":"loaded",
			"content_library":[{"id":"c1","name":"Doc 1","type":"Article"}]
		}`))
	}))
	defer server.Close()

	svc := NewSCXService(apiKey, server.URL)
	resp, err := svc.Content().Library(ctx)
	if err != nil {
		t.Fatalf("content library returned unexpected error: %v", err)
	}
	if resp.Status != "ok" || resp.Message != "loaded" {
		t.Fatalf("unexpected base response: %+v", resp)
	}
	if resp.ContentLibrary == nil {
		t.Fatal("expected content library to be present")
	}
	if len(*resp.ContentLibrary) != 1 {
		t.Fatalf("expected one content item, got %d", len(*resp.ContentLibrary))
	}
	item := (*resp.ContentLibrary)[0]
	if item.ID != "c1" || item.Name != "Doc 1" || item.Type != ContentTypeArticle {
		t.Fatalf("unexpected content item: %+v", item)
	}
}

func TestContentLibrary_NonOKStatus(t *testing.T) {
	ctx := context.Background()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusBadRequest)
	}))
	defer server.Close()

	svc := NewSCXService("content-key", server.URL)
	_, err := svc.Content().Library(ctx)
	if err == nil {
		t.Fatal("expected error for non-200 content library response")
	}
	if !strings.Contains(err.Error(), "content list failed") {
		t.Fatalf("expected content list failure error, got %v", err)
	}
}

func TestContentGet_PathEscapingAndSuccess(t *testing.T) {
	ctx := context.Background()
	apiKey := "content-get-key"
	id := "id/with space?"
	expectedPath := "/content/get/" + url.PathEscape(id)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.EscapedPath() != expectedPath {
			http.Error(w, "bad escaped path", http.StatusNotFound)
			return
		}
		if r.Header.Get("Authorization") != "Bearer "+apiKey {
			http.Error(w, "missing auth", http.StatusUnauthorized)
			return
		}

		_, _ = w.Write([]byte(`{
			"api_status":"ok",
			"api_message":"found",
			"content_library":[{"id":"id/with space?","name":"Doc 1","type":"Article"}]
		}`))
	}))
	defer server.Close()

	svc := NewSCXService(apiKey, server.URL)
	resp, err := svc.Content().Get(ctx, id)
	if err != nil {
		t.Fatalf("content get returned unexpected error: %v", err)
	}
	if resp.Status != "ok" || resp.Message != "found" {
		t.Fatalf("unexpected base response: %+v", resp)
	}
}

func TestContentManage_Success(t *testing.T) {
	ctx := context.Background()
	apiKey := "content-manage-key"
	request := ContentManageRequest{
		ID:              "id/with space?",
		Name:            "Name",
		Title:           "Title",
		Body:            "Body",
		SameAs:          []string{"https://example.com/a"},
		DifferentFrom:   []string{"https://example.com/b"},
		Context:         "Context",
		Language:        "en",
		Type:            ContentTypeArticle,
		ImportSourceURL: "https://example.com/src",
	}

	expectedPath := "/content/" + url.PathEscape(request.ID)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.EscapedPath() != expectedPath {
			http.Error(w, "bad escaped path", http.StatusNotFound)
			return
		}
		if r.Header.Get("Authorization") != "Bearer "+apiKey {
			http.Error(w, "missing auth", http.StatusUnauthorized)
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "missing content type", http.StatusBadRequest)
			return
		}

		var gotRequest ContentManageRequest
		if err := json.NewDecoder(r.Body).Decode(&gotRequest); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}

		if gotRequest.ID != request.ID ||
			gotRequest.Name != request.Name ||
			gotRequest.Title != request.Title ||
			gotRequest.Body != request.Body ||
			gotRequest.Context != request.Context ||
			gotRequest.Language != request.Language ||
			gotRequest.Type != request.Type ||
			gotRequest.ImportSourceURL != request.ImportSourceURL {
			http.Error(w, "invalid payload fields", http.StatusBadRequest)
			return
		}
		if len(gotRequest.SameAs) != 1 || gotRequest.SameAs[0] != request.SameAs[0] {
			http.Error(w, "invalid same_as payload", http.StatusBadRequest)
			return
		}
		if len(gotRequest.DifferentFrom) != 1 || gotRequest.DifferentFrom[0] != request.DifferentFrom[0] {
			http.Error(w, "invalid different_from payload", http.StatusBadRequest)
			return
		}

		_, _ = w.Write([]byte(`{"id":"id/with space?","status":2}`))
	}))
	defer server.Close()

	svc := NewSCXService(apiKey, server.URL)
	resp, err := svc.Content().Manage(ctx, request)
	if err != nil {
		t.Fatalf("content manage returned unexpected error: %v", err)
	}
	if resp.ID != request.ID {
		t.Fatalf("expected response id %q, got %q", request.ID, resp.ID)
	}
	if resp.Status != ContentManageStatusUpdated {
		t.Fatalf("expected status %d, got %d", ContentManageStatusUpdated, resp.Status)
	}
}

func TestContentManage_NonOKStatus(t *testing.T) {
	ctx := context.Background()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "update failed", http.StatusInternalServerError)
	}))
	defer server.Close()

	svc := NewSCXService("content-manage-key", server.URL)
	_, err := svc.Content().Manage(ctx, ContentManageRequest{ID: "123"})
	if err == nil {
		t.Fatal("expected error for non-200 content manage response")
	}
	if !strings.Contains(err.Error(), "content manage failed") {
		t.Fatalf("expected content manage failure error, got %v", err)
	}
}
