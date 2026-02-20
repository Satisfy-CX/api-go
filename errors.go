package scxapi

import "fmt"

// APIError is returned when the API responds with a non-200 status code.
type APIError struct {
	Operation  string
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s failed with status %d: %s", e.Operation, e.StatusCode, e.Body)
}

// RequestEncodeError is returned when request payload encoding fails.
type RequestEncodeError struct {
	Operation string
	Err       error
}

func (e *RequestEncodeError) Error() string {
	return fmt.Sprintf("%s request encoding failed: %v", e.Operation, e.Err)
}

func (e *RequestEncodeError) Unwrap() error { return e.Err }

// RequestBuildError is returned when request creation fails.
type RequestBuildError struct {
	Operation string
	Err       error
}

func (e *RequestBuildError) Error() string {
	return fmt.Sprintf("%s request build failed: %v", e.Operation, e.Err)
}

func (e *RequestBuildError) Unwrap() error { return e.Err }

// RequestExecuteError is returned when the HTTP client cannot complete a request.
type RequestExecuteError struct {
	Operation string
	Err       error
}

func (e *RequestExecuteError) Error() string {
	return fmt.Sprintf("%s request execution failed: %v", e.Operation, e.Err)
}

func (e *RequestExecuteError) Unwrap() error { return e.Err }

// ResponseReadError is returned when reading the response body fails.
type ResponseReadError struct {
	Operation string
	Err       error
}

func (e *ResponseReadError) Error() string {
	return fmt.Sprintf("%s response read failed: %v", e.Operation, e.Err)
}

func (e *ResponseReadError) Unwrap() error { return e.Err }

// ResponseDecodeError is returned when JSON response decoding fails.
type ResponseDecodeError struct {
	Operation string
	Err       error
}

func (e *ResponseDecodeError) Error() string {
	return fmt.Sprintf("%s response decode failed: %v", e.Operation, e.Err)
}

func (e *ResponseDecodeError) Unwrap() error { return e.Err }
