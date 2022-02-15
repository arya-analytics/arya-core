package mock

import (
	"context"
	api "github.com/arya-analytics/aryacore/pkg/cluster/gen/proto/go/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"google.golang.org/grpc"
	"io"
)

type Server struct {
	api.UnimplementedChannelChunkServiceServer
	CreatedChunks *model.Reflect
	DeletedChunks model.PKChain
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
		s.CreatedChunks.ChainAppend(model.NewReflect(cc.Chunk))
	}
	if err := stream.SendAndClose(&api.ChannelChunkServiceCreateReplicasResponse{}); err != nil {
		return err
	}
	return nil
}

func (s *Server) RetrieveReplicas(req *api.ChannelChunkServiceRetrieveReplicasRequest, stream api.ChannelChunkService_RetrieveReplicasServer) error {
	for _, pk := range req.Id {
		if err := stream.Send(&api.ChannelChunkServiceRetrieveReplicasResponse{Chunk: &api.ChannelChunkReplica{
			Id:    pk,
			Telem: []byte{1, 2, 3, 4},
		}}); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) DeleteReplicas(ctx context.Context, req *api.ChannelChunkServiceDeleteReplicasRequest) (*api.ChannelChunkServiceDeleteReplicasResponse, error) {
	s.DeletedChunks = append(s.DeletedChunks, model.NewPKChain(req.Id)...)
	return &api.ChannelChunkServiceDeleteReplicasResponse{}, nil
}
