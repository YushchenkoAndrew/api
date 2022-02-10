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
	Method string
	Path   string
}

var operations map[string]Handler

func GetOperation(key string) (Handler, bool) {
	value, ok := operations[key]
	return value, ok
}

type operation struct {
	Cfg []struct {
		Name   string `mapstructure:"name"`
		Method string `mapstructure:"method"`
		Path   string `mapstructure:"path"`
	} `mapstructure:"cfg"`
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
	operations = make(map[string]Handler)
	for _, cfg := range c.operations.Cfg {
		operations[cfg.Name] = Handler{Method: cfg.Method, Path: cfg.Path}
	}
}
