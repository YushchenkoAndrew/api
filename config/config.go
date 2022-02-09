package config

import "api/interfaces"

type config struct {
	cfg []interfaces.Config
}

func NewConfig(handlers []func() interfaces.Config) interfaces.Config {
	var cfg []interfaces.Config
	for _, handler := range handlers {
		cfg = append(cfg, handler())
	}

	return &config{cfg: cfg}
}

func (c *config) Init() {
	for _, cfg := range c.cfg {
		cfg.Init()
	}
}
