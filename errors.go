package telnyx

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIError represents a Telnyx API error response.
type APIError struct {
	StatusCode int
	Errors     []ErrorDetail `json:"errors"`
}

// ErrorDetail is a single error object returned by the API.
type ErrorDetail struct {
	Code   string      `json:"code"`
	Title  string      `json:"title"`
	Detail string      `json:"detail"`
	Source ErrorSource `json:"source"`
}

// ErrorSource identifies where in the request the error occurred.
type ErrorSource struct {
	Pointer   string `json:"pointer"`
	Parameter string `json:"parameter"`
}

func (e *APIError) Error() string {
	if len(e.Errors) == 0 {
		return fmt.Sprintf("telnyx: HTTP %d", e.StatusCode)
	}
	return fmt.Sprintf("telnyx: HTTP %d: %s – %s", e.StatusCode, e.Errors[0].Title, e.Errors[0].Detail)
}

func parseAPIError(resp *http.Response) error {
	var apiErr APIError
	apiErr.StatusCode = resp.StatusCode
	_ = json.NewDecoder(resp.Body).Decode(&apiErr)
	return &apiErr
}

// ErrSignatureVerification is returned when a webhook signature is invalid.
type ErrSignatureVerification struct {
	Reason string
}

func (e *ErrSignatureVerification) Error() string {
	return fmt.Sprintf("telnyx: webhook signature verification failed: %s", e.Reason)
}
