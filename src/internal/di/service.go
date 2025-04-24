package di

import (
		"context"
		"log"
		"net/http"
		"krestikinoliki/internal/web"
		"go.uber.org/fx"
)

func StartHTTPServer(lc fx.Lifecycle, handler *web.GameHandler, user_handler *web.UserHandler) {
	// Регистрируем маршруты
	mux := http.NewServeMux()
	web.RegisterRoutes(mux, handler, user_handler) // новый метод, если нужно

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux, // важно!
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Server started on http://localhost:8080")
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down server...")
			return server.Close()
		},
	})
}


// func main() {
// 	// Репозиторий — пока можно мок, или in-memory реализацию
// 	repo := datasource.NewGameStorage()

// 	// Сервис
// 	service := datasource.NewGameServiceImpl(repo)

// 	// Хендлер
// 	handler := web.NewGameHandler(service)

// 	// Регистрируем маршруты
// 	web.RegisterRoutes(handler)

// 	log.Println("Server started on http://localhost:8080")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }


// Эта функция будет вызвана автоматически и поднимет HTTP-сервер