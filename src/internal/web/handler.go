package web

import (
	"encoding/json"
	"fmt"
	"io"
	"krestikinoliki/internal/app"
	"krestikinoliki/internal/datasource"
	"net/http"
	"strings"

)

type GameHandler struct {
	service app.GameService //интерфейс
	repo datasource.GameRepository //интерфейс
}

func NewGameHandler(service app.GameService, repo datasource.GameRepository) *GameHandler {
	return &GameHandler{service: service, repo: repo}
}

// POST /game/{id}
func (h *GameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path == "/game"{
		h.CreateGame(w,r)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/connect/"){
		h.Connect(w,r)
		return
	}

	if r.URL.Path == "/getgames"{
		h.GetGames(w,r)
		return
	}

	h.PlayGame(w,r)
}

func (h *GameHandler) GetGames(w http.ResponseWriter, r *http.Request){
	games := h.repo.GetGames()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(games); err != nil {
		http.Error(w, "Failed to encode games", http.StatusInternalServerError)
		return
	}
}	

func (h *GameHandler) PlayGame(w http.ResponseWriter, r *http.Request){
	// Извлечение ID из URL: /game/{id}
	id := strings.TrimPrefix(r.URL.Path, "/game/")
	if id == "" {
		http.Error(w, "Game ID is required", http.StatusBadRequest)
		return
	}
	

	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Can't read body", http.StatusBadRequest)
		return
	}
	// Декодируем JSON в DTO
	var dto GameDTO
	if err := json.Unmarshal(body, &dto); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Маппим DTO в доменную модель
	game, err := FromWeb(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Валидируем ход
	check, err := h.service.FieldValidation(game)
	if  !check {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Проверяем победителя
	isgameover := h.service.GameIsOver(game)
	if isgameover {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Game is over!\n"))
		return
	}


	// Выполняем ход
	updatedGame, err := h.service.NextMove(game)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to process move: %v", err), http.StatusBadRequest)
		return
	}

	// Маппим обратно в DTO
	responseDTO := ToWeb(updatedGame)

	// Отправляем JSON обратно
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseDTO)
}

// POST /game — создать новую игру
func (h *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	type vs_computer struct{
		Computer bool `json:"vs_computer"`
	}
	var req vs_computer
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req.Computer = true
	}
	cookie, err := r.Cookie("user_id")
	if err != nil{
		return
	}
	newGame := h.service.NewGame(req.Computer, cookie.Value)

    dto := ToWeb(newGame) // маппер domain -> web
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(dto)
}

// POST /game — создать новую игру
func (h *GameHandler) Connect(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/connect/")
	if id == "" {
		http.Error(w, "Game ID is required", http.StatusBadRequest)
		return
	}
	cookie, err := r.Cookie("user_id")
	if err != nil{
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	var game *app.CurrentGame
	game = h.service.Connect(game, id, cookie.Value)

    dto := ToWeb(game) // маппер domain -> web
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(dto)
}