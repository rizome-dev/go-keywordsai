package main

import (
	"context"
	"fmt"
	"log"

	"github.com/rizome-dev/go-keywordsai"
)

func main() {
	// Create SDK with all services
	// Automatically uses KEYWORDSAI_API_KEY env var
	sdk := keywordsai.New()
	ctx := context.Background()

	fmt.Println("KeywordsAI Go SDK Demo")
	fmt.Println("=====================")

	// Example: Log the same request as before, but using the SDK
	fmt.Println("\nLogging a request...")
	err := sdk.Logs.Create(ctx, &keywordsai.RequestLog{
		Model: "gpt-4o-mini",
		PromptMessages: []keywordsai.Message{
			{
				Role:    "user",
				Content: "Hi",
			},
		},
		CompletionMessage: &keywordsai.Message{
			Role:    "assistant",
			Content: "Hi, how can I assist you today?",
		},
	})

	if err != nil {
		log.Printf("Failed to log request: %v", err)
	} else {
		fmt.Println("✓ Successfully logged request")
	}

	// Example: List available models
	fmt.Println("\nFetching available models...")
	models, err := sdk.Models.List(ctx)
	if err != nil {
		log.Printf("Failed to list models: %v", err)
	} else {
		fmt.Printf("✓ Found %d models\n", len(models))
		if len(models) > 0 && len(models) <= 5 {
			for _, model := range models {
				fmt.Printf("  - %s (%s)\n", model.ID, model.Provider)
			}
		} else if len(models) > 5 {
			fmt.Printf("  (showing first 5 models)\n")
			for i := 0; i < 5; i++ {
				fmt.Printf("  - %s (%s)\n", models[i].ID, models[i].Provider)
			}
		}
	}

	fmt.Println("\nDemo complete!")
}
