package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrInvalidToken = errors.New("invalid token")

type HMACProvider struct {
	secret []byte
}

func NewHMACProvider(secret string) (*HMACProvider, error) {
	if secret == "" {
		return nil, errors.New("jwt secret is required")
	}

	return &HMACProvider{secret: []byte(secret)}, nil
}

func (p *HMACProvider) UserIDFromToken(tokenString string) (uuid.UUID, error) {
	claims, err := p.parseClaims(tokenString)
	if err != nil {
		return uuid.Nil, err
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

func (p *HMACProvider) RoleFromToken(tokenString string) (string, error) {
	claims, err := p.parseClaims(tokenString)
	if err != nil {
		return "", err
	}

	role, ok := claims["role"].(string)
	if !ok || role == "" {
		return "", ErrInvalidToken
	}

	return role, nil
}

func (p *HMACProvider) parseClaims(tokenString string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return p.secret, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
