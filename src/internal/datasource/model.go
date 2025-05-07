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
}