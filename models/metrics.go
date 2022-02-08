package models

import "time"

const (
	MILLI = "milli"
	MICRO = "micro"
	NANO  = "nano"
)

type Metrics struct {
	ID            uint32    `gorm:"type:bigint;primary_key,AUTO_INCREMENT" json:"id" xml:"id"`
	CreatedAt     time.Time `gorm:"default:CURRENT_DATE;unique" json:"created_at" xml:"created_at" example:"2021-08-06"`
	Name          string    `gorm:"not null" json:"name" xml:"name" example:"void-deployment-8985bd57d-k9n5g"`
	Namespace     string    `gorm:"not null" json:"namespace" xml:"namespace" example:"void-deployment-8985bd57d-k9n5g"`
	ContainerName string    `gorm:"not null" json:"container_name" xml:"container_name" example:"void"`
	CPU           int64     `gorm:"not null" json:"cpu" xml:"cpu" example:"690791"`
	CpuScale      uint8     `gorm:"not null" json:"cpu_scale" xml:"cpu_scale" example:"3"`
	Memory        int64     `gorm:"not null" json:"memory" xml:"memory" example:"690791"`
	MemoryScale   uint8     `gorm:"not null" json:"memory_scale" xml:"memory_scale" example:"6"`
}
