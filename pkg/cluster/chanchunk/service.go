package chanchunk

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster/internal"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
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

func (s *Service) CanHandle(q *internal.QueryRequest) bool {
	return catalog().Contains(q.Model.Pointer())
}

func (s *Service) Exec(ctx context.Context, qr *internal.QueryRequest) error {
	return internal.SwitchQueryRequestVariant(ctx, qr, internal.QueryRequestVariantOperations{
		internal.QueryVariantCreate:   s.createReplica,
		internal.QueryVariantRetrieve: s.retrieveReplica,
		internal.QueryVariantDelete:   s.deleteReplica,
		internal.QueryVariantUpdate:   s.updateReplica,
	})
}

// |||| REPLICA ||||

const (
	RangeReplicaIDField = "RangeReplicaID"
	RangeReplicaField   = "RangeReplica"
	NodeIsHostField     = "RangeReplica.Node.IsHost"
	NodeField           = "RangeReplica.Node"
)

func (s *Service) createReplica(ctx context.Context, qr *internal.QueryRequest) error {
	// CLARIFICATION: Retrieves information about the rng replicas and nodes model belongs to.
	// It will bind the results to qr.Model itself.
	if err := s.local.RetrieveRangeReplica(
		ctx,
		qr.Model.FieldsByName(RangeReplicaField).ToReflect().Pointer(),
		qr.Model.FieldsByName(RangeReplicaIDField).ToPKChain(),
	); err != nil {
		return err
	}
	// CLARIFICATION: Now that we have the RangeReplicas.Node.IsHost field populated, we can switch on it.
	return replicaNodeIsHostSwitch(
		qr.Model,
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

func (s *Service) retrieveReplica(ctx context.Context, qr *internal.QueryRequest) error {
	baseOpts := LocalRetrieveOpts{NodeRelations: true, Fields: retrieveRequiredFields()}
	PKC, pkOK := internal.PKQueryOpt(qr)
	if pkOK {
		baseOpts.PKC = PKC
	}

	whereFldsOpt, whereFldsOK := internal.WhereFieldsQueryOpt(qr)
	if whereFldsOK {
		baseOpts.WhereFields = whereFldsOpt
	}

	fldsOpt, fldsOptOK := internal.RetrieveFieldsQueryOpt(qr)
	if fldsOptOK {
		baseOpts.Fields = fldsOpt.AllExcept(BulkTelemField).Append(retrieveRequiredFields()...)
	}

	// CLARIFICATION: Retrieves information about the rng replicas and nodes model belongs to.
	// It will bind the results to qr.Model itself.
	if err := s.local.Retrieve(ctx, qr.Model.Pointer(), baseOpts); err != nil {
		return err
	}

	// CLARIFICATION: If we specified a fields query opt, and it doesn't contain the telem field, we don't
	// need to fetch bulk, so we can just return here.
	if fldsOptOK && !fldsOpt.ContainsAny(BulkTelemField) {
		return nil
	}

	// CLARIFICATION: Now that we have the RangeReplicas.Node.IsHost field populated, we can switch on it.
	return replicaNodeIsHostSwitch(
		qr.Model,
		func(m *model.Reflect) error {
			return s.local.Retrieve(ctx, m.Pointer(), LocalRetrieveOpts{PKC: m.PKChain()})
		},
		func(m *model.Reflect) error { return s.remote.Retrieve(ctx, m.Pointer(), remRetrieveOpts(m)) },
	)
}

func (s *Service) deleteReplica(ctx context.Context, qr *internal.QueryRequest) error {
	PKC, ok := internal.PKQueryOpt(qr)
	if !ok {
		panic("delete queries require a primary key!")
	}
	// CLARIFICATION: Retrieves information about the rng replicas and nodes model belongs to.
	// It will bind the results to qr.Model itself.
	if err := s.local.Retrieve(ctx, qr.Model.Pointer(), LocalRetrieveOpts{PKC: PKC, NodeRelations: true}); err != nil {
		return err
	}
	// CLARIFICATION: Now that we have the RangeReplicas.Node.IsHost field populated, we can switch on it.
	return replicaNodeIsHostSwitch(
		qr.Model,
		func(m *model.Reflect) error { return s.local.Delete(ctx, LocalDeleteOpts{m.PKChain()}) },
		func(m *model.Reflect) error { return s.remote.Delete(ctx, remDeleteOpts(m)) },
	)
}

func (s *Service) updateReplica(ctx context.Context, qr *internal.QueryRequest) error {
	opts := LocalUpdateOpts{Bulk: internal.BulkUpdateQueryOpt(qr)}
	PKC, pkOk := internal.PKQueryOpt(qr)
	if pkOk {
		if len(PKC) > 1 {
			panic("update queries can't have more than one primary key")
		}
		opts.PK = PKC[0]
	}
	fieldsOpt, ok := internal.RetrieveFieldsQueryOpt(qr)
	if ok {
		opts.Fields = fieldsOpt
	}
	if !qr.Model.FieldsByName("Telem").AllNonZero() {
		log.
			WithFields(log.Fields{"ID": qr.Model.PKChain().Raw()}).
			Warn("can't perform update on replica's telemetry, but was still specified!")
	}
	return s.local.Update(ctx, qr.Model.Pointer(), opts)
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
