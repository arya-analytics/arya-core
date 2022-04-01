package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/pool"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/tasks"
)

// Engine opens connections and execute queries with a roach database.
// implements the storage.EngineMD interface.
type Engine struct {
	query.AssembleBase
	Driver Driver
	Pool   *pool.Pool[*Engine]
}

// New creates a new roach.Engine with the provided Driver and storage.Pool.
func New(driver Driver) *Engine {
	e := &Engine{Driver: driver, Pool: pool.New[*Engine]()}
	e.Pool.Factory = e
	e.AssembleBase = query.NewAssemble(e.Exec)
	return e
}

// Exec implements query.Execute.
func (e *Engine) Exec(ctx context.Context, p *query.Pack) error {
	if !e.shouldHandle(p) {
		return nil
	}
	a, err := e.Pool.Acquire(e)
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
	e.Pool.Release(a)
	return newErrorConvert().Exec(err)
}

// NewAdapt opens a new connection with the data store and returns a storage.Adapter.
func (e *Engine) NewAdapt(*Engine) (pool.Adapt[*Engine], error) {
	return newAdapter(e.Driver)
}

func (e *Engine) NewTasks(opts ...tasks.ScheduleOpt) (tasks.Schedule, error) {
	a, err := e.Pool.Acquire(e)
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
