package web

import (
	"krestikinoliki/internal/app"
)

type GameDTO struct{
	Field app.Field	`json:"Field"`
	ID string `json:"Id"`
}