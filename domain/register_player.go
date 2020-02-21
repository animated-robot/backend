package domain

type RegisterPlayer struct {
	SessionCode string `json:"sessionCode"`
	Player Player `json:"player"`
}
