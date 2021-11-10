package db

import (
	"api/config"
	"api/models"
)

func Link() {
	if config.ENV.ForceMigrate {
		DB.Migrator().DropTable(&models.Link{})
	}
	DB.AutoMigrate(&models.Link{})
}
