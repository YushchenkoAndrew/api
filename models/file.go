package models

import "time"

type File struct {
	ID        uint32    `gorm:"primaryKey" json:"ID" xml:"ID"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"UpdatedAt" xml:"UpdatedAt" example:"2021-08-27T16:17:53.119571+03:00"`
	Name      string    `gorm:"not null" json:"Name" xml:"Name" example:"index.js"`
	Type      string    `gorm:"not null" json:"Type" xml:"Type" example:"js"`
	Role      string    `gorm:"not null" json:"Role" xml:"Role" example:"src"`
	ProjectID uint32    `gorm:"foreignKey:ProjectID;not null" json:"ProjectID" xml:"ProjectID" example:"1"`
	// Project   Project   `gorm:""`
}

func (*File) TableName() string {
	return "file"
}

type ReqFile struct {
	// ID        uint32    `json:"id" xml:"id"`
	Name string `json:"Name" xml:"Name"`
	Type string `json:"Type" xml:"Type"`
	Role string `json:"Role" xml:"Role"`
}
