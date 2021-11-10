package db

import (
	"api/config"
	"api/models"
	"fmt"
)

func GeoIpBlocks() {
	if config.ENV.ForceMigrate {
		DB.Migrator().DropTable(&models.GeoIpBlocks{})
	}

	DB.AutoMigrate(&models.GeoIpBlocks{})

	var nSize int64
	if DB.Model(&models.GeoIpBlocks{}).Count(&nSize); nSize == 0 {

		// The most quick and easiest way !!!
		DB.Exec(fmt.Sprintf("copy geo_ip_blocks from '%s/GeoLite2-Country-Blocks.csv' delimiter ',' csv header;", config.ENV.MigrationPath))
		DB.Exec("CREATE INDEX geo_ip_blocks_network_idx ON geo_ip_blocks USING gist (network inet_ops);")
	}
}
