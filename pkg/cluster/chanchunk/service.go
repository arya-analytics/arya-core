package chanchunk

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/batch"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"reflect"
)

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

func catalog() model.Catalog {
	return model.Catalog{
		&storage.ChannelChunk{},
		&storage.ChannelChunkReplica{},
	}
}

func (s *Service) CanHandle(q *cluster.QueryRequest) bool {
	return catalog().Contains(q.Model.Pointer())
}

func (s *Service) Exec(ctx context.Context, qr *cluster.QueryRequest) error {
	switch qr.Variant {
	case cluster.QueryVariantCreate:
		s.switchModel(ctx, qr, s.createChunk, s.createReplica)
	case cluster.QueryVariantRetrieve:
		s.switchModel(ctx, qr, s.retrieveChunk, s.retrieveReplica)
	case cluster.QueryVariantDelete:
		s.switchModel(ctx, qr, s.deleteChunk, s.deleteReplica)
	}
	return s.error()
}

type serviceOperation func(ctx context.Context, qr *cluster.QueryRequest)

func (s *Service) switchModel(ctx context.Context, qr *cluster.QueryRequest, soChunk, soReplica serviceOperation) {
	switch qr.Model.Type() {
	case reflect.TypeOf(storage.ChannelChunk{}):
		soChunk(ctx, qr)
	case reflect.TypeOf(storage.ChannelChunkReplica{}):
		soReplica(ctx, qr)
	}
}

func (s *Service) catchExec(actionFunc errutil.ActionFunc) {
	s.catcher.Exec(actionFunc)
}

func (s *Service) error() error {
	return s.catcher.Error()
}

// |||| CHUNK ||||

func (s *Service) createChunk(ctx context.Context, qr *cluster.QueryRequest) {
	s.catchExec(func() error { return s.local.CreateChunk(ctx, qr.Model.Pointer()) })
}

func (s *Service) retrieveChunk(ctx context.Context, qr *cluster.QueryRequest) {
	PKC, ok := cluster.PKQueryOpt(qr)
	if !ok {
		panic("retrieve queries require a primary key!")
	}
	s.catchExec(func() error { return s.local.RetrieveChunk(ctx, qr.Model.Pointer(), LocalChunkRetrieveOpts{PKC: PKC}) })
}

func (s *Service) deleteChunk(ctx context.Context, qr *cluster.QueryRequest) {
	PKC, ok := cluster.PKQueryOpt(qr)
	if !ok {
		panic("delete queries require a primary key!")
	}
	s.catcher.Exec(func() error { return s.local.DeleteChunk(ctx, LocalChunkDeleteOpts{PKC: PKC}) })
}

// |||| REPLICA ||||

func (s *Service) createReplica(ctx context.Context, qr *cluster.QueryRequest) {
	rrPKs := model.NewPKChain(qr.Model.FieldsByName("RangeReplicaID").Raw())
	opts := LocalRangeReplicaRetrieveOpts{PKC: rrPKs}
	rrS := qr.Model.FieldsByName("RangeReplica").ToReflect()
	s.catchExec(func() error { return s.local.RetrieveRangeReplica(ctx, rrS.Pointer(), opts) })
	s.replicaExec(
		qr.Model,
		func(m *model.Reflect) error { return s.local.CreateReplica(ctx, m.Pointer()) },
		func(m *model.Reflect) error { return s.remote.CreateReplica(ctx, buildRemoteReplicaCreateOpts(m)) },
	)
}

func (s *Service) retrieveReplica(ctx context.Context, qr *cluster.QueryRequest) {
	PKC, ok := cluster.PKQueryOpt(qr)
	if !ok {
		panic("retrieve queries require a primary key!")
	}
	opts := LocalReplicaRetrieveOpts{PKC: PKC, OmitBulk: true, Relations: true}
	s.catchExec(func() error { return s.local.RetrieveReplica(ctx, qr.Model.Pointer(), opts) })
	s.replicaExec(
		qr.Model,
		func(m *model.Reflect) error {
			return s.local.RetrieveReplica(ctx, m.Pointer(), LocalReplicaRetrieveOpts{PKC: m.PKChain()})
		},
		func(m *model.Reflect) error {
			return s.remote.RetrieveReplica(ctx, m.Pointer(), buildRemoteReplicaRetrieveOpts(m))
		},
	)
}

func (s *Service) deleteReplica(ctx context.Context, qr *cluster.QueryRequest) {
	PKC, ok := cluster.PKQueryOpt(qr)
	if !ok {
		panic("delete queries require a primary key!")
	}
	s.catcher.Exec(func() error {
		return s.local.RetrieveReplica(ctx, qr.Model.Pointer(), LocalReplicaRetrieveOpts{PKC: PKC})
	})
	s.replicaExec(
		qr.Model,
		func(m *model.Reflect) error { return s.local.DeleteReplica(ctx, LocalReplicaDeleteOpts{m.PKChain()}) },
		func(m *model.Reflect) error { return s.remote.DeleteReplica(ctx, buildRemoteReplicaDeleteOpts(m)) },
	)
}

func (s *Service) replicaExec(
	mRfl *model.Reflect,
	localF func(m *model.Reflect) error,
	remoteF func(m *model.Reflect) error) {
	isLocal := batch.NewModel(mRfl).Exec("RangeReplica.Node.IsHost")
	if localCCR, ok := isLocal[true]; ok {
		s.catchExec(func() error { return localF(localCCR) })
	}
	if remoteCCr, ok := isLocal[false]; ok {
		s.catchExec(func() error { return remoteF(remoteCCr) })
	}

}

func buildRemoteReplicaRetrieveOpts(remoteCCR *model.Reflect) (qp []RemoteReplicaRetrieveOpts) {
	addrMap := batch.NewModel(remoteCCR).Exec("RangeReplica.Node.Address")
	for addr, m := range addrMap {
		qp = append(qp, RemoteReplicaRetrieveOpts{
			Addr: addr.(string),
			PKC:  m.PKChain(),
		})
	}
	return qp
}

func buildRemoteReplicaCreateOpts(remoteCCR *model.Reflect) (qp []RemoterReplicaCreateOpts) {
	addrMap := batch.NewModel(remoteCCR).Exec("RangeReplica.Node.Address")
	for addr, m := range addrMap {
		qp = append(qp, RemoterReplicaCreateOpts{
			Addr:         addr.(string),
			ChunkReplica: m.Pointer(),
		})
	}
	return qp
}

func buildRemoteReplicaDeleteOpts(remoteCCR *model.Reflect) (qp []RemoteReplicaDeleteOpts) {
	addrMap := batch.NewModel(remoteCCR).Exec("RangeReplica.Node.Address")
	for addr, m := range addrMap {
		qp = append(qp, RemoteReplicaDeleteOpts{
			Addr: addr.(string),
			PKC:  m.PKChain(),
		})
	}
	return qp
}
