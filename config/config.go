package config

import "github.com/spf13/viper"

type Env struct {
	Host     string `mapstructure:"HOST"`
	Port     string `mapstructure:"PORT"`
	BasePath string `mapstructure:"BASE_PATH"`

	DBType string `mapstructure:"DB_TYPE"`
	DBName string `mapstructure:"DB_NAME"`
	DBHost string `mapstructure:"DB_HOST"`
	DBPort string `mapstructure:"DB_PORT"`
	DBUser string `mapstructure:"DB_USER"`
	DBPass string `mapstructure:"DB_PASS"`

	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort string `mapstructure:"REDIS_PORT"`
	RedisPass string `mapstructure:"REDIS_PASS"`
}

type Config struct{}

func (o *Config) Init(path string) (env Env) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	o.CheckOnErr(&err, "Failed on reading .env file")

	err = viper.Unmarshal(&env)
	o.CheckOnErr(&err, "Failed on reading .env file")
	return
}

func (*Config) CheckOnErr(err *error, msg string) {
	if *err != nil {
		panic(msg)
	}
}
