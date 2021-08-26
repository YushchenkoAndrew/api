package db

import (
	"api/config"
	models "api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	DB, err = gorm.Open(postgres.Open(
		"host="+config.ENV.DBHost+
			" user="+config.ENV.DBUser+
			" password="+config.ENV.DBPass+
			" port="+config.ENV.DBPort+
			" dbname="+config.ENV.DBName), &gorm.Config{})
	config.CheckOnErr(&err, "Failed on db connection")
}

func MigrateTables(bForce bool) {
	if bForce {
		DB.Migrator().DropTable(
			&models.Info{},
			&models.World{},
		)
	}

	DB.AutoMigrate(
		&models.Info{},
		&models.World{},
	)
}
