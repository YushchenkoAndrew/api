package models

import (
	"api/config"
	"api/interfaces"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type GeoIpBlocks struct {
	Network                     string `csv:"network" gorm:"type:cidr"`
	GeonameId                   int64  `csv:"geoname_id"`
	RegisteredCountryGeonameId  int64  `csv:"registered_country_geoname_id"`
	RepresentedCountryGeonameId int64  `csv:"represented_country_geoname_id"`
	IsAnonymousProxy            bool   `csv:"is_anonymous_proxy"`
	IsSatelliteProvider         bool   `csv:"is_satellite_provider"`
}

func NewGeoIpBlocks() interfaces.Table {
	return &GeoIpBlocks{}
}

func (*GeoIpBlocks) TableName() string {
	return "geo_ip_blocks"
}

func (c *GeoIpBlocks) Migrate(db *gorm.DB, forced bool) {
	if forced {
		db.Migrator().DropTable(c)
	}

	db.AutoMigrate(c)

	var nSize int64
	if db.Model(c).Count(&nSize); nSize == 0 {

		// The most quick and easiest way !!!
		db.Exec("\\copy geo_ip_blocks from '" + config.ENV.MigrationPath + "/GeoLite2-Country-Blocks.csv' delimiter ',' csv header;")
		db.Exec("CREATE INDEX geo_ip_blocks_network_idx ON geo_ip_blocks USING gist (network inet_ops);")
	}
}

func (*GeoIpBlocks) Redis(*gorm.DB, *redis.Client) error {
	return nil
}
