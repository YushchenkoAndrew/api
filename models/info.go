package models

import "time"

type Info struct {
	ID        uint32    `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"default:CURRENT_DATE"`
	Countries string
	Views     uint16 `gorm:"default:0"`
	Clicks    uint16 `gorm:"default:0"`
	Media     uint16 `gorm:"default:0"`
	Visitors  uint16 `gorm:"default:0"`
}

func (*Info) TableName() string {
	return "info"
}
