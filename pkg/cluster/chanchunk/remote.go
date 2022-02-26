package chanchunk

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	api "github.com/arya-analytics/aryacore/pkg/cluster/gen/proto/go/chanchunk/v1"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/rpc"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"io"
)

type ServiceRemote interface {
	Retrieve(ctx context.Context, chunkReplica interface{}, qp []RemoteRetrieveOpts) error
	Create(ctx context.Context, qp []RemoteCreateOpts) error
	Delete(ctx context.Context, qp []RemoteDeleteOpts) error
}

type RemoteRetrieveOpts struct {
	Node *models.Node
	PKC  model.PKChain
}

type RemoteCreateOpts struct {
	Node         *models.Node
	ChunkReplica interface{}
}

type RemoteDeleteOpts struct {
	Node *models.Node
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
	pool *cluster.NodeRPCPool
}

func NewServiceRemoteRPC(pool *cluster.NodeRPCPool) ServiceRemote {
	return &ServiceRemoteRPC{pool: pool}
}

func (s *ServiceRemoteRPC) client(node *models.Node) (api.ChannelChunkServiceClient, error) {
	conn, err := s.pool.Retrieve(node)
	if err != nil {
		return nil, err
	}
	return api.NewChannelChunkServiceClient(conn), nil
}

func (s *ServiceRemoteRPC) Retrieve(ctx context.Context, chunkReplica interface{}, qp []RemoteRetrieveOpts) error {
	exc := newExchange(chunkReplica)
	for _, params := range qp {
		rq := &api.ChannelChunkServiceRetrieveReplicasRequest{PKC: params.PKC.Strings()}
		client, err := s.client(params.Node)
		if err != nil {
			return err
		}
		stream, err := client.RetrieveReplicas(ctx, rq)
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
			inRfl := model.NewReflect(in.CCR)
			exc.Dest.ChainAppend(inRfl)
		}
	}
	exc.ToSource()
	return nil
}

func (s *ServiceRemoteRPC) Create(ctx context.Context, qp []RemoteCreateOpts) error {
	for _, params := range qp {
		exc := newExchange(params.ChunkReplica)
		exc.ToDest()
		client, err := s.client(params.Node)
		if err != nil {
			return err
		}
		stream, err := client.CreateReplicas(ctx)
		if err != nil {
			return err
		}

		var sErr error
		exc.Dest.ForEach(func(rfl *model.Reflect, i int) {
			req := &api.ChannelChunkServiceCreateReplicasRequest{CCR: rfl.Pointer().(*api.ChannelChunkReplica)}
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

func (s *ServiceRemoteRPC) Delete(ctx context.Context, qp []RemoteDeleteOpts) error {
	for _, params := range qp {
		req := &api.ChannelChunkServiceDeleteReplicasRequest{PKC: params.PKC.Strings()}
		client, err := s.client(params.Node)
		if err != nil {
			return err
		}
		if _, err := client.DeleteReplicas(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// || RPC SERVER SCAFFOLD ||

type ServiceRemoteRPCServerScaffold struct {
	api.UnimplementedChannelChunkServiceServer
}
