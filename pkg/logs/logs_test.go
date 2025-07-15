package logs

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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/request-logs/create/" {
			t.Errorf("Expected path /api/request-logs/create/, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		var log types.RequestLog
		if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if log.Model != "gpt-4" {
			t.Errorf("Expected model gpt-4, got %s", log.Model)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := NewService(c)

	log := &types.RequestLog{
		Model: "gpt-4",
		PromptMessages: []types.Message{
			{Role: "user", Content: "Hello"},
		},
	}

	err := s.Create(context.Background(), log)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
}

func TestBatchCreate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/request-logs/batch/create" {
			t.Errorf("Expected path /api/request-logs/batch/create, got %s", r.URL.Path)
		}

		var payload types.BatchRequestLogsPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if len(payload.Logs) != 2 {
			t.Errorf("Expected 2 logs, got %d", len(payload.Logs))
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := NewService(c)

	logs := []types.RequestLog{
		{
			Model: "gpt-4",
			PromptMessages: []types.Message{
				{Role: "user", Content: "Hello"},
			},
		},
		{
			Model: "gpt-3.5-turbo",
			PromptMessages: []types.Message{
				{Role: "user", Content: "Hi"},
			},
		},
	}

	err := s.BatchCreate(context.Background(), logs)
	if err != nil {
		t.Fatalf("BatchCreate() error = %v", err)
	}
}

func TestBatchCreateTooManyLogs(t *testing.T) {
	c := client.New("test-key")
	s := NewService(c)

	// Create 5001 logs (exceeds limit)
	logs := make([]types.RequestLog, 5001)
	
	err := s.BatchCreate(context.Background(), logs)
	if err == nil {
		t.Fatal("Expected error for too many logs, got nil")
	}
}

func TestList(t *testing.T) {
	now := time.Now()
	expectedLogs := []types.RequestLog{
		{
			Model:     "gpt-4",
			Timestamp: &now,
		},
		{
			Model:     "gpt-3.5-turbo",
			Timestamp: &now,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := types.LogsResponse{
			Logs:       expectedLogs,
			TotalCount: 2,
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := NewService(c)

	result, err := s.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if result.TotalCount != 2 {
		t.Errorf("Expected total count 2, got %d", result.TotalCount)
	}
	if len(result.Logs) != 2 {
		t.Errorf("Expected 2 logs, got %d", len(result.Logs))
	}
}

func TestGet(t *testing.T) {
	expectedLog := types.RequestLog{
		Model: "gpt-4",
		PromptMessages: []types.Message{
			{Role: "user", Content: "Hello"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/request-logs/test-log-id" {
			t.Errorf("Expected path /api/request-logs/test-log-id, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(expectedLog)
	}))
	defer server.Close()

	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := NewService(c)

	result, err := s.Get(context.Background(), "test-log-id")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if result.Model != expectedLog.Model {
		t.Errorf("Expected model %s, got %s", expectedLog.Model, result.Model)
	}
}

func TestUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/request-logs/test-log-id" {
			t.Errorf("Expected path /api/request-logs/test-log-id, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH method, got %s", r.Method)
		}

		var updates map[string]interface{}
		json.NewDecoder(r.Body).Decode(&updates)
		
		if updates["tags"] == nil {
			t.Error("Expected tags in update payload")
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := NewService(c)

	updates := map[string]interface{}{
		"tags": []string{"important", "reviewed"},
	}

	err := s.Update(context.Background(), "test-log-id", updates)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
}

func TestListThreads(t *testing.T) {
	expectedThreads := []types.Thread{
		{
			ID:                 "thread-1",
			CustomerIdentifier: "customer-123",
			Messages:          []types.Message{},
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/threads" {
			t.Errorf("Expected path /api/threads, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(expectedThreads)
	}))
	defer server.Close()

	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := NewService(c)

	result, err := s.ListThreads(context.Background(), "")
	if err != nil {
		t.Fatalf("ListThreads() error = %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 thread, got %d", len(result))
	}
	if result[0].ID != "thread-1" {
		t.Errorf("Expected thread ID thread-1, got %s", result[0].ID)
	}
}