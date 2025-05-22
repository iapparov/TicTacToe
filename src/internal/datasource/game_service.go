package datasource

import (
	"krestikinoliki/internal/app"

	"github.com/google/uuid"
)


type GameServiceImpl struct{
	core *app.TicTacToeService // структура
	repo GameRepository //интерфейс
}

func (s *GameServiceImpl) LoadGame(id uuid.UUID) (*app.CurrentGame, error) {
	return s.repo.LoadGame(id)
}

func NewGameServiceImpl(repo GameRepository) *GameServiceImpl{ //принимает структуру подходяющую под интерфейс
	return &GameServiceImpl{
		repo: repo,
		core: &app.TicTacToeService{},
	}
}

// NextMove реализует логику хода.
func (s *GameServiceImpl) NextMove(game *app.CurrentGame) (*app.CurrentGame, error) {

	updated, err := s.core.NextMove(game)
	if err != nil {
		return nil, err
	}
	// Сохраняем игру в репозитории после хода.
	if err := s.repo.SaveGame(updated); err != nil {
		return nil, err
	}
	return updated, nil
}

// FieldValidation проверяет валидность игрового поля.
func (s *GameServiceImpl) FieldValidation(game *app.CurrentGame) (bool, error) {
	return s.core.FieldValidation(game)
}

// GameIsOver проверяет, завершена ли игра.
func (s *GameServiceImpl) GameIsOver(game *app.CurrentGame) bool  {
	return s.core.GameIsOver(game)
}

func (s *GameServiceImpl) Connect(game *app.CurrentGame, Uuidgame string, Uuidplayero string) (*app.CurrentGame){
	tmp, err := uuid.Parse(Uuidgame)
	if err != nil {
		return game
	}
	game, err = s.repo.LoadGame(tmp)
	if err !=nil {
		return nil
	}
	game = s.core.Connect(game, Uuidgame, Uuidplayero)
	s.repo.SaveGame(game)
	return game
}

func (s *GameServiceImpl) NewGame(Computer bool, Uuid string) (*app.CurrentGame) {
	game := s.core.NewGame(Computer, Uuid)
	s.repo.SaveGame(game)
	return game
}

//Добавить сохранение юзер id в модель