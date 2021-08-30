package models

import "time"

type Info struct {
	ID        uint32    `gorm:"primaryKey" json:"ID" xml:"ID"`
	CreatedAt time.Time `gorm:"default:CURRENT_DATE;unique" json:"CreatedAt" xml:"CreatedAt" example:"2021-08-06"`
	Countries string    `json:"Countries" xml:"Contries" example:"UK,US"`
	Views     uint16    `gorm:"default:0" json:"Views" xml:"Views" example:"1"`
	Clicks    uint16    `gorm:"default:0" json:"Clicks" xml:"Clicks" example:"2"`
	Media     uint16    `gorm:"default:0" json:"Media" xml:"Media" example:"3"`
	Visitors  uint16    `gorm:"default:0" json:"Visitors" xml:"Visitors" example:"4"`
}

func (*Info) TableName() string {
	return "info"
}

type ReqInfo struct {
	// ID        uint32    `json:"id" xml:"id"`
	// CreatedAt *time.Time `json:"CreatedAt" xml:"CreatedAt"`
	Countries string  `json:"Countries" xml:"Contries"`
	Views     *uint16 `json:"Views,omitempty" xml:"Views,omitempty"`
	Clicks    *uint16 `json:"Clicks,omitempty" xml:"Clicks,omitempty"`
	Media     *uint16 `json:"Media,omitempty" xml:"Media,omitempty"`
	Visitors  *uint16 `json:"Visitors,omitempty" xml:"Visitors,omitempty"`
}

type StatInfo struct {
	Views    uint16
	Clicks   uint16
	Media    uint16
	Visitors uint16
}
