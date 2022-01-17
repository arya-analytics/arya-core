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
	EngineRoleMD = iota
	EngineRoleCache
	EngineRoleBulk
)

// |||| ENGINE ||||

type BaseEngine interface {
	NewAdapter() Adapter
	IsAdapter(a Adapter) bool
}

// || META DATA ||

type MDEngine interface {
	BaseEngine
	NewRetrieve(a Adapter) MDRetrieveQuery
	NewCreate(a Adapter) MDCreateQuery
	NewDelete(a Adapter) MDDeleteQuery
	NewMigrate(a Adapter) MDMigrateQuery
}

// |||| QUERY ||||

type BaseQuery interface {
	Exec(ctx context.Context) error
}

// || META DATA ||

type MDBaseQuery interface {
	BaseQuery
}

type MDRetrieveQuery interface {
	MDBaseQuery
	Model(model interface{}) MDRetrieveQuery
	Where(query string, args ...interface{}) MDRetrieveQuery
	WherePK(pk interface{}) MDRetrieveQuery
	WherePKs(pks interface{}) MDRetrieveQuery
}

type MDCreateQuery interface {
	MDBaseQuery
	UpdateOnConflict() MDCreateQuery
	Model(model interface{}) MDCreateQuery
}

type MDDeleteQuery interface {
	MDBaseQuery
	WherePK(pk interface{}) MDDeleteQuery
	WherePKs(pks interface{}) MDDeleteQuery
	Model(model interface{}) MDDeleteQuery
}

type MDMigrateQuery interface {
	MDBaseQuery
	Verify(ctx context.Context) error
	Exec(ctx context.Context) error
}
