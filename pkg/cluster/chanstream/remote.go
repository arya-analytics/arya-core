package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	api "github.com/arya-analytics/aryacore/pkg/cluster/gen/proto/go/chanstream/v1"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/rpc"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/tsquery"
)

// |||| RPC IMPLEMENTATION ||||

func catalogRemoteRPC() model.Catalog {
	return model.Catalog{
		&api.ChannelSample{},
	}
}

func newExchange(m interface{}) *model.Exchange {
	return rpc.NewModelExchange(m, catalogRemoteRPC().New(m))
}

type RemoteRPC struct {
	pool         *cluster.NodeRPCPool
	retrievePool map[int]api.ChannelStreamService_RetrieveClient
	createPool   map[int]api.ChannelStreamService_CreateClient
}

func NewRemoteRPC(pool *cluster.NodeRPCPool) *RemoteRPC {
	return &RemoteRPC{
		pool:         pool,
		retrievePool: map[int]api.ChannelStreamService_RetrieveClient{},
		createPool:   map[int]api.ChannelStreamService_CreateClient{},
	}
}

func (r *RemoteRPC) exec(ctx context.Context, p *query.Pack) error {
	return query.Switch(ctx, p, query.Ops{
		&tsquery.Create{}:   r.create,
		&tsquery.Retrieve{}: r.retrieve,
	})

}

func (r *RemoteRPC) client(node *models.Node) (api.ChannelStreamServiceClient, error) {
	conn, err := r.pool.Retrieve(node)
	if err != nil {
		return nil, err
	}
	return api.NewChannelStreamServiceClient(conn), nil
}

func (r *RemoteRPC) retrieveCreateStream(ctx context.Context, node *models.Node) (api.ChannelStreamService_CreateClient, error) {
	stream, ok := r.createPool[node.ID]
	if ok {
		return stream, nil
	}
	c, err := r.client(node)
	if err != nil {
		return nil, err
	}
	return c.Create(ctx)
}

func (r *RemoteRPC) create(ctx context.Context, p *query.Pack) error {
	goExecOpt, ok := tsquery.RetrieveGoExecOpt(p)
	if !ok {
		panic("go exec")
	}
	errors := goExecOpt.Errors
	for {
		rfl, cOk := p.Model().ChanRecv()
		if !cOk {
			break
		}
		stream, err := r.retrieveCreateStream(ctx, rfl.StructFieldByName(csFieldNode).Interface().(*models.Node))
		if err != nil {
			errors <- err
			break
		}
		exc := newExchange(rfl)
		exc.ToDest()
		if sErr := stream.Send(&api.CreateRequest{CCR: exc.Dest().Pointer().(*api.ChannelSample)}); sErr != nil {
			errors <- sErr
		}
	}
	return nil
}

func (r *RemoteRPC) retrieve(ctx context.Context, p *query.Pack) error {
	return nil
}
