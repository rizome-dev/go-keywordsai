package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	defaultBaseURL = "https://api.keywordsai.co"
	defaultTimeout = 30 * time.Second
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

type Option func(*Client)

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// New creates a new KeywordsAI client.
// Usage:
//   client.New()                    // Uses KEYWORDSAI_API_KEY env var
//   client.New("api-key")           // Uses provided API key
//   client.New(client.WithTimeout(30*time.Second))  // With options only
//   client.New("api-key", client.WithTimeout(30*time.Second))  // With API key and options
func New(params ...interface{}) *Client {
	var apiKey string
	var opts []Option

	// Parse parameters - first string is API key, rest are options
	for _, param := range params {
		switch v := param.(type) {
		case string:
			if apiKey == "" {
				apiKey = v
			}
		case Option:
			opts = append(opts, v)
		case func(*Client):
			opts = append(opts, Option(v))
		}
	}

	c := &Client{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		baseURL: defaultBaseURL,
		apiKey:  apiKey,
	}

	// Use environment variable if no API key provided
	if c.apiKey == "" {
		c.apiKey = os.Getenv("KEYWORDSAI_API_KEY")
	}

	// Use environment variable for base URL if not overridden
	if baseURL := os.Getenv("KEYWORDS_BASE_URL"); baseURL != "" && c.baseURL == defaultBaseURL {
		c.baseURL = baseURL
	}

	// Apply options
	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) do(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

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

func (c *Client) Get(ctx context.Context, path string, result interface{}) error {
	return c.do(ctx, http.MethodGet, path, nil, result)
}

func (c *Client) GetWithQuery(ctx context.Context, path string, query interface{}, result interface{}) error {
	if query != nil {
		queryString := BuildQueryString(query)
		if queryString != "" {
			path = path + "?" + queryString
		}
	}
	return c.do(ctx, http.MethodGet, path, nil, result)
}

func (c *Client) Post(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.do(ctx, http.MethodPost, path, body, result)
}

func (c *Client) Put(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.do(ctx, http.MethodPut, path, body, result)
}

func (c *Client) Patch(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.do(ctx, http.MethodPatch, path, body, result)
}

func (c *Client) Delete(ctx context.Context, path string, result interface{}) error {
	return c.do(ctx, http.MethodDelete, path, nil, result)
}

func (c *Client) BaseURL() string {
	return c.baseURL
}

func (c *Client) APIKey() string {
	return c.apiKey
}

func (c *Client) HTTPClient() *http.Client {
	return c.httpClient
}