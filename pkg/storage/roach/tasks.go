package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	tasks2 "github.com/arya-analytics/aryacore/pkg/util/tasks"
	"github.com/uptrace/bun"
	"time"
)

const taskTickInterval = 3 * time.Second

func newTaskScheduler(db *bun.DB) *tasks2.Scheduler {
	return tasks2.NewScheduler(
		[]tasks2.Task{
			{
				Action:   syncNodesAction(db),
				Interval: syncNodesInterval,
			},
		},
		taskTickInterval,
		tasks2.ScheduleWithName("roach tasks"),
	)
}

// |||| NODE SYNCING |||

const syncNodesInterval = 3 * time.Second

func syncNodesAction(db *bun.DB) tasks2.Action {
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
