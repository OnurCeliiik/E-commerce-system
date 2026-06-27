package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OnurCeliiik/ecommerce/services/user/dto"
	"github.com/OnurCeliiik/ecommerce/services/user/handlers"
	"github.com/OnurCeliiik/ecommerce/services/user/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type mockUserService struct {
	register     func(ctx context.Context, req dto.RegisterUserRequest) (*dto.RegisterUserResponse, error)
	login        func(ctx context.Context, req dto.LoginUserRequest) (*dto.LoginUserResponse, error)
	me           func(ctx context.Context, userID uuid.UUID) (*dto.MeResponse, error)
	getUserEmail func(ctx context.Context, userID uuid.UUID) (string, error)
}

func (m *mockUserService) Register(ctx context.Context, req dto.RegisterUserRequest) (*dto.RegisterUserResponse, error) {
	return m.register(ctx, req)
}

func (m *mockUserService) Login(ctx context.Context, req dto.LoginUserRequest) (*dto.LoginUserResponse, error) {
	return m.login(ctx, req)
}

func (m *mockUserService) Me(ctx context.Context, userID uuid.UUID) (*dto.MeResponse, error) {
	return m.me(ctx, userID)
}

func (m *mockUserService) GetUserEmail(ctx context.Context, userID uuid.UUID) (string, error) {
	return m.getUserEmail(ctx, userID)
}

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRegister_Conflict(t *testing.T) {
	handler := handlers.NewUserHandler(&mockUserService{
		register: func(ctx context.Context, req dto.RegisterUserRequest) (*dto.RegisterUserResponse, error) {
			return nil, service.ErrEmailAlreadyExists
		},
	})

	body := `{"first_name":"Ada","last_name":"Lovelace","email":"ada@example.com","password":"secret123"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Register(c)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w.Code)
	}
}

func TestLogin_Unauthorized(t *testing.T) {
	handler := handlers.NewUserHandler(&mockUserService{
		login: func(ctx context.Context, req dto.LoginUserRequest) (*dto.LoginUserResponse, error) {
			return nil, service.ErrInvalidCredentials
		},
	})

	body := `{"email":"ada@example.com","password":"wrong"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestMe_Success(t *testing.T) {
	userID := uuid.New()
	handler := handlers.NewUserHandler(&mockUserService{
		me: func(ctx context.Context, id uuid.UUID) (*dto.MeResponse, error) {
			return &dto.MeResponse{
				ID:        userID.String(),
				FirstName: "Ada",
				LastName:  "Lovelace",
				Email:     "ada@example.com",
			}, nil
		},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	c.Set("userID", userID)

	handler.Me(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp dto.MeResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Email != "ada@example.com" {
		t.Fatalf("expected ada@example.com, got %s", resp.Email)
	}
}

func TestMe_MissingUserID(t *testing.T) {
	handler := handlers.NewUserHandler(&mockUserService{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)

	handler.Me(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
