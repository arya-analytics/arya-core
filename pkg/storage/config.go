package storage

import (
	"fmt"
)

type Config interface {
	Args() interface{}
	Type() EngineType
	Role() EngineRole
}

type ConfigChain []Config

type ConfigError struct {
	Op string
	Et EngineType
}

func (e ConfigError) Error() string {
	return fmt.Sprintf("%s %v", e.Op, e.Et)
}

func (cc ConfigChain) Retrieve(et EngineType) (Config, error) {
	for _, cfg := range cc {
		if cfg.Type() == et {
			return cfg, nil
		}
	}
	return nil, ConfigError{Op: "config not found in config chain", Et: et}
}

type MDStubConfig struct{}

func (a MDStubConfig) Args() interface{} {
	return []string{""}
}

func (a MDStubConfig) Role() EngineRole {
	return EngineRoleMetaData
}

func (a MDStubConfig) Type() EngineType {
	return EngineTypeMDStub
}

type CacheStubConfig struct {}

func (a CacheStubConfig) Args() interface{} {
	return []string{""}
}

func (a CacheStubConfig) Role() EngineRole {
	return EngineRoleCache
}

func (a CacheStubConfig) Type() EngineType {
	return EngineTypeCacheStub
}