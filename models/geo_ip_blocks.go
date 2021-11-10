package models

type GeoIpBlocks struct {
	Network                     string `csv:"network" gorm:"type:cidr"`
	GeonameId                   int64  `csv:"geoname_id"`
	RegisteredCountryGeonameId  int64  `csv:"registered_country_geoname_id"`
	RepresentedCountryGeonameId int64  `csv:"represented_country_geoname_id"`
	IsAnonymousProxy            bool   `csv:"is_anonymous_proxy"`
	IsSatelliteProvider         bool   `csv:"is_satellite_provider"`
}

func (*GeoIpBlocks) TableName() string {
	return "geo_ip_blocks"
}
