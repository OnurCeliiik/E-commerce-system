package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OnurCeliiik/ecommerce/services/user/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type mockTokenValidator struct {
	userID uuid.UUID
	err    error
}

func (m *mockTokenValidator) UserIDFromToken(token string) (uuid.UUID, error) {
	return m.userID, m.err
}

func init() {
	gin.SetMode(gin.TestMode)
}

func TestAuth(t *testing.T) {
	validUserID := uuid.New()

	tests := []struct {
		name       string
		header     string
		validator  middleware.TokenValidator
		wantStatus int
	}{
		{
			name:       "valid token",
			header:     "Bearer valid-token",
			validator:  &mockTokenValidator{userID: validUserID},
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing token",
			header:     "",
			validator:  &mockTokenValidator{},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid token",
			header:     "Bearer invalid-token",
			validator:  &mockTokenValidator{err: errors.New("invalid token")},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid prefix",
			header:     "Bearer: valid-token",
			validator:  &mockTokenValidator{userID: validUserID},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "bearer with extra spaces",
			header:     "Bearer    valid-token",
			validator:  &mockTokenValidator{userID: validUserID},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/", middleware.Auth(tt.validator), func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Fatalf("expected %d, got %d", tt.wantStatus, rr.Code)
			}
		})
	}
}

func TestAuth_SetsUserIDOnContext(t *testing.T) {
	userID := uuid.New()
	validator := &mockTokenValidator{userID: userID}

	router := gin.New()
	router.GET("/", middleware.Auth(validator), func(c *gin.Context) {
		got, ok := middleware.UserIDFromContext(c)
		if !ok {
			t.Fatal("expected user id in context")
		}
		if got != userID {
			t.Fatalf("expected %s, got %s", userID, got)
		}
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}
