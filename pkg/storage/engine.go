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
	NewMigrate(a Adapter) MigrateQuery
}

// || META DATA ||

type MDEngine interface {
	BaseEngine
	NewRetrieve(a Adapter) MDRetrieveQuery
	NewCreate(a Adapter) MDCreateQuery
	NewDelete(a Adapter) MDDeleteQuery
}

// |||| QUERY ||||

type BaseQuery interface {
	Exec(ctx context.Context) error
}

// || META DATA ||

type MDRetrieveQuery interface {
	BaseQuery
	Model(model interface{}) MDRetrieveQuery
	Where(query string, args ...interface{}) MDRetrieveQuery
	WhereID(id interface{}) MDRetrieveQuery
}

type MDCreateQuery interface {
	BaseQuery
	Model(model interface{}) MDCreateQuery
}

type MDDeleteQuery interface {
	BaseQuery
	WhereID(id interface{}) MDDeleteQuery
	Model(model interface{}) MDDeleteQuery
}

// ||| MIGRATE |||

type MigrateQuery interface {
	Verify(ctx context.Context) error
	Exec(ctx context.Context) error
}
