package models

type Ping struct {
	Stat    string `json:"stat" example:"OK"`
	Message string `json:"message" example:"pong"`
}
