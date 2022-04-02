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
)

// |||| QUERY EXECUTE ||||

type streamRetrieve struct {
	pool *streamRetrievePool
}

func newStreamRetrieve(pool *streamRetrievePool) *streamRetrieve {
	return &streamRetrieve{pool: pool}
}

func (sr *streamRetrieve) exec(ctx context.Context, p *query.Pack) error {
	qStream, nodes, pkc := stream(p), nodeOpt(p), pkOpt(p)
	for _, n := range nodes {
		s, err := sr.pool.Acquire(n)
		if err != nil {
			return err
		}
		d := &deltaOutlet{
			d:           s.delta,
			pkc:         pkc,
			qStream:     qStream,
			oValStream:  *query.ConcreteModel[*chan *models.ChannelSample](p),
			inValStream: make(chan *models.ChannelSample, len(pkc)),
		}
		d.Start(ctx)
	}
	return nil
}

// |||| OUTLET ||||

type outletContext struct {
	pkc model.PKChain
}

type deltaOutlet struct {
	d           *route.Delta[*models.ChannelSample, outletContext]
	pkc         model.PKChain
	oValStream  chan *models.ChannelSample
	qStream     *streamq.Stream
	inValStream chan *models.ChannelSample
}

func (d *deltaOutlet) SendError() chan<- error {
	return d.qStream.Errors
}

func (d *deltaOutlet) Send() chan<- *models.ChannelSample {
	return d.inValStream
}

func (d *deltaOutlet) Context() outletContext {
	return outletContext{pkc: d.pkc}
}

func (d *deltaOutlet) Start(ctx context.Context) {
	d.qStream.Segment(func() {
		d.d.AddOutlet(d)
		defer d.d.RemoveOutlet(d)
		for v := range d.inValStream {
			if route.CtxDone(ctx) {
				break
			}
			if d.pkc.Contains(v.ChannelConfigID) {
				d.oValStream <- v
			}
		}
	})
}

// |||| POOL ||||

type streamRetrievePool struct {
	*pool.Pool[*models.Node]
}

func newStreamRetrievePool(ctx context.Context, rpcPool *cluster.NodeRPCPool) *streamRetrievePool {
	p := &streamRetrievePool{pool.New[*models.Node]()}
	p.Factory = &streamRetrieveFactory{ctx: context.Background(), rpcPool: rpcPool}
	return p
}

func (r *streamRetrievePool) Acquire(n *models.Node) (*streamRetrieveAdapter, error) {
	a, err := r.Pool.Acquire(n)
	return a.(*streamRetrieveAdapter), err
}

// |||| ADAPTER ||||

type streamRetrieveAdapter struct {
	nodePK    int
	rpcStream api.ChannelStreamService_RetrieveClient
	delta     *route.Delta[*models.ChannelSample, outletContext]
}

func (r *streamRetrieveAdapter) Acquire() {

}

func (r *streamRetrieveAdapter) Healthy() bool {
	return true
}

func (r *streamRetrieveAdapter) Release() {

}

func (r *streamRetrieveAdapter) Match(n *models.Node) bool {
	return r.nodePK == n.ID
}

// || DELTA INLET ||

type deltaInlet struct {
	rpcStream    api.ChannelStreamService_RetrieveClient
	errStream    chan error
	sampleStream chan *models.ChannelSample
}

func (d *deltaInlet) Next() <-chan *models.ChannelSample {
	return d.sampleStream
}

func (d *deltaInlet) Errors() <-chan error {
	return d.errStream
}

func (d *deltaInlet) Update(dCtx route.DeltaContext[*models.ChannelSample, outletContext]) {
	pkc := parsePKC(dCtx)
	d.sampleStream = make(chan *models.ChannelSample, len(pkc))
	if err := d.rpcStream.SendMsg(&api.RetrieveRequest{PKC: pkc.Strings()}); err != nil {
		d.errStream <- err
	}
}

func (d *deltaInlet) Start() {
	for {
		resp, err := d.rpcStream.Recv()
		if err != nil {
			d.errStream <- err
			return
		}
		s := &models.ChannelSample{}
		rpc.NewModelExchange(resp.Sample, s).ToDest()
		d.sampleStream <- s
	}
}

func parsePKC(dCtx route.DeltaContext[*models.ChannelSample, outletContext]) (pkc model.PKChain) {
	for o := range dCtx.Outlets {
		pkc = append(pkc, o.Context().pkc...)
	}
	return pkc.Unique()
}

// || ADAPTER ||

// || RETRIEVE STREAM FACTORY ||

type streamRetrieveFactory struct {
	ctx     context.Context
	rpcPool *cluster.NodeRPCPool
}

func (r *streamRetrieveFactory) NewAdapt(n *models.Node) (pool.Adapt[*models.Node], error) {
	c, err := r.rpcPool.Retrieve(n)
	if err != nil {
		return nil, err
	}
	s, err := api.NewChannelStreamServiceClient(c).Retrieve(r.ctx)
	if err != nil {
		return nil, err
	}
	di := &deltaInlet{
		rpcStream:    s,
		errStream:    make(chan error, 10),
		sampleStream: make(chan *models.ChannelSample, 1),
	}
	go di.Start()
	d := route.NewDelta[*models.ChannelSample, outletContext](di)
	go d.Start()
	return &streamRetrieveAdapter{nodePK: n.ID, rpcStream: s, delta: d}, nil
}

func (r *streamRetrieveFactory) Match(*models.Node) bool {
	return true
}
