package storage

import (
	"animated-robot/domain"
	"github.com/google/uuid"
)

func NewPlayer(name string) domain.Player {
	playerUuid, _ := uuid.NewRandom()
	playerId := playerUuid.String()
	return domain.Player{
		Id:   playerId,
		Name: name,
	}
}
