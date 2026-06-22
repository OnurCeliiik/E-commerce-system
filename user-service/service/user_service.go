package service

import (
	"context"
	"errors"

	"github.com/OnurCeliiik/ecommerce/services/user/dto"
	"github.com/OnurCeliiik/ecommerce/services/user/model"
	"github.com/OnurCeliiik/ecommerce/services/user/repository"
	"github.com/OnurCeliiik/ecommerce/services/user/utils"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, userID uuid.UUID) (*model.User, error)
}

type TokenGenerator interface {
	Generate(userID uuid.UUID) (string, error)
}

type userService struct {
	repo   UserRepository
	tokens TokenGenerator
}

func NewUserService(repo UserRepository, tokens TokenGenerator) *userService {
	return &userService{repo: repo, tokens: tokens}
}

func (s *userService) Register(ctx context.Context, req dto.RegisterUserRequest) (*dto.RegisterUserResponse, error) {
	_, err := s.repo.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrEmailAlreadyExists
	}
	if !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, ErrInvalidInput
	}

	user := &model.User{
		ID:           uuid.New(),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PasswordHash: hash,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &dto.RegisterUserResponse{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *userService) Login(ctx context.Context, req dto.LoginUserRequest) (*dto.LoginUserResponse, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := utils.ComparePassword(req.Password, user.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.tokens.Generate(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.LoginUserResponse{
		Token: token,
	}, nil
}

func (s *userService) Me(ctx context.Context, userID uuid.UUID) (*dto.MeResponse, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &dto.MeResponse{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}, nil
}
