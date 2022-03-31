package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	api "github.com/arya-analytics/aryacore/pkg/cluster/gen/proto/go/chanstream/v1"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/rpc"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	"io"
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
		&streamq.TSCreate{}:   r.create,
		&streamq.TSRetrieve{}: r.retrieve,
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
	s, ok := r.createPool[node.ID]
	if ok {
		return s, nil
	}
	c, err := r.client(node)
	if err != nil {
		return nil, err
	}
	s, err = c.Create(ctx)
	if err != nil {
		return nil, err
	}
	r.createPool[node.ID] = s
	return s, nil
}

func (r *RemoteRPC) retrieveRetrieveStream(ctx context.Context, node *models.Node) (api.ChannelStreamService_RetrieveClient, error) {
	s, ok := r.retrievePool[node.ID]
	if ok {
		return s, nil
	}
	c, err := r.client(node)
	if err != nil {
		return nil, err
	}
	s, err = c.Retrieve(ctx)
	if err != nil {
		return nil, err
	}
	r.retrievePool[node.ID] = s
	return s, nil
}

func (r *RemoteRPC) create(ctx context.Context, p *query.Pack) error {
	s := stream(p)
	s.Segment(func() {
		for {
			rfl, cOk := p.Model().ChanRecv()
			if !cOk {
				break
			}
			stream, err := r.retrieveCreateStream(ctx, rfl.StructFieldByName(csFieldNode).Interface().(*models.Node))
			if err != nil {
				s.Errors <- err
				break
			}
			exc := newExchange(rfl)
			exc.ToDest()
			if sErr := stream.Send(&api.CreateRequest{CCR: exc.Dest().Pointer().(*api.ChannelSample)}); sErr != nil {
				s.Errors <- sErr
			}
		}
	})
	return nil
}

func (r *RemoteRPC) retrieve(ctx context.Context, p *query.Pack) error {
	goe, nodes, pkc := stream(p), nodeOpt(p), pkOpt(p)
	for _, n := range nodes {
		s, err := r.retrieveRetrieveStream(ctx, n)
		if err != nil {
			return err
		}
		err = s.Send(&api.RetrieveRequest{PKC: pkc.Strings()})
		if err != nil {
			return err
		}
		go func() {
			for {
				res, sErr := s.Recv()
				if sErr == io.EOF {
					break
				}
				select {
				case <-ctx.Done():
					return
				default:
				}
				if sErr != nil {
					goe.Errors <- sErr
					break
				}
				exc := rpc.NewModelExchange(res.CCR, &models.ChannelSample{})
				exc.ToDest()
				p.Model().ChanSend(exc.Dest())
			}
		}()
	}
	return nil
}
