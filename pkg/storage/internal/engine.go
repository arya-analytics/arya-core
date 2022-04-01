package internal

import (
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	"github.com/arya-analytics/aryacore/pkg/util/tasks"
)

// |||| ENGINE ||||

// Engine is a set of general interfaces that each engine variant must meet.
//
// Assigning Telem Responsibility
//
// Each engine variant is responsible for storing specific data types.
// These responsibilities are assigned in the model struct using the storage.re key.
// If no responsibility is assigned, EngineMD is assumed responsible.
type Engine interface {
	query.AssembleExec
}

// || META DATA ||

// EngineMD or the Metadata Engine is responsible for storing lightweight,
// strongly consistent data across the cluster.
type EngineMD interface {
	Engine
	query.Assemble
	NewTasks(opts ...tasks.ScheduleOpt) (tasks.Schedule, error)
}

// || OBJECT ||

// EngineObject is responsible for storing bulktelem data to node localstorage data storage.
type EngineObject interface {
	Engine
	query.AssembleCreate
	query.AssembleRetrieve
	query.AssembleDelete
	query.AssembleMigrate
}

// || CACHE ||

// EngineCache is responsible for storing and serving lightweight,
// ephemeral data at high speeds.
type EngineCache interface {
	Engine
	streamq.AssembleTSCreate
	streamq.AssembleTSRetrieve
}
