package chanchunk

import (
	"context"
	api "github.com/arya-analytics/aryacore/pkg/cluster/gen/proto/go/chanchunk/v1"
	"github.com/arya-analytics/aryacore/pkg/rpc"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"io"
)

/// |||| INTERFACE ||||

type ServiceRemote interface {
	// |||| REPLICA ||||

	RetrieveReplica(ctx context.Context, chunkReplica interface{}, qp []RemoteReplicaRetrieveOpts) error
	CreateReplica(ctx context.Context, qp []RemoterReplicaCreateOpts) error
	DeleteReplica(ctx context.Context, qp []RemoteReplicaDeleteOpts) error
}

type RemoteReplicaRetrieveOpts struct {
	Addr string
	PKC  model.PKChain
}

type RemoterReplicaCreateOpts struct {
	Addr         string
	ChunkReplica interface{}
}

type RemoteReplicaDeleteOpts struct {
	Addr string
	PKC  model.PKChain
}

// |||| RPC IMPLEMENTATION ||||

func catalogRemoteRPC() model.Catalog {
	return model.Catalog{
		&api.ChannelChunkReplica{},
	}

}
func newExchange(m interface{}) *model.Exchange {
	return rpc.NewModelExchange(m, catalogRemoteRPC().New(m))
}

type ServiceRemoteRPC struct {
	pool rpc.Pool
}

func NewServiceRemoteRPC(pool rpc.Pool) ServiceRemote {
	return &ServiceRemoteRPC{pool: pool}
}

func (s *ServiceRemoteRPC) client(addr string) api.ChannelChunkServiceClient {
	return api.NewChannelChunkServiceClient(s.pool.Retrieve(addr))
}

func (s *ServiceRemoteRPC) RetrieveReplica(ctx context.Context, chunkReplica interface{}, qp []RemoteReplicaRetrieveOpts) error {
	exc := newExchange(chunkReplica)
	for _, params := range qp {
		rq := &api.ChannelChunkServiceRetrieveReplicasRequest{Id: params.PKC.Strings()}
		stream, err := s.client(params.Addr).RetrieveReplicas(ctx, rq)
		if err != nil {
			return err
		}
		for {
			in, sErr := stream.Recv()
			if sErr == io.EOF {
				break
			}
			if sErr != nil {
				return sErr
			}
			inRfl := model.NewReflect(in.Chunk)
			exc.Dest.ChainAppend(inRfl)
		}
	}
	exc.ToSource()
	return nil
}

func (s *ServiceRemoteRPC) CreateReplica(ctx context.Context, qp []RemoterReplicaCreateOpts) error {
	for _, params := range qp {
		exc := newExchange(params.ChunkReplica)
		exc.ToDest()

		stream, err := s.client(params.Addr).CreateReplicas(ctx)
		if err != nil {
			return err
		}

		var sErr error
		exc.Dest.ForEach(func(rfl *model.Reflect, i int) {
			req := &api.ChannelChunkServiceCreateReplicasRequest{Chunk: rfl.Pointer().(*api.ChannelChunkReplica)}
			if sErr == nil {
				sErr = stream.Send(req)
			}
		})

		if sErr != nil {
			return sErr
		}

		if _, cErr := stream.CloseAndRecv(); cErr != nil {
			return cErr
		}
	}
	return nil
}

func (s *ServiceRemoteRPC) DeleteReplica(ctx context.Context, qp []RemoteReplicaDeleteOpts) error {
	for _, params := range qp {
		req := &api.ChannelChunkServiceDeleteReplicasRequest{Id: params.PKC.Strings()}
		_, err := s.client(params.Addr).DeleteReplicas(ctx, req)
		if err != nil {
			return err
		}
	}
	return nil
}

// || RPC SERVER SCAFFOLD ||

type ServiceRemoteRPCServerScaffold struct {
	api.UnimplementedChannelChunkServiceServer
}
