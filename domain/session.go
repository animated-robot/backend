package domain

import "github.com/google/uuid"

type Session struct {
	FrontSocketId string `json:"socketId"`
	Code string	`json:"sessionCode"`
	PlayersIds []uuid.UUID `json:"playersIds"`
}

