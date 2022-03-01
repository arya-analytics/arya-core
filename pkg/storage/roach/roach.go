package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/tasks"
	"github.com/uptrace/bun"
)

// |||| CONFIG ||||

type Driver interface {
	Connect() (*bun.DB, error)
}

// |||| ENGINE ||||

// Engine opens connections and execute queries with a roach database.
// implements the storage.EngineMD interface.
type Engine struct {
	query.AssembleBase
	driver Driver
	pool   *storage.Pool
}

func New(driver Driver, pool *storage.Pool) *Engine {
	e := &Engine{driver: driver, pool: pool}
	e.AssembleBase = query.NewAssemble(e.Exec)
	return e
}

func (e *Engine) Exec(ctx context.Context, p *query.Pack) error {
	db := conn(e.pool.Retrieve(e))
	return query.Switch(ctx, p, query.Ops{
		Create:   newCreate(db).Exec,
		Retrieve: newRetrieve(db).Exec,
		Delete:   newDelete(db).Exec,
		Update:   newUpdate(db).Exec,
	})
}

// NewAdapter opens a new connection with the data store and returns a storage.Adapter.
func (e *Engine) NewAdapter() storage.Adapter {
	return newAdapter(e.driver)
}

// IsAdapter checks if the provided adapter is a roach adapter.
func (e *Engine) IsAdapter(a storage.Adapter) bool {
	_, ok := bindAdapter(a)
	return ok
}

func (e *Engine) ShouldHandle(m interface{}, _ ...string) bool {
	return catalog().Contains(m)
}

// NewMigrate opens a new migrateExec with the provided storage.Adapter.
func (e *Engine) NewMigrate() storage.QueryMDMigrate {
	return newMigrate(conn(e.pool.Retrieve(e)), e.driver)
}

func (e *Engine) NewTasks(opts ...tasks.ScheduleOpt) tasks.Schedule {
	return newTaskScheduler(conn(e.pool.Retrieve(e)), opts...)
}
