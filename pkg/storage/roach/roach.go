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
	DemandCap() int
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
	a := e.pool.Retrieve(e)
	db := conn(a)
	return query.Switch(ctx, p, query.Ops{
		Create:   newCreate(db).exec,
		Retrieve: newRetrieve(db).exec,
		Delete:   newDelete(db).exec,
		Update:   newUpdate(db).exec,
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
	a := e.pool.Retrieve(e)
	return newMigrate(conn(a), e.driver)
}

func (e *Engine) NewTasks(opts ...tasks.ScheduleOpt) tasks.Schedule {
	return newTaskScheduler(conn(e.pool.Retrieve(e)), opts...)
}
