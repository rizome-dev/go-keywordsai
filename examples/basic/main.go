package main

import (
	"context"
	"fmt"
	"log"

	"github.com/rizome-dev/go-keywordsai/pkg/client"
	"github.com/rizome-dev/go-keywordsai/pkg/logs"
	"github.com/rizome-dev/go-keywordsai/pkg/models"
	"github.com/rizome-dev/go-keywordsai/pkg/prompts"
	"github.com/rizome-dev/go-keywordsai/pkg/types"
)

func main() {
	// Create a new client (uses KEYWORDSAI_API_KEY env var automatically)
	c := client.New()
	logsService := logs.NewService(c)
	modelsService := models.NewService(c)
	promptsService := prompts.NewService(c)

	ctx := context.Background()

	// Example 1: Log a single request
	fmt.Println("=== Logging a single request ===")
	logRequest(ctx, logsService)

	// Example 2: List available models
	fmt.Println("\n=== Listing available models ===")
	listModels(ctx, modelsService)

	// Example 3: Create and manage prompts
	fmt.Println("\n=== Managing prompts ===")
	managePrompts(ctx, promptsService)
}

func logRequest(ctx context.Context, logsService *logs.Service) {
	requestLog := &types.RequestLog{
		Model: "gpt-4",
		PromptMessages: []types.Message{
			{
				Role:    "system",
				Content: "You are a helpful assistant.",
			},
			{
				Role:    "user",
				Content: "What is the capital of France?",
			},
		},
		CompletionMessage: &types.Message{
			Role:    "assistant",
			Content: "The capital of France is Paris.",
		},
		PromptTokens:     intPtr(25),
		CompletionTokens: intPtr(8),
		Cost:             floatPtr(0.0033),
		CustomerParams: &types.CustomerParams{
			CustomerIdentifier: "user-123",
			Metadata: map[string]interface{}{
				"session_id": "abc-123",
			},
		},
		Tags: []string{"geography", "factual"},
	}

	err := logsService.Create(ctx, requestLog)
	if err != nil {
		log.Printf("Failed to create log: %v", err)
	} else {
		fmt.Println("Successfully logged request")
	}
}

func listModels(ctx context.Context, modelsService *models.Service) {
	models, err := modelsService.List(ctx)
	if err != nil {
		log.Printf("Failed to list models: %v", err)
		return
	}

	fmt.Printf("Found %d models:\n", len(models))
	for _, model := range models {
		fmt.Printf("  - %s (%s): Input=$%.4f/1K, Output=$%.4f/1K\n",
			model.ID, model.Provider, model.InputCost*1000, model.OutputCost*1000)
	}
}

func managePrompts(ctx context.Context, promptsService *prompts.Service) {
	// Create a prompt
	description := "A prompt for answering geography questions"
	prompt, err := promptsService.Create(ctx, "Geography Assistant", &description)
	if err != nil {
		log.Printf("Failed to create prompt: %v", err)
		return
	}
	fmt.Printf("Created prompt: %s (ID: %s)\n", prompt.Name, prompt.ID)

	// Create a version
	version := &types.PromptVersion{
		Name:     "Version 1.0",
		Template: "You are a geography expert. Answer the following question: {{question}}",
		Model:    stringPtr("gpt-4"),
		Parameters: map[string]interface{}{
			"temperature": 0.7,
			"max_tokens":  200,
		},
		IsActive: true,
	}

	createdVersion, err := promptsService.CreateVersion(ctx, prompt.ID, version)
	if err != nil {
		log.Printf("Failed to create prompt version: %v", err)
		return
	}
	fmt.Printf("Created prompt version: %s (ID: %s)\n", createdVersion.Name, createdVersion.ID)

	// List all prompts
	prompts, err := promptsService.List(ctx)
	if err != nil {
		log.Printf("Failed to list prompts: %v", err)
		return
	}
	fmt.Printf("\nTotal prompts: %d\n", len(prompts))
}

// For convenience, you can also import utils package for pointer helpers:
// import "github.com/rizome-dev/go-keywordsai/pkg/utils"
// Then use utils.String("value"), utils.Int(123), etc.

func intPtr(i int) *int {
	return &i
}

func floatPtr(f float64) *float64 {
	return &f
}

func stringPtr(s string) *string {
	return &s
}