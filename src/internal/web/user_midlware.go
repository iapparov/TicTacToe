package web

import (
	"krestikinoliki/internal/app"
	"net/http"
)

type UserAuthenticator struct {
	Next http.Handler
	Repo app.UserRepository
}

func (ua *UserAuthenticator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/register" || r.URL.Path == "/login" {
		ua.Next.ServeHTTP(w, r)
		return
	}

	cookie, err := r.Cookie("user_id")
	if err != nil || cookie.Value == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if !ua.Repo.FindByUUID(cookie.Value) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// todo: можно проверить uuid.Parse(cookie.Value), если нужно
	ua.Next.ServeHTTP(w, r)
}