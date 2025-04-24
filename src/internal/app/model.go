package app

import "github.com/google/uuid"


type Field [][]int

type CurrentGame struct{
	Field Field	
	UUID uuid.UUID
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