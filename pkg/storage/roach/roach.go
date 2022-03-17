package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
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
	pool   *internal.Pool
}

func New(driver Driver, pool *internal.Pool) *Engine {
	e := &Engine{driver: driver, pool: pool}
	e.AssembleBase = query.NewAssemble(e.Exec)
	return e
}

func (e *Engine) Exec(ctx context.Context, p *query.Pack) error {
	a, err := e.pool.Retrieve(e)
	if err != nil {
		return err
	}
	db := conn(a)
	return query.Switch(ctx, p, query.Ops{
		Create:   newCreate(db).exec,
		Retrieve: newRetrieve(db).exec,
		Delete:   newDelete(db).exec,
		Update:   newUpdate(db).exec,
		Migrate:  newMigrate(db).exec,
	})
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
	a, err := e.pool.Retrieve(e)
	return newTaskScheduler(conn(a), opts...), err
}
