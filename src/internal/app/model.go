package app

import "github.com/google/uuid"


type Field [][]int

type State int
const (
	Wait = 0
	Move = 1
	Draw = 2
	Win = 3
)

type CurrentGame struct{
	Field Field	
	UUID uuid.UUID
	Status State	
	Computer bool
	PlayerX uuid.UUID
	PlayerO uuid.UUID
}

type User struct{
	UUID uuid.UUID
	Login string
	Password string
}

type SignUpRequest struct{
	Login string
	Password string
}