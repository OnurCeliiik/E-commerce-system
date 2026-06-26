package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// HTTPProductClient fetches product data from product-service over HTTP.
// It satisfies service.ProductCatalog when passed to NewOrderService.
type HTTPProductClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHTTPProductClient(baseURL string) *HTTPProductClient {
	return &HTTPProductClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *HTTPProductClient) GetUnitPrice(ctx context.Context, productID uuid.UUID) (float64, error) {
	url := fmt.Sprintf("%s/api/v1/products/%s", c.baseURL, productID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// Decode only what we need — do not import product-service packages.
		var product struct {
			Price float64 `json:"price"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
			return 0, fmt.Errorf("decode response: %w", err)
		}
		return product.Price, nil
	case http.StatusNotFound:
		return 0, ErrProductNotFound
	default:
		return 0, fmt.Errorf("product service returned status %d", resp.StatusCode)
	}
}
