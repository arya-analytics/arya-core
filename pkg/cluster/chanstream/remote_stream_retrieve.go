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

type remoteStreamRetrieve struct {
	pool *remoteStreamRetrievePool
}

func newRemoteStreamRetrieve(pool *remoteStreamRetrievePool) *remoteStreamRetrieve {
	return &remoteStreamRetrieve{pool: pool}
}

func (rsr *remoteStreamRetrieve) exec(ctx context.Context, p *query.Pack) error {
	qStream, _ := streamq.StreamOpt(p, query.PanicIfOptNotPresent())
	pkc, _ := query.PKOpt(p, query.PanicIfOptNotPresent())
	nodes := nodeOpt(p)
	for _, n := range nodes {
		s, err := rsr.pool.Acquire(n)
		if err != nil {
			return err
		}
		ldi := &remoteDeltaOutlet{
			delta:    s.delta,
			pkc:      pkc,
			qStream:  qStream,
			oSamples: *query.ConcreteModel[*chan *models.ChannelSample](p),
			iSamples: make(chan *models.ChannelSample, len(pkc)),
		}
		ldi.Start(ctx)
	}
	return nil
}

// |||| OUTLET ||||

type outletContext struct {
	pkc model.PKChain
}

type remoteDeltaOutlet struct {
	delta    *route.Delta[*models.ChannelSample, outletContext]
	pkc      model.PKChain
	oSamples chan *models.ChannelSample
	qStream  *streamq.Stream
	iSamples chan *models.ChannelSample
}

func (rdo *remoteDeltaOutlet) SendError() chan<- error {
	return rdo.qStream.Errors
}

func (rdo *remoteDeltaOutlet) Send() chan<- *models.ChannelSample {
	return rdo.iSamples
}

func (rdo *remoteDeltaOutlet) Context() outletContext {
	return outletContext{pkc: rdo.pkc}
}

func (rdo *remoteDeltaOutlet) Start(ctx context.Context) {
	rdo.qStream.Segment(func() {
		rdo.delta.AddOutlet(rdo)
		defer rdo.delta.RemoveOutlet(rdo)
		for v := range rdo.iSamples {
			if route.CtxDone(ctx) {
				break
			}
			if rdo.pkc.Contains(v.ChannelConfigID) {
				rdo.oSamples <- v
			}
		}
	}, streamq.WithSegmentName("cluster.telemstream.remoteDeltaOutlet"))
}

// |||| POOL ||||

type remoteStreamRetrievePool struct {
	*pool.Pool[*models.Node]
}

func newStreamRetrievePool(rpcPool *cluster.NodeRPCPool) *remoteStreamRetrievePool {
	p := &remoteStreamRetrievePool{pool.New[*models.Node]()}
	p.Factory = &remoteStreamRetrieveFactory{ctx: context.Background(), rpcPool: rpcPool}
	return p
}

func (rp *remoteStreamRetrievePool) Acquire(n *models.Node) (*remoteStreamRetrieveAdapter, error) {
	a, err := rp.Pool.Acquire(n)
	return a.(*remoteStreamRetrieveAdapter), err
}

// |||| ADAPTER ||||

type remoteStreamRetrieveAdapter struct {
	nodePK    int
	rpcStream api.ChannelStreamService_RetrieveClient
	delta     *route.Delta[*models.ChannelSample, outletContext]
	demand    pool.Demand
}

func (ra *remoteStreamRetrieveAdapter) Acquire() {
	ra.demand.Increment()
}

func (ra *remoteStreamRetrieveAdapter) Healthy() bool {
	return true
}

func (ra *remoteStreamRetrieveAdapter) Release() {
	ra.demand.Decrement()
}

func (ra *remoteStreamRetrieveAdapter) Match(n *models.Node) bool {
	return ra.nodePK == n.ID
}

// || DELTA INLET ||

type remoteDeltaInlet struct {
	rpcStream    api.ChannelStreamService_RetrieveClient
	errStream    chan error
	sampleStream chan *models.ChannelSample
}

func (rdi *remoteDeltaInlet) Next() <-chan *models.ChannelSample {
	return rdi.sampleStream
}

func (rdi *remoteDeltaInlet) Errors() <-chan error {
	return rdi.errStream
}

func (rdi *remoteDeltaInlet) Update(dCtx route.DeltaContext[*models.ChannelSample, outletContext]) {
	pkc := parsePKC(dCtx)
	rdi.sampleStream = make(chan *models.ChannelSample, len(pkc))
	if err := rdi.rpcStream.SendMsg(&api.RetrieveRequest{PKC: pkc.Strings()}); err != nil {
		rdi.errStream <- err
	}
}

func (rdi *remoteDeltaInlet) Start() {
	for {
		resp, err := rdi.rpcStream.Recv()
		if err != nil {
			rdi.errStream <- err
			return
		}
		s := &models.ChannelSample{}
		rpc.NewModelExchange(resp.Sample, s).ToDest()
		rdi.sampleStream <- s
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

type remoteStreamRetrieveFactory struct {
	ctx     context.Context
	rpcPool *cluster.NodeRPCPool
}

func (rf *remoteStreamRetrieveFactory) NewAdapt(n *models.Node) (pool.Adapt[*models.Node], error) {
	c, err := rf.rpcPool.Retrieve(n)
	if err != nil {
		return nil, err
	}
	s, err := api.NewChannelStreamServiceClient(c).Retrieve(rf.ctx)
	if err != nil {
		return nil, err
	}
	di := &remoteDeltaInlet{
		rpcStream:    s,
		errStream:    make(chan error, 10),
		sampleStream: make(chan *models.ChannelSample, 1),
	}
	go di.Start()
	d := route.NewDelta[*models.ChannelSample, outletContext](di)
	go d.Start()
	return &remoteStreamRetrieveAdapter{nodePK: n.ID, rpcStream: s, delta: d}, nil
}

func (rf *remoteStreamRetrieveFactory) Match(*models.Node) bool {
	return true
}
