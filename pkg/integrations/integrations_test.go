package integrations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rizome-dev/go-keywordsai/pkg/client"
	"github.com/rizome-dev/go-keywordsai/pkg/types"
)

func TestTextToSpeech(t *testing.T) {
	expectedAudio := []byte("fake audio data")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/audio/speech" {
			t.Errorf("Expected path /api/audio/speech, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		var req types.TTSRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Model != "tts-1" {
			t.Errorf("Expected model tts-1, got %s", req.Model)
		}
		if req.Input != "Hello, world!" {
			t.Errorf("Expected input 'Hello, world!', got %s", req.Input)
		}
		if req.Voice != "alloy" {
			t.Errorf("Expected voice alloy, got %s", req.Voice)
		}

		w.Header().Set("Content-Type", "audio/mpeg")
		w.Write(expectedAudio)
	}))
	defer server.Close()

	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := NewService(c)

	req := &types.TTSRequest{
		Model: "tts-1",
		Input: "Hello, world!",
		Voice: "alloy",
	}

	result, err := s.TextToSpeech(context.Background(), req)
	if err != nil {
		t.Fatalf("TextToSpeech() error = %v", err)
	}

	if !bytes.Equal(result, expectedAudio) {
		t.Errorf("Expected audio data %v, got %v", expectedAudio, result)
	}
}

func TestSpeechToText(t *testing.T) {
	expectedResponse := types.STTResponse{
		Text:     "Hello, world!",
		Language: "en",
		Duration: 1.5,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/audio/transcriptions" {
			t.Errorf("Expected path /api/audio/transcriptions, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Parse multipart form
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			t.Fatalf("Failed to parse multipart form: %v", err)
		}

		model := r.FormValue("model")
		if model != "whisper-1" {
			t.Errorf("Expected model whisper-1, got %s", model)
		}

		// Check file was uploaded
		file, _, err := r.FormFile("file")
		if err != nil {
			t.Errorf("Expected file upload, got error: %v", err)
		}
		file.Close()

		json.NewEncoder(w).Encode(expectedResponse)
	}))
	defer server.Close()

	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := NewService(c)

	audioData := []byte("fake audio data")
	req := &types.STTRequest{
		Model: "whisper-1",
	}

	result, err := s.SpeechToText(context.Background(), audioData, req)
	if err != nil {
		t.Fatalf("SpeechToText() error = %v", err)
	}

	if result.Text != expectedResponse.Text {
		t.Errorf("Expected text %s, got %s", expectedResponse.Text, result.Text)
	}
	if result.Language != expectedResponse.Language {
		t.Errorf("Expected language %s, got %s", expectedResponse.Language, result.Language)
	}
}

func TestSpeechToTextWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			t.Fatalf("Failed to parse multipart form: %v", err)
		}

		// Check all optional fields
		if lang := r.FormValue("language"); lang != "en" {
			t.Errorf("Expected language 'en', got %s", lang)
		}
		if format := r.FormValue("response_format"); format != "json" {
			t.Errorf("Expected response_format 'json', got %s", format)
		}
		if temp := r.FormValue("temperature"); temp != "0.500000" {
			t.Errorf("Expected temperature '0.500000', got %s", temp)
		}
		if prompt := r.FormValue("prompt"); prompt != "Test prompt" {
			t.Errorf("Expected prompt 'Test prompt', got %s", prompt)
		}

		json.NewEncoder(w).Encode(types.STTResponse{Text: "Test"})
	}))
	defer server.Close()

	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := NewService(c)

	audioData := []byte("fake audio data")
	responseFormat := "json"
	language := "en"
	temperature := 0.5
	prompt := "Test prompt"
	
	req := &types.STTRequest{
		Model:          "whisper-1",
		ResponseFormat: &responseFormat,
		Language:       &language,
		Temperature:    &temperature,
		Prompt:         &prompt,
	}

	_, err := s.SpeechToText(context.Background(), audioData, req)
	if err != nil {
		t.Fatalf("SpeechToText() error = %v", err)
	}
}

func TestCreateEmbeddings(t *testing.T) {
	expectedResponse := types.EmbeddingResponse{
		Object: "list",
		Data: []types.EmbeddingData{
			{
				Object:    "embedding",
				Embedding: []float64{0.1, 0.2, 0.3},
				Index:     0,
			},
		},
		Model: "text-embedding-ada-002",
		Usage: types.Usage{
			PromptTokens: 5,
			TotalTokens:  5,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/embeddings" {
			t.Errorf("Expected path /api/embeddings, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		var req types.EmbeddingRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Model != "text-embedding-ada-002" {
			t.Errorf("Expected model text-embedding-ada-002, got %s", req.Model)
		}

		json.NewEncoder(w).Encode(expectedResponse)
	}))
	defer server.Close()

	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := NewService(c)

	req := &types.EmbeddingRequest{
		Model: "text-embedding-ada-002",
		Input: "Hello, world!",
	}

	result, err := s.CreateEmbeddings(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateEmbeddings() error = %v", err)
	}

	if result.Model != expectedResponse.Model {
		t.Errorf("Expected model %s, got %s", expectedResponse.Model, result.Model)
	}
	if len(result.Data) != 1 {
		t.Errorf("Expected 1 embedding, got %d", len(result.Data))
	}
	if result.Usage.PromptTokens != expectedResponse.Usage.PromptTokens {
		t.Errorf("Expected prompt tokens %d, got %d", expectedResponse.Usage.PromptTokens, result.Usage.PromptTokens)
	}
}

func TestCreateEmbeddingsWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req types.EmbeddingRequest
		json.NewDecoder(r.Body).Decode(&req)

		if *req.EncodingFormat != "float" {
			t.Errorf("Expected encoding format 'float', got %s", *req.EncodingFormat)
		}
		if *req.Dimensions != 1536 {
			t.Errorf("Expected dimensions 1536, got %d", *req.Dimensions)
		}

		json.NewEncoder(w).Encode(types.EmbeddingResponse{
			Object: "list",
			Data:   []types.EmbeddingData{},
			Model:  "text-embedding-ada-002",
		})
	}))
	defer server.Close()

	c := client.New("test-key", client.WithBaseURL(server.URL))
	s := NewService(c)

	encodingFormat := "float"
	dimensions := 1536
	
	req := &types.EmbeddingRequest{
		Model:          "text-embedding-ada-002",
		Input:          []string{"Hello", "World"},
		EncodingFormat: &encodingFormat,
		Dimensions:     &dimensions,
	}

	_, err := s.CreateEmbeddings(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateEmbeddings() error = %v", err)
	}
}

// Helper to compare multipart content
func parseMultipartBody(contentType string, body io.Reader) (map[string]string, error) {
	_, params, found := strings.Cut(contentType, "boundary=")
	if !found {
		return nil, fmt.Errorf("boundary not found in content type")
	}
	
	reader := multipart.NewReader(body, strings.TrimPrefix(params, "boundary="))
	fields := make(map[string]string)
	
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		
		data, err := io.ReadAll(part)
		if err != nil {
			return nil, err
		}
		
		fields[part.FormName()] = string(data)
	}
	
	return fields, nil
}