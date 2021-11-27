package models

import "time"

type Project struct {
	ID        uint32    `gorm:"primaryKey" json:"id" xml:"id"`
	CreatedAt time.Time `gorm:"default:CURRENT_DATE" json:"created_at" xml:"created_at" example:"2021-08-06"`
	Name      string    `gorm:"not null;unique" json:"name" xml:"name" example:"CodeRain"`
	Title     string    `gorm:"not null" json:"title" xml:"title" example:"Code Rain"`
	Flag      string    `json:"flag" xml:"flag" example:"js"`
	Desc      string    `json:"desc" xml:"desc" example:"Take the blue pill and the sit will close, or take the red pill and I show how deep the rebbit hole goes"`
	Note      string    `json:"note" xml:"note" example:"Creating a 'Code Rain' effect from Matrix. As funny joke you can put any text to display at the end."`
	Files     []File    `gorm:"foreignKey:ProjectID" json:"files" xml:"files"`
	Links     []Link    `gorm:"foreignKey:ProjectID" json:"links" xml:"links"`
}

func (*Project) TableName() string {
	return "project"
}

type ReqProject struct {
	// ID        uint32    `json:"id" xml:"id"`
	Name  string `json:"name" xml:"name"`
	Title string `json:"title" xml:"title"`
	Flag  string `json:"flag" xml:"flag"`
	Desc  string `json:"desc" xml:"desc"`
	Note  string `json:"note" xml:"note"`
	Files []File `json:"files,omitempty" xml:"files,omitempty"`
	Links []Link `json:"links,omitempty" xml:"links,omitempty"`
}
