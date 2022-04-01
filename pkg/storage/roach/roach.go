package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"github.com/arya-analytics/aryacore/pkg/util/pool"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/tasks"
)

// Engine opens connections and execute queries with a roach database.
// implements the storage.EngineMD interface.
type Engine struct {
	query.AssembleBase
	driver Driver
	pool   *pool.Pool[internal.Engine]
}

// New creates a new roach.Engine with the provided Driver and storage.Pool.
func New(driver Driver, pool *pool.Pool[internal.Engine]) *Engine {
	e := &Engine{driver: driver, pool: pool}
	e.AssembleBase = query.NewAssemble(e.Exec)
	return e
}

// Exec implements query.Execute.
func (e *Engine) Exec(ctx context.Context, p *query.Pack) error {
	if !e.shouldHandle(p) {
		return nil
	}
	a, err := e.pool.Acquire(e)
	if err != nil {
		return newErrorConvert().Exec(err)
	}
	db := UnsafeDB(a)
	err = query.Switch(ctx, p, query.Ops{
		&query.Create{}:   newCreate(db).exec,
		&query.Retrieve{}: newRetrieve(db).exec,
		&query.Delete{}:   newDelete(db).exec,
		&query.Update{}:   newUpdate(db).exec,
		&query.Migrate{}:  newMigrate(db).exec,
	}, query.SwitchWithoutPanic())
	e.pool.Release(a)
	return newErrorConvert().Exec(err)
}

// NewAdapt opens a new connection with the data store and returns a storage.Adapter.
func (e *Engine) NewAdapt() (pool.Adapt[internal.Engine], error) {
	return newAdapter(e.driver)
}

func (e *Engine) Match(ce internal.Engine) bool {
	_, ok := ce.(*Engine)
	return ok
}

func (e *Engine) NewTasks(opts ...tasks.ScheduleOpt) (tasks.Schedule, error) {
	a, err := e.pool.Acquire(e)
	return newTaskScheduler(UnsafeDB(a), opts...), err
}

func (e *Engine) shouldHandle(p *query.Pack) bool {
	switch p.Query().(type) {
	case *query.Migrate:
		return true
	default:
		return catalog().Contains(p.Model())
	}
}
