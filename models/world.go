package models

import "time"

type World struct {
	ID        uint32    `gorm:"primaryKey"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Country   string    `gorm:"size:2;unique"`
	Visitors  uint16    `gorm:"default:0"`
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
