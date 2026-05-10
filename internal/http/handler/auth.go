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

type RegiserRequest struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type AuthHandler struct {
	users repo.UserRepository
}

func NewAuthHandler(users repo.UserRepository) *AuthHandler {
	return &AuthHandler{users: users}
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

	w.Write([]byte("login success"))
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegiserRequest
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

	err = h.users.Create(r.Context(), req.Email, req.UserName, hash)
	if err != nil {
		slog.Error("failed to create user from database", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("regiser success"))
}
