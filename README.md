# KeywordsAI Go SDK

[![GoDoc](https://pkg.go.dev/badge/github.com/rizome-dev/go-keywordsai)](https://pkg.go.dev/github.com/rizome-dev/go-keywordsai)
[![Go Report Card](https://goreportcard.com/badge/github.com/rizome-dev/go-keywordsai)](https://goreportcard.com/report/github.com/rizome-dev/go-keywordsai)
[![CI](https://github.com/rizome-dev/go-keywordsai/actions/workflows/ci.yml/badge.svg)](https://github.com/rizome-dev/go-keywordsai/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A full-featured Go SDK for [KeywordsAI](https://keywordsai.co), providing comprehensive access to all KeywordsAI API endpoints for LLM observability, logging, and management.

built by [rizome labs](https://rizome.dev) | contact: [hi@rizome.dev](mailto:hi@rizome.dev)

## Installation

```bash
go get github.com/rizome-dev/go-keywordsai
```

## Quick Start

### Option 1: All-in-One SDK (Recommended)

```go
package main

import (
    "context"
    "log"
    
    "github.com/rizome-dev/go-keywordsai"
)

func main() {
    // Create SDK with all services (uses KEYWORDSAI_API_KEY env var)
    sdk := keywordsai.New()
    
    ctx := context.Background()
    
    // Use any service directly
    err := sdk.Logs.Create(ctx, &keywordsai.RequestLog{
        Model: "gpt-4",
        PromptMessages: []keywordsai.Message{
            {Role: "user", Content: "Hello!"},
        },
        CompletionMessage: &keywordsai.Message{
            Role: "assistant",
            Content: "Hi! How can I help you?",
        },
    })
    
    if err != nil {
        log.Fatal(err)
    }
    
    // Access other services
    models, _ := sdk.Models.List(ctx)
    prompts, _ := sdk.Prompts.List(ctx)
}
```

### Option 2: Individual Services

```go
package main

import (
    "context"
    "log"
    
    "github.com/rizome-dev/go-keywordsai/pkg/client"
    "github.com/rizome-dev/go-keywordsai/pkg/logs"
    "github.com/rizome-dev/go-keywordsai/pkg/types"
)

func main() {
    // Create client and only the services you need
    c := client.New()
    logsService := logs.NewService(c)
    
    ctx := context.Background()
    
    // Log a request
    err := logsService.Create(ctx, &types.RequestLog{
        Model: "gpt-4",
        PromptMessages: []types.Message{
            {Role: "user", Content: "Hello!"},
        },
        CompletionMessage: &types.Message{
            Role: "assistant",
            Content: "Hi! How can I help you?",
        },
    })
    
    if err != nil {
        log.Fatal(err)
    }
}
```

## Configuration

### Environment Variables

- `KEYWORDSAI_API_KEY` - Your KeywordsAI API key
- `KEYWORDS_BASE_URL` - Custom API base URL (default: https://api.keywordsai.co)

### Client Options

```go
import "net/http"
import "time"

// Custom HTTP client
httpClient := &http.Client{
    Timeout: 60 * time.Second,
}

// All in one New method - supports all combinations
c := client.New("api-key",
    client.WithHTTPClient(httpClient),
    client.WithBaseURL("https://custom.api.url"),
    client.WithTimeout(30 * time.Second),
)

// Or just with options (uses KEYWORDSAI_API_KEY env var)
c := client.New(
    client.WithTimeout(30 * time.Second),
    client.WithBaseURL("https://custom.api.url"),
)

// Then create services
logsService := logs.NewService(c)
```

### Utility Functions

The SDK provides helper functions for creating pointers to basic types:

```go
// Using the all-in-one SDK package
import "github.com/rizome-dev/go-keywordsai"

// Utilities are available directly
namePtr := keywordsai.String("example")
count := keywordsai.Int(42)
cost := keywordsai.Float64(0.0033)
enabled := keywordsai.Bool(true)

// Or import utils package directly
import "github.com/rizome-dev/go-keywordsai/pkg/utils"

namePtr := utils.String("example")
count := utils.Int(42)
```

## API Reference

### Logging

#### Single Request Logging

```go
import (
    "github.com/rizome-dev/go-keywordsai/pkg/client"
    "github.com/rizome-dev/go-keywordsai/pkg/logs"
    "github.com/rizome-dev/go-keywordsai/pkg/types"
)

c := client.New() // Uses KEYWORDSAI_API_KEY env var
logsService := logs.NewService(c)

err := logsService.Create(ctx, &types.RequestLog{
    Model: "gpt-4",
    PromptMessages: []types.Message{
        {Role: "system", Content: "You are a helpful assistant."},
        {Role: "user", Content: "What is the capital of France?"},
    },
    CompletionMessage: &types.Message{
        Role: "assistant",
        Content: "The capital of France is Paris.",
    },
    PromptTokens: intPtr(25),
    CompletionTokens: intPtr(8),
    Cost: floatPtr(0.0033),
    CustomerParams: &types.CustomerParams{
        CustomerIdentifier: "user-123",
        Metadata: map[string]interface{}{
            "session_id": "abc-123",
        },
    },
    Tags: []string{"geography", "factual"},
})
```

#### Batch Request Logging

```go
logs := []types.RequestLog{
    // ... multiple request logs
}

// Maximum 5000 logs per batch
err := logsService.BatchCreate(ctx, logs)
```

#### Query Logs

```go
filter := &types.LogFilter{
    Model: stringPtr("gpt-4"),
    Category: stringPtr("technical"),
    StartTime: &startTime,
    EndTime: &endTime,
    Limit: intPtr(100),
}

response, err := logsService.List(ctx, filter)
```

### Prompt Management

#### Create Prompt

```go
import (
    "github.com/rizome-dev/go-keywordsai/pkg/client"
    "github.com/rizome-dev/go-keywordsai/pkg/prompts"
)

c := client.New() // Uses KEYWORDSAI_API_KEY env var
promptsService := prompts.NewService(c)

description := "A helpful assistant prompt"
prompt, err := promptsService.Create(ctx, "My Assistant", &description)
```

#### Create Prompt Version

```go
version := &types.PromptVersion{
    Name: "v1.0",
    Template: "You are a {{role}}. {{task}}",
    Model: stringPtr("gpt-4"),
    Parameters: map[string]interface{}{
        "temperature": 0.7,
        "max_tokens": 200,
    },
    IsActive: true,
}

createdVersion, err := promptsService.CreateVersion(ctx, prompt.ID, version)
```

### Integrations

#### Text-to-Speech

```go
import (
    "github.com/rizome-dev/go-keywordsai/pkg/client"
    "github.com/rizome-dev/go-keywordsai/pkg/integrations"
    "github.com/rizome-dev/go-keywordsai/pkg/types"
)

c := client.New() // Uses KEYWORDSAI_API_KEY env var
integrationsService := integrations.NewService(c)

ttsRequest := &types.TTSRequest{
    Model: "tts-1",
    Input: "Hello, world!",
    Voice: "alloy",
    ResponseFormat: stringPtr("mp3"),
    Speed: floatPtr(1.0),
}

audioData, err := integrationsService.TextToSpeech(ctx, ttsRequest)
// audioData contains the audio file bytes
```

#### Speech-to-Text

```go
audioData, _ := os.ReadFile("speech.wav")

sttRequest := &types.STTRequest{
    Model: "whisper-1",
    ResponseFormat: stringPtr("json"),
    Language: stringPtr("en"),
}

result, err := integrationsService.SpeechToText(ctx, audioData, sttRequest)
fmt.Println(result.Text)
```

#### Embeddings

```go
embeddingRequest := &types.EmbeddingRequest{
    Model: "text-embedding-ada-002",
    Input: "Keywords AI is great!",
    EncodingFormat: stringPtr("float"),
}

result, err := integrationsService.CreateEmbeddings(ctx, embeddingRequest)
```

### Temporary API Keys

```go
import (
    "time"
    "github.com/rizome-dev/go-keywordsai/pkg/client"
    "github.com/rizome-dev/go-keywordsai/pkg/keys"
)

c := client.New() // Uses KEYWORDSAI_API_KEY env var
keysService := keys.NewService(c)

// Create a temporary key
keyRequest := &keys.CreateKeyRequest{
    Name: stringPtr("Dev Key"),
    ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
    UsageLimit: intPtr(10000),
    AllowedModels: []string{"gpt-3.5-turbo", "gpt-4"},
}

key, err := keysService.Create(ctx, keyRequest)

// List keys
keys, err := keysService.List(ctx)

// Update key
updates := map[string]interface{}{
    "name": "Updated Key Name",
    "is_active": false,
}
updatedKey, err := keysService.Update(ctx, key.ID, updates)

// Delete key
err = keysService.Delete(ctx, key.ID)
```

### Models

```go
import (
    "github.com/rizome-dev/go-keywordsai/pkg/client"
    "github.com/rizome-dev/go-keywordsai/pkg/models"
)

c := client.New() // Uses KEYWORDSAI_API_KEY env var
modelsService := models.NewService(c)

// List available models
models, err := modelsService.List(ctx)
for _, model := range models {
    fmt.Printf("%s: $%.4f/1K input, $%.4f/1K output\n", 
        model.ID, model.InputCost*1000, model.OutputCost*1000)
}
```

## Examples

See the [examples](./examples) directory for complete working examples:

- [Basic Usage](./examples/basic/main.go) - Simple logging and prompt management
- [Batch Logging](./examples/batch-logging/main.go) - Batch operations and querying
- [Integrations](./examples/integrations/main.go) - TTS, STT, and embeddings
- [API Keys](./examples/api-keys/main.go) - Temporary key management

## Development

```bash
# Install tools
make setup

# Run tests
make test

# Run tests with coverage
go test -cover ./...

# Lint
make lint

# All checks
make ci
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This SDK is distributed under the MIT License. See [LICENSE](LICENSE) file for details.

## Support

- Email: hi@rizome.dev
- KeywordsAI Documentation: https://docs.keywordsai.co
- Issues: https://github.com/rizome-dev/go-keywordsai/issues
