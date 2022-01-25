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
	EngineRoleObject
	EngineRoleCache
)

// |||| ENGINE ||||

type BaseEngine interface {
	NewAdapter() Adapter
	IsAdapter(a Adapter) bool
	InCatalog(m interface{}) bool
}

// || META DATA ||

type MDEngine interface {
	BaseEngine
	// NewRetrieve opens a new MDRetrieveQuery.
	NewRetrieve(a Adapter) MDRetrieveQuery
	// NewCreate opens a new MDCreateQuery.
	NewCreate(a Adapter) MDCreateQuery
	// NewDelete opens a new MDDeleteQuery.
	NewDelete(a Adapter) MDDeleteQuery
	// NewMigrate opens a new MDMigrateQuery.
	NewMigrate(a Adapter) MDMigrateQuery
	// NewUpdate opens a new MDUpdateQuery.
	NewUpdate(a Adapter) MDUpdateQuery
}

// || OBJECT ||

type ObjectEngine interface {
	BaseEngine
	NewRetrieve(a Adapter) ObjectRetrieveQuery
	NewCreate(a Adapter) ObjectCreateQuery
	NewDelete(a Adapter) ObjectDeleteQuery
	NewMigrate(a Adapter) ObjectMigrateQuery
}

// |||| QUERY ||||

type BaseQuery interface {
	Exec(ctx context.Context) error
}

// || META DATA ||

type MDBaseQuery interface {
	BaseQuery
}

// MDRetrieveQuery is for retrieving items from metadata storage.
type MDRetrieveQuery interface {
	MDBaseQuery
	Model(model interface{}) MDRetrieveQuery
	Where(query string, args ...interface{}) MDRetrieveQuery
	WherePK(pk interface{}) MDRetrieveQuery
	WherePKs(pks interface{}) MDRetrieveQuery
}

// MDCreateQuery is for creating items in metadata storage.
type MDCreateQuery interface {
	MDBaseQuery
	Model(model interface{}) MDCreateQuery
}

// MDUpdateQuery is for updating items in metadata storage.
type MDUpdateQuery interface {
	MDBaseQuery
	Model(model interface{}) MDUpdateQuery
	Where(query string, args ...interface{}) MDUpdateQuery
	WherePK(pk interface{}) MDUpdateQuery
}

// MDDeleteQuery is for deleting items in metadata storage.
type MDDeleteQuery interface {
	MDBaseQuery
	WherePK(pk interface{}) MDDeleteQuery
	WherePKs(pks interface{}) MDDeleteQuery
	Model(model interface{}) MDDeleteQuery
}

// MDMigrateQuery applies migration changes to metadata storage.
type MDMigrateQuery interface {
	MDBaseQuery
	Verify(ctx context.Context) error
}

// || BULK ||

type ObjectBaseQuery interface {
	BaseQuery
}

type ObjectCreateQuery interface {
	ObjectBaseQuery
	Model(model interface{}) ObjectCreateQuery
}

type ObjectRetrieveQuery interface {
	ObjectBaseQuery
	Model(model interface{}) ObjectRetrieveQuery
	WherePK(pk interface{}) ObjectRetrieveQuery
	WherePKs(pks interface{}) ObjectRetrieveQuery
}

type ObjectDeleteQuery interface {
	ObjectBaseQuery
	Model(model interface{}) ObjectDeleteQuery
	WherePK(pk interface{}) ObjectDeleteQuery
	WherePKs(pks interface{}) ObjectDeleteQuery
}

type ObjectMigrateQuery interface {
	ObjectBaseQuery
	Verify(ctx context.Context) error
}
