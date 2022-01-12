package storage

import (
	"context"
	"github.com/google/uuid"
)

type Adapter interface {
	ID() uuid.UUID
}

type EngineRole int

const (
	EngineRoleMetaData = iota
	EngineRoleCache
	EngineRoleBulk
)

type Engine interface {
	NewAdapter() Adapter
	IsAdapter(Adapter) bool
}

type MetaDataEngine interface {
	NewRetrieve(a Adapter) MetaDataRetrieve
}

type MetaDataRetrieve interface {
	Model(model interface{}) MetaDataRetrieve
	Where(query string, args ...interface{}) MetaDataRetrieve
	Exec(ctx context.Context) error
}