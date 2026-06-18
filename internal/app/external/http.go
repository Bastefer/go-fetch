package external

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"service-parser/internal/app/dto"
)

type SourceClient interface {
	GetProducts(ctx context.Context, url string) ([]dto.SourceProduct, error)
	GetClients(ctx context.Context, url string) ([]dto.SourceClient, error)
}

type HTTPClient struct {
	client *http.Client
}

func NewHTTPClient(client *http.Client) *HTTPClient {
	return &HTTPClient{
		client: client,
	}
}

func (c *HTTPClient) GetProducts(
	ctx context.Context,
	url string,
) ([]dto.SourceProduct, error) {
	const op = "internal/app/external/http/GetProducts"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	var products []dto.SourceProduct

	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return products, nil
}

func (c *HTTPClient) GetClients(
	ctx context.Context,
	url string,
) ([]dto.SourceClient, error) {
	const op = "internal/app/external/http/GetClients"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	var clients []dto.SourceClient

	if err := json.NewDecoder(resp.Body).Decode(&clients); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return clients, nil
}
