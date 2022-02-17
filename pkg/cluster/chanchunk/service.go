package chanchunk

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster/internal"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/route"
	"reflect"
)

// ||||

type Service struct {
	remote  ServiceRemote
	local   ServiceLocal
	catcher *errutil.Catcher
}

func NewService(local ServiceLocal, remote ServiceRemote) *Service {
	return &Service{
		remote:  remote,
		local:   local,
		catcher: &errutil.Catcher{},
	}
}

func (s *Service) CanHandle(q *internal.QueryRequest) bool {
	return catalog().Contains(q.Model.Pointer())
}

func (s *Service) Exec(ctx context.Context, qr *internal.QueryRequest) error {
	s.catcher.Reset()
	switch qr.Model.Type() {
	case reflect.TypeOf(storage.ChannelChunk{}):
		s.switchVariant(ctx, qr, s.createChunk, s.retrieveChunk, s.deleteChunk)
	case reflect.TypeOf(storage.ChannelChunkReplica{}):
		s.switchVariant(ctx, qr, s.createReplica, s.retrieveReplica, s.deleteReplica)
	}
	return s.error()
}

func (s *Service) switchVariant(
	ctx context.Context,
	qr *internal.QueryRequest,
	createOp, retrieveOp, deleteOp internal.ServiceOperation) {
	switch qr.Variant {
	case internal.QueryVariantCreate:
		createOp(ctx, qr)
	case internal.QueryVariantRetrieve:
		retrieveOp(ctx, qr)
	case internal.QueryVariantDelete:
		deleteOp(ctx, qr)
	}
}

func (s *Service) catchExec(actionFunc errutil.ActionFunc) {
	s.catcher.Exec(actionFunc)
}

func (s *Service) error() error {
	return s.catcher.Error()
}

// |||| CHUNK ||||

func (s *Service) createChunk(ctx context.Context, qr *internal.QueryRequest) {
	s.catchExec(func() error { return s.local.CreateChunk(ctx, qr.Model.Pointer()) })
}

func (s *Service) retrieveChunk(ctx context.Context, qr *internal.QueryRequest) {
	PKC, ok := internal.PKQueryOpt(qr)
	if !ok {
		panic("retrieve queries require a primary key!")
	}
	s.catchExec(func() error { return s.local.RetrieveChunk(ctx, qr.Model.Pointer(), LocalChunkRetrieveOpts{PKC: PKC}) })
}

func (s *Service) deleteChunk(ctx context.Context, qr *internal.QueryRequest) {
	PKC, ok := internal.PKQueryOpt(qr)
	if !ok {
		panic("delete queries require a primary key!")
	}
	s.catcher.Exec(func() error { return s.local.DeleteChunk(ctx, LocalChunkDeleteOpts{PKC: PKC}) })
}

// |||| REPLICA ||||

const (
	RangeReplicaIDField = "RangeReplicaID"
	RangeReplicaField   = "RangeReplica"
	NodeIsHostField     = "RangeReplica.Node.IsHost"
	NodeAddressField    = "RangeReplica.Node.Address"
)

func (s *Service) createReplica(ctx context.Context, qr *internal.QueryRequest) {
	rrPKs := model.NewPKChain(qr.Model.FieldsByName(RangeReplicaIDField).Raw())
	opts := LocalRangeReplicaRetrieveOpts{PKC: rrPKs}
	rrS := qr.Model.FieldsByName(RangeReplicaField).ToReflect()
	s.catchExec(func() error { return s.local.RetrieveRangeReplica(ctx, rrS.Pointer(), opts) })
	replicaIsHostSwitch(
		qr.Model,
		func(_ bool, m *model.Reflect) {
			s.catchExec(func() error { return s.local.CreateReplica(ctx, m.Pointer()) })
		},
		func(_ bool, m *model.Reflect) {
			s.catchExec(func() error { return s.remote.CreateReplica(ctx, buildRemoteReplicaCreateOpts(m)) })
		},
	)
}

func (s *Service) retrieveReplica(ctx context.Context, qr *internal.QueryRequest) {
	PKC, ok := internal.PKQueryOpt(qr)
	if !ok {
		panic("retrieve queries require a primary key!")
	}
	opts := LocalReplicaRetrieveOpts{PKC: PKC, OmitBulk: true, Relations: true}
	s.catchExec(func() error { return s.local.RetrieveReplica(ctx, qr.Model.Pointer(), opts) })
	replicaIsHostSwitch(
		qr.Model,
		func(_ bool, m *model.Reflect) {
			s.catchExec(func() error {
				return s.local.RetrieveReplica(ctx, m.Pointer(), LocalReplicaRetrieveOpts{PKC: m.PKChain()})
			})
		},
		func(_ bool, m *model.Reflect) {
			s.catchExec(func() error {
				return s.remote.RetrieveReplica(ctx, m.Pointer(), buildRemoteReplicaRetrieveOpts(m))
			})
		},
	)
}

func (s *Service) deleteReplica(ctx context.Context, qr *internal.QueryRequest) {
	PKC, ok := internal.PKQueryOpt(qr)
	if !ok {
		panic("delete queries require a primary key!")
	}
	s.catchExec(func() error {
		return s.local.RetrieveReplica(ctx, qr.Model.Pointer(), LocalReplicaRetrieveOpts{PKC: PKC})
	})
	replicaIsHostSwitch(
		qr.Model,
		func(_ bool, m *model.Reflect) {
			s.catchExec(func() error { return s.local.DeleteReplica(ctx, LocalReplicaDeleteOpts{m.PKChain()}) })
		},
		func(_ bool, m *model.Reflect) {
			s.catchExec(func() error { return s.remote.DeleteReplica(ctx, buildRemoteReplicaDeleteOpts(m)) })
		},
	)
}

// |||| ROUTING ||||

func replicaIsHostSwitch(mRfl *model.Reflect, localF, remoteF func(_ bool, m *model.Reflect)) {
	route.ModelSwitchBoolean(mRfl, NodeIsHostField, localF, remoteF)

}

func replicaAddrSwitch(rfl *model.Reflect, action func(addr string, rfl *model.Reflect)) {
	route.ModelSwitchIter(rfl, NodeAddressField, action)
}

/// |||| REMOTE OPTION BUILDING ||||

func buildRemoteReplicaRetrieveOpts(remoteCCR *model.Reflect) (qp []RemoteReplicaRetrieveOpts) {
	replicaAddrSwitch(remoteCCR, func(addr string, m *model.Reflect) {
		qp = append(qp, RemoteReplicaRetrieveOpts{Addr: addr, PKC: m.PKChain()})
	})
	return qp

}

func buildRemoteReplicaCreateOpts(remoteCCR *model.Reflect) (qp []RemoterReplicaCreateOpts) {
	replicaAddrSwitch(remoteCCR, func(addr string, m *model.Reflect) {
		qp = append(qp, RemoterReplicaCreateOpts{Addr: addr, ChunkReplica: m.Pointer()})
	})
	return qp
}

func buildRemoteReplicaDeleteOpts(remoteCCR *model.Reflect) (qp []RemoteReplicaDeleteOpts) {
	replicaAddrSwitch(remoteCCR, func(addr string, m *model.Reflect) {
		qp = append(qp, RemoteReplicaDeleteOpts{Addr: addr, PKC: m.PKChain()})
	})
	return qp
}

// |||| CATALOG ||||

func catalog() model.Catalog {
	return model.Catalog{
		&storage.ChannelChunk{},
		&storage.ChannelChunkReplica{},
	}
}
