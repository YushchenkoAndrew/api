package models

import (
	"api/interfaces"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
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

	ProjectID      uint32 `gorm:"foreignKey:ProjectID;not null" json:"project_id" xml:"project_id" example:"1"`
	SubscriptionID uint32 `gorm:"foreignKey:SubscriptionID;not null" json:"subscription_id" xml:"subscription_id" example:"1"`
}

func NewMetrics() interfaces.Table {
	return &Metrics{}
}

func (*Metrics) TableName() string {
	return "metrics"
}

func (c *Metrics) Migrate(db *gorm.DB, forced bool) {
	if forced {
		db.Migrator().DropTable(c)
	}
	db.AutoMigrate(c)
}

func (*Metrics) Redis(*gorm.DB, *redis.Client) error {
	return nil
}
