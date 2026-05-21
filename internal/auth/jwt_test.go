package auth

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestJWTManagerGenerateAndParse(t *testing.T) {
	manager, err := NewJWTManager("test-secret", time.Hour)
	if err != nil {
		t.Fatalf("NewJWTManager returned error: %v", err)
	}

	token, err := manager.Generate(42)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if len(strings.Split(token, ".")) != 3 {
		t.Fatalf("token = %q, want three jwt parts", token)
	}

	claims, err := manager.Parse(token)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if claims.UserID != 42 {
		t.Fatalf("user id = %d, want %d", claims.UserID, int64(42))
	}
	if claims.ExpiresAt.Before(time.Now().UTC()) {
		t.Fatal("expires_at is in the past")
	}
}

func TestJWTManagerParseRejectsTamperedToken(t *testing.T) {
	manager, err := NewJWTManager("test-secret", time.Hour)
	if err != nil {
		t.Fatalf("NewJWTManager returned error: %v", err)
	}

	token, err := manager.Generate(42)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	token = token[:len(token)-1] + "x"

	_, err = manager.Parse(token)
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidToken)
	}
}

func TestJWTManagerParseRejectsExpiredToken(t *testing.T) {
	manager, err := NewJWTManager("test-secret", time.Nanosecond)
	if err != nil {
		t.Fatalf("NewJWTManager returned error: %v", err)
	}

	token, err := manager.Generate(42)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	time.Sleep(time.Millisecond)

	_, err = manager.Parse(token)
	if !errors.Is(err, ErrExpiredToken) {
		t.Fatalf("error = %v, want %v", err, ErrExpiredToken)
	}
}
