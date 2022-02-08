package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
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

const syncNodesInterval = 5 * time.Second

func syncNodesAction(db *bun.DB) tasks.Action {
	return func(ctx context.Context) error {
		gnc, nc, err := nodeCounts(db, ctx)
		if err != nil {
			return err
		}
		if gnc != nc {
			var (
				gnIds []int
				nodes []*storage.Node
			)
			if gnErr := db.NewSelect().Table(crdbGossipNodes).Column("node_id").Scan(ctx,
				&gnIds); gnErr != nil {
				return gnErr
			}

			if nErr := newRetrieve(db).Model(&nodes).Exec(ctx); nErr != nil {
				return nErr
			}

			nodesRfl := model.NewReflect(&nodes)

			if gnc > nc {
				for _, gnId := range gnIds {
					pk := model.NewPK(gnId)
					if _, ok := nodesRfl.ValueByPK(pk); !ok {
						nodeRfl := nodesRfl.NewStruct()
						nodeRfl.PKField().Set(pk.Value())
						if cErr := newCreate(db).Model(nodeRfl.Pointer()).Exec(
							ctx); cErr != nil {
							return cErr
						}
					}
				}
			}
		}
		return nil
	}
}

func nodeCounts(db *bun.DB, ctx context.Context) (int, int, error) {
	gnc, err := db.NewSelect().Column("node_id").Count(ctx)
	if err != nil {
		return 0, 0, err
	}
	nc, err := newRetrieve(db).Model(&storage.Node{}).Count(ctx)
	if err != nil {
		return 0, 0, err
	}
	return gnc, nc, nil
}
