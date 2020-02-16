package domain

type Direction struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type InputContext struct {
	PlayerId string `json:"playerId"`
	SessionCode string `json: "sessionCode"`
	Direction Direction `json: "direction"`
	Attack bool `json: "attack"`
}
