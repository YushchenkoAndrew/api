package interfaces

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Table interface {
	Migrate(db *gorm.DB, forced bool)
	Redis(db *gorm.DB, client *redis.Client) error
}
