package db

import (
	"api/config"
	"api/models"
)

func Info() {
	if config.ENV.ForceMigrate {
		DB.Migrator().DropTable(&models.Info{})
	}
	DB.AutoMigrate(&models.Info{})
}
