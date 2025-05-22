package web

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, gameHandler *GameHandler, userHandler *UserHandler) {
	mux.HandleFunc("/register", userHandler.Register)
	mux.HandleFunc("/login", userHandler.Login)
	mux.HandleFunc("/refresh-access", userHandler.RefreshAccessToken)
	mux.HandleFunc("/refresh-refresh", userHandler.RefreshRefreshToken)

	protected := &UserAuthenticator{
		Next: gameHandler, 
		Repo: userHandler.userRepo,
		JwtProvider: gameHandler.jwt,
	}
	mux.Handle("/game/", protected)
	mux.Handle("/game", protected)
	mux.Handle("/connect/", protected)
	mux.Handle("/getgames", protected)
	mux.Handle("/currentgame", protected)
	mux.Handle("/userinfo", protected)
	mux.Handle("/getendedgames", protected)
	mux.Handle("/getleaderboard", protected)
	
}
