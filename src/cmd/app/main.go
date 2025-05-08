package main

import (
	"krestikinoliki/internal/datasource"
	"krestikinoliki/internal/app"
	"krestikinoliki/internal/web"
	"krestikinoliki/internal/di"
	"go.uber.org/fx"
)



func main() {
	app := fx.New(
		// Провайдеры зависимостей
		fx.Provide(
			datasource.ConnectDB, // коннект к бд
			datasource.NewPostgresGameRepo, //БД + Игра
			datasource.NewPostgresUserRepo, //БД + ЮЗЕР
			app.NewUserServiceImpl, // Юзерсервис
		
			// Интерфейс GameRepository
			func(repo *datasource.PostgresGameRepo) datasource.GameRepository {
				return repo
			},
		
			// Интерфейс GameService
			func(repo datasource.GameRepository) app.GameService {
				return datasource.NewGameServiceImpl(repo)
			},
		
			// Интерфейс UserRepository
			func(repo *datasource.PostgresUserRepo) app.UserRepository {
				return repo
			},
		
			// Интерфейс UserService
			func(s *app.UserServiceImpl) app.UserService {
				return s
			},
		
			web.NewUserHandler,
			web.NewGameHandler,
		),

		// Регистрация HTTP-сервера
		fx.Invoke(di.StartHTTPServer),
	)

	app.Run()
}