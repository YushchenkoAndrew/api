package models

import (
	"api/interfaces"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type SubscribeDto struct {
	CronTime  string `json:"cron_time" xml:"cron_time" example:"00 00 00 */1 * *"`
	Operation string `json:"operation" xml:"operation" example:"metrics"`
}

type Subscription struct {
	ID        uint32    `gorm:"primaryKey" json:"id" xml:"id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at" xml:"created_at" example:"2021-08-06"`
	CronID    string    `json:"cron_id" xml:"cron_id" example:"d266389ebf09e1e8a95a5b4286b504b2"`
	CronTime  string    `json:"cron_time" xml:"cron_time" example:"00 00 00 */1 * *"`
	Operation string    `json:"operation" xml:"operation" example:"metrics"`
	Metrics   []Metrics `gorm:"foreignKey:SubscriptionID" json:"metrics" xml:"metrics"`
}

func NewSubscription() interfaces.Table {
	return &Subscription{}
}

func (*Subscription) TableName() string {
	return "subscription"
}

func (c *Subscription) Migrate(db *gorm.DB, forced bool) {
	if forced {
		db.Migrator().DropTable(c)
	}
	db.AutoMigrate(c)
}

func (*Subscription) Redis(*gorm.DB, *redis.Client) error {
	return nil
}
