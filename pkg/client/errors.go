package client

import (
	"encoding/json"
	"fmt"
)

// APIError represents an error response from the KeywordsAI API
type APIError struct {
	StatusCode int    `json:"status_code"`
	ErrorText  string `json:"error"`
	Message    string `json:"message"`
	Details    string `json:"details"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("KeywordsAI API error (status %d): %s", e.StatusCode, e.Message)
	}
	if e.ErrorText != "" {
		return fmt.Sprintf("KeywordsAI API error (status %d): %s", e.StatusCode, e.ErrorText)
	}
	return fmt.Sprintf("KeywordsAI API error (status %d)", e.StatusCode)
}

// parseAPIError attempts to parse an API error response
func parseAPIError(statusCode int, body []byte) error {
	var apiErr APIError
	apiErr.StatusCode = statusCode
	
	if err := json.Unmarshal(body, &apiErr); err != nil {
		// If we can't parse as JSON, use the raw body as the error message
		return &APIError{
			StatusCode: statusCode,
			ErrorText:  string(body),
		}
	}
	
	return &apiErr
}