package db

import (
	"api/config"
	"api/models"
)

func Project() {
	if config.ENV.ForceMigrate {
		DB.Migrator().DropTable(&models.Project{})
	}
	DB.AutoMigrate(&models.Project{})
}
