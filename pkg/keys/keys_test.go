package keys

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rizome-dev/go-keywordsai/pkg/client"
	"github.com/rizome-dev/go-keywordsai/pkg/types"
)

func TestCreate(t *testing.T) {
	expiresAt := time.Now().Add(24 * time.Hour)
	expectedKey := types.TemporaryKey{
		ID:        "key-123",
		Key:       "tmp_abc123xyz",
		Name:      stringPtr("Test Key"),
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		IsActive:  true,
		UsageLimit: intPtr(1000),
		UsageCount: 0,
		AllowedModels: []string{"gpt-4", "gpt-3.5-turbo"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/temporary-keys" {
			t.Errorf("Expected path /api/temporary-keys, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		var req CreateKeyRequest
		json.NewDecoder(r.Body).Decode(&req)

		if *req.Name != "Test Key" {
			t.Errorf("Expected name 'Test Key', got %s", *req.Name)
		}

		json.NewEncoder(w).Encode(expectedKey)
	}))
	defer server.Close()

	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := NewService(c)

	req := &CreateKeyRequest{
		Name:          stringPtr("Test Key"),
		ExpiresAt:     expiresAt,
		UsageLimit:    intPtr(1000),
		AllowedModels: []string{"gpt-4", "gpt-3.5-turbo"},
	}

	result, err := s.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if result.ID != expectedKey.ID {
		t.Errorf("Expected ID %s, got %s", expectedKey.ID, result.ID)
	}
	if result.Key != expectedKey.Key {
		t.Errorf("Expected key %s, got %s", expectedKey.Key, result.Key)
	}
}

func TestList(t *testing.T) {
	expectedKeys := []types.TemporaryKey{
		{
			ID:        "key-1",
			Key:       "tmp_abc123",
			IsActive:  true,
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(24 * time.Hour),
		},
		{
			ID:        "key-2",
			Key:       "tmp_xyz789",
			IsActive:  false,
			CreatedAt: time.Now().Add(-48 * time.Hour),
			ExpiresAt: time.Now().Add(-24 * time.Hour),
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/temporary-keys" {
			t.Errorf("Expected path /api/temporary-keys, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(expectedKeys)
	}))
	defer server.Close()

	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := NewService(c)

	result, err := s.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(result))
	}
}

func TestGet(t *testing.T) {
	expectedKey := types.TemporaryKey{
		ID:         "key-123",
		Key:        "tmp_abc123xyz",
		IsActive:   true,
		UsageCount: 42,
		UsageLimit: intPtr(1000),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/temporary-keys/key-123" {
			t.Errorf("Expected path /api/temporary-keys/key-123, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(expectedKey)
	}))
	defer server.Close()

	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := NewService(c)

	result, err := s.Get(context.Background(), "key-123")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if result.ID != expectedKey.ID {
		t.Errorf("Expected ID %s, got %s", expectedKey.ID, result.ID)
	}
	if result.UsageCount != expectedKey.UsageCount {
		t.Errorf("Expected usage count %d, got %d", expectedKey.UsageCount, result.UsageCount)
	}
}

func TestUpdate(t *testing.T) {
	updatedKey := types.TemporaryKey{
		ID:       "key-123",
		Key:      "tmp_abc123xyz",
		Name:     stringPtr("Updated Key"),
		IsActive: false,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/temporary-keys/key-123" {
			t.Errorf("Expected path /api/temporary-keys/key-123, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH method, got %s", r.Method)
		}

		var updates map[string]interface{}
		json.NewDecoder(r.Body).Decode(&updates)

		if updates["is_active"] != false {
			t.Errorf("Expected is_active false, got %v", updates["is_active"])
		}

		json.NewEncoder(w).Encode(updatedKey)
	}))
	defer server.Close()

	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := NewService(c)

	updates := map[string]interface{}{
		"name":      "Updated Key",
		"is_active": false,
	}

	result, err := s.Update(context.Background(), "key-123", updates)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if *result.Name != "Updated Key" {
		t.Errorf("Expected name 'Updated Key', got %s", *result.Name)
	}
	if result.IsActive != false {
		t.Errorf("Expected is_active false, got %v", result.IsActive)
	}
}

func TestDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/temporary-keys/key-123" {
			t.Errorf("Expected path /api/temporary-keys/key-123, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := NewService(c)

	err := s.Delete(context.Background(), "key-123")
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}