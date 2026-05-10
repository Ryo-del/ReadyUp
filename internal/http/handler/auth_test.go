package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ReadyUp/internal/auth"
	"ReadyUp/internal/model"
)

type fakeUserRepository struct {
	createFunc     func(ctx context.Context, email, username, passwordHash string) error
	getByEmailFunc func(ctx context.Context, email string) (*model.User, error)
}

func (r fakeUserRepository) Create(ctx context.Context, email, username, passwordHash string) error {
	if r.createFunc == nil {
		return nil
	}
	return r.createFunc(ctx, email, username, passwordHash)
}

func (r fakeUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	if r.getByEmailFunc == nil {
		return nil, nil
	}
	return r.getByEmailFunc(ctx, email)
}

func TestAuthHandlerRegister(t *testing.T) {
	t.Run("creates user with hashed password", func(t *testing.T) {
		var gotEmail, gotUsername, gotHash string
		handler := NewAuthHandler(fakeUserRepository{
			createFunc: func(ctx context.Context, email, username, passwordHash string) error {
				gotEmail = email
				gotUsername = username
				gotHash = passwordHash
				return nil
			},
		})

		req := httptest.NewRequest(
			http.MethodPost,
			"/auth/register",
			strings.NewReader(`{"username":"artem","email":"artem@example.com","password":"secret"}`),
		)
		rec := httptest.NewRecorder()

		handler.Register(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body: %q", rec.Code, http.StatusOK, rec.Body.String())
		}
		if gotEmail != "artem@example.com" {
			t.Fatalf("email = %q, want %q", gotEmail, "artem@example.com")
		}
		if gotUsername != "artem" {
			t.Fatalf("username = %q, want %q", gotUsername, "artem")
		}
		if gotHash == "" {
			t.Fatal("password hash is empty")
		}
		if gotHash == "secret" {
			t.Fatal("plain password was passed to repository")
		}
		if err := auth.CheckPassword(gotHash, "secret"); err != nil {
			t.Fatalf("password hash does not match original password: %v", err)
		}
	})

	t.Run("rejects invalid json", func(t *testing.T) {
		handler := NewAuthHandler(fakeUserRepository{
			createFunc: func(ctx context.Context, email, username, passwordHash string) error {
				t.Fatal("Create should not be called")
				return nil
			},
		})
		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(`{`))
		rec := httptest.NewRecorder()

		handler.Register(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("returns internal server error when repository fails", func(t *testing.T) {
		handler := NewAuthHandler(fakeUserRepository{
			createFunc: func(ctx context.Context, email, username, passwordHash string) error {
				return errors.New("insert failed")
			},
		})
		req := httptest.NewRequest(
			http.MethodPost,
			"/auth/register",
			strings.NewReader(`{"username":"artem","email":"artem@example.com","password":"secret"}`),
		)
		rec := httptest.NewRecorder()

		handler.Register(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
		}
	})
}

func TestAuthHandlerLogin(t *testing.T) {
	t.Run("returns ok for valid credentials", func(t *testing.T) {
		passwordHash, err := auth.HashPassword("secret")
		if err != nil {
			t.Fatalf("HashPassword returned error: %v", err)
		}
		handler := NewAuthHandler(fakeUserRepository{
			getByEmailFunc: func(ctx context.Context, email string) (*model.User, error) {
				if email != "artem@example.com" {
					t.Fatalf("email = %q, want %q", email, "artem@example.com")
				}
				return &model.User{Email: email, PasswordHash: passwordHash}, nil
			},
		})
		req := httptest.NewRequest(
			http.MethodPost,
			"/auth/login",
			strings.NewReader(`{"email":"artem@example.com","password":"secret"}`),
		)
		rec := httptest.NewRecorder()

		handler.Login(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body: %q", rec.Code, http.StatusOK, rec.Body.String())
		}
	})

	t.Run("rejects invalid json", func(t *testing.T) {
		handler := NewAuthHandler(fakeUserRepository{
			getByEmailFunc: func(ctx context.Context, email string) (*model.User, error) {
				t.Fatal("GetByEmail should not be called")
				return nil, nil
			},
		})
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{`))
		rec := httptest.NewRecorder()

		handler.Login(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("rejects wrong password", func(t *testing.T) {
		passwordHash, err := auth.HashPassword("secret")
		if err != nil {
			t.Fatalf("HashPassword returned error: %v", err)
		}
		handler := NewAuthHandler(fakeUserRepository{
			getByEmailFunc: func(ctx context.Context, email string) (*model.User, error) {
				return &model.User{Email: email, PasswordHash: passwordHash}, nil
			},
		})
		req := httptest.NewRequest(
			http.MethodPost,
			"/auth/login",
			strings.NewReader(`{"email":"artem@example.com","password":"bad-secret"}`),
		)
		rec := httptest.NewRecorder()

		handler.Login(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("returns internal server error when repository fails", func(t *testing.T) {
		handler := NewAuthHandler(fakeUserRepository{
			getByEmailFunc: func(ctx context.Context, email string) (*model.User, error) {
				return nil, errors.New("select failed")
			},
		})
		req := httptest.NewRequest(
			http.MethodPost,
			"/auth/login",
			strings.NewReader(`{"email":"artem@example.com","password":"secret"}`),
		)
		rec := httptest.NewRecorder()

		handler.Login(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
		}
	})
}
