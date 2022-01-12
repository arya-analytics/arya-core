package storage

import "github.com/google/uuid"

type Adapter interface {
	ID() uuid.UUID
}

type EngineRole int

const (
	EngineRoleMetaData = iota
	EngineRoleCache
	EngineRoleBulk
)

type EngineBase interface {
	NewAdapter() Adapter
	IsAdapter(interface{}) bool
	Role() EngineRole
}
