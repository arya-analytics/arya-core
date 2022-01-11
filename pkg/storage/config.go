package storage

import (
	"fmt"
)

type Config interface {
	Args() interface{}
	DSN() string
	Type() EngineType
	Role() EngineRole
}

type ConfigChain []Config

func (cc ConfigChain) Retrieve(et EngineRole) (Config, error) {
	for _, cfg := range cc {
		if cfg.Role() == et {
			return cfg, nil
		}
	}
	return nil, fmt.Errorf("config with type %v not found in config chain", et)
}

type MDStubConfig struct {}

func (a MDStubConfig) Args() interface{} {
	return []string{""}
}

func (a MDStubConfig) DSN() string {
	return ""
}

func (a MDStubConfig) Role() EngineRole {
	return EngineRoleMetaData
}

func (a MDStubConfig) Type() EngineType {
	return EngineTypeMDStub
}