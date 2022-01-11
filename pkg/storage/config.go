package storage

import (
	"fmt"
)

type Config interface {
	Args() []interface{}
	DSN() string
	Type() EngineType
	Role() EngineRole
}

type ConfigChain struct {
	chain []Config
}

func (cc ConfigChain) Retrieve(et EngineRole) (Config, error) {
	for _, cfg := range cc.chain {
		if cfg.Role() == et {
			return cfg, nil
		}
	}
	return nil, fmt.Errorf("config with type %v not found in config chain", et)
}
