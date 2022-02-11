package config

import (
	"api/interfaces"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const (
	TYPE = "yaml"
)

type Handler struct {
	Name     string   `mapstructure:"name"`
	Method   string   `mapstructure:"method"`
	Path     string   `mapstructure:"path"`
	Required []string `mapstructure:"required"`
}

type operation struct {
	Cfg []Handler `mapstructure:"cfg"`
}

var operations = make(map[string]Handler)

func GetOperation(key string) (*Handler, bool) {
	value, ok := operations[key]
	return &value, ok
}

type operationConfig struct {
	path, name string
	operations operation
}

func NewOperationConfig(path, name string) func() interfaces.Config {
	return func() interfaces.Config {
		return &operationConfig{path: path, name: name}
	}
}

func (c *operationConfig) Init() {
	viper.AddConfigPath(c.path)
	viper.SetConfigName(c.name)
	viper.SetConfigType(TYPE)

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic("Failed on reading operations file")
	}

	if err := viper.Unmarshal(&c.operations); err != nil {
		panic("Failed on reading operation file")
	}

	if path, err := os.Getwd(); err == nil {
		fmt.Println(path) // for example /home/user
	}

	// Form map
	for _, cfg := range c.operations.Cfg {
		operations[cfg.Name] = cfg
	}
}
