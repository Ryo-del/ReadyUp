package http

import (
	"ReadyUp/internal/auth"
	"context"
	nethttp "net/http"
	"strings"
)

type contextKey string

const userIDContextKey contextKey = "user_id"

func AuthMiddleware(jwtManager *auth.JWTManager) func(nethttp.Handler) nethttp.Handler {
	return func(next nethttp.Handler) nethttp.Handler {
		return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
			token, ok := bearerToken(r.Header.Get("Authorization"))
			if !ok {
				nethttp.Error(w, "missing bearer token", nethttp.StatusUnauthorized)
				return
			}

			claims, err := jwtManager.Parse(token)
			if err != nil {
				nethttp.Error(w, "invalid bearer token", nethttp.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDContextKey).(int64)
	return userID, ok
}

func bearerToken(header string) (string, bool) {
	scheme, token, ok := strings.Cut(header, " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") || strings.TrimSpace(token) == "" {
		return "", false
	}
	return strings.TrimSpace(token), true
}
