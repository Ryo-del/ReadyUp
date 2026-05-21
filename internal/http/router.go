package http

import (
	"ReadyUp/internal/auth"
	"ReadyUp/internal/http/handler"
	"ReadyUp/internal/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(userRepo repository.UserRepository, jwtManager *auth.JWTManager) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)

	authhandler := handler.NewAuthHandler(userRepo, jwtManager)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", authhandler.Login)
		r.Post("/register", authhandler.Register)
	})
	return r
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
