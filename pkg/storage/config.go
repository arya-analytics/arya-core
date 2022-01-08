package storage

import (
	"fmt"
)

type Config interface {
	Args() []interface{}
	DSN() string
	Type() EngineType
}

type ConfigChain struct {
	chain []Config
}

func (cc ConfigChain) Retrieve(et EngineType) (Config, error) {
	for _, cfg := range cc.chain {
		if cfg.Type() == et {
			return cfg, nil
		}
	}
	return nil, fmt.Errorf("config with type %v not found in config chain", et)
}
