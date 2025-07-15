package models

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rizome-dev/go-keywordsai/pkg/client"
	"github.com/rizome-dev/go-keywordsai/pkg/types"
)

func TestList(t *testing.T) {
	expectedModels := []types.Model{
		{
			ID:            "gpt-4",
			Name:          "GPT-4",
			Provider:      "openai",
			InputCost:     0.03,
			OutputCost:    0.06,
			MaxTokens:     8192,
			ContextWindow: 8192,
			SupportedModes: []string{"chat", "completion"},
			IsAvailable:   true,
		},
		{
			ID:            "claude-3-opus",
			Name:          "Claude 3 Opus",
			Provider:      "anthropic",
			InputCost:     0.015,
			OutputCost:    0.075,
			MaxTokens:     4096,
			ContextWindow: 200000,
			SupportedModes: []string{"chat"},
			IsAvailable:   true,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/models" {
			t.Errorf("Expected path /api/models, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		
		json.NewEncoder(w).Encode(expectedModels)
	}))
	defer server.Close()

	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := NewService(c)

	result, err := s.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 models, got %d", len(result))
	}

	for i, model := range result {
		if model.ID != expectedModels[i].ID {
			t.Errorf("Expected model ID %s, got %s", expectedModels[i].ID, model.ID)
		}
		if model.Provider != expectedModels[i].Provider {
			t.Errorf("Expected provider %s, got %s", expectedModels[i].Provider, model.Provider)
		}
		if model.InputCost != expectedModels[i].InputCost {
			t.Errorf("Expected input cost %f, got %f", expectedModels[i].InputCost, model.InputCost)
		}
		if model.OutputCost != expectedModels[i].OutputCost {
			t.Errorf("Expected output cost %f, got %f", expectedModels[i].OutputCost, model.OutputCost)
		}
		if model.IsAvailable != expectedModels[i].IsAvailable {
			t.Errorf("Expected availability %v, got %v", expectedModels[i].IsAvailable, model.IsAvailable)
		}
	}
}