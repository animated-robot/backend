package storage

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MyGenerator struct {
	Code string
}

func (mg MyGenerator) Generate() string {
	return mg.Code
}

func TestJson(t *testing.T) {
	type Player map[string]interface{}
	str := `{ 
		"name": "joao",
		"color": "black",
		"height": 1.85
	}`

	var player Player
	err := json.Unmarshal([]byte(str), &player)
	if err != nil {
		fmt.Println(err.Error())
	}
	player["id"] = "lçkajsflkasjlksafj"

	p, err := json.Marshal(player)
	if err != nil {
		fmt.Println(err.Error())
	}
	playerJson := string(p)
	fmt.Printf("player: %s", playerJson)
}

func TestNewSessionStoreInMemory(t *testing.T) {
	code := "my test"
	playerId := uuid.New()
	socketId := "front socket id"

	myGenerator := MyGenerator{
		Code: code,
	}

	sessionStore := NewSessionStoreInMemory(myGenerator)

	createdSessionCode, _ := sessionStore.Create(socketId)

	assert.Equal(t, code, createdSessionCode)

	session, _ := sessionStore.Get(code)

	assert.NotNil(t, session)
	assert.Equal(t, session.SessionCode, code)

	sessionStore.AddPlayer(code, playerId)

	sessionStore.RemovePlayer(code, playerId)

	sessionStore.Delete(code)

	_, notFoundSession := sessionStore.Get(code)

	assert.NotNil(t, notFoundSession)
}