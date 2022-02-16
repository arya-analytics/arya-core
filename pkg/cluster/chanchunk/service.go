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
	remote  ClientRemote
	local   ClientLocal
	catcher *errutil.Catcher
}

func NewService(local ClientLocal, remote ClientRemote) *Service {
	return &Service{
		remote:  remote,
		local:   local,
		catcher: &errutil.Catcher{},
	}
}

type RemoteReplicaRetrieveParams struct {
	Addr string
	PKC  model.PKChain
}

type RemoteReplicaCreateParams struct {
	Addr  string
	Model *model.Reflect
}

type RemoteReplicaDeleteParams struct {
	Addr string
	PKC  model.PKChain
}

type ClientRemote interface {
	// |||| REPLICA ||||

	RetrieveReplicas(ctx context.Context, ccr *model.Reflect, qp []RemoteReplicaRetrieveParams) error
	CreateReplicas(ctx context.Context, qp []RemoteReplicaCreateParams) error
	DeleteReplicas(ctx context.Context, qp []RemoteReplicaDeleteParams) error
}

type ClientLocal interface {
	// |||| CHUNK ||||

	Create(ctx context.Context, cc *model.Reflect) error
	Retrieve(ctx context.Context, cc *model.Reflect, ccPKC model.PKChain) error
	Delete(ctx context.Context, ccPKC model.PKChain) error

	// |||| REPLICA ||||

	CreateReplicas(ctx context.Context, ccr *model.Reflect) error
	RetrieveReplicas(ctx context.Context, ccr *model.Reflect, ccrPKC model.PKChain, omitBulk bool) error
	DeleteReplicas(ctx context.Context, ccrPKC model.PKChain) error

	// |||| RANGE REPLICA ||||

	RetrieveRangeReplicas(ctx context.Context, rr *model.Reflect, rrPKC model.PKChain) error
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
