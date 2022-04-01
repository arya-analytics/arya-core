package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	api "github.com/arya-analytics/aryacore/pkg/cluster/gen/proto/go/chanstream/v1"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/rpc"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/pool"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	"github.com/arya-analytics/aryacore/pkg/util/route"
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
	srp *retrievePool
	scp *createPool
}

func NewRemoteRPC(rpcPool *cluster.NodeRPCPool) *RemoteRPC {
	return &RemoteRPC{
		srp: newRetrievePool(context.Background(), rpcPool),
		scp: newCreatePool(context.Background(), rpcPool),
	}
}

func (r *RemoteRPC) exec(ctx context.Context, p *query.Pack) error {
	return query.Switch(ctx, p, query.Ops{
		&streamq.TSCreate{}:   r.create,
		&streamq.TSRetrieve{}: r.retrieve,
	})
}
func (r *RemoteRPC) create(ctx context.Context, p *query.Pack) error {
	s := stream(p)
	s.Segment(func() {
		for {
			rfl, cOk := p.Model().ChanRecv()
			if !cOk || route.CtxDone(ctx) {
				break
			}
			stream, err := r.scp.Acquire(rfl.StructFieldByName(csFieldNode).Interface().(*models.Node))
			if err != nil {
				s.Errors <- err
				break
			}
			exc := newExchange(rfl)
			exc.ToDest()
			if sErr := stream.Send(&api.CreateRequest{Sample: exc.Dest().Pointer().(*api.ChannelSample)}); sErr != nil {
				s.Errors <- sErr
			}
			r.scp.Release(stream)
		}
	})
	return nil
}

func (r *RemoteRPC) retrieve(ctx context.Context, p *query.Pack) error {
	goe, nodes, pkc := stream(p), nodeOpt(p), pkOpt(p)
	for _, n := range nodes {
		s, err := r.srp.Acquire(n)
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
				if sErr == io.EOF || route.CtxDone(ctx) {
					break
				}
				if sErr != nil {
					goe.Errors <- sErr
					break
				}
				exc := rpc.NewModelExchange(res.Sample, &models.ChannelSample{})
				exc.ToDest()
				p.Model().ChanSend(exc.Dest())
			}
		}()
	}
	return nil
}

// |||| POOL ||||

// || RETRIEVE ADAPTER ||

type retrievePool struct {
	*pool.Pool[*models.Node]
}

func newRetrievePool(ctx context.Context, rpcPool *cluster.NodeRPCPool) *retrievePool {
	p := &retrievePool{pool.New[*models.Node]()}
	p.Factory = &retrieveStreamFactory{ctx: context.Background(), rpcPool: rpcPool}
	return p
}

func (r *retrievePool) Acquire(n *models.Node) (*retrieveAdapter, error) {
	a, err := r.Pool.Acquire(n)
	return a.(*retrieveAdapter), err
}

type retrieveAdapter struct {
	nodePK int
	api.ChannelStreamService_RetrieveClient
}

func (r *retrieveAdapter) Acquire() {

}

func (r *retrieveAdapter) Healthy() bool {
	return true
}

func (r *retrieveAdapter) Release() {

}

func (r *retrieveAdapter) Match(n *models.Node) bool {
	return r.nodePK == n.ID
}

// || RETRIEVE STREAM FACTORY ||

type retrieveStreamFactory struct {
	ctx     context.Context
	rpcPool *cluster.NodeRPCPool
}

func (r *retrieveStreamFactory) NewAdapt(n *models.Node) (pool.Adapt[*models.Node], error) {
	c, err := r.rpcPool.Retrieve(n)
	if err != nil {
		return nil, err
	}
	s, err := api.NewChannelStreamServiceClient(c).Retrieve(r.ctx)
	if err != nil {
		return nil, err
	}
	return &retrieveAdapter{
		nodePK:                              n.ID,
		ChannelStreamService_RetrieveClient: s,
	}, nil
}

func (r *retrieveStreamFactory) Match(*models.Node) bool {
	return true
}

// || CREATE POOL ||

type createPool struct {
	*pool.Pool[*models.Node]
}

func (c *createPool) Acquire(n *models.Node) (*createAdapter, error) {
	a, err := c.Pool.Acquire(n)
	return a.(*createAdapter), err
}

func newCreatePool(ctx context.Context, rpcPool *cluster.NodeRPCPool) *createPool {
	p := &createPool{pool.New[*models.Node]()}
	p.Factory = &createStreamFactory{ctx: context.Background(), rpcPool: rpcPool}
	return p
}

// || CREATE ADAPTER ||

type createAdapter struct {
	nodePK int
	api.ChannelStreamService_CreateClient
	demand pool.Demand
}

func (c *createAdapter) Acquire() {
	c.demand.Increment()
}

func (c *createAdapter) Healthy() bool {
	return true
}

func (c *createAdapter) Release() {
	c.demand.Decrement()
}

func (c *createAdapter) Match(n *models.Node) bool {
	return c.nodePK == n.ID
}

// || CREATE STREAM FACTORY ||

type createStreamFactory struct {
	ctx     context.Context
	rpcPool *cluster.NodeRPCPool
}

func (c *createStreamFactory) NewAdapt(n *models.Node) (pool.Adapt[*models.Node], error) {
	client, err := c.rpcPool.Retrieve(n)
	if err != nil {
		return nil, err
	}
	s, err := api.NewChannelStreamServiceClient(client).Create(c.ctx)
	if err != nil {
		return nil, err
	}
	return &createAdapter{
		nodePK:                            n.ID,
		ChannelStreamService_CreateClient: s,
	}, nil
}

func (c *createStreamFactory) Match(node *models.Node) bool {
	return true
}
