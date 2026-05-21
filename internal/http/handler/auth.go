package handler

import (
	"ReadyUp/internal/auth"
	repo "ReadyUp/internal/repository"
	"encoding/json"
	"log/slog"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
	ExpiresIn int64  `json:"expires_in"`
}

type AuthHandler struct {
	users      repo.UserRepository
	jwtManager *auth.JWTManager
}

func NewAuthHandler(users repo.UserRepository, jwtManager *auth.JWTManager) *AuthHandler {
	return &AuthHandler{users: users, jwtManager: jwtManager}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("failed to parse login request", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	user, err := h.users.GetByEmail(r.Context(), req.Email)
	if err != nil {
		slog.Error("failed to get user from database", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	err = auth.CheckPassword(
		user.PasswordHash,
		req.Password,
	)

	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	h.writeToken(w, user.ID)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("failed to parse login request", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	userID, err := h.users.Create(r.Context(), req.Email, req.UserName, hash)
	if err != nil {
		slog.Error("failed to create user from database", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	h.writeToken(w, userID)
}

func (h *AuthHandler) writeToken(w http.ResponseWriter, userID int64) {
	token, err := h.jwtManager.Generate(userID)
	if err != nil {
		slog.Error("failed to generate jwt token", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(AuthResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: int64(h.jwtManager.TTL().Seconds()),
	})
	if err != nil {
		slog.Error("failed to encode auth response", "error", err)
	}
}
