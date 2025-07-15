package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetWithQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check query parameters
		if r.URL.Query().Get("model") != "gpt-4" {
			t.Errorf("expected model=gpt-4 in query, got %s", r.URL.Query().Get("model"))
		}
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("expected limit=10 in query, got %s", r.URL.Query().Get("limit"))
		}
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	c := New("test-key", WithBaseURL(server.URL))
	
	type testQuery struct {
		Model string `url:"model"`
		Limit int    `url:"limit"`
	}
	
	query := testQuery{
		Model: "gpt-4",
		Limit: 10,
	}
	
	var result map[string]string
	err := c.GetWithQuery(context.Background(), "/test", query, &result)
	if err != nil {
		t.Fatalf("GetWithQuery() error = %v", err)
	}
	
	if result["status"] != "ok" {
		t.Errorf("expected status=ok, got %v", result["status"])
	}
}

func TestAccessorMethods(t *testing.T) {
	baseURL := "https://custom.api.com"
	apiKey := "test-api-key"
	httpClient := &http.Client{Timeout: 60 * time.Second}
	
	c := New(apiKey, WithBaseURL(baseURL), WithHTTPClient(httpClient))
	
	// Test BaseURL()
	if got := c.BaseURL(); got != baseURL {
		t.Errorf("BaseURL() = %v, want %v", got, baseURL)
	}
	
	// Test APIKey()
	if got := c.APIKey(); got != apiKey {
		t.Errorf("APIKey() = %v, want %v", got, apiKey)
	}
	
	// Test HTTPClient()
	if got := c.HTTPClient(); got != httpClient {
		t.Errorf("HTTPClient() = %v, want %v", got, httpClient)
	}
}

func TestAPIError(t *testing.T) {
	tests := []struct {
		name     string
		apiError APIError
		want     string
	}{
		{
			name: "with message",
			apiError: APIError{
				StatusCode: 400,
				Message:    "Bad request",
			},
			want: "KeywordsAI API error (status 400): Bad request",
		},
		{
			name: "with error text",
			apiError: APIError{
				StatusCode: 500,
				ErrorText:  "Internal server error",
			},
			want: "KeywordsAI API error (status 500): Internal server error",
		},
		{
			name: "with both message and error text",
			apiError: APIError{
				StatusCode: 400,
				Message:    "Bad request",
				ErrorText:  "Invalid input",
			},
			want: "KeywordsAI API error (status 400): Bad request",
		},
		{
			name: "with only status code",
			apiError: APIError{
				StatusCode: 404,
			},
			want: "KeywordsAI API error (status 404)",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.apiError.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseAPIError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       []byte
		wantErr    string
	}{
		{
			name:       "valid JSON error",
			statusCode: 400,
			body:       []byte(`{"error": "Bad request", "message": "Invalid input"}`),
			wantErr:    "KeywordsAI API error (status 400): Invalid input",
		},
		{
			name:       "invalid JSON",
			statusCode: 500,
			body:       []byte(`Internal server error`),
			wantErr:    "KeywordsAI API error (status 500): Internal server error",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parseAPIError(tt.statusCode, tt.body)
			if err == nil {
				t.Fatal("expected error")
			}
			if err.Error() != tt.wantErr {
				t.Errorf("parseAPIError() error = %v, want %v", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestPostMultipart(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check content type
		contentType := r.Header.Get("Content-Type")
		if !bytes.HasPrefix([]byte(contentType), []byte("multipart/form-data")) {
			t.Errorf("expected multipart/form-data content type, got %s", contentType)
		}
		
		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20) // 10MB
		if err != nil {
			t.Fatalf("failed to parse multipart form: %v", err)
		}
		
		// Check form values
		if r.FormValue("model") != "whisper-1" {
			t.Errorf("expected model=whisper-1, got %s", r.FormValue("model"))
		}
		
		// Check file
		file, header, err := r.FormFile("audio")
		if err != nil {
			t.Fatalf("failed to get file: %v", err)
		}
		defer file.Close()
		
		if header.Filename != "test.wav" {
			t.Errorf("expected filename=test.wav, got %s", header.Filename)
		}
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"text": "Hello world"})
	}))
	defer server.Close()
	
	c := New("test-key", WithBaseURL(server.URL))
	
	fields := []MultipartField{
		{Name: "model", Value: "whisper-1"},
		{Name: "audio", IsFile: true, FileName: "test.wav", Data: []byte("test audio data")},
	}
	
	var result map[string]string
	err := c.PostMultipart(context.Background(), "/transcribe", fields, &result)
	if err != nil {
		t.Fatalf("PostMultipart() error = %v", err)
	}
	
	if result["text"] != "Hello world" {
		t.Errorf("expected text='Hello world', got %v", result["text"])
	}
}

func TestPostMultipartError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid file"})
	}))
	defer server.Close()
	
	c := New("test-key", WithBaseURL(server.URL))
	
	fields := []MultipartField{
		{Name: "model", Value: "whisper-1"},
	}
	
	var result map[string]string
	err := c.PostMultipart(context.Background(), "/transcribe", fields, &result)
	if err == nil {
		t.Fatal("expected error")
	}
	
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	
	if apiErr.StatusCode != 400 {
		t.Errorf("expected status code 400, got %d", apiErr.StatusCode)
	}
}