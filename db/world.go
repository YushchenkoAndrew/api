package db

import (
	"api/config"
	"api/models"
)

func World() {
	if config.ENV.ForceMigrate {
		DB.Migrator().DropTable(&models.World{})
	}
	DB.AutoMigrate(&models.World{})
}
