package keys

import (
	"context"
	"fmt"
	"time"

	"github.com/rizome-dev/go-keywordsai/pkg/client"
	"github.com/rizome-dev/go-keywordsai/pkg/types"
)

type Service struct {
	client *client.Client
}

func NewService(client *client.Client) *Service {
	return &Service{client: client}
}

type CreateKeyRequest struct {
	Name             *string                `json:"name,omitempty"`
	ExpiresAt        time.Time              `json:"expires_at"`
	UsageLimit       *int                   `json:"usage_limit,omitempty"`
	AllowedModels    []string               `json:"allowed_models,omitempty"`
	AllowedEndpoints []string               `json:"allowed_endpoints,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

func (s *Service) Create(ctx context.Context, req *CreateKeyRequest) (*types.TemporaryKey, error) {
	var result types.TemporaryKey
	return &result, s.client.Post(ctx, "/api/temporary-keys", req, &result)
}

func (s *Service) List(ctx context.Context) ([]types.TemporaryKey, error) {
	var result []types.TemporaryKey
	return result, s.client.Get(ctx, "/api/temporary-keys", &result)
}

func (s *Service) Get(ctx context.Context, keyID string) (*types.TemporaryKey, error) {
	var result types.TemporaryKey
	path := fmt.Sprintf("/api/temporary-keys/%s", keyID)
	return &result, s.client.Get(ctx, path, &result)
}

func (s *Service) Update(ctx context.Context, keyID string, updates map[string]interface{}) (*types.TemporaryKey, error) {
	var result types.TemporaryKey
	path := fmt.Sprintf("/api/temporary-keys/%s", keyID)
	return &result, s.client.Patch(ctx, path, updates, &result)
}

func (s *Service) Delete(ctx context.Context, keyID string) error {
	path := fmt.Sprintf("/api/temporary-keys/%s", keyID)
	return s.client.Delete(ctx, path, nil)
}