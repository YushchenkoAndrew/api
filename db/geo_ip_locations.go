package db

import (
	"api/config"
	"api/models"
	"fmt"
)

func GeoIpLocations() {
	if config.ENV.ForceMigrate {
		DB.Migrator().DropTable(&models.GeoIpLocations{})
	}

	DB.AutoMigrate(&models.GeoIpLocations{})

	var nSize int64
	if DB.Model(&models.GeoIpLocations{}).Count(&nSize); nSize == 0 {

		// The most quick and easiest way !!!
		DB.Exec(fmt.Sprintf("copy geo_ip_locations from '%s/GeoLite2-Country-Locations-en.csv' delimiter ',' csv header;", config.ENV.MigrationPath))
	}
}
