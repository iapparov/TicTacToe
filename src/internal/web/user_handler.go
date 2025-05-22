package web

import (
	"encoding/json"
	"krestikinoliki/internal/app"
	"net/http"
	"strings"
)

type UserHandler struct {
	userService app.UserService
	userRepo app.UserRepository
	JwtProvider *app.JwtProvider
}

func NewUserHandler(service app.UserService, repo app.UserRepository, jwt *app.JwtProvider) *UserHandler {
	return &UserHandler{
		userService: service,
		userRepo: repo,
		JwtProvider: jwt,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req app.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.userService.RegisterUser(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.userRepo.Save(user)

	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req app.JwtRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	resp, err := h.userService.LoginJwt(req, *h.JwtProvider)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful\n"))
}

func (h *UserHandler) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
        return
    }
    var req app.RefreshJwtRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    resp, err := h.userService.RefreshAccessToken(req)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) RefreshRefreshToken(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
        return
    }
    var req app.RefreshJwtRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    // Получаем старый accessToken из заголовка Authorization
    accessToken := r.Header.Get("Authorization")
    accessToken = strings.TrimPrefix(accessToken, "Bearer ")
    resp, err := h.userService.RefreshRefreshToken(req, accessToken)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}