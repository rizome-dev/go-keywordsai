package models

import (
	"context"

	"github.com/rizome-dev/go-keywordsai/pkg/client"
	"github.com/rizome-dev/go-keywordsai/pkg/types"
)

type Service struct {
	client *client.Client
}

func NewService(client *client.Client) *Service {
	return &Service{client: client}
}

func (s *Service) List(ctx context.Context) ([]types.Model, error) {
	var result []types.Model
	return result, s.client.Get(ctx, "/api/models", &result)
}