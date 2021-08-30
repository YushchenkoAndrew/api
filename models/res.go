package models

type Success struct {
	Status     string      `json:"status" xml:"status" example:"OK"`
	Result     interface{} `json:"result" xml:"result"`
	Page       int         `json:"page,omitempty" xml:"page,omitempty" example:"0"`
	Limit      int         `json:"limit,omitempty" xml:"limit,omitempty" example:"20"`
	Items      int64       `json:"items" xml:"items" example:"1"`
	TotalItems int64       `json:"totalItems" xml:"totalItems" example:"20"`
}

type Error struct {
	Status  string   `json:"status" xml:"status" example:"ERR"`
	Result  []string `json:"result" xml:"result" example:""`
	Message string   `json:"message" xml:"message" example:"Server side error: Something went wrong"`
}
