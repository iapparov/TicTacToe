package web

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, gameHandler *GameHandler, userHandler *UserHandler) {
	mux.HandleFunc("/register", userHandler.Register)
	mux.HandleFunc("/login", userHandler.Login)

	protected := &UserAuthenticator{
		Next: gameHandler, 
		Repo: userHandler.userRepo,
	}
	mux.Handle("/game/", protected)
	mux.Handle("/game", protected)
	mux.Handle("/connect/", protected)
	mux.Handle("/getgames", protected)
}
