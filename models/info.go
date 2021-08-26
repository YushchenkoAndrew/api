package models

import "time"

type Info struct {
	ID        uint32    `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"default:CURRENT_DATE;unique"`
	Countries string
	Views     uint16 `gorm:"default:0"`
	Clicks    uint16 `gorm:"default:0"`
	Media     uint16 `gorm:"default:0"`
	Visitors  uint16 `gorm:"default:0"`
}

func (*Info) TableName() string {
	return "info"
}

type ReqInfo struct {
	// ID        uint32    `json:"id" xml:"id"`
	CreatedAt *time.Time `json:"CreatedAt" xml:"CreatedAt"`
	Countries string     `json:"Countries" xml:"Contries"`
	Views     *uint16    `json:"Views,omitempty" xml:"Views"`
	Clicks    *uint16    `json:"Clicks,omitempty" xml:"Clicks"`
	Media     *uint16    `json:"Media,omitempty" xml:"Media"`
	Visitors  *uint16    `json:"Visitors,omitempty" xml:"Visitors"`
}

type StatInfo struct {
	Views    uint16
	Clicks   uint16
	Media    uint16
	Visitors uint16
}
