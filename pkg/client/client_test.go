package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		apiKey    string
		envAPIKey string
		envURL    string
		wantKey   string
		wantURL   string
	}{
		{
			name:    "Direct API key",
			apiKey:  "test-key",
			wantKey: "test-key",
			wantURL: defaultBaseURL,
		},
		{
			name:      "Env API key",
			apiKey:    "",
			envAPIKey: "env-test-key",
			wantKey:   "env-test-key",
			wantURL:   defaultBaseURL,
		},
		{
			name:    "Custom base URL",
			apiKey:  "test-key",
			envURL:  "https://custom.api.com",
			wantKey: "test-key",
			wantURL: "https://custom.api.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envAPIKey != "" {
				t.Setenv("KEYWORDSAI_API_KEY", tt.envAPIKey)
			}
			if tt.envURL != "" {
				t.Setenv("KEYWORDS_BASE_URL", tt.envURL)
			}

			var c *Client
			if tt.apiKey != "" {
				c = New(tt.apiKey)
			} else {
				c = New()
			}
			if c.apiKey != tt.wantKey {
				t.Errorf("New() apiKey = %v, want %v", c.apiKey, tt.wantKey)
			}
			if c.baseURL != tt.wantURL {
				t.Errorf("New() baseURL = %v, want %v", c.baseURL, tt.wantURL)
			}
		})
	}
}

func TestClientOptions(t *testing.T) {
	customClient := &http.Client{Timeout: 60 * time.Second}
	c := New("test-key",
		WithHTTPClient(customClient),
		WithBaseURL("https://custom.api.com"),
		WithTimeout(45*time.Second),
	)

	if c.httpClient.Timeout != 45*time.Second {
		t.Errorf("WithTimeout() timeout = %v, want %v", c.httpClient.Timeout, 45*time.Second)
	}
	if c.baseURL != "https://custom.api.com" {
		t.Errorf("WithBaseURL() baseURL = %v, want %v", c.baseURL, "https://custom.api.com")
	}
}

func TestClientDo(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		serverResponse interface{}
		serverStatus   int
		wantErr        bool
	}{
		{
			name:           "Successful GET",
			method:         http.MethodGet,
			path:           "/test",
			serverResponse: map[string]string{"status": "ok"},
			serverStatus:   http.StatusOK,
			wantErr:        false,
		},
		{
			name:   "Successful POST",
			method: http.MethodPost,
			path:   "/test",
			body:   map[string]string{"name": "test"},
			serverResponse: map[string]string{"id": "123"},
			serverStatus:   http.StatusOK,
			wantErr:        false,
		},
		{
			name:         "Error response",
			method:       http.MethodGet,
			path:         "/test",
			serverStatus: http.StatusBadRequest,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Authorization") != "Bearer test-key" {
					t.Errorf("Missing or incorrect Authorization header")
				}
				if r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("Missing or incorrect Content-Type header")
				}

				w.WriteHeader(tt.serverStatus)
				if tt.serverResponse != nil {
					json.NewEncoder(w).Encode(tt.serverResponse)
				}
			}))
			defer server.Close()

			c := New("test-key", WithBaseURL(server.URL))
			
			var result map[string]string
			err := c.do(context.Background(), tt.method, tt.path, tt.body, &result)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("do() error = %v, wantErr %v", err, tt.wantErr)
			}
			
			if !tt.wantErr && tt.serverResponse != nil {
				expected := tt.serverResponse.(map[string]string)
				for k, v := range expected {
					if result[k] != v {
						t.Errorf("do() result[%s] = %v, want %v", k, result[k], v)
					}
				}
			}
		})
	}
}

func TestClientMethods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Echo back the method for verification
		response := map[string]string{"method": r.Method}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	c := New("test-key", WithBaseURL(server.URL))
	ctx := context.Background()

	tests := []struct {
		name       string
		fn         func() error
		wantMethod string
	}{
		{
			name: "GET",
			fn: func() error {
				var result map[string]string
				return c.Get(ctx, "/test", &result)
			},
			wantMethod: "GET",
		},
		{
			name: "POST",
			fn: func() error {
				var result map[string]string
				return c.Post(ctx, "/test", map[string]string{"test": "data"}, &result)
			},
			wantMethod: "POST",
		},
		{
			name: "PUT",
			fn: func() error {
				var result map[string]string
				return c.Put(ctx, "/test", map[string]string{"test": "data"}, &result)
			},
			wantMethod: "PUT",
		},
		{
			name: "PATCH",
			fn: func() error {
				var result map[string]string
				return c.Patch(ctx, "/test", map[string]string{"test": "data"}, &result)
			},
			wantMethod: "PATCH",
		},
		{
			name: "DELETE",
			fn: func() error {
				var result map[string]string
				return c.Delete(ctx, "/test", &result)
			},
			wantMethod: "DELETE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fn(); err != nil {
				t.Fatalf("Method %s failed: %v", tt.name, err)
			}
		})
	}
}