package scxapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestTypedErrorTypeHelpers(t *testing.T) {
	inner := errors.New("boom")

	type helperCheck struct {
		name string
		err  error
		is   func(error) bool
	}

	checks := []helperCheck{
		{
			name: "api",
			err:  fmt.Errorf("wrapped: %w", &APIError{Operation: "content.get", StatusCode: http.StatusBadRequest, Body: "bad"}),
			is:   IsAPIError,
		},
		{
			name: "request encode",
			err:  fmt.Errorf("wrapped: %w", &RequestEncodeError{Operation: "content.manage", Err: inner}),
			is:   IsRequestEncodeError,
		},
		{
			name: "request build",
			err:  fmt.Errorf("wrapped: %w", &RequestBuildError{Operation: "content.get", Err: inner}),
			is:   IsRequestBuildError,
		},
		{
			name: "request execute",
			err:  fmt.Errorf("wrapped: %w", &RequestExecuteError{Operation: "content.get", Err: inner}),
			is:   IsRequestExecuteError,
		},
		{
			name: "response read",
			err:  fmt.Errorf("wrapped: %w", &ResponseReadError{Operation: "content.get", Err: inner}),
			is:   IsResponseReadError,
		},
		{
			name: "response decode",
			err:  fmt.Errorf("wrapped: %w", &ResponseDecodeError{Operation: "content.get", Err: inner}),
			is:   IsResponseDecodeError,
		},
	}

	for _, tc := range checks {
		if !tc.is(tc.err) {
			t.Fatalf("expected %s helper to detect wrapped error", tc.name)
		}
	}
}

func TestAPIStatusAndOperationHelpers(t *testing.T) {
	err := fmt.Errorf("wrapped: %w", &APIError{
		Operation:  "content.get",
		StatusCode: http.StatusNotFound,
		Body:       `{"detail":"missing"}`,
	})

	status, ok := StatusCode(err)
	if !ok || status != http.StatusNotFound {
		t.Fatalf("expected status code %d, got %d (ok=%v)", http.StatusNotFound, status, ok)
	}

	body, ok := ResponseBody(err)
	if !ok || body == "" {
		t.Fatalf("expected response body to be extracted")
	}

	op, ok := ErrorOperation(err)
	if !ok || op != "content.get" {
		t.Fatalf("expected operation content.get, got %q (ok=%v)", op, ok)
	}

	if !IsStatusCode(err, http.StatusNotFound) || !IsNotFound(err) {
		t.Fatalf("expected not found helpers to match")
	}
	if !IsClientError(err) {
		t.Fatalf("expected 404 to be client error")
	}
	if IsServerError(err) {
		t.Fatalf("did not expect 404 to be server error")
	}
	if IsRetryable(err) {
		t.Fatalf("did not expect 404 to be retryable")
	}
}

func TestRetryAndTimeoutHelpers(t *testing.T) {
	rateLimited := &APIError{Operation: "content.library", StatusCode: http.StatusTooManyRequests, Body: "rate limit"}
	if !IsRateLimited(rateLimited) {
		t.Fatalf("expected rate limited helper to match")
	}
	if !IsRetryable(rateLimited) {
		t.Fatalf("expected 429 to be retryable")
	}

	serverErr := &APIError{Operation: "content.library", StatusCode: http.StatusBadGateway, Body: "upstream error"}
	if !IsServerError(serverErr) || !IsRetryable(serverErr) {
		t.Fatalf("expected 502 to be server and retryable")
	}

	timeoutErr := &RequestExecuteError{Operation: "content.library", Err: context.DeadlineExceeded}
	if !IsTimeout(timeoutErr) {
		t.Fatalf("expected deadline exceeded to be timeout")
	}
	if !IsRetryable(timeoutErr) {
		t.Fatalf("expected timeout to be retryable")
	}

	readErr := &ResponseReadError{Operation: "content.get", Err: errors.New("short read")}
	if !IsRetryable(readErr) {
		t.Fatalf("expected response read errors to be retryable")
	}

	encodeErr := &RequestEncodeError{Operation: "content.manage", Err: errors.New("json: unsupported value")}
	if IsRetryable(encodeErr) {
		t.Fatalf("did not expect request encoding errors to be retryable")
	}
}
