package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rizome-dev/go-keywordsai/pkg/client"
	"github.com/rizome-dev/go-keywordsai/pkg/keys"
)

func main() {
	c := client.New() // Uses KEYWORDSAI_API_KEY env var automatically
	keysService := keys.NewService(c)
	ctx := context.Background()

	// Example 1: Create a temporary API key
	fmt.Println("=== Creating Temporary API Key ===")
	keyID := createTemporaryKey(ctx, keysService)

	// Example 2: List all temporary keys
	fmt.Println("\n=== Listing Temporary Keys ===")
	listKeys(ctx, keysService)

	// Example 3: Update a key
	if keyID != "" {
		fmt.Println("\n=== Updating Temporary Key ===")
		updateKey(ctx, keysService, keyID)
	}

	// Example 4: Get key details
	if keyID != "" {
		fmt.Println("\n=== Getting Key Details ===")
		getKeyDetails(ctx, keysService, keyID)
	}

	// Example 5: Delete a key
	if keyID != "" {
		fmt.Println("\n=== Deleting Temporary Key ===")
		deleteKey(ctx, keysService, keyID)
	}
}

func createTemporaryKey(ctx context.Context, keysService *keys.Service) string {
	// Create a key that expires in 30 days
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	
	keyRequest := &keys.CreateKeyRequest{
		Name:      stringPtr("Development API Key"),
		ExpiresAt: expiresAt,
		UsageLimit: intPtr(10000), // 10,000 request limit
		AllowedModels: []string{
			"gpt-3.5-turbo",
			"gpt-4",
			"claude-3-opus",
			"claude-3-sonnet",
		},
		AllowedEndpoints: []string{
			"/api/request-logs/create/",
			"/api/embeddings",
		},
		Metadata: map[string]interface{}{
			"environment": "development",
			"team":        "engineering",
			"project":     "demo-app",
		},
	}

	key, err := keysService.Create(ctx, keyRequest)
	if err != nil {
		log.Printf("Failed to create temporary key: %v", err)
		return ""
	}

	fmt.Printf("Created temporary key:\n")
	fmt.Printf("  ID: %s\n", key.ID)
	fmt.Printf("  Key: %s\n", key.Key)
	fmt.Printf("  Expires: %s\n", key.ExpiresAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("  Usage Limit: %d\n", *key.UsageLimit)
	
	return key.ID
}

func listKeys(ctx context.Context, keysService *keys.Service) {
	keys, err := keysService.List(ctx)
	if err != nil {
		log.Printf("Failed to list keys: %v", err)
		return
	}

	fmt.Printf("Found %d temporary keys:\n", len(keys))
	for _, key := range keys {
		status := "Active"
		if !key.IsActive {
			status = "Inactive"
		}
		
		name := "Unnamed"
		if key.Name != nil {
			name = *key.Name
		}

		fmt.Printf("  - %s (%s): %s, Usage: %d",
			name, key.ID, status, key.UsageCount)
		
		if key.UsageLimit != nil {
			fmt.Printf("/%d", *key.UsageLimit)
		}
		
		fmt.Printf(", Expires: %s\n", key.ExpiresAt.Format("2006-01-02"))
	}
}

func updateKey(ctx context.Context, keysService *keys.Service, keyID string) {
	updates := map[string]interface{}{
		"name": "Updated Development Key",
		"metadata": map[string]interface{}{
			"environment": "staging",
			"team":        "engineering",
			"project":     "demo-app",
			"updated_at":  time.Now().Format(time.RFC3339),
		},
		"usage_limit": 20000, // Increase limit to 20,000
	}

	updatedKey, err := keysService.Update(ctx, keyID, updates)
	if err != nil {
		log.Printf("Failed to update key: %v", err)
		return
	}

	fmt.Printf("Updated key:\n")
	fmt.Printf("  Name: %s\n", *updatedKey.Name)
	if updatedKey.UsageLimit != nil {
		fmt.Printf("  New Usage Limit: %d\n", *updatedKey.UsageLimit)
	}
	fmt.Printf("  Metadata: %v\n", updatedKey.Metadata)
}

func getKeyDetails(ctx context.Context, keysService *keys.Service, keyID string) {
	key, err := keysService.Get(ctx, keyID)
	if err != nil {
		log.Printf("Failed to get key details: %v", err)
		return
	}

	fmt.Printf("Key Details:\n")
	fmt.Printf("  ID: %s\n", key.ID)
	if key.Name != nil {
		fmt.Printf("  Name: %s\n", *key.Name)
	}
	fmt.Printf("  Created: %s\n", key.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("  Expires: %s\n", key.ExpiresAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("  Status: %s\n", map[bool]string{true: "Active", false: "Inactive"}[key.IsActive])
	fmt.Printf("  Usage: %d", key.UsageCount)
	if key.UsageLimit != nil {
		fmt.Printf(" / %d", *key.UsageLimit)
	}
	fmt.Println()
	
	if len(key.AllowedModels) > 0 {
		fmt.Printf("  Allowed Models: %v\n", key.AllowedModels)
	}
	if len(key.AllowedEndpoints) > 0 {
		fmt.Printf("  Allowed Endpoints: %v\n", key.AllowedEndpoints)
	}
}

func deleteKey(ctx context.Context, keysService *keys.Service, keyID string) {
	err := keysService.Delete(ctx, keyID)
	if err != nil {
		log.Printf("Failed to delete key: %v", err)
		return
	}
	
	fmt.Printf("Successfully deleted key: %s\n", keyID)
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}