package datasource

import (
	"krestikinoliki/internal/app"
)

func ToEntity(CurrentGame *app.CurrentGame) *GameEntity{
	return &GameEntity{
		ID: CurrentGame.UUID,
		Field: CurrentGame.Field,
		Status: CurrentGame.Status,
		Computer: CurrentGame.Computer,
		PlayerX: CurrentGame.PlayerX,
		PlayerO: CurrentGame.PlayerO,
	}
}

func FromEntity(Entity *GameEntity) *app.CurrentGame{
	return &app.CurrentGame{
		Field: Entity.Field,
		UUID: Entity.ID,
		Status: Entity.Status,
		Computer: Entity.Computer,
		PlayerX: Entity.PlayerX,
		PlayerO: Entity.PlayerO,
	}
}