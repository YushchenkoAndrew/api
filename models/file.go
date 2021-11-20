package models

import "time"

type File struct {
	ID        uint32    `gorm:"primaryKey" json:"id" xml:"id"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at" xml:"updated_at" example:"2021-08-27T16:17:53.119571+03:00"`
	Name      string    `gorm:"not null" json:"name" xml:"name" example:"index.js"`
	Path      string    `json:"path" xml:"path" example:"/test"`
	Type      string    `gorm:"not null" json:"type" xml:"type" example:"js"`
	Role      string    `gorm:"not null" json:"role" xml:"role" example:"src"`
	ProjectID uint32    `gorm:"foreignKey:ProjectID;not null" json:"project_id" xml:"project_id" example:"1"`
	// Project   Project   `gorm:""`
}

func (*File) TableName() string {
	return "file"
}

type ReqFile struct {
	// ID        uint32    `json:"id" xml:"id"`
	Name string `json:"name" xml:"name"`
	Path string `json:"path" xml:"path"`
	Type string `json:"type" xml:"type"`
	Role string `json:"role" xml:"role"`
}
