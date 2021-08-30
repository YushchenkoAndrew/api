package models

import "time"

type World struct {
	ID        uint32    `gorm:"primaryKey" json:"ID" xml:"ID"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"UpdatedAt" xml:"UpdatedAt" example:"2021-08-27T16:17:53.119571+03:00"`
	Country   string    `gorm:"size:2;unique" json:"Country" xml:"Country" example:"UK"`
	Visitors  uint16    `gorm:"default:0" json:"Visitors" xml:"Visitors" example:"5"`
}

func (*World) TableName() string {
	return "world"
}

type ReqWorld struct {
	// ID        uint32
	// UpdatedAt time.Time
	Country  string  `json:"Country" xml:"Country"`
	Visitors *uint16 `json:"Visitors" xml:"Visitors"`
}
