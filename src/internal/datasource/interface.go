package datasource

import (
	"krestikinoliki/internal/app"
	"github.com/google/uuid"
)

type GameRepository interface{
	SaveGame(currentgame *app.CurrentGame) error
	LoadGame(ID uuid.UUID) (*app.CurrentGame, error)
	GetGames() []string
	CurrentGame(Userid string) []string
	GetEndedGames(uuid string) ([]string)
	GetLeaderBoard(count int) ([]app.LeaderBoard, error)
}