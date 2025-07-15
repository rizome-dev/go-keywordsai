// Package keywordsai provides a complete SDK for the KeywordsAI API.
package keywordsai

import (
	"github.com/rizome-dev/go-keywordsai/pkg/client"
	"github.com/rizome-dev/go-keywordsai/pkg/integrations"
	"github.com/rizome-dev/go-keywordsai/pkg/keys"
	"github.com/rizome-dev/go-keywordsai/pkg/logs"
	"github.com/rizome-dev/go-keywordsai/pkg/models"
	"github.com/rizome-dev/go-keywordsai/pkg/prompts"
	"github.com/rizome-dev/go-keywordsai/pkg/types"
	"github.com/rizome-dev/go-keywordsai/pkg/utils"
)

// Re-export commonly used types for convenience
type (
	RequestLog        = types.RequestLog
	Message           = types.Message
	LogFilter         = types.LogFilter
	Prompt            = types.Prompt
	PromptVersion     = types.PromptVersion
	Model             = types.Model
	TemporaryKey      = types.TemporaryKey
	TTSRequest        = types.TTSRequest
	STTRequest        = types.STTRequest
	STTResponse       = types.STTResponse
	EmbeddingRequest  = types.EmbeddingRequest
	EmbeddingResponse = types.EmbeddingResponse
)

// Re-export pointer utilities
var (
	String  = utils.String
	Int     = utils.Int
	Int64   = utils.Int64
	Float64 = utils.Float64
	Bool    = utils.Bool
)

// SDK provides a convenient all-in-one client with all services
type SDK struct {
	Client       *client.Client
	Logs         *logs.Service
	Prompts      *prompts.Service
	Models       *models.Service
	Keys         *keys.Service
	Integrations *integrations.Service
}

// New creates a new KeywordsAI SDK instance with all services initialized.
// This is a convenience wrapper for users who want everything in one place.
// Usage:
//   sdk := keywordsai.New()                    // Uses KEYWORDSAI_API_KEY env var
//   sdk := keywordsai.New("api-key")           // Uses provided API key
//   sdk := keywordsai.New(client.WithTimeout(30*time.Second))  // With options
func New(params ...interface{}) *SDK {
	c := client.New(params...)
	
	return &SDK{
		Client:       c,
		Logs:         logs.NewService(c),
		Prompts:      prompts.NewService(c),
		Models:       models.NewService(c),
		Keys:         keys.NewService(c),
		Integrations: integrations.NewService(c),
	}
}