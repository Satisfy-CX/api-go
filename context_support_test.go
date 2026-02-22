package scxapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestContext(t *testing.T) {
	if requestContext(nil) == nil {
		t.Fatalf("expected nil context to normalize to background context")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if got := requestContext(ctx); got != ctx {
		t.Fatalf("expected non-nil context to be preserved")
	}
}

func TestHealthCheckAcceptsContextAndNilContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ping" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"api_status":"error","api_message":"unexpected path"}`))
			return
		}
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"api_status":"error","api_message":"unexpected method"}`))
			return
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"api_status":"error","api_message":"missing auth"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"api_status":"ok","api_message":"ok"}`))
	}))
	defer server.Close()

	svc := NewSCXService("test-key", server.URL)
	svc.Client = server.Client()

	if _, err := svc.Health().Check(context.Background()); err != nil {
		t.Fatalf("expected context check to succeed: %v", err)
	}
	if _, err := svc.Health().Check(nil); err != nil {
		t.Fatalf("expected nil-context check to succeed: %v", err)
	}
}

func TestHealthCheckHonorsCanceledContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"api_status":"ok","api_message":"ok"}`))
	}))
	defer server.Close()

	svc := NewSCXService("test-key", server.URL)
	svc.Client = server.Client()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := svc.Health().Check(ctx)
	if err == nil {
		t.Fatalf("expected canceled context to return an error")
	}
	if !IsRequestExecuteError(err) {
		t.Fatalf("expected canceled context to be wrapped as RequestExecuteError, got %T", err)
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected error chain to include context canceled")
	}
}

func TestContentLibraryAcceptsNilContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/content/library" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"api_status":"error","api_message":"unexpected path"}`))
			return
		}
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"api_status":"error","api_message":"unexpected method"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"api_status":"ok","api_message":"ok","content_library":[]}`))
	}))
	defer server.Close()

	svc := NewSCXService("test-key", server.URL)
	svc.Client = server.Client()

	if _, err := svc.Content().Library(nil); err != nil {
		t.Fatalf("expected nil-context content library request to succeed: %v", err)
	}
}
