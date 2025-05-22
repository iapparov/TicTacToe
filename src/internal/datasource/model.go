package datasource

import (
	"krestikinoliki/internal/app"
	"github.com/google/uuid"
)


type GameEntity struct{
	Field app.Field	`db:"Field"`
	ID uuid.UUID `db:"Id"`
	Status app.State `db:"Status"`
	Computer bool `db:"Computer"`
	PlayerX uuid.UUID `db:"playerx"`
	PlayerO uuid.UUID `db:"playero"`
	CreatedAt int64 `db:"createdat"`
}