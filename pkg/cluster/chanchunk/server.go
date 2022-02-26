package chanchunk

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	api "github.com/arya-analytics/aryacore/pkg/cluster/gen/proto/go/chanchunk/v1"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/rpc"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"io"
)

type ServerRPC struct {
	cluster cluster.Cluster
}

func NewServerRPC(clust cluster.Cluster) *ServerRPC {
	return &ServerRPC{cluster: clust}
}

func (s *ServerRPC) BindTo(srv *grpc.Server) {
	api.RegisterChannelChunkServiceServer(srv, s)
}

func (s *ServerRPC) CreateReplicas(stream api.ChannelChunkService_CreateReplicasServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		ccr := &models.ChannelChunkReplica{}
		rpc.NewModelExchange(ccr, req.CCR).ToSource()
		if cErr := s.cluster.NewCreate().Model(ccr).Exec(stream.Context()); cErr != nil {
			return cErr
		}
	}
	return stream.SendAndClose(&api.ChannelChunkServiceCreateReplicasResponse{})
}

func (s *ServerRPC) RetrieveReplicas(req *api.ChannelChunkServiceRetrieveReplicasRequest, stream api.ChannelChunkService_RetrieveReplicasServer) error {
	PKC := parsePKC(req.PKC)
	for _, PK := range PKC {
		ccr := &models.ChannelChunkReplica{}
		if err := s.cluster.NewRetrieve().Model(ccr).WherePK(PK).Exec(stream.Context()); err != nil {
			return err
		}
		res := &api.ChannelChunkServiceRetrieveReplicasResponse{CCR: &api.ChannelChunkReplica{}}
		rpc.NewModelExchange(res.CCR, ccr).ToSource()
		if err := stream.Send(res); err != nil {
			return err
		}
	}
	return nil
}

func (s *ServerRPC) DeleteReplicas(ctx context.Context, req *api.ChannelChunkServiceDeleteReplicasRequest) (*api.ChannelChunkServiceDeleteReplicasResponse, error) {
	err := s.cluster.NewDelete().WherePKs(parsePKC(req.PKC).Raw()).Exec(ctx)
	return &api.ChannelChunkServiceDeleteReplicasResponse{}, err
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
