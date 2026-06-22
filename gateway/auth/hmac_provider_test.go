package auth_test

import (
	"testing"
	"time"

	"github.com/OnurCeliiik/ecommerce/gateway/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestHMACProvider_ValidateToken(t *testing.T) {
	secret := "gateway-test-secret"
	provider, err := auth.NewHMACProvider(secret)
	if err != nil {
		t.Fatal(err)
	}

	_, err = provider.UserIDFromToken("invalid-token")
	if err == nil {
		t.Fatal("expected invalid token error")
	}
}

func TestHMACProvider_RequiresSecret(t *testing.T) {
	_, err := auth.NewHMACProvider("")
	if err == nil {
		t.Fatal("expected error for empty secret")
	}
}

func TestHMACProvider_ValidUserServiceToken(t *testing.T) {
	secret := "shared-secret"
	userID := uuid.New()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatal(err)
	}

	provider, err := auth.NewHMACProvider(secret)
	if err != nil {
		t.Fatal(err)
	}

	gotID, err := provider.UserIDFromToken(tokenString)
	if err != nil {
		t.Fatal(err)
	}
	if gotID != userID {
		t.Fatalf("expected %s, got %s", userID, gotID)
	}
}
