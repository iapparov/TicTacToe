package web

import (
	"krestikinoliki/internal/app"
	"errors"
	"github.com/google/uuid"
)

func ToWeb(CurrentGame *app.CurrentGame) *GameDTO {

	Idstr := CurrentGame.UUID.String()
	return &GameDTO{
		Field: CurrentGame.Field,
		ID: Idstr,
	}
}

func FromWeb(CurrentGame *GameDTO) (*app.CurrentGame, error){

	IDuuid, err := uuid.Parse(CurrentGame.ID)
	if err != nil{
		return nil, errors.New(CurrentGame.ID)
	}
	return &app.CurrentGame{
		Field: CurrentGame.Field,
		UUID: IDuuid,
	}, nil
}

