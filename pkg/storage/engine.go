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
	Migrate(ctx context.Context, adapter Adapter) error
}

type MetaDataEngine interface {
	Engine
	NewRetrieve(a Adapter) MetaDataRetrieve
	NewCreate(a Adapter) MetaDataCreate
}

type MetaDataQuery interface {}

type MetaDataRetrieve interface {
	MetaDataQuery
	Model(model interface{}) MetaDataRetrieve
	Where(query string, args ...interface{}) MetaDataRetrieve
	WhereID(id interface{}) MetaDataRetrieve
	Exec(ctx context.Context) error
}

type MetaDataCreate interface {
	MetaDataQuery
	Model(model interface{}) MetaDataCreate
	Exec(ctx context.Context) error
}
