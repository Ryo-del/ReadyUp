package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("failed to parse login request", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

}

func Register(w http.ResponseWriter, r *http.Request) {

}
