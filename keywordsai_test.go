package keywordsai

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/rizome-dev/go-keywordsai/pkg/client"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		params []interface{}
		check  func(*testing.T, *SDK)
	}{
		{
			name:   "no parameters",
			params: []interface{}{},
			check: func(t *testing.T, sdk *SDK) {
				if sdk == nil {
					t.Fatal("expected non-nil SDK")
				}
				if sdk.Client == nil {
					t.Fatal("expected non-nil Client")
				}
				if sdk.Logs == nil {
					t.Fatal("expected non-nil Logs service")
				}
				if sdk.Prompts == nil {
					t.Fatal("expected non-nil Prompts service")
				}
				if sdk.Models == nil {
					t.Fatal("expected non-nil Models service")
				}
				if sdk.Keys == nil {
					t.Fatal("expected non-nil Keys service")
				}
				if sdk.Integrations == nil {
					t.Fatal("expected non-nil Integrations service")
				}
			},
		},
		{
			name:   "with API key",
			params: []interface{}{"test-api-key"},
			check: func(t *testing.T, sdk *SDK) {
				if sdk.Client.APIKey() != "test-api-key" {
					t.Errorf("expected API key to be 'test-api-key', got %s", sdk.Client.APIKey())
				}
			},
		},
		{
			name:   "with options",
			params: []interface{}{client.WithTimeout(30 * time.Second)},
			check: func(t *testing.T, sdk *SDK) {
				// Just verify SDK is created properly
				if sdk == nil {
					t.Fatal("expected non-nil SDK")
				}
			},
		},
		{
			name: "with API key and options",
			params: []interface{}{
				"test-api-key",
				client.WithBaseURL("https://custom.api.url"),
				client.WithTimeout(45 * time.Second),
			},
			check: func(t *testing.T, sdk *SDK) {
				if sdk.Client.APIKey() != "test-api-key" {
					t.Errorf("expected API key to be 'test-api-key', got %s", sdk.Client.APIKey())
				}
				if sdk.Client.BaseURL() != "https://custom.api.url" {
					t.Errorf("expected base URL to be 'https://custom.api.url', got %s", sdk.Client.BaseURL())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sdk := New(tt.params...)
			tt.check(t, sdk)
		})
	}
}

func TestTypeAliases(t *testing.T) {
	// Verify that type aliases work correctly
	var _ RequestLog = RequestLog{}
	var _ Message = Message{}
	var _ LogFilter = LogFilter{}
	var _ Prompt = Prompt{}
	var _ PromptVersion = PromptVersion{}
	var _ Model = Model{}
	var _ TemporaryKey = TemporaryKey{}
	var _ TTSRequest = TTSRequest{}
	var _ STTRequest = STTRequest{}
	var _ STTResponse = STTResponse{}
	var _ EmbeddingRequest = EmbeddingRequest{}
	var _ EmbeddingResponse = EmbeddingResponse{}
}

func TestPointerUtilities(t *testing.T) {
	// Test String
	s := "test"
	sp := String(s)
	if sp == nil || *sp != s {
		t.Errorf("String() failed")
	}

	// Test Int
	i := 42
	ip := Int(i)
	if ip == nil || *ip != i {
		t.Errorf("Int() failed")
	}

	// Test Int64
	i64 := int64(42)
	i64p := Int64(i64)
	if i64p == nil || *i64p != i64 {
		t.Errorf("Int64() failed")
	}

	// Test Float64
	f := 3.14
	fp := Float64(f)
	if fp == nil || *fp != f {
		t.Errorf("Float64() failed")
	}

	// Test Bool
	b := true
	bp := Bool(b)
	if bp == nil || *bp != b {
		t.Errorf("Bool() failed")
	}
}

func TestNewWithEnvironmentVariable(t *testing.T) {
	// Save original env var
	originalAPIKey := os.Getenv("KEYWORDSAI_API_KEY")
	defer func() {
		if originalAPIKey != "" {
			os.Setenv("KEYWORDSAI_API_KEY", originalAPIKey)
		} else {
			os.Unsetenv("KEYWORDSAI_API_KEY")
		}
	}()

	// Set test API key
	os.Setenv("KEYWORDSAI_API_KEY", "env-api-key")

	sdk := New()
	if sdk.Client.APIKey() != "env-api-key" {
		t.Errorf("expected API key from env var, got %s", sdk.Client.APIKey())
	}
}

func TestNewWithCustomHTTPClient(t *testing.T) {
	customClient := &http.Client{
		Timeout: 60 * time.Second,
	}

	sdk := New(
		"test-api-key",
		client.WithHTTPClient(customClient),
	)

	if sdk.Client.HTTPClient() != customClient {
		t.Error("expected custom HTTP client to be used")
	}
}