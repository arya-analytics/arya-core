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

type ConfigError struct {
	Op string
	Et EngineType
}

func NewConfigError(op string, et EngineType) ConfigError {
	return ConfigError{Op: op, Et: et}
}

func (e ConfigError) Error() string {
	return fmt.Sprintf("%s %v", e.Op, e.Et)
}

func (cc ConfigChain) Retrieve(t EngineType) (Config, error) {
	for _, cfg := range cc {
		if cfg.Type() == t {
			return cfg, nil
		}
	}
	return nil, NewConfigError("Config not found in config chain", t)
}

type MDStubConfig struct{}

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

type CacheStubConfig struct {}

func (a CacheStubConfig) Args() interface{} {
	return []string{""}
}

func (a CacheStubConfig) DSN() string {
	return ""
}

func (a CacheStubConfig) Role() EngineRole {
	return EngineRoleCache
}

func (a CacheStubConfig) Type() EngineType {
	return EngineTypeCacheStub
}