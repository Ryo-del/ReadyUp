package auth

import (
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPasswordCreatesBcryptHash(t *testing.T) {
	password := "strong-password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if hash == "" {
		t.Fatal("HashPassword returned empty hash")
	}

	if hash == password {
		t.Fatal("HashPassword returned the plain password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		t.Fatalf("hash does not match original password: %v", err)
	}
}

func TestCheckPassword(t *testing.T) {
	hash, err := HashPassword("correct-password")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	tests := []struct {
		name        string
		password    string
		wantErr     bool
		wantErrType error
	}{
		{
			name:     "matching password",
			password: "correct-password",
		},
		{
			name:        "wrong password",
			password:    "wrong-password",
			wantErr:     true,
			wantErrType: bcrypt.ErrMismatchedHashAndPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPassword(hash, tt.password)
			if tt.wantErr {
				if err == nil {
					t.Fatal("CheckPassword returned nil error")
				}
				if tt.wantErrType != nil && !errors.Is(err, tt.wantErrType) {
					t.Fatalf("CheckPassword error = %v, want %v", err, tt.wantErrType)
				}
				return
			}

			if err != nil {
				t.Fatalf("CheckPassword returned error: %v", err)
			}
		})
	}
}
