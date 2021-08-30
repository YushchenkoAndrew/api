package config

import (
	"github.com/spf13/viper"
)

type EnvType struct {
	Host     string `mapstructure:"HOST"`
	Port     string `mapstructure:"PORT"`
	BasePath string `mapstructure:"BASE_PATH"`

	DBType string `mapstructure:"DB_TYPE"`
	DBName string `mapstructure:"DB_NAME"`
	DBHost string `mapstructure:"DB_HOST"`
	DBPort string `mapstructure:"DB_PORT"`
	DBUser string `mapstructure:"DB_USER"`
	DBPass string `mapstructure:"DB_PASS"`

	ForceMigrate bool `mapstructure:"FORCE_MIGRATE"`

	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort string `mapstructure:"REDIS_PORT"`
	RedisPass string `mapstructure:"REDIS_PASS"`

	AccessSecret  string `mapstructure:"ACCESS_SECRET"`
	RefreshSecret string `mapstructure:"REFRESH_SECRET"`

	ID     string `mapstructure:"ID"`
	User   string `mapstructure:"USER"`
	Pass   string `mapstructure:"PASS"`
	Pepper string `mapstructure:"PEPPER"`

	LiveTime int64 `mapstructure:"LIVE_TIME"`
	Items    int   `mapstructure:"ITEMS"`
	Limit    int   `mapstructure:"LIMIT"`
}

var ENV EnvType

func LoadEnv(path string) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic("Failed on reading .env file")
	}

	if err := viper.Unmarshal(&ENV); err != nil {
		panic("Failed on reading .env file")
	}
}
