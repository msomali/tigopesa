package config

import (
	"github.com/techcraftt/tigosdk/push"
	"github.com/techcraftt/tigosdk/ussd/aw"
	"github.com/techcraftt/tigosdk/ussd/wa"
)

var (
	_ Validator = (*Config)(nil)
)

type (
	
	Config struct {
		
	}

	//Validator validates the configurations and return nil if all
	//is good or return an error
	Validator interface {
		Validate() error
	}
)

func (c *Config) Validate() error {
	panic("implement me")
}

func (c *Config) Split() (push.Config,wa.Config,aw.Config,error) {
	return push.Config{}, wa.Config{}, aw.Config{}, nil
}

func Merge(config push.Config, config2 wa.Config, config3 aw.Config) (Config,error) {
	return Config{}, nil
}
