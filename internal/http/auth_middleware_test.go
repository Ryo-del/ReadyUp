package http

import (
	"ReadyUp/internal/auth"
	nethttp "net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthMiddlewareAllowsValidBearerToken(t *testing.T) {
	manager := testJWTManager(t)
	token, err := manager.Generate(42)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	handler := AuthMiddleware(manager)(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		userID, ok := UserIDFromContext(r.Context())
		if !ok {
			t.Fatal("user id is missing from context")
		}
		if userID != 42 {
			t.Fatalf("user id = %d, want %d", userID, int64(42))
		}
		w.WriteHeader(nethttp.StatusNoContent)
	}))

	req := httptest.NewRequest(nethttp.MethodGet, "/private", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, nethttp.StatusNoContent)
	}
}

func TestAuthMiddlewareRejectsMissingBearerToken(t *testing.T) {
	handler := AuthMiddleware(testJWTManager(t))(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		t.Fatal("next handler should not be called")
	}))
	req := httptest.NewRequest(nethttp.MethodGet, "/private", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, nethttp.StatusUnauthorized)
	}
}

func TestAuthMiddlewareRejectsInvalidBearerToken(t *testing.T) {
	handler := AuthMiddleware(testJWTManager(t))(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		t.Fatal("next handler should not be called")
	}))
	req := httptest.NewRequest(nethttp.MethodGet, "/private", nil)
	req.Header.Set("Authorization", "Bearer invalid")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, nethttp.StatusUnauthorized)
	}
}

func testJWTManager(t *testing.T) *auth.JWTManager {
	t.Helper()

	manager, err := auth.NewJWTManager("test-secret", time.Hour)
	if err != nil {
		t.Fatalf("NewJWTManager returned error: %v", err)
	}
	return manager
}
