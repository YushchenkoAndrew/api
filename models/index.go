package models

type Ping struct {
	Status  string `json:"status" example:"OK"`
	Message string `json:"message" example:"pong"`
}

type DefultRes struct {
	Status  string   `json:"status" example:"OK"`
	Message string   `json:"message" example:"pong"`
	Result  []string `json:"result" example:"[]"`
}
