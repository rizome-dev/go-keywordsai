package integrations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rizome-dev/go-keywordsai/pkg/client"
	"github.com/rizome-dev/go-keywordsai/pkg/types"
)

type Service struct {
	client *client.Client
}

func NewService(client *client.Client) *Service {
	return &Service{client: client}
}

func (s *Service) TextToSpeech(ctx context.Context, req *types.TTSRequest) ([]byte, error) {
	// TTS returns audio data, so we need custom handling
	var buf bytes.Buffer
	body, _ := json.Marshal(req)
	
	httpReq, err := http.NewRequestWithContext(ctx, "POST", s.client.BaseURL()+"/api/audio/speech", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	
	httpReq.Header.Set("Authorization", "Bearer "+s.client.APIKey())
	httpReq.Header.Set("Content-Type", "application/json")
	
	resp, err := s.client.HTTPClient().Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: status=%d, body=%s", resp.StatusCode, string(bodyBytes))
	}
	
	_, err = io.Copy(&buf, resp.Body)
	return buf.Bytes(), err
}

func (s *Service) SpeechToText(ctx context.Context, audioData []byte, req *types.STTRequest) (*types.STTResponse, error) {
	// Build multipart fields
	fields := []client.MultipartField{
		{Name: "model", Value: req.Model},
		{Name: "file", IsFile: true, FileName: "audio.wav", Data: audioData},
	}

	// Add optional fields
	if req.ResponseFormat != nil {
		fields = append(fields, client.MultipartField{Name: "response_format", Value: *req.ResponseFormat})
	}
	if req.Language != nil {
		fields = append(fields, client.MultipartField{Name: "language", Value: *req.Language})
	}
	if req.Temperature != nil {
		fields = append(fields, client.MultipartField{Name: "temperature", Value: fmt.Sprintf("%f", *req.Temperature)})
	}
	if req.Prompt != nil {
		fields = append(fields, client.MultipartField{Name: "prompt", Value: *req.Prompt})
	}

	var result types.STTResponse
	err := s.client.PostMultipart(ctx, "/api/audio/transcriptions", fields, &result)
	return &result, err
}

func (s *Service) CreateEmbeddings(ctx context.Context, req *types.EmbeddingRequest) (*types.EmbeddingResponse, error) {
	var result types.EmbeddingResponse
	return &result, s.client.Post(ctx, "/api/embeddings", req, &result)
}