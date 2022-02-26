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
	DeleteReplicas(ctx context.Context, pkc model.PKChain) error
}

type ServerRPC struct {
	persist ServerRPCPersist
}

func NewServerRPC(p ServerRPCPersist) *ServerRPC {
	return &ServerRPC{persist: p}
}

func (s *ServerRPC) BindTo(srv *grpc.Server) {
	api.RegisterChannelChunkServiceServer(srv, s)
}

func (s *ServerRPC) CreateReplicas(stream api.ChannelChunkService_CreateReplicasServer) error {
	c := &errutil.Catcher{}
	for {
		var req *api.ChannelChunkServiceCreateReplicasRequest
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
	return stream.SendAndClose(&api.ChannelChunkServiceCreateReplicasResponse{})
}

func (s *ServerRPC) RetrieveReplicas(req *api.ChannelChunkServiceRetrieveReplicasRequest, stream api.ChannelChunkService_RetrieveReplicasServer) error {
	pkc := parsePKC(req.PKC)
	c := &errutil.Catcher{}
	for _, pk := range pkc {
		res := &api.ChannelChunkServiceRetrieveReplicasResponse{CCR: &api.ChannelChunkReplica{}}
		c.Exec(func() error { return s.persist.RetrieveReplica(stream.Context(), res.CCR, pk) })
		c.Exec(func() error { return stream.Send(res) })
	}
	return c.Error()
}

func (s *ServerRPC) DeleteReplicas(ctx context.Context, req *api.ChannelChunkServiceDeleteReplicasRequest) (*api.ChannelChunkServiceDeleteReplicasResponse, error) {
	return &api.ChannelChunkServiceDeleteReplicasResponse{}, s.persist.DeleteReplicas(ctx, parsePKC(req.PKC))
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
	cluster cluster.Cluster
}

func (sp *ServerRPCPersistCluster) RetrieveReplica(ctx context.Context, ccr *api.ChannelChunkReplica, pk model.PK) error {
	exc := rpc.NewModelExchange(&models.ChannelChunkReplica{}, ccr)
	if err := sp.cluster.NewRetrieve().Model(exc.Source).WherePK(pk.Raw()).Exec(ctx); err != nil {
		return err
	}
	exc.ToDest()
	return nil
}
func (sp *ServerRPCPersistCluster) CreateReplica(ctx context.Context, ccr *api.ChannelChunkReplica) error {
	exc := rpc.NewModelExchange(ccr, &models.ChannelChunkReplica{})
	exc.ToDest()
	return sp.cluster.NewCreate().Model(exc.Dest).Exec(ctx)
}

func (sp *ServerRPCPersistCluster) DeleteReplica(ctx context.Context, pkc model.PKChain) error {
	return sp.cluster.NewDelete().Model(&models.ChannelChunkReplica{}).WherePKs(pkc.Raw()).Exec(ctx)
}
