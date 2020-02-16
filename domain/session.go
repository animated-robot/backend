package domain

type Session struct {
	FrontSocketId string
	Code string
	Players []Player
}

