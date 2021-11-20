package models

import "time"

type Info struct {
	ID        uint32    `gorm:"primaryKey" json:"id" xml:"id"`
	CreatedAt time.Time `gorm:"default:CURRENT_DATE;unique" json:"created_at" xml:"created_at" example:"2021-08-06"`
	Countries string    `json:"countries" xml:"contries" example:"UK,US"`
	Views     uint16    `gorm:"default:0" json:"views" xml:"views" example:"1"`
	Clicks    uint16    `gorm:"default:0" json:"clicks" xml:"clicks" example:"2"`
	Media     uint16    `gorm:"default:0" json:"media" xml:"media" example:"3"`
	Visitors  uint16    `gorm:"default:0" json:"visitors" xml:"visitors" example:"4"`
}

func (*Info) TableName() string {
	return "info"
}

type ReqInfo struct {
	// ID        uint32    `json:"id" xml:"id"`
	// CreatedAt *time.Time `json:"CreatedAt" xml:"CreatedAt"`
	Countries string  `json:"countries" xml:"contries"`
	Views     *uint16 `json:"views,omitempty" xml:"views,omitempty"`
	Clicks    *uint16 `json:"clicks,omitempty" xml:"clicks,omitempty"`
	Media     *uint16 `json:"media,omitempty" xml:"media,omitempty"`
	Visitors  *uint16 `json:"visitors,omitempty" xml:"visitors,omitempty"`
}

type StatInfo struct {
	Views    uint16
	Clicks   uint16
	Media    uint16
	Visitors uint16
}
