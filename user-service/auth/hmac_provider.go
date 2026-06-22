package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrInvalidToken = errors.New("invalid token")

type HMACProvider struct {
	secret []byte
	ttl    time.Duration
}

func NewHMACProvider(secret string, ttl time.Duration) (*HMACProvider, error) {
	if secret == "" {
		return nil, errors.New("jwt secret is required")
	}

	return &HMACProvider{
		secret: []byte(secret),
		ttl:    ttl,
	}, nil
}

func (p *HMACProvider) Generate(userID uuid.UUID) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID.String(),
		"exp": now.Add(p.ttl).Unix(),
		"iat": now.Unix(),
	})

	return token.SignedString(p.secret)
}

func (p *HMACProvider) UserIDFromToken(tokenString string) (uuid.UUID, error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return p.secret, nil
	})
	if err != nil {
		return uuid.Nil, ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return uuid.Nil, ErrInvalidToken
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, ErrInvalidToken
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, ErrInvalidToken
	}

	return userID, nil
}
