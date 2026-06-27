package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

// HTTPUserClient fetches user data from user-service over HTTP.
type HTTPUserClient struct {
	baseURL    string
	secret     string
	httpClient *http.Client
}

func NewHTTPUserClient(baseURL string) *HTTPUserClient {
	return &HTTPUserClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		secret:  os.Getenv("INTERNAL_SERVICE_SECRET"),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *HTTPUserClient) GetUserEmail(ctx context.Context, userID uuid.UUID) (string, error) {
	url := fmt.Sprintf("%s/api/v1/internal/users/%s", c.baseURL, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("X-Internal-Secret", c.secret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var user struct {
			Email string `json:"email"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			return "", fmt.Errorf("decode response: %w", err)
		}
		return user.Email, nil
	case http.StatusNotFound:
		return "", ErrUserNotFound
	default:
		return "", fmt.Errorf("user service returned status %d", resp.StatusCode)
	}
}
