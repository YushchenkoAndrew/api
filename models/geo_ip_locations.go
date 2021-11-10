package models

type GeoIpLocations struct {
	GeonameId         int64  `csv:"geoname_id"`
	LocaleCode        string `csv:"locale_code" gorm:"size:2"`
	ContinentCode     string `csv:"continent_code" gorm:"size:2"`
	ContinentName     string `csv:"continent_name" gorm:"size:255"`
	CountryIsoCode    string `csv:"country_iso_code" gorm:"size:2"`
	CountryName       string `csv:"country_name" gorm:"size:255"`
	IsInEuropeanUnion bool   `csv:"is_in_european_union"`
}

func (*GeoIpLocations) TableName() string {
	return "geo_ip_locations"
}
