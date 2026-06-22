package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/OnurCeliiik/ecommerce/services/user/dto"
	"github.com/OnurCeliiik/ecommerce/services/user/model"
	"github.com/OnurCeliiik/ecommerce/services/user/repository"
	"github.com/OnurCeliiik/ecommerce/services/user/service"
	"github.com/OnurCeliiik/ecommerce/services/user/utils"
	"github.com/google/uuid"
)

type mockUserRepository struct {
	findByEmail func(ctx context.Context, email string) (*model.User, error)
	findByID    func(ctx context.Context, userID uuid.UUID) (*model.User, error)
	create      func(ctx context.Context, user *model.User) error
}

func (m *mockUserRepository) Create(ctx context.Context, user *model.User) error {
	return m.create(ctx, user)
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return m.findByEmail(ctx, email)
}

func (m *mockUserRepository) FindByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	return m.findByID(ctx, userID)
}

type mockTokenGenerator struct {
	generate func(userID uuid.UUID) (string, error)
}

func (m *mockTokenGenerator) Generate(userID uuid.UUID) (string, error) {
	return m.generate(userID)
}

func TestRegister_Success(t *testing.T) {
	repo := &mockUserRepository{
		findByEmail: func(ctx context.Context, email string) (*model.User, error) {
			return nil, repository.ErrUserNotFound
		},
		create: func(ctx context.Context, user *model.User) error {
			user.CreatedAt = time.Now()
			return nil
		},
	}

	svc := service.NewUserService(repo, &mockTokenGenerator{})
	resp, err := svc.Register(context.Background(), dto.RegisterUserRequest{
		FirstName: "Ada",
		LastName:  "Lovelace",
		Email:     "ada@example.com",
		Password:  "secret123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Email != "ada@example.com" {
		t.Fatalf("expected email ada@example.com, got %s", resp.Email)
	}
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	repo := &mockUserRepository{
		findByEmail: func(ctx context.Context, email string) (*model.User, error) {
			return &model.User{Email: email}, nil
		},
	}

	svc := service.NewUserService(repo, &mockTokenGenerator{})
	_, err := svc.Register(context.Background(), dto.RegisterUserRequest{
		FirstName: "Ada",
		LastName:  "Lovelace",
		Email:     "ada@example.com",
		Password:  "secret123",
	})
	if !errors.Is(err, service.ErrEmailAlreadyExists) {
		t.Fatalf("expected ErrEmailAlreadyExists, got %v", err)
	}
}

func TestLogin_Success(t *testing.T) {
	userID := uuid.New()
	hash, err := utils.HashPassword("secret123")
	if err != nil {
		t.Fatal(err)
	}

	repo := &mockUserRepository{
		findByEmail: func(ctx context.Context, email string) (*model.User, error) {
			return &model.User{
				ID:           userID,
				Email:        email,
				PasswordHash: hash,
			}, nil
		},
	}
	tokens := &mockTokenGenerator{
		generate: func(id uuid.UUID) (string, error) {
			if id != userID {
				t.Fatalf("expected user id %s, got %s", userID, id)
			}
			return "jwt-token", nil
		},
	}

	svc := service.NewUserService(repo, tokens)
	resp, err := svc.Login(context.Background(), dto.LoginUserRequest{
		Email:    "ada@example.com",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Token != "jwt-token" {
		t.Fatalf("expected jwt-token, got %s", resp.Token)
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	repo := &mockUserRepository{
		findByEmail: func(ctx context.Context, email string) (*model.User, error) {
			return nil, repository.ErrUserNotFound
		},
	}

	svc := service.NewUserService(repo, &mockTokenGenerator{})
	_, err := svc.Login(context.Background(), dto.LoginUserRequest{
		Email:    "missing@example.com",
		Password: "secret123",
	})
	if !errors.Is(err, service.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestMe_UserNotFound(t *testing.T) {
	repo := &mockUserRepository{
		findByID: func(ctx context.Context, userID uuid.UUID) (*model.User, error) {
			return nil, repository.ErrUserNotFound
		},
	}

	svc := service.NewUserService(repo, &mockTokenGenerator{})
	_, err := svc.Me(context.Background(), uuid.New())
	if !errors.Is(err, service.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
