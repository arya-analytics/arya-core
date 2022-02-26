package mock

import (
	"context"
	api "github.com/arya-analytics/aryacore/pkg/cluster/gen/proto/go/chanchunk/v1"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"io"
)

type Server struct {
	api.UnimplementedChannelChunkServiceServer
	CreatedChunks       *model.Reflect
	DeletedChunkPKChain model.PKChain
}

func NewServer() *Server {
	return &Server{
		CreatedChunks: model.NewReflect(&[]*api.ChannelChunkReplica{}),
	}
}

func (s *Server) BindTo(srv *grpc.Server) {
	api.RegisterChannelChunkServiceServer(srv, s)
}

func (s *Server) CreateReplicas(stream api.ChannelChunkService_CreateReplicasServer) error {
	for {
		cc, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		s.CreatedChunks.ChainAppend(model.NewReflect(cc.CCR))
	}
	if err := stream.SendAndClose(&api.ChannelChunkServiceCreateReplicasResponse{}); err != nil {
		return err
	}
	return nil
}

func (s *Server) RetrieveReplicas(req *api.ChannelChunkServiceRetrieveReplicasRequest, stream api.ChannelChunkService_RetrieveReplicasServer) error {
	for _, id := range req.PKC {
		if err := stream.Send(&api.ChannelChunkServiceRetrieveReplicasResponse{CCR: &api.ChannelChunkReplica{
			ID:    id,
			Telem: []byte{1, 2, 3, 4},
		}}); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) DeleteReplicas(ctx context.Context, req *api.ChannelChunkServiceDeleteReplicasRequest) (*api.ChannelChunkServiceDeleteReplicasResponse, error) {
	pkC := model.NewPKChain([]uuid.UUID{})
	for _, id := range req.PKC {
		pk, err := model.NewPK(uuid.New()).NewFromString(id)
		if err != nil {
			panic(err)
		}
		pkC = append(pkC, pk)
	}
	s.DeletedChunkPKChain = append(s.DeletedChunkPKChain, pkC...)
	return &api.ChannelChunkServiceDeleteReplicasResponse{}, nil
}
