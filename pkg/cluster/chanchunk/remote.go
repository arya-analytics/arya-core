package chanchunk

import (
	"context"
	api "github.com/arya-analytics/aryacore/pkg/cluster/gen/proto/go/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/rpc"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"io"
)

/// |||| INTERFACE ||||

type ServiceRemote interface {
	// |||| REPLICA ||||

	RetrieveReplicas(ctx context.Context, ccr *model.Reflect, qp []RemoteReplicaRetrieveParams) error
	CreateReplicas(ctx context.Context, qp []RemoteReplicaCreateParams) error
	DeleteReplicas(ctx context.Context, qp []RemoteReplicaDeleteParams) error
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

// |||| RPC IMPLEMENTATION ||||

func catalogRemoteRPC() model.Catalog {
	return model.Catalog{
		&api.ChannelChunkReplica{},
	}

}
func newExchange(m *model.Reflect) *model.Exchange {
	return rpc.NewModelExchange(m.Pointer(), catalogRemoteRPC().New(m.Pointer()))
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

func (s *ServiceRemoteRPC) RetrieveReplicas(ctx context.Context, ccr *model.Reflect, qp []RemoteReplicaRetrieveParams) error {
	exc := newExchange(ccr)
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

func (s *ServiceRemoteRPC) CreateReplicas(ctx context.Context, qp []RemoteReplicaCreateParams) error {
	for _, params := range qp {
		exc := newExchange(params.Model)
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

func (s *ServiceRemoteRPC) DeleteReplicas(ctx context.Context, qp []RemoteReplicaDeleteParams) error {
	for _, params := range qp {
		req := &api.ChannelChunkServiceDeleteReplicasRequest{Id: params.PKC.Strings()}
		_, err := s.client(params.Addr).DeleteReplicas(ctx, req)
		if err != nil {
			return err
		}
	}
	return nil
}
