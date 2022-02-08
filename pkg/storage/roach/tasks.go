package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/tasks"
	"github.com/uptrace/bun"
	"time"
)

const (
	taskTickInterval = 5 * time.Second
	tasksName        = "roach_tasks"
)

func newTaskScheduler(db *bun.DB, opts ...tasks.SchedulerOpt) *tasks.Scheduler {
	opts = append(opts, tasks.ScheduleWithName(tasksName))
	return tasks.NewScheduler(
		[]tasks.Task{
			{
				Action:   syncNodesAction(db),
				Interval: syncNodesInterval,
			},
		},
		taskTickInterval,
		opts...,
	)
}

// |||| NODE SYNCING |||

const (
	syncNodesInterval  = 5 * time.Second
	gossipNodeIDColumn = "node_id"
)

func syncNodesAction(db *bun.DB) tasks.Action {
	sn := &syncNodes{db: db, catcher: &errutil.Catcher{}}
	return func(ctx context.Context) error { return sn.exec(ctx) }
}

type syncNodes struct {
	ctx     context.Context
	db      *bun.DB
	catcher *errutil.Catcher
}

func (sn *syncNodes) exec(ctx context.Context) error {
	sn.ctx = ctx
	gnc, nc := sn.countNodeImbalance()
	if gnc > nc {
		sn.createMissingNodes()
	}
	return sn.catcher.Error()
}

func (sn *syncNodes) countNodeImbalance() (gnc int, nc int) {
	sn.catcher.Exec(func() (err error) {
		if gnc, err = sn.db.NewSelect().
			Column(gossipNodeIDColumn).
			Count(sn.ctx); err != nil {
			return err
		}
		if nc, err = newRetrieve(sn.db).
			Model(&storage.Node{}).
			Count(sn.ctx); err != nil {
			return err
		}
		return nil
	})
	return gnc, nc
}

func (sn *syncNodes) createMissingNodes() {
	gnIDs, nodesRfl := sn.retrieveGossipNodeIDs(), sn.retrieveNodesRfl()
	for _, gnID := range gnIDs {
		pk := model.NewPK(gnID)
		if _, ok := nodesRfl.ValueByPK(pk); !ok {
			sn.createNodeWithPK(pk, nodesRfl)
		}
	}

}

func (sn *syncNodes) createNodeWithPK(pk model.PK, nodesRfl *model.Reflect) {
	nodeRfl := nodesRfl.NewStruct()
	nodeRfl.PKField().Set(pk.Value())
	sn.catcher.Exec(func() error {
		return newCreate(sn.db).Model(nodeRfl.Pointer()).Exec(sn.ctx)
	})
}

func (sn *syncNodes) retrieveGossipNodeIDs() (gnIDs []int) {
	sn.catcher.Exec(func() error {
		return sn.db.NewSelect().
			Table(crdbGossipNodes).
			Column(gossipNodeIDColumn).
			Scan(sn.ctx, &gnIDs)
	})
	return gnIDs
}

func (sn *syncNodes) retrieveNodesRfl() *model.Reflect {
	nodesRfl := model.NewReflect(&[]*storage.Node{})
	sn.catcher.Exec(func() error {
		return newRetrieve(sn.db).Model(nodesRfl.Pointer()).Exec(sn.ctx)
	})
	return nodesRfl
}
