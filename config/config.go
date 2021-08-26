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

	LiveTime int64 `mapstructure:"LIVE_TIME"`
	Items    int   `mapstructure:"ITEMS"`
	Limit    int   `mapstructure:"LIMIT"`
}

var ENV EnvType

func LoadEnv(path string) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	CheckOnErr(&err, "Failed on reading .env file")

	err = viper.Unmarshal(&ENV)
	CheckOnErr(&err, "Failed on reading .env file")
}

func CheckOnErr(err *error, msg string) {
	if *err != nil {
		panic(msg)
	}
}
