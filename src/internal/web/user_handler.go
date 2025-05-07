package web

import (
	"encoding/json"
	"krestikinoliki/internal/app"
	"net/http"
)

type UserHandler struct {
	userService app.UserService
	userRepo app.UserRepository
}

func NewUserHandler(service app.UserService, repo app.UserRepository) *UserHandler {
	return &UserHandler{
		userService: service,
		userRepo: repo,
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

	var req app.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userID, err := h.userService.Authenticate(req.Login, req.Password)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Простая кука авторизации
	http.SetCookie(w, &http.Cookie{
		Name:  "user_id",
		Value: userID.String(),
		Path:  "/",
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful\n"))
}