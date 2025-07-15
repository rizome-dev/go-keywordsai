package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/rizome-dev/go-keywordsai/pkg/client"
	"github.com/rizome-dev/go-keywordsai/pkg/integrations"
	"github.com/rizome-dev/go-keywordsai/pkg/types"
)

func main() {
	c := client.New() // Uses KEYWORDSAI_API_KEY env var automatically
	integrationsService := integrations.NewService(c)
	ctx := context.Background()

	// Example 1: Text-to-Speech
	fmt.Println("=== Text-to-Speech Example ===")
	textToSpeech(ctx, integrationsService)

	// Example 2: Speech-to-Text
	fmt.Println("\n=== Speech-to-Text Example ===")
	speechToText(ctx, integrationsService)

	// Example 3: Embeddings
	fmt.Println("\n=== Embeddings Example ===")
	createEmbeddings(ctx, integrationsService)
}

func textToSpeech(ctx context.Context, integrationsService *integrations.Service) {
	ttsRequest := &types.TTSRequest{
		Model: "tts-1",
		Input: "Hello! This is a test of the Keywords AI text-to-speech integration.",
		Voice: "alloy",
		ResponseFormat: stringPtr("mp3"),
		Speed: floatPtr(1.0),
	}

	audioData, err := integrationsService.TextToSpeech(ctx, ttsRequest)
	if err != nil {
		log.Printf("Failed to generate speech: %v", err)
		return
	}

	// Save the audio file
	err = os.WriteFile("output.mp3", audioData, 0644)
	if err != nil {
		log.Printf("Failed to save audio file: %v", err)
		return
	}

	fmt.Printf("Generated audio file: output.mp3 (%d bytes)\n", len(audioData))
}

func speechToText(ctx context.Context, integrationsService *integrations.Service) {
	// In a real application, you would read actual audio data from a file
	// For this example, we'll use placeholder data
	audioData := []byte("audio data would go here")

	// Note: In production, you would read from an actual audio file like:
	// audioData, err := os.ReadFile("speech.wav")

	sttRequest := &types.STTRequest{
		Model:          "whisper-1",
		ResponseFormat: stringPtr("json"),
		Language:       stringPtr("en"),
		Temperature:    floatPtr(0.2),
		Prompt:         stringPtr("This is a conversation about AI technology."),
	}

	result, err := integrationsService.SpeechToText(ctx, audioData, sttRequest)
	if err != nil {
		log.Printf("Failed to transcribe speech: %v", err)
		return
	}

	fmt.Printf("Transcription: %s\n", result.Text)
	if result.Language != "" {
		fmt.Printf("Detected language: %s\n", result.Language)
	}
	if result.Duration > 0 {
		fmt.Printf("Audio duration: %.2f seconds\n", result.Duration)
	}
}

func createEmbeddings(ctx context.Context, integrationsService *integrations.Service) {
	// Example 1: Single text embedding
	embeddingRequest := &types.EmbeddingRequest{
		Model: "text-embedding-ada-002",
		Input: "Keywords AI is a powerful platform for LLM observability and monitoring.",
		EncodingFormat: stringPtr("float"),
	}

	result, err := integrationsService.CreateEmbeddings(ctx, embeddingRequest)
	if err != nil {
		log.Printf("Failed to create embeddings: %v", err)
		return
	}

	fmt.Printf("Created %d embedding(s)\n", len(result.Data))
	if len(result.Data) > 0 {
		fmt.Printf("Embedding dimension: %d\n", len(result.Data[0].Embedding))
		fmt.Printf("First 5 values: %.4f, %.4f, %.4f, %.4f, %.4f...\n",
			result.Data[0].Embedding[0],
			result.Data[0].Embedding[1],
			result.Data[0].Embedding[2],
			result.Data[0].Embedding[3],
			result.Data[0].Embedding[4])
	}
	fmt.Printf("Tokens used: %d\n", result.Usage.TotalTokens)

	// Example 2: Multiple text embeddings
	multiEmbeddingRequest := &types.EmbeddingRequest{
		Model: "text-embedding-ada-002",
		Input: []string{
			"First text to embed",
			"Second text to embed",
			"Third text to embed",
		},
		EncodingFormat: stringPtr("float"),
		Dimensions:     intPtr(1536),
	}

	multiResult, err := integrationsService.CreateEmbeddings(ctx, multiEmbeddingRequest)
	if err != nil {
		log.Printf("Failed to create multiple embeddings: %v", err)
		return
	}

	fmt.Printf("\nCreated %d embeddings for multiple inputs\n", len(multiResult.Data))
	for i, embedding := range multiResult.Data {
		fmt.Printf("  Embedding %d: dimension=%d\n", i+1, len(embedding.Embedding))
	}
}

func stringPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}

func intPtr(i int) *int {
	return &i
}