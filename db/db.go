package db

import (
	"api/config"
	"api/interfaces"
	"api/logs"
	"api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// FIXME: Should fix this one day
var DB *gorm.DB

func ConnectToDB(tables []interfaces.Table) {
	var err error
	DB, err = gorm.Open(postgres.Open(
		"host="+config.ENV.DBHost+
			" user="+config.ENV.DBUser+
			" password="+config.ENV.DBPass+
			" port="+config.ENV.DBPort+
			" dbname="+config.ENV.DBName), &gorm.Config{})

	if err != nil {
		logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			File:    "/db/db.go",
			Message: "Bruhhh, did you even start the Postgres ???",
			Desc:    err.Error(),
		})
		panic("Failed on db connection")
	}

	for _, table := range tables {
		table.Migrate(DB, config.ENV.ForceMigrate)
	}
}
