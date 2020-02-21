package domain

import "github.com/google/uuid"

type Session struct {
	FrontSocketId string      `json:"socketId"`
	SessionCode   string      `json:"sessionCode"`
	PlayersIds    []uuid.UUID `json:"playersIds"`
}

