package app

import "github.com/google/uuid"


type Field [][]int

type State int
const (
	Wait = 0
	MoveX = 1
	MoveO = 2
	Draw = 3
	WinX = 4
	WinO = 5
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