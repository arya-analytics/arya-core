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

// |||| ENGINE ||||

type Engine interface {
	NewAdapter() Adapter
	IsAdapter(a Adapter) bool
	NewMigrate(a Adapter) Migrate
}

// || MIGRATE ||

type Migrate interface {
	Verify(ctx context.Context) error
	Exec(ctx context.Context) error
}

// || META DATA ||

type MetaDataEngine interface {
	Engine
	NewRetrieve(a Adapter) MetaDataRetrieve
	NewCreate(a Adapter) MetaDataCreate
	NewDelete(a Adapter) MetaDataDelete
}

type MetaDataRetrieve interface {
	Model(model interface{}) MetaDataRetrieve
	Where(query string, args ...interface{}) MetaDataRetrieve
	WhereID(id interface{}) MetaDataRetrieve
	Exec(ctx context.Context) error
}

type MetaDataCreate interface {
	Model(model interface{}) MetaDataCreate
	Exec(ctx context.Context) error
}

type MetaDataDelete interface {
	WhereID(id interface{}) MetaDataDelete
	Model(model interface{}) MetaDataDelete
	Exec(ctx context.Context) error
}
