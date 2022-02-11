package helper

import (
	"api/config"
	"api/db"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func PrecacheResult(key string, client *gorm.DB, model interface{}) error {
	ctx := context.Background()

	// Check if cache have requested data
	if data, err := db.Redis.Get(ctx, key).Result(); err == nil {
		json.Unmarshal([]byte(data), model)
		go db.Redis.Expire(ctx, key, time.Duration(config.ENV.LiveTime)*time.Second)
	} else {
		result := client.Find(model)
		if result.Error != nil {
			return fmt.Errorf("Server side error: Something went wrong - %v", result.Error)
		}

		// Encode json to strO
		if str, err := json.Marshal(model); err == nil {
			go db.Redis.Set(ctx, key, str, time.Duration(config.ENV.LiveTime)*time.Second)
		}
	}
	return nil
}
