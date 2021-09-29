package models

import "time"

type Project struct {
	ID        uint32    `gorm:"primaryKey" json:"ID" xml:"ID"`
	CreatedAt time.Time `gorm:"default:CURRENT_DATE" json:"CreatedAt" xml:"CreatedAt" example:"2021-08-06"`
	Name      string    `gorm:"not null" json:"Name" xml:"Name" example:"Code Rain"`
	Title     string    `json:"Title" xml:"Title" example:"Take the blue pill and the sit will close, or take the red pill and I show how deep the rebbit hole goes"`
	Desc      string    `json:"Desc" xml:"Desc" example:"Creating a 'Code Rain' effect from Matrix. As funny joke you can put any text to display at the end."`
	Files     []File    `gorm:"foreignKey:ProjectID"`
}

func (*Project) TableName() string {
	return "project"
}

type ReqProject struct {
	// ID        uint32    `json:"id" xml:"id"`
	Name  string `json:"Name" xml:"Name"`
	Title string `json:"Title" xml:"Title"`
	Desc  string `json:"Desc" xml:"Desc"`
	Files []File `json:"Files,omitempty" xml:"Files,omitempty"`
}
