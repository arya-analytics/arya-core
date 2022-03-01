package chanchunk

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/route"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	remote ServiceRemote
	local  Local
}

func NewService(local Local, remote ServiceRemote) *Service {
	return &Service{remote: remote, local: local}
}

func (s *Service) CanHandle(p *query.Pack) bool {
	return catalog().Contains(p.Model().Pointer())
}

func (s *Service) Exec(ctx context.Context, p *query.Pack) error {
	return query.Switch(ctx, p, query.Ops{
		Create:   s.createReplica,
		Retrieve: s.retrieveReplica,
		Delete:   s.deleteReplica,
		Update:   s.updateReplica,
	})
}

// |||| REPLICA ||||

const (
	RangeReplicaIDField = "RangeReplicaID"
	RangeReplicaField   = "RangeReplica"
	NodeIsHostField     = "RangeReplica.Node.IsHost"
	NodeField           = "RangeReplica.Node"
)

func (s *Service) createReplica(ctx context.Context, p *query.Pack) error {
	// CLARIFICATION: Retrieves information about the rng replicas and nodes model belongs to.
	// It will bind the results to p.Model itself.
	if err := s.local.RetrieveRangeReplica(
		ctx,
		p.Model().FieldsByName(RangeReplicaField).ToReflect().Pointer(),
		p.Model().FieldsByName(RangeReplicaIDField).ToPKChain(),
	); err != nil {
		return err
	}
	// CLARIFICATION: Now that we have the RangeReplicas.Node.IsHost field populated, we can switch on it.
	return replicaNodeIsHostSwitch(
		p.Model(),
		func(m *model.Reflect) error { return s.local.Create(ctx, m.Pointer()) },
		func(m *model.Reflect) error { return s.remote.Create(ctx, remCreateOpts(m)) },
	)
}

const BulkTelemField = "Telem"

// retrieveRequiredFields returns the minimum set of fields we need to complete a channel chunk replica retrieve
// request. We need this info to resolve the node that the replica belongs to.
func retrieveRequiredFields() []string {
	return []string{"ID", "ChannelChunkID", "RangeReplicaID"}
}

func (s *Service) retrieveReplica(ctx context.Context, p *query.Pack) error {
	baseOpts := LocalRetrieveOpts{NodeRelations: true, Fields: retrieveRequiredFields()}
	PKC, pkOK := query.PKOpt(p)
	if pkOK {
		baseOpts.PKC = PKC
	}

	whereFldsOpt, whereFldsOK := query.WhereFieldsOpt(p)
	if whereFldsOK {
		baseOpts.WhereFields = whereFldsOpt
	}

	fldsOpt, fldsOptOK := query.RetrieveFieldsOpt(p)
	if fldsOptOK {
		baseOpts.Fields = fldsOpt.AllExcept(BulkTelemField).Append(retrieveRequiredFields()...)
	}

	// CLARIFICATION: Retrieves information about the rng replicas and nodes model belongs to.
	// It will bind the results to p .Model itself.
	if err := s.local.Retrieve(ctx, p.Model().Pointer(), baseOpts); err != nil {
		return err
	}

	// CLARIFICATION: If we specified a fields query opt, and it doesn't contain the telem field, we don't
	// need to fetch bulk, so we can just return here.
	if fldsOptOK && !fldsOpt.ContainsAny(BulkTelemField) {
		return nil
	}

	// CLARIFICATION: Now that we have the RangeReplicas.Node.IsHost field populated, we can switch on it.
	return replicaNodeIsHostSwitch(
		p.Model(),
		func(m *model.Reflect) error {
			return s.local.Retrieve(ctx, m.Pointer(), LocalRetrieveOpts{PKC: m.PKChain()})
		},
		func(m *model.Reflect) error { return s.remote.Retrieve(ctx, m.Pointer(), remRetrieveOpts(m)) },
	)
}

func (s *Service) deleteReplica(ctx context.Context, p *query.Pack) error {
	PKC, ok := query.PKOpt(p)
	if !ok {
		panic("delete queries require a primary key!")
	}
	// CLARIFICATION: Retrieves information about the rng replicas and nodes model belongs to.
	// It will bind the results to p .Model itself.
	if err := s.local.Retrieve(ctx, p.Model().Pointer(), LocalRetrieveOpts{PKC: PKC, NodeRelations: true}); err != nil {
		return err
	}
	// CLARIFICATION: Now that we have the RangeReplicas.Node.IsHost field populated, we can switch on it.
	return replicaNodeIsHostSwitch(
		p.Model(),
		func(m *model.Reflect) error { return s.local.Delete(ctx, LocalDeleteOpts{m.PKChain()}) },
		func(m *model.Reflect) error { return s.remote.Delete(ctx, remDeleteOpts(m)) },
	)
}

func (s *Service) updateReplica(ctx context.Context, p *query.Pack) error {
	opts := LocalUpdateOpts{Bulk: query.BulkUpdateOpt(p)}
	PKC, pkOk := query.PKOpt(p)
	if pkOk {
		if len(PKC) > 1 {
			panic("update queries can't have more than one primary key")
		}
		opts.PK = PKC[0]
	}
	fieldsOpt, ok := query.RetrieveFieldsOpt(p)
	if ok {
		opts.Fields = fieldsOpt
	}
	if !p.Model().FieldsByName("Telem").AllNonZero() {
		log.
			WithFields(log.Fields{"ID": p.Model().PKChain().Raw()}).
			Warn("can't perform update on replica's telemetry, but was still specified!")
	}
	return s.local.Update(ctx, p.Model().Pointer(), opts)
}

// |||| ROUTING ||||

func replicaNodeIsHostSwitch(mRfl *model.Reflect, localF, remoteF func(m *model.Reflect) error) error {
	c := errutil.CatchSimple{}
	route.ModelSwitchBoolean(mRfl,
		NodeIsHostField,
		func(_ bool, m *model.Reflect) { c.Exec(func() error { return localF(m) }) },
		func(_ bool, m *model.Reflect) { c.Exec(func() error { return remoteF(m) }) },
	)
	return c.Error()
}

func replicaNodeSwitch(rfl *model.Reflect, action func(node *models.Node, rfl *model.Reflect)) {
	route.ModelSwitchIter(rfl, NodeField, action)
}

/// |||| REMOTE OPTION BUILDING ||||

func remRetrieveOpts(ccr *model.Reflect) (opts []RemoteRetrieveOpts) {
	replicaNodeSwitch(ccr, func(node *models.Node, m *model.Reflect) {
		opts = append(opts, RemoteRetrieveOpts{Node: node, PKC: m.PKChain()})
	})
	return opts

}

func remCreateOpts(ccr *model.Reflect) (opts []RemoteCreateOpts) {
	replicaNodeSwitch(ccr, func(node *models.Node, m *model.Reflect) {
		opts = append(opts, RemoteCreateOpts{Node: node, ChunkReplica: m.Pointer()})
	})
	return opts
}

func remDeleteOpts(ccr *model.Reflect) (opts []RemoteDeleteOpts) {
	replicaNodeSwitch(ccr, func(node *models.Node, m *model.Reflect) {
		opts = append(opts, RemoteDeleteOpts{Node: node, PKC: m.PKChain()})
	})
	return opts
}

// |||| CATALOG ||||

func catalog() model.Catalog {
	return model.Catalog{&models.ChannelChunkReplica{}}
}
