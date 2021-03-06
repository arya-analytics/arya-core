package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/tasks"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"time"
)

const (
	tasksName = "roach_tasks"
)

func newTaskScheduler(db *bun.DB, opts ...tasks.ScheduleOpt) tasks.Schedule {
	opts = append(opts, tasks.ScheduleWithName(tasksName))
	return tasks.NewScheduleSimple(
		[]tasks.Task{
			{
				Action:   syncNodesAction(db),
				Interval: syncNodesInterval,
			},
		},
		opts...,
	)
}

// |||| NODE SYNCING |||

const (
	syncNodesInterval  = 5 * time.Second
	gossipNodeIDColumn = "node_id"
	nodesTable         = "nodes"
	nodesIDColumn      = "id"
)

// syncNodesAction scans the cockroach internal node table,
// and updates the arya nodes table to add/remove nodes that have
// joined/exited the cluster.
func syncNodesAction(db *bun.DB) tasks.Action {
	return func(ctx context.Context, cfg tasks.ScheduleConfig) error {
		sn := &syncNodes{db: db, catcher: errutil.NewCatchSimple(), cfg: cfg}
		return sn.exec(ctx)
	}
}

type syncNodes struct {
	ctx     context.Context
	db      *bun.DB
	catcher *errutil.CatchSimple
	cfg     tasks.ScheduleConfig
}

func (sn *syncNodes) exec(ctx context.Context) error {
	sn.ctx = ctx
	gnPKC, nodePKC := sn.retrieveGossipNodePKChain(), sn.retrieveNodePKChain()
	sn.runNodeAction(gnPKC, nodePKC, sn.createNodeWithPK)
	sn.runNodeAction(nodePKC, gnPKC, sn.deleteNodeWithPK)
	return newErrorConvert().Exec(sn.catcher.Error())
}

func (sn *syncNodes) runNodeAction(sourcePKC model.PKChain, destPKC model.PKChain,
	action func(pk model.PK)) {
	for _, sPK := range sourcePKC {
		found := false
		for _, dPK := range destPKC {
			if sPK.Equals(dPK) {
				found = true
			}
		}
		if !found {
			action(sPK)
		}
	}
}

func (sn *syncNodes) createNodeWithPK(pk model.PK) {
	fld := log.Fields{
		"pk": pk.Raw(),
	}
	if !sn.cfg.Silent {
		log.WithFields(fld).Info("A new node joined the cluster. Creating table entry.")
	}
	newNode := &models.Node{ID: pk.Raw().(int)}
	sn.catcher.Exec(func() error {
		if err := query.NewCreate().BindExec(newCreate(sn.db).exec).Model(newNode).Exec(sn.ctx); err != nil {
			if sErr, ok := err.(query.Error); !ok || sErr.Type != query.ErrorTypeUniqueViolation {
				return err
			}
		}
		return nil
	})
}

func (sn *syncNodes) deleteNodeWithPK(pk model.PK) {
	if !sn.cfg.Silent {
		log.WithFields(log.Fields{"pk": pk.Raw()}).Info("A node left the cluster. Removing table entry.")
	}
	sn.catcher.Exec(func() error {
		_, err := sn.db.NewDelete().
			Table(nodesTable).
			Where("ID = ?", pk.Raw()).
			Exec(sn.ctx)
		return err
	})
}

func (sn *syncNodes) retrieveGossipNodePKChain() model.PKChain {
	var gnIDs []int
	sn.catcher.Exec(func() error {
		return sn.db.NewSelect().
			Table(crdbGossipNodes).
			Column(gossipNodeIDColumn).
			Scan(sn.ctx, &gnIDs)
	})
	return model.NewPKChain(gnIDs)
}

func (sn *syncNodes) retrieveNodePKChain() model.PKChain {
	var nodeIDs []int
	sn.catcher.Exec(func() error {
		return sn.db.
			NewSelect().
			Table(nodesTable).
			Column(nodesIDColumn).
			Scan(sn.ctx, &nodeIDs)
	})
	return model.NewPKChain(nodeIDs)
}
