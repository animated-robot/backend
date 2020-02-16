package domain

type RegisterPlayer struct {
	SessionCode string `json:"sessionCode"`
	PlayerName string `json:"playerName"`
}
