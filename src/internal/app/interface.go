package app

import "github.com/google/uuid"

type GameService interface {
	NextMove(game *CurrentGame) (*CurrentGame, error)
	FieldValidation(game *CurrentGame) (bool, error)
	GameIsOver(game *CurrentGame) bool
	NewGame(Computer bool, uuid string) (*CurrentGame)
	Connect(game *CurrentGame, Uuidgame string, Uuidplayero string) (*CurrentGame)
}

type UserService interface {
	RegisterUser(req SignUpRequest) (User, error)
	Authenticate(login, password string) (uuid.UUID, error)
}

type UserRepository interface {
	Save(user User) error
	FindByLogin(login string) (User, error)
	FindByUUID(uuid string) (bool, []string)
}