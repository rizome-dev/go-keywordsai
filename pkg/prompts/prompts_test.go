package prompts

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
	expectedPrompt := types.Prompt{
		ID:          "prompt-123",
		Name:        "Test Prompt",
		Description: ptr("A test prompt"),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/prompts/" {
			t.Errorf("Expected path /api/prompts/, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		var payload map[string]interface{}
		json.NewDecoder(r.Body).Decode(&payload)

		if payload["name"] != "Test Prompt" {
			t.Errorf("Expected name 'Test Prompt', got %v", payload["name"])
		}

		json.NewEncoder(w).Encode(expectedPrompt)
	}))
	defer server.Close()

	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := NewService(c)

	desc := "A test prompt"
	result, err := s.Create(context.Background(), "Test Prompt", &desc)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if result.ID != expectedPrompt.ID {
		t.Errorf("Expected ID %s, got %s", expectedPrompt.ID, result.ID)
	}
}

func TestList(t *testing.T) {
	expectedPrompts := []types.Prompt{
		{
			ID:   "prompt-1",
			Name: "Prompt 1",
		},
		{
			ID:   "prompt-2",
			Name: "Prompt 2",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/prompts/" {
			t.Errorf("Expected path /api/prompts/, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(expectedPrompts)
	}))
	defer server.Close()

	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := NewService(c)

	result, err := s.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 prompts, got %d", len(result))
	}
}

func TestGet(t *testing.T) {
	expectedPrompt := types.Prompt{
		ID:   "prompt-123",
		Name: "Test Prompt",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/prompts/prompt-123" {
			t.Errorf("Expected path /api/prompts/prompt-123, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(expectedPrompt)
	}))
	defer server.Close()

	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := NewService(c)

	result, err := s.Get(context.Background(), "prompt-123")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if result.ID != expectedPrompt.ID {
		t.Errorf("Expected ID %s, got %s", expectedPrompt.ID, result.ID)
	}
}

func TestUpdate(t *testing.T) {
	updatedPrompt := types.Prompt{
		ID:          "prompt-123",
		Name:        "Updated Prompt",
		Description: ptr("Updated description"),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/prompts/prompt-123" {
			t.Errorf("Expected path /api/prompts/prompt-123, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH method, got %s", r.Method)
		}
		json.NewEncoder(w).Encode(updatedPrompt)
	}))
	defer server.Close()

	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := NewService(c)

	updates := map[string]interface{}{
		"name":        "Updated Prompt",
		"description": "Updated description",
	}

	result, err := s.Update(context.Background(), "prompt-123", updates)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if result.Name != updatedPrompt.Name {
		t.Errorf("Expected name %s, got %s", updatedPrompt.Name, result.Name)
	}
}

func TestDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/prompts/prompt-123" {
			t.Errorf("Expected path /api/prompts/prompt-123, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := NewService(c)

	err := s.Delete(context.Background(), "prompt-123")
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
}

func TestCreateVersion(t *testing.T) {
	expectedVersion := types.PromptVersion{
		ID:       "version-123",
		PromptID: "prompt-123",
		Version:  1,
		Name:     "Version 1",
		Template: "Hello {{name}}",
		IsActive: true,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/prompts/prompt-123/versions" {
			t.Errorf("Expected path /api/prompts/prompt-123/versions, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(expectedVersion)
	}))
	defer server.Close()

	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := NewService(c)

	version := &types.PromptVersion{
		Name:     "Version 1",
		Template: "Hello {{name}}",
		IsActive: true,
	}

	result, err := s.CreateVersion(context.Background(), "prompt-123", version)
	if err != nil {
		t.Fatalf("CreateVersion() error = %v", err)
	}

	if result.ID != expectedVersion.ID {
		t.Errorf("Expected ID %s, got %s", expectedVersion.ID, result.ID)
	}
}

func TestListVersions(t *testing.T) {
	expectedVersions := []types.PromptVersion{
		{
			ID:       "version-1",
			PromptID: "prompt-123",
			Version:  1,
		},
		{
			ID:       "version-2",
			PromptID: "prompt-123",
			Version:  2,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/prompts/prompt-123/versions" {
			t.Errorf("Expected path /api/prompts/prompt-123/versions, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(expectedVersions)
	}))
	defer server.Close()

	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := NewService(c)

	result, err := s.ListVersions(context.Background(), "prompt-123")
	if err != nil {
		t.Fatalf("ListVersions() error = %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 versions, got %d", len(result))
	}
}

func ptr(s string) *string {
	return &s
}

func TestGetVersion(t *testing.T) {
	expectedVersion := &types.PromptVersion{
		ID:       "version-123",
		PromptID: "prompt-123",
		Version:  1,
		Name:     "v1.0",
		Template: "Test template",
		IsActive: true,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/prompts/prompt-123/versions/version-123" {
			t.Errorf("Expected path /api/prompts/prompt-123/versions/version-123, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(expectedVersion)
	}))
	defer server.Close()

	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := NewService(c)

	result, err := s.GetVersion(context.Background(), "prompt-123", "version-123")
	if err != nil {
		t.Fatalf("GetVersion() error = %v", err)
	}

	if result.ID != expectedVersion.ID {
		t.Errorf("Expected ID %s, got %s", expectedVersion.ID, result.ID)
	}
}

func TestUpdateVersion(t *testing.T) {
	updates := map[string]interface{}{
		"name":     "Updated Version",
		"template": "Updated template",
	}

	expectedVersion := &types.PromptVersion{
		ID:       "version-123",
		PromptID: "prompt-123",
		Version:  1,
		Name:     "Updated Version",
		Template: "Updated template",
		IsActive: true,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH method, got %s", r.Method)
		}
		if r.URL.Path != "/api/prompts/prompt-123/versions/version-123" {
			t.Errorf("Expected path /api/prompts/prompt-123/versions/version-123, got %s", r.URL.Path)
		}

		var receivedUpdates map[string]interface{}
		json.NewDecoder(r.Body).Decode(&receivedUpdates)
		
		if receivedUpdates["name"] != updates["name"] {
			t.Errorf("Expected name update %v, got %v", updates["name"], receivedUpdates["name"])
		}

		json.NewEncoder(w).Encode(expectedVersion)
	}))
	defer server.Close()

	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := NewService(c)

	result, err := s.UpdateVersion(context.Background(), "prompt-123", "version-123", updates)
	if err != nil {
		t.Fatalf("UpdateVersion() error = %v", err)
	}

	if result.Name != expectedVersion.Name {
		t.Errorf("Expected name %s, got %s", expectedVersion.Name, result.Name)
	}
}

func TestDeleteVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}
		if r.URL.Path != "/api/prompts/prompt-123/versions/version-123" {
			t.Errorf("Expected path /api/prompts/prompt-123/versions/version-123, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := NewService(c)

	err := s.DeleteVersion(context.Background(), "prompt-123", "version-123")
	if err != nil {
		t.Fatalf("DeleteVersion() error = %v", err)
	}
}