package datasource

import (
	"krestikinoliki/internal/app"
)

func ToEntity(CurrentGame *app.CurrentGame) *GameEntity{
	return &GameEntity{
		ID:CurrentGame.UUID,
		Field: CurrentGame.Field,
	}
}

func FromEntity(Entity *GameEntity) *app.CurrentGame{
	return &app.CurrentGame{
		Field: Entity.Field,
		UUID: Entity.ID,
	}
}