package auth_test

import (
	"testing"
	"time"

	"github.com/OnurCeliiik/ecommerce/services/user/auth"
	"github.com/google/uuid"
)

func TestHMACProvider_GenerateAndValidate(t *testing.T) {
	provider, err := auth.NewHMACProvider("test-secret", time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	userID := uuid.New()
	token, err := provider.Generate(userID)
	if err != nil {
		t.Fatal(err)
	}

	gotID, err := provider.UserIDFromToken(token)
	if err != nil {
		t.Fatal(err)
	}
	if gotID != userID {
		t.Fatalf("expected %s, got %s", userID, gotID)
	}
}

func TestHMACProvider_InvalidToken(t *testing.T) {
	provider, err := auth.NewHMACProvider("test-secret", time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	_, err = provider.UserIDFromToken("not-a-jwt")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestHMACProvider_RequiresSecret(t *testing.T) {
	_, err := auth.NewHMACProvider("", time.Hour)
	if err == nil {
		t.Fatal("expected error when secret is empty")
	}
}
