package datasource

import (
	"errors"
	"krestikinoliki/internal/app"
	"sync"
	"github.com/google/uuid"
)

type GameStorage struct{
	data sync.Map
}

func NewGameStorage() *GameStorage{
	return &GameStorage{}
}

func (s *GameStorage) SaveGame(currentgame *app.CurrentGame) error {
	entity := ToEntity(currentgame)
	s.data.Store(entity.ID.String(), entity)
	return nil
}

func (s* GameStorage) LoadGame(ID uuid.UUID) (*app.CurrentGame, error){

	value, err := s.data.Load(ID.String())

	if err{
		return nil, errors.New("did not found game")
	}
	entity := value.(*GameEntity)
	return FromEntity(entity), nil
}