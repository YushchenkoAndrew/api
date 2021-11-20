package models

import "time"

type Link struct {
	ID        uint32    `gorm:"primaryKey" json:"id" xml:"id"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at" xml:"updated_at" example:"2021-08-27T16:17:53.119571+03:00"`
	Name      string    `gorm:"not null" json:"name" xml:"name" example:"main"`
	Link      string    `gorm:"not null" json:"link" xml:"link" example:"https://github.com/YushchenkoAndrew/template"`
	ProjectID uint32    `gorm:"foreignKey:ProjectID;not null" json:"project_id" xml:"project_id" example:"1"`
	// Project   Project   `gorm:""`
}

func (*Link) TableName() string {
	return "link"
}

type ReqLink struct {
	// ID        uint32    `json:"id" xml:"id"`
	Name string `json:"name" xml:"name"`
	Link string `json:"link" xml:"link"`
}
