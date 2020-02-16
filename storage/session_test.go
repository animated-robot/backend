package storage

import (
	"animated-robot/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MyGenerator struct {
	Code string
}

func (mg MyGenerator) Generate() string {
	return mg.Code
}

func TestNewSessionStoreInMemory(t *testing.T) {
	code := "my test"
	playerId := "player id"
	socketId := "front socket id"

	myGenerator := MyGenerator{
		Code: code,
	}

	sessionStore := NewSessionStoreInMemory(myGenerator)

	createdSessionCode, _ := sessionStore.Create(socketId)

	assert.Equal(t, code, createdSessionCode)

	session, _ := sessionStore.Get(code)

	assert.NotNil(t, session)
	assert.Equal(t, session.Code, code)

	sessionStore.AddPlayer(code, domain.Player{
		Id:   playerId,
		Name: "Joaozinho",
	})

	player, _ := sessionStore.GetPlayer(code, playerId)

	assert.NotNil(t, player)
	assert.Equal(t, playerId, player.Id)

	sessionStore.RemovePlayer(code, playerId)

	_, notFoundPlayer := sessionStore.GetPlayer(code, playerId)

	assert.NotNil(t, notFoundPlayer)

	sessionStore.Delete(code)

	_, notFoundSession := sessionStore.Get(code)

	assert.NotNil(t, notFoundSession)
}