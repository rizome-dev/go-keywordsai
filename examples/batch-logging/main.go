package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rizome-dev/go-keywordsai/pkg/client"
	"github.com/rizome-dev/go-keywordsai/pkg/logs"
	"github.com/rizome-dev/go-keywordsai/pkg/types"
)

func main() {
	c := client.New() // Uses KEYWORDSAI_API_KEY env var automatically
	logsService := logs.NewService(c)
	ctx := context.Background()

	// Example 1: Batch logging multiple requests
	fmt.Println("=== Batch logging requests ===")
	batchLogRequests(ctx, logsService)

	// Example 2: Query and filter logs
	fmt.Println("\n=== Querying logs ===")
	queryLogs(ctx, logsService)

	// Example 3: Working with threads
	fmt.Println("\n=== Working with threads ===")
	workWithThreads(ctx, logsService)
}

func batchLogRequests(ctx context.Context, logsService *logs.Service) {
	// Create multiple logs for batch submission
	now := time.Now()
	logs := []types.RequestLog{
		{
			Model: "gpt-4",
			PromptMessages: []types.Message{
				{Role: "user", Content: "What is machine learning?"},
			},
			CompletionMessage: &types.Message{
				Role:    "assistant",
				Content: "Machine learning is a subset of artificial intelligence...",
			},
			PromptTokens:     intPtr(10),
			CompletionTokens: intPtr(50),
			Cost:             floatPtr(0.0060),
			Timestamp:        &now,
			CustomerParams: &types.CustomerParams{
				CustomerIdentifier: "user-456",
			},
			Tags:     []string{"ml", "educational"},
			Category: stringPtr("technical"),
		},
		{
			Model: "claude-3-opus",
			PromptMessages: []types.Message{
				{Role: "user", Content: "Explain quantum computing"},
			},
			CompletionMessage: &types.Message{
				Role:    "assistant",
				Content: "Quantum computing is a revolutionary approach to computation...",
			},
			PromptTokens:     intPtr(8),
			CompletionTokens: intPtr(100),
			Cost:             floatPtr(0.0158),
			Timestamp:        &now,
			CustomerParams: &types.CustomerParams{
				CustomerIdentifier: "user-789",
			},
			Tags:     []string{"quantum", "physics", "educational"},
			Category: stringPtr("technical"),
		},
		{
			Model: "gpt-3.5-turbo",
			PromptMessages: []types.Message{
				{Role: "user", Content: "Hello!"},
			},
			CompletionMessage: &types.Message{
				Role:    "assistant",
				Content: "Hello! How can I help you today?",
			},
			PromptTokens:     intPtr(3),
			CompletionTokens: intPtr(10),
			Cost:             floatPtr(0.0000195),
			Timestamp:        &now,
			Failed:           boolPtr(false),
			CustomerParams: &types.CustomerParams{
				CustomerIdentifier: "user-456",
			},
			Tags:     []string{"greeting"},
			Category: stringPtr("chat"),
		},
	}

	err := logsService.BatchCreate(ctx, logs)
	if err != nil {
		log.Printf("Failed to batch create logs: %v", err)
		return
	}
	fmt.Printf("Successfully batch logged %d requests\n", len(logs))
}

func queryLogs(ctx context.Context, logsService *logs.Service) {
	// Query logs with filters
	now := time.Now()
	startTime := now.Add(-24 * time.Hour)
	
	filter := &types.LogFilter{
		Category:  stringPtr("technical"),
		StartTime: &startTime,
		EndTime:   &now,
		Limit:     intPtr(10),
	}

	logsResponse, err := logsService.List(ctx, filter)
	if err != nil {
		log.Printf("Failed to query logs: %v", err)
		return
	}

	fmt.Printf("Found %d logs (total: %d)\n", len(logsResponse.Logs), logsResponse.TotalCount)
	for i, logEntry := range logsResponse.Logs {
		fmt.Printf("  %d. Model: %s, Cost: $%.6f, Tags: %v\n",
			i+1, logEntry.Model, *logEntry.Cost, logEntry.Tags)
	}

	// Get a specific log (if we have any)
	if len(logsResponse.Logs) > 0 && logsResponse.Logs[0].Model != "" {
		// Note: This assumes the API returns an ID field, which wasn't in the types
		// In a real implementation, you'd need to get the actual log ID
		fmt.Println("\nNote: Individual log retrieval would require log IDs from the API")
	}
}

func workWithThreads(ctx context.Context, logsService *logs.Service) {
	// List threads for a customer
	threads, err := logsService.ListThreads(ctx, "user-456")
	if err != nil {
		log.Printf("Failed to list threads: %v", err)
		return
	}

	fmt.Printf("Found %d threads for customer user-456\n", len(threads))
	for _, thread := range threads {
		fmt.Printf("  Thread %s: %d messages, created %s\n",
			thread.ID, len(thread.Messages), thread.CreatedAt.Format("2006-01-02 15:04:05"))
	}
}

func intPtr(i int) *int {
	return &i
}

func floatPtr(f float64) *float64 {
	return &f
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}