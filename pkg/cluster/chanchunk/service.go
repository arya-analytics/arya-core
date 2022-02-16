package chanchunk

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/batch"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
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
		storage.ChannelChunk{},
		storage.ChannelChunkReplica{},
	}
}

func (s *Service) CanHandle(q *cluster.QueryRequest) bool {
	return catalog().Contains(q.Model)
}

func (s *Service) Exec(ctx context.Context, q *cluster.QueryRequest) error {
	switch q.Variant {
	case cluster.QueryVariantCreate:
		s.createReplicas(ctx, q)
	}
	return s.catcher.Error()
}

func (s *Service) createReplicas(ctx context.Context, q *cluster.QueryRequest) {
	rrPKs := model.NewPKChain(q.Model.FieldsByName("RangeReplicaID").Raw())
	rrS := q.Model.FieldsByName("RangeReplica").ToReflect()
	s.catcher.Exec(func() error {
		return s.local.RetrieveRangeReplicas(ctx, rrS, rrPKs)
	})
	isLocal := batch.NewModel(q.Model).Exec("RangeReplica.Node.IsHost")
	localCC, ok := isLocal[true]
	if ok {
		s.catcher.Exec(func() error { return s.local.CreateReplicas(ctx, localCC) })
	}
	remoteCC, ok := isLocal[false]
	if ok {
		s.catcher.Exec(func() error {
			return s.remote.CreateReplicas(ctx, buildRemoteCreateParams(remoteCC))
		})
	}
}

//func (s *Service) retrieveReplicas(ctx context.Context, q *cluster.QueryRequest) {
//	if !ok {
//		panic("replica retrieve queries require a PK")
//	}
//	ccrPKC := model.NewPKChain(pkOpt())
//	s.catcher.Exec(func() error {
//		s.local.RetrieveReplicas(ctx, q.Model, ccrPKC, false)
//	})
//}

func buildRemoteCreateParams(remoteCC *model.Reflect) (qp []RemoteReplicaCreateParams) {
	addrMap := batch.NewModel(remoteCC).Exec("RangeReplica.Node.Address")
	for addr, m := range addrMap {
		qp = append(qp, RemoteReplicaCreateParams{
			Addr:  addr.(string),
			Model: m,
		})
	}
	return qp

}
