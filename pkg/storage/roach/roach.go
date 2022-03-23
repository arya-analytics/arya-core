package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/tasks"
)

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
	a, err := e.pool.Acquire(e)
	if err != nil {
		return newErrorConvert().Exec(err)
	}
	db := conn(a)
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

// NewAdapter opens a new connection with the data store and returns a storage.Adapter.
func (e *Engine) NewAdapter() (internal.Adapter, error) {
	return newAdapter(e.driver)
}

// IsAdapter checks if the provided adapter is a roach adapter.
func (e *Engine) IsAdapter(a internal.Adapter) bool {
	_, ok := bindAdapter(a)
	return ok
}

func (e *Engine) ShouldHandle(m interface{}, _ ...string) bool {
	return catalog().Contains(m)
}

func (e *Engine) NewTasks(opts ...tasks.ScheduleOpt) (tasks.Schedule, error) {
	a, err := e.pool.Acquire(e)
	return newTaskScheduler(conn(a), opts...), err
}
