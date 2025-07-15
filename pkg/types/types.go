package types

import (
	"time"
)

type Message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
	Name    *string     `json:"name,omitempty"`
}

type ToolCall struct {
	ID       string         `json:"id"`
	Type     string         `json:"type"`
	Function FunctionCall   `json:"function"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens,omitempty"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
	TotalTokens      int `json:"total_tokens,omitempty"`
}

type CustomerParams struct {
	CustomerIdentifier string                 `json:"customer_identifier"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

type RequestLog struct {
	Model                  string                 `json:"model"`
	PromptMessages        []Message              `json:"prompt_messages"`
	CompletionMessage     *Message               `json:"completion_message,omitempty"`
	CustomerParams        *CustomerParams        `json:"customer_params,omitempty"`
	PromptTokens          *int                   `json:"prompt_tokens,omitempty"`
	CompletionTokens      *int                   `json:"completion_tokens,omitempty"`
	Cost                  *float64               `json:"cost,omitempty"`
	Latency               *int                   `json:"latency,omitempty"`
	Failed                *bool                  `json:"failed,omitempty"`
	StatusCode            *int                   `json:"status_code,omitempty"`
	Error                 *string                `json:"error,omitempty"`
	Timestamp             *time.Time             `json:"timestamp,omitempty"`
	Usage                 *Usage                 `json:"usage,omitempty"`
	ToolCalls             []ToolCall             `json:"tool_calls,omitempty"`
	Metadata              map[string]interface{} `json:"metadata,omitempty"`
	ExtraHeaders          map[string]string      `json:"extra_headers,omitempty"`
	RequestParams         map[string]interface{} `json:"request_params,omitempty"`
	Provider              *string                `json:"provider,omitempty"`
	Stream                *bool                  `json:"stream,omitempty"`
	Category              *string                `json:"category,omitempty"`
	Tags                  []string               `json:"tags,omitempty"`
}

type BatchRequestLogsPayload struct {
	Logs []RequestLog `json:"logs"`
}

type LogFilter struct {
	Model              *string    `json:"model,omitempty"`
	Failed             *bool      `json:"failed,omitempty"`
	Category           *string    `json:"category,omitempty"`
	CustomerIdentifier *string    `json:"customer_identifier,omitempty"`
	StartTime          *time.Time `json:"start_time,omitempty"`
	EndTime            *time.Time `json:"end_time,omitempty"`
	Tags               []string   `json:"tags,omitempty"`
	Limit              *int       `json:"limit,omitempty"`
	Offset             *int       `json:"offset,omitempty"`
}

type LogsResponse struct {
	Logs       []RequestLog `json:"logs"`
	TotalCount int          `json:"total_count"`
	NextOffset *int         `json:"next_offset,omitempty"`
}

type Thread struct {
	ID                 string                 `json:"id"`
	CustomerIdentifier string                 `json:"customer_identifier"`
	Messages          []Message              `json:"messages"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

type Prompt struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PromptVersion struct {
	ID             string                 `json:"id"`
	PromptID       string                 `json:"prompt_id"`
	Version        int                    `json:"version"`
	Name           string                 `json:"name"`
	Template       string                 `json:"template"`
	Model          *string                `json:"model,omitempty"`
	Parameters     map[string]interface{} `json:"parameters,omitempty"`
	IsActive       bool                   `json:"is_active"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

type Model struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Provider        string   `json:"provider"`
	InputCost       float64  `json:"input_cost"`
	OutputCost      float64  `json:"output_cost"`
	MaxTokens       int      `json:"max_tokens"`
	ContextWindow   int      `json:"context_window"`
	SupportedModes  []string `json:"supported_modes"`
	IsAvailable     bool     `json:"is_available"`
}

type TemporaryKey struct {
	ID               string                 `json:"id"`
	Key              string                 `json:"key"`
	Name             *string                `json:"name,omitempty"`
	ExpiresAt        time.Time              `json:"expires_at"`
	CreatedAt        time.Time              `json:"created_at"`
	IsActive         bool                   `json:"is_active"`
	UsageLimit       *int                   `json:"usage_limit,omitempty"`
	UsageCount       int                    `json:"usage_count"`
	AllowedModels    []string               `json:"allowed_models,omitempty"`
	AllowedEndpoints []string               `json:"allowed_endpoints,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

type TTSRequest struct {
	Model          string  `json:"model"`
	Input          string  `json:"input"`
	Voice          string  `json:"voice"`
	ResponseFormat *string `json:"response_format,omitempty"`
	Speed          *float64 `json:"speed,omitempty"`
}

type STTRequest struct {
	Model           string  `json:"model"`
	ResponseFormat  *string `json:"response_format,omitempty"`
	Language        *string `json:"language,omitempty"`
	Temperature     *float64 `json:"temperature,omitempty"`
	Prompt          *string `json:"prompt,omitempty"`
}

type STTResponse struct {
	Text     string                 `json:"text"`
	Language string                 `json:"language,omitempty"`
	Duration float64                `json:"duration,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type EmbeddingRequest struct {
	Model          string      `json:"model"`
	Input          interface{} `json:"input"`
	EncodingFormat *string     `json:"encoding_format,omitempty"`
	Dimensions     *int        `json:"dimensions,omitempty"`
}

type EmbeddingData struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

type EmbeddingResponse struct {
	Object string          `json:"object"`
	Data   []EmbeddingData `json:"data"`
	Model  string          `json:"model"`
	Usage  Usage           `json:"usage"`
}