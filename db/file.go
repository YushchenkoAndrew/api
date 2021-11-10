package db

import (
	"api/config"
	"api/models"
)

func File() {
	if config.ENV.ForceMigrate {
		DB.Migrator().DropTable(&models.File{})
	}
	DB.AutoMigrate(&models.File{})
}
