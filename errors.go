package scxapi

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
)

// APIError is returned when the API responds with a non-200 status code.
type APIError struct {
	Operation  string
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s failed with status %d: %s", e.Operation, e.StatusCode, e.Body)
}

// Op returns the API operation name that produced this error.
func (e *APIError) Op() string {
	if e == nil {
		return ""
	}
	return e.Operation
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

// Op returns the API operation name that produced this error.
func (e *RequestEncodeError) Op() string {
	if e == nil {
		return ""
	}
	return e.Operation
}

// RequestBuildError is returned when request creation fails.
type RequestBuildError struct {
	Operation string
	Err       error
}

func (e *RequestBuildError) Error() string {
	return fmt.Sprintf("%s request build failed: %v", e.Operation, e.Err)
}

func (e *RequestBuildError) Unwrap() error { return e.Err }

// Op returns the API operation name that produced this error.
func (e *RequestBuildError) Op() string {
	if e == nil {
		return ""
	}
	return e.Operation
}

// RequestExecuteError is returned when the HTTP client cannot complete a request.
type RequestExecuteError struct {
	Operation string
	Err       error
}

func (e *RequestExecuteError) Error() string {
	return fmt.Sprintf("%s request execution failed: %v", e.Operation, e.Err)
}

func (e *RequestExecuteError) Unwrap() error { return e.Err }

// Op returns the API operation name that produced this error.
func (e *RequestExecuteError) Op() string {
	if e == nil {
		return ""
	}
	return e.Operation
}

// ResponseReadError is returned when reading the response body fails.
type ResponseReadError struct {
	Operation string
	Err       error
}

func (e *ResponseReadError) Error() string {
	return fmt.Sprintf("%s response read failed: %v", e.Operation, e.Err)
}

func (e *ResponseReadError) Unwrap() error { return e.Err }

// Op returns the API operation name that produced this error.
func (e *ResponseReadError) Op() string {
	if e == nil {
		return ""
	}
	return e.Operation
}

// ResponseDecodeError is returned when JSON response decoding fails.
type ResponseDecodeError struct {
	Operation string
	Err       error
}

func (e *ResponseDecodeError) Error() string {
	return fmt.Sprintf("%s response decode failed: %v", e.Operation, e.Err)
}

func (e *ResponseDecodeError) Unwrap() error { return e.Err }

// Op returns the API operation name that produced this error.
func (e *ResponseDecodeError) Op() string {
	if e == nil {
		return ""
	}
	return e.Operation
}

// AsAPIError extracts an APIError from err.
func AsAPIError(err error) (*APIError, bool) {
	var out *APIError
	if errors.As(err, &out) {
		return out, true
	}
	return nil, false
}

// IsAPIError reports whether err is or wraps APIError.
func IsAPIError(err error) bool {
	_, ok := AsAPIError(err)
	return ok
}

// AsRequestEncodeError extracts a RequestEncodeError from err.
func AsRequestEncodeError(err error) (*RequestEncodeError, bool) {
	var out *RequestEncodeError
	if errors.As(err, &out) {
		return out, true
	}
	return nil, false
}

// IsRequestEncodeError reports whether err is or wraps RequestEncodeError.
func IsRequestEncodeError(err error) bool {
	_, ok := AsRequestEncodeError(err)
	return ok
}

// AsRequestBuildError extracts a RequestBuildError from err.
func AsRequestBuildError(err error) (*RequestBuildError, bool) {
	var out *RequestBuildError
	if errors.As(err, &out) {
		return out, true
	}
	return nil, false
}

// IsRequestBuildError reports whether err is or wraps RequestBuildError.
func IsRequestBuildError(err error) bool {
	_, ok := AsRequestBuildError(err)
	return ok
}

// AsRequestExecuteError extracts a RequestExecuteError from err.
func AsRequestExecuteError(err error) (*RequestExecuteError, bool) {
	var out *RequestExecuteError
	if errors.As(err, &out) {
		return out, true
	}
	return nil, false
}

// IsRequestExecuteError reports whether err is or wraps RequestExecuteError.
func IsRequestExecuteError(err error) bool {
	_, ok := AsRequestExecuteError(err)
	return ok
}

// AsResponseReadError extracts a ResponseReadError from err.
func AsResponseReadError(err error) (*ResponseReadError, bool) {
	var out *ResponseReadError
	if errors.As(err, &out) {
		return out, true
	}
	return nil, false
}

// IsResponseReadError reports whether err is or wraps ResponseReadError.
func IsResponseReadError(err error) bool {
	_, ok := AsResponseReadError(err)
	return ok
}

// AsResponseDecodeError extracts a ResponseDecodeError from err.
func AsResponseDecodeError(err error) (*ResponseDecodeError, bool) {
	var out *ResponseDecodeError
	if errors.As(err, &out) {
		return out, true
	}
	return nil, false
}

// IsResponseDecodeError reports whether err is or wraps ResponseDecodeError.
func IsResponseDecodeError(err error) bool {
	_, ok := AsResponseDecodeError(err)
	return ok
}

// ErrorOperation extracts the operation name from any SDK typed error.
func ErrorOperation(err error) (string, bool) {
	var opErr interface {
		Op() string
	}
	if errors.As(err, &opErr) {
		return opErr.Op(), true
	}
	return "", false
}

// StatusCode extracts an HTTP status code from APIError.
func StatusCode(err error) (int, bool) {
	apiErr, ok := AsAPIError(err)
	if !ok {
		return 0, false
	}
	return apiErr.StatusCode, true
}

// ResponseBody extracts an API response body from APIError.
func ResponseBody(err error) (string, bool) {
	apiErr, ok := AsAPIError(err)
	if !ok {
		return "", false
	}
	return apiErr.Body, true
}

// IsStatusCode reports whether err is APIError with the given status code.
func IsStatusCode(err error, code int) bool {
	status, ok := StatusCode(err)
	return ok && status == code
}

// IsClientError reports whether err is APIError with status 4xx.
func IsClientError(err error) bool {
	status, ok := StatusCode(err)
	return ok && status >= http.StatusBadRequest && status <= 499
}

// IsServerError reports whether err is APIError with status 5xx.
func IsServerError(err error) bool {
	status, ok := StatusCode(err)
	return ok && status >= http.StatusInternalServerError && status <= 599
}

// IsNotFound reports whether err is APIError with status 404.
func IsNotFound(err error) bool { return IsStatusCode(err, http.StatusNotFound) }

// IsUnauthorized reports whether err is APIError with status 401.
func IsUnauthorized(err error) bool { return IsStatusCode(err, http.StatusUnauthorized) }

// IsForbidden reports whether err is APIError with status 403.
func IsForbidden(err error) bool { return IsStatusCode(err, http.StatusForbidden) }

// IsConflict reports whether err is APIError with status 409.
func IsConflict(err error) bool { return IsStatusCode(err, http.StatusConflict) }

// IsRateLimited reports whether err is APIError with status 429.
func IsRateLimited(err error) bool { return IsStatusCode(err, http.StatusTooManyRequests) }

// IsTimeout reports whether err represents a timeout.
func IsTimeout(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

// IsRetryable reports whether the operation can be retried.
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}
	if IsTimeout(err) {
		return true
	}
	if status, ok := StatusCode(err); ok {
		return status == http.StatusRequestTimeout ||
			status == http.StatusTooManyRequests ||
			(status >= http.StatusInternalServerError && status <= 599)
	}
	// Network execution and stream read issues are often transient.
	return IsRequestExecuteError(err) || IsResponseReadError(err)
}
