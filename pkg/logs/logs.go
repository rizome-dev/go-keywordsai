package logs

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

func (s *Service) Create(ctx context.Context, log *types.RequestLog) error {
	return s.client.Post(ctx, "/api/request-logs/create/", log, nil)
}

func (s *Service) BatchCreate(ctx context.Context, logs []types.RequestLog) error {
	if len(logs) > 5000 {
		return fmt.Errorf("batch size exceeds maximum of 5000 logs")
	}

	payload := types.BatchRequestLogsPayload{
		Logs: logs,
	}
	return s.client.Post(ctx, "/api/request-logs/batch/create", payload, nil)
}

func (s *Service) List(ctx context.Context, filter *types.LogFilter) (*types.LogsResponse, error) {
	var result types.LogsResponse
	path := "/api/request-logs"
	return &result, s.client.GetWithQuery(ctx, path, filter, &result)
}

func (s *Service) Get(ctx context.Context, logID string) (*types.RequestLog, error) {
	var result types.RequestLog
	path := fmt.Sprintf("/api/request-logs/%s", logID)
	return &result, s.client.Get(ctx, path, &result)
}

func (s *Service) Update(ctx context.Context, logID string, updates map[string]interface{}) error {
	path := fmt.Sprintf("/api/request-logs/%s", logID)
	return s.client.Patch(ctx, path, updates, nil)
}

func (s *Service) ListThreads(ctx context.Context, customerIdentifier string) ([]types.Thread, error) {
	var result []types.Thread
	path := "/api/threads"
	if customerIdentifier != "" {
		// TODO: Add query parameter support
		payload := map[string]string{"customer_identifier": customerIdentifier}
		return result, s.client.Post(ctx, path, payload, &result)
	}
	return result, s.client.Get(ctx, path, &result)
}