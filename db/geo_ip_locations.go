package db

import (
	"api/config"
	"api/models"
)

func GeoIpLocations() {
	if config.ENV.ForceMigrate {
		DB.Migrator().DropTable(&models.GeoIpLocations{})
	}

	DB.AutoMigrate(&models.GeoIpLocations{})

	var nSize int64
	if DB.Model(&models.GeoIpLocations{}).Count(&nSize); nSize == 0 {

		// The most quick and easiest way !!!
		DB.Exec("copy geo_ip_locations from '" + config.ENV.MigrationPath + "/GeoLite2-Country-Locations-en.csv' delimiter ',' csv header;")
	}
}
