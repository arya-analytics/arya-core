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

type ServiceRemoteRPC struct {
	pool         *cluster.NodeRPCPool
	retrievePool map[int]api.ChannelStreamService_RetrieveClient
	createPool   map[int]api.ChannelStreamService_CreateClient
}

func NewServiceRemoteRPC(pool *cluster.NodeRPCPool) *ServiceRemoteRPC {
	return &ServiceRemoteRPC{
		pool:         pool,
		retrievePool: map[int]api.ChannelStreamService_RetrieveClient{},
		createPool:   map[int]api.ChannelStreamService_CreateClient{},
	}
}

func (s *ServiceRemoteRPC) client(node *models.Node) (api.ChannelStreamServiceClient, error) {
	conn, err := s.pool.Retrieve(node)
	if err != nil {
		return nil, err
	}
	return api.NewChannelStreamServiceClient(conn), nil
}

func (s *ServiceRemoteRPC) retrieveCreateStream(ctx context.Context, node *models.Node) (api.ChannelStreamService_CreateClient, error) {
	stream, ok := s.createPool[node.ID]
	if ok {
		return stream, nil
	}
	c, err := s.client(node)
	if err != nil {
		return nil, err
	}
	return c.Create(ctx)
}

func (s *ServiceRemoteRPC) Create(ctx context.Context, p *query.Pack) error {
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
		stream, err := s.retrieveCreateStream(ctx, rfl.StructFieldByName(csFieldNode).Interface().(*models.Node))
		if err != nil {
			errors <- err
		}
		exc := newExchange(rfl)
		exc.ToDest()
		if err := stream.Send(&api.CreateRequest{CCR: exc.Dest().Pointer().(*api.ChannelSample)}); err != nil {
			errors <- err
		}
	}
	return nil
}
