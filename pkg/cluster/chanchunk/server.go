package chanchunk

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	api "github.com/arya-analytics/aryacore/pkg/cluster/gen/proto/go/chanchunk/v1"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/rpc"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"io"
)

type ServerRPCPersist interface {
	CreateReplica(ctx context.Context, ccr *api.ChannelChunkReplica) error
	RetrieveReplica(ctx context.Context, ccr *api.ChannelChunkReplica, pk model.PK) error
	DeleteReplica(ctx context.Context, pkc model.PKChain) error
}

type ServerRPC struct {
	api.UnimplementedChannelChunkServiceServer
	persist ServerRPCPersist
}

func NewServerRPC(p ServerRPCPersist) *ServerRPC {
	return &ServerRPC{persist: p}
}

func (s *ServerRPC) BindTo(srv *grpc.Server) {
	api.RegisterChannelChunkServiceServer(srv, s)
}

func (s *ServerRPC) CreateReplicas(stream api.ChannelChunkService_CreateReplicasServer) error {
	c := errutil.NewCatchSimple()
	for {
		var req *api.CreateReplicasRequest
		c.Exec(func() (err error) {
			req, err = stream.Recv()
			return err
		})
		c.Exec(func() error { return s.persist.CreateReplica(stream.Context(), req.CCR) })
		if c.Error() != nil {
			if c.Error() == io.EOF {
				break
			}
			return c.Error()
		}
	}
	return stream.SendAndClose(&api.CreateReplicasResponse{})
}

func (s *ServerRPC) RetrieveReplicas(req *api.RetrieveReplicasRequest, stream api.ChannelChunkService_RetrieveReplicasServer) error {
	pkc := parsePKC(req.PKC)
	c := errutil.NewCatchSimple()
	for _, pk := range pkc {
		res := &api.RetrieveReplicasResponse{CCR: &api.ChannelChunkReplica{}}
		c.Exec(func() error { return s.persist.RetrieveReplica(stream.Context(), res.CCR, pk) })
		c.Exec(func() error { return stream.Send(res) })
	}
	return c.Error()
}

func (s *ServerRPC) DeleteReplicas(ctx context.Context, req *api.DeleteReplicasRequest) (*api.DeleteReplicasResponse, error) {
	return &api.DeleteReplicasResponse{}, s.persist.DeleteReplica(ctx, parsePKC(req.PKC))
}

func parsePKC(strPKC []string) model.PKChain {
	PKC := model.NewPKChain([]uuid.UUID{})
	for _, strPK := range strPKC {
		pk, err := model.NewPK(uuid.New()).NewFromString(strPK)
		if err != nil {
			panic(err)
		}
		PKC = append(PKC, pk)
	}
	return PKC
}

type ServerRPCPersistCluster struct {
	Cluster cluster.Cluster
}

func (sp *ServerRPCPersistCluster) RetrieveReplica(ctx context.Context, ccr *api.ChannelChunkReplica, pk model.PK) error {
	exc := rpc.NewModelExchange(&models.ChannelChunkReplica{}, ccr)
	if err := sp.Cluster.NewRetrieve().Model(exc.Source()).WherePK(pk.Raw()).Exec(ctx); err != nil {
		return err
	}
	exc.ToDest()
	return nil
}
func (sp *ServerRPCPersistCluster) CreateReplica(ctx context.Context, ccr *api.ChannelChunkReplica) error {
	mCCR := &models.ChannelChunkReplica{}
	rpc.NewModelExchange(ccr, mCCR).ToDest()
	return sp.Cluster.NewCreate().Model(mCCR).Exec(ctx)
}

func (sp *ServerRPCPersistCluster) DeleteReplica(ctx context.Context, pkc model.PKChain) error {
	return sp.Cluster.NewDelete().Model(&models.ChannelChunkReplica{}).WherePKs(pkc.Raw()).Exec(ctx)
}
