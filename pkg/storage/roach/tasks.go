package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/tasks"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/uptrace/bun"
	"time"
)

const taskTickInterval = 3 * time.Second

func newTaskScheduler(db *bun.DB) *tasks.Scheduler {
	return tasks.NewScheduler(
		[]tasks.Task{
			{
				Action:   syncNodesAction(db),
				Interval: syncNodesInterval,
			},
		},
		taskTickInterval,
		tasks.SchedulerWithName("roach tasks"),
	)
}

// |||| NODE SYNCING |||

const syncNodesInterval = 3 * time.Second

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
					_, ok := nodesRfl.ValueByPK(pk)
					if !ok {
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
	gnc, err := db.NewSelect().Table(crdbGossipNodes).ScanAndCount(ctx)
	if err != nil {
		return 0, 0, err
	}
	nc, err := newRetrieve(db).Model(&storage.Node{}).Count(ctx)
	if err != nil {
		return 0, 0, err
	}
	return gnc, nc, nil
}
