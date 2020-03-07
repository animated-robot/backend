package domain

type Session struct {
	FrontSocketId string      `json:"socketId"`
	SessionCode   string      `json:"sessionCode"`
	Players       []Player    `json:"players"`
}
