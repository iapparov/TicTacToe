package web

import (
	"krestikinoliki/internal/app"
	"net/http"
	"strings"
)

type UserAuthenticator struct {
	Next http.Handler
	Repo app.UserRepository
	JwtProvider *app.JwtProvider
}

func (ua *UserAuthenticator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/register" || r.URL.Path == "/login" {
		ua.Next.ServeHTTP(w, r)
		return
	}

    authHeader := r.Header.Get("Authorization")
    if !strings.HasPrefix(authHeader, "Bearer ") {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
	
    token := strings.TrimPrefix(authHeader, "Bearer ")
    id, err := ua.JwtProvider.ValidateAccessToken(token)
	flag, _ := ua.Repo.FindByUUID(id.String())
    if err != nil || !flag{
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

	// todo: можно проверить uuid.Parse(cookie.Value), если нужно
	ua.Next.ServeHTTP(w, r)
}