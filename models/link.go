package models

import "time"

type Link struct {
	ID        uint32    `gorm:"primaryKey" json:"ID" xml:"ID"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"UpdatedAt" xml:"UpdatedAt" example:"2021-08-27T16:17:53.119571+03:00"`
	Name      string    `gorm:"not null" json:"Name" xml:"Name" example:"main"`
	Link      string    `gorm:"not null" json:"Link" xml:"Link" example:"https://github.com/YushchenkoAndrew/template"`
	ProjectID uint32    `gorm:"foreignKey:ProjectID;not null" json:"ProjectID" xml:"ProjectID" example:"1"`
	// Project   Project   `gorm:""`
}

func (*Link) TableName() string {
	return "link"
}

type ReqLink struct {
	// ID        uint32    `json:"id" xml:"id"`
	Name string `json:"Name" xml:"Name"`
	Link string `json:"Link" xml:"Link"`
}
