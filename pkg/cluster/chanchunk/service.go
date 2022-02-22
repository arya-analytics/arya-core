package chanchunk

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster/internal"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/route"
	"reflect"
)

type Service struct {
	remote ServiceRemote
	local  ServiceLocal
}

func NewService(local ServiceLocal, remote ServiceRemote) *Service {
	return &Service{remote: remote, local: local}
}

func (s *Service) CanHandle(q *internal.QueryRequest) bool {
	return catalog().Contains(q.Model.Pointer())
}

func (s *Service) Exec(ctx context.Context, qr *internal.QueryRequest) error {
	switch qr.Model.Type() {
	case reflect.TypeOf(models.ChannelChunk{}):
		return internal.SwitchQueryRequestVariant(ctx, qr, internal.QueryRequestVariantOperations{
			internal.QueryVariantCreate:   s.createChunk,
			internal.QueryVariantRetrieve: s.retrieveChunk,
			internal.QueryVariantUpdate:   s.updateChunk,
			internal.QueryVariantDelete:   s.deleteChunk,
		})
	case reflect.TypeOf(models.ChannelChunkReplica{}):
		return internal.SwitchQueryRequestVariant(ctx, qr, internal.QueryRequestVariantOperations{
			internal.QueryVariantCreate:   s.createReplica,
			internal.QueryVariantRetrieve: s.retrieveReplica,
			internal.QueryVariantDelete:   s.deleteReplica,
		})
	default:
		panic("channel chunk service received an unknown model type!")
	}
}

// |||| CHUNK ||||

func (s *Service) createChunk(ctx context.Context, qr *internal.QueryRequest) error {
	return s.local.CreateChunk(ctx, qr.Model.Pointer())
}

func (s *Service) retrieveChunk(ctx context.Context, qr *internal.QueryRequest) error {
	PKC, ok := internal.PKQueryOpt(qr)
	if !ok {
		panic("retrieve queries require a primary key!")
	}
	return s.local.RetrieveChunk(ctx, qr.Model.Pointer(), LocalChunkRetrieveOpts{PKC: PKC})
}

func (s *Service) updateChunk(ctx context.Context, qr *internal.QueryRequest) error {
	PKC, ok := internal.PKQueryOpt(qr)
	if !ok {
		panic("update queries require a primary key!")
	}
	if len(PKC) > 1 {
		panic("update query cannot have more than one primary key!")
	}
	return s.local.UpdateChunk(ctx, qr.Model.Pointer(), LocalChunkUpdateOpts{PK: PKC[0]})
}

func (s *Service) deleteChunk(ctx context.Context, qr *internal.QueryRequest) error {
	PKC, ok := internal.PKQueryOpt(qr)
	if !ok {
		panic("delete queries require a primary key!")
	}
	return s.local.DeleteChunk(ctx, LocalChunkDeleteOpts{PKC: PKC})
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
		LocalRangeReplicaRetrieveOpts{PKC: qr.Model.FieldsByName(RangeReplicaIDField).ToPKChain()},
	); err != nil {
		return err
	}
	// CLARIFICATION: Now that we have the RangeReplica.Node.IsHost field populated, we can switch on it.
	return replicaNodeIsHostSwitch(
		qr.Model,
		func(m *model.Reflect) error { return s.local.CreateReplica(ctx, m.Pointer()) },
		func(m *model.Reflect) error { return s.remote.CreateReplica(ctx, buildRemoteReplicaCreateOpts(m)) },
	)
}

func (s *Service) retrieveReplica(ctx context.Context, qr *internal.QueryRequest) error {
	PKC, ok := internal.PKQueryOpt(qr)
	if !ok {
		panic("retrieve queries require a primary key!")
	}
	// CLARIFICATION: Retrieves information about the rng replicas and nodes model belongs to.
	// It will bind the results to qr.Model itself.
	if err := s.local.RetrieveReplica(
		ctx,
		qr.Model.Pointer(),
		LocalReplicaRetrieveOpts{PKC: PKC, OmitBulk: true, Relations: true}); err != nil {
		return err
	}
	// CLARIFICATION: Now that we have the RangeReplica.Node.IsHost field populated, we can switch on it.
	return replicaNodeIsHostSwitch(
		qr.Model,
		func(m *model.Reflect) error {
			return s.local.RetrieveReplica(ctx, m.Pointer(), LocalReplicaRetrieveOpts{PKC: m.PKChain()})
		},
		func(m *model.Reflect) error {
			return s.remote.RetrieveReplica(ctx, m.Pointer(), buildRemoteReplicaRetrieveOpts(m))
		},
	)
}

func (s *Service) deleteReplica(ctx context.Context, qr *internal.QueryRequest) error {
	PKC, ok := internal.PKQueryOpt(qr)
	if !ok {
		panic("delete queries require a primary key!")
	}
	// CLARIFICATION: Retrieves information about the rng replicas and nodes model belongs to.
	// It will bind the results to qr.Model itself.
	if err := s.local.RetrieveReplica(ctx, qr.Model.Pointer(), LocalReplicaRetrieveOpts{PKC: PKC}); err != nil {
		return err
	}
	// CLARIFICATION: Now that we have the RangeReplica.Node.IsHost field populated, we can switch on it.
	return replicaNodeIsHostSwitch(
		qr.Model,
		func(m *model.Reflect) error { return s.local.DeleteReplica(ctx, LocalReplicaDeleteOpts{m.PKChain()}) },
		func(m *model.Reflect) error { return s.remote.DeleteReplica(ctx, buildRemoteReplicaDeleteOpts(m)) },
	)
}

// |||| ROUTING ||||

func replicaNodeIsHostSwitch(mRfl *model.Reflect, localF, remoteF func(m *model.Reflect) error) (err error) {
	route.ModelSwitchBoolean(mRfl,
		NodeIsHostField,
		func(_ bool, m *model.Reflect) {
			if lErr := localF(m); lErr != nil {
				err = lErr
			}
		},
		func(_ bool, m *model.Reflect) {
			if rErr := remoteF(m); rErr != nil {
				err = rErr
			}
		})
	return err
}

func replicaNodeSwitch(rfl *model.Reflect, action func(node *models.Node, rfl *model.Reflect)) {
	route.ModelSwitchIter(rfl, NodeField, action)
}

/// |||| REMOTE OPTION BUILDING ||||

func buildRemoteReplicaRetrieveOpts(remoteCCR *model.Reflect) (opts []RemoteReplicaRetrieveOpts) {
	replicaNodeSwitch(remoteCCR, func(node *models.Node, m *model.Reflect) {
		opts = append(opts, RemoteReplicaRetrieveOpts{Node: node, PKC: m.PKChain()})
	})
	return opts

}

func buildRemoteReplicaCreateOpts(remoteCCR *model.Reflect) (opts []RemoterReplicaCreateOpts) {
	replicaNodeSwitch(remoteCCR, func(node *models.Node, m *model.Reflect) {
		opts = append(opts, RemoterReplicaCreateOpts{Node: node, ChunkReplica: m.Pointer()})
	})
	return opts
}

func buildRemoteReplicaDeleteOpts(remoteCCR *model.Reflect) (opts []RemoteReplicaDeleteOpts) {
	replicaNodeSwitch(remoteCCR, func(node *models.Node, m *model.Reflect) {
		opts = append(opts, RemoteReplicaDeleteOpts{Node: node, PKC: m.PKChain()})
	})
	return opts
}

// |||| CATALOG ||||

func catalog() model.Catalog {
	return model.Catalog{&models.ChannelChunkReplica{}}
}
