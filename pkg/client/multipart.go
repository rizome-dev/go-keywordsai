package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// MultipartField represents a field in a multipart form
type MultipartField struct {
	Name     string
	Value    string
	IsFile   bool
	FileName string
	Data     []byte
}

// PostMultipart sends a multipart form request
func (c *Client) PostMultipart(ctx context.Context, path string, fields []MultipartField, result interface{}) error {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	for _, field := range fields {
		if field.IsFile {
			part, err := writer.CreateFormFile(field.Name, field.FileName)
			if err != nil {
				return fmt.Errorf("failed to create form file: %w", err)
			}
			if _, err := part.Write(field.Data); err != nil {
				return fmt.Errorf("failed to write file data: %w", err)
			}
		} else {
			if err := writer.WriteField(field.Name, field.Value); err != nil {
				return fmt.Errorf("failed to write field: %w", err)
			}
		}
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, &buf)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return parseAPIError(resp.StatusCode, bodyBytes)
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}