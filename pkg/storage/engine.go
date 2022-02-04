package storage

import (
	"context"
	"github.com/google/uuid"
)

type Adapter interface {
	ID() uuid.UUID
}

// |||| ENGINE ||||

// Engine is a set of general interfaces that each engine variant must meet.
//
// Assigning Data Responsibility
//
// Each engine variant is responsible for storing specific data types.
// These responsibilities are assigned in the model struct using the storage.re key.
// If no responsibility is assigned, MDEngine is assumed responsible.
type Engine interface {
	NewAdapter() Adapter
	IsAdapter(a Adapter) bool
	InCatalog(model interface{}) bool
}

// || META DATA ||

// MDEngine or the Metadata Engine is responsible for storing lightweight,
// strongly consistent data across the cluster.
type MDEngine interface {
	Engine
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

// ObjectEngine is responsible for storing bulk data to node local data storage.
type ObjectEngine interface {
	Engine
	// NewRetrieve opens a new ObjectRetrieveQuery.
	NewRetrieve(a Adapter) ObjectRetrieveQuery
	// NewCreate opens a new ObjectCreateQuery.
	NewCreate(a Adapter) ObjectCreateQuery
	// NewDelete opens a new ObjectDeleteQuery.
	NewDelete(a Adapter) ObjectDeleteQuery
	// NewMigrate opens a new ObjectMigrateQuery.
	NewMigrate(a Adapter) ObjectMigrateQuery
}

// || CACHE ||

// CacheEngine is responsible for storing and serving lightweight,
// ephemeral data at high speeds.
type CacheEngine interface {
	Engine
	// NewTSRetrieve opens a new CacheTSRetrieveQuery.
	NewTSRetrieve(a Adapter) CacheTSRetrieveQuery
	// NewTSCreate opens a new CacheTSCreateQuery.
	NewTSCreate(a Adapter) CacheTSCreateQuery
}

// |||| QUERY ||||

type Query interface {
	Exec(ctx context.Context) error
}

// || META DATA ||

type MDBaseQuery interface {
	Query
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
	Model(model interface{}) MDDeleteQuery
	WherePK(pk interface{}) MDDeleteQuery
	WherePKs(pks interface{}) MDDeleteQuery
}

// MDMigrateQuery applies migration changes to metadata storage.
type MDMigrateQuery interface {
	MDBaseQuery
	Verify(ctx context.Context) error
}

// || OBJECT ||

type ObjectBaseQuery interface {
	Query
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

// || TS CACHE ||

type CacheBaseQuery interface {
	Query
}

type CacheCreateQuery interface {
	CacheBaseQuery
	Model(model interface{}) CacheCreateQuery
}

type CacheTSRetrieveQuery interface {
	CacheBaseQuery
	SeriesExists(ctx context.Context, pk interface{}) (bool, error)
	Model(model interface{}) CacheTSRetrieveQuery
	WherePK(pk interface{}) CacheTSRetrieveQuery
	WherePKs(pks interface{}) CacheTSRetrieveQuery
	AllTimeRange() CacheTSRetrieveQuery
	WhereTimeRange(fromTS int64, toTS int64) CacheTSRetrieveQuery
}

type CacheTSCreateQuery interface {
	CacheBaseQuery
	Model(model interface{}) CacheTSCreateQuery
	Series() CacheTSCreateQuery
	Sample() CacheTSCreateQuery
}
