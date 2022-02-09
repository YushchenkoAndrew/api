package config

import (
	"api/interfaces"

	"github.com/spf13/viper"
)

type EnvType struct {
	Host     string `mapstructure:"HOST"`
	Port     string `mapstructure:"PORT"`
	BasePath string `mapstructure:"BASE_PATH"`

	// DataBase
	DBType string `mapstructure:"DB_TYPE"`
	DBName string `mapstructure:"DB_NAME"`
	DBHost string `mapstructure:"DB_HOST"`
	DBPort string `mapstructure:"DB_PORT"`
	DBUser string `mapstructure:"DB_USER"`
	DBPass string `mapstructure:"DB_PASS"`

	// Redis
	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort string `mapstructure:"REDIS_PORT"`
	RedisPass string `mapstructure:"REDIS_PASS"`

	// JWT
	AccessSecret  string `mapstructure:"ACCESS_SECRET"`
	RefreshSecret string `mapstructure:"REFRESH_SECRET"`

	// Root User Login + Pass & Pepper
	ID     string `mapstructure:"API_ID"`
	URL    string `mapstructure:"API_URL"`
	User   string `mapstructure:"API_USER"`
	Pass   string `mapstructure:"API_PASS"`
	Pepper string `mapstructure:"API_PEPPER"`

	// Pagination setting
	LiveTime int64 `mapstructure:"LIVE_TIME"`
	Items    int   `mapstructure:"ITEMS"`
	Limit    int   `mapstructure:"LIMIT"`

	// Rate Info
	RateLimit int `mapstructure:"RATE_LIMIT"`
	RateTime  int `mapstructure:"RATE_TIME"`

	BotUrl    string `mapstructure:"BOT_URL"`
	BotKey    string `mapstructure:"BOT_KEY"`
	BotPepper string `mapstructure:"BOT_Pepper"`

	// Migration Settings
	ForceMigrate  bool   `mapstructure:"FORCE_MIGRATE"`
	MigrationPath string `mapstructure:"MIGRATION_PATH"`

	// Metrics
	Metrics int `mapstructure:"METRICS_COUNT"`
}

// FIXME: I should fix this one day
var ENV EnvType

type envConfig struct {
	path string
}

func NewEnvConfig(path string) func() interfaces.Config {
	return func() interfaces.Config {
		return &envConfig{path: path}
	}
}

func (c *envConfig) Init() {
	viper.AddConfigPath(c.path)
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic("Failed on reading .env file")
	}

	if err := viper.Unmarshal(&ENV); err != nil {
		panic("Failed on reading .env file")
	}
}
