package prompts

import (
	"context"
	"fmt"

	"github.com/rizome-dev/go-keywordsai/pkg/client"
	"github.com/rizome-dev/go-keywordsai/pkg/types"
)

type Service struct {
	client *client.Client
}

func NewService(client *client.Client) *Service {
	return &Service{client: client}
}

func (s *Service) Create(ctx context.Context, name string, description *string) (*types.Prompt, error) {
	payload := map[string]interface{}{
		"name": name,
	}
	if description != nil {
		payload["description"] = *description
	}

	var result types.Prompt
	return &result, s.client.Post(ctx, "/api/prompts/", payload, &result)
}

func (s *Service) List(ctx context.Context) ([]types.Prompt, error) {
	var result []types.Prompt
	return result, s.client.Get(ctx, "/api/prompts/", &result)
}

func (s *Service) Get(ctx context.Context, promptID string) (*types.Prompt, error) {
	var result types.Prompt
	path := fmt.Sprintf("/api/prompts/%s", promptID)
	return &result, s.client.Get(ctx, path, &result)
}

func (s *Service) Update(ctx context.Context, promptID string, updates map[string]interface{}) (*types.Prompt, error) {
	var result types.Prompt
	path := fmt.Sprintf("/api/prompts/%s", promptID)
	return &result, s.client.Patch(ctx, path, updates, &result)
}

func (s *Service) Delete(ctx context.Context, promptID string) error {
	path := fmt.Sprintf("/api/prompts/%s", promptID)
	return s.client.Delete(ctx, path, nil)
}

func (s *Service) CreateVersion(ctx context.Context, promptID string, version *types.PromptVersion) (*types.PromptVersion, error) {
	var result types.PromptVersion
	path := fmt.Sprintf("/api/prompts/%s/versions", promptID)
	return &result, s.client.Post(ctx, path, version, &result)
}

func (s *Service) ListVersions(ctx context.Context, promptID string) ([]types.PromptVersion, error) {
	var result []types.PromptVersion
	path := fmt.Sprintf("/api/prompts/%s/versions", promptID)
	return result, s.client.Get(ctx, path, &result)
}

func (s *Service) GetVersion(ctx context.Context, promptID string, versionID string) (*types.PromptVersion, error) {
	var result types.PromptVersion
	path := fmt.Sprintf("/api/prompts/%s/versions/%s", promptID, versionID)
	return &result, s.client.Get(ctx, path, &result)
}

func (s *Service) UpdateVersion(ctx context.Context, promptID string, versionID string, updates map[string]interface{}) (*types.PromptVersion, error) {
	var result types.PromptVersion
	path := fmt.Sprintf("/api/prompts/%s/versions/%s", promptID, versionID)
	return &result, s.client.Patch(ctx, path, updates, &result)
}

func (s *Service) DeleteVersion(ctx context.Context, promptID string, versionID string) error {
	path := fmt.Sprintf("/api/prompts/%s/versions/%s", promptID, versionID)
	return s.client.Delete(ctx, path, nil)
}