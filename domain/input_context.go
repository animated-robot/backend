package domain

type Direction struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type InputContext struct {
	Player        Player    `json:"player"`
	SessionCode   string    `json: "sessionCode"`
	Direction     Direction `json: "direction"`
	ActiveActions []string      `json: "activeActions"`
}
