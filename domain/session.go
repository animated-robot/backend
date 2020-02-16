package domain

type Session struct {
	FrontSocketId string `json:"socketId"`
	Code string	`json:"sessionCode"`
	Players []Player `json:"players"`
}

