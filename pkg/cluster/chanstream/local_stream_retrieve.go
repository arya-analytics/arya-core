package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	"github.com/arya-analytics/aryacore/pkg/util/route"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"time"
)

// |||| QUERY EXECUTE ||||

type localStreamRetrieve struct {
	delta *route.Delta[*models.ChannelSample, outletContext]
}

func newLocalStreamRetrieve(delta *route.Delta[*models.ChannelSample, outletContext]) *localStreamRetrieve {
	return &localStreamRetrieve{delta: delta}
}

func (lsr *localStreamRetrieve) exec(ctx context.Context, p *query.Pack) error {
	qStream, pkc := stream(p), pkOpt(p)
	ldo := &localDeltaOutlet{
		d:           lsr.delta,
		pkc:         pkc,
		qStream:     qStream,
		oValStream:  *query.ConcreteModel[*chan *models.ChannelSample](p),
		inValStream: make(chan *models.ChannelSample, len(pkc)),
	}
	go ldo.Start(ctx)
	return nil
}

// |||| OUTLET ||||

type localDeltaOutlet struct {
	d           *route.Delta[*models.ChannelSample, outletContext]
	pkc         model.PKChain
	oValStream  chan *models.ChannelSample
	qStream     *streamq.Stream
	inValStream chan *models.ChannelSample
}

func (ldo *localDeltaOutlet) Send() chan<- *models.ChannelSample {
	return ldo.inValStream
}

func (ldo *localDeltaOutlet) SendError() chan<- error {
	return ldo.qStream.Errors
}

func (ldo *localDeltaOutlet) Context() outletContext {
	return outletContext{pkc: ldo.pkc}
}

func (ldo *localDeltaOutlet) Start(ctx context.Context) {
	ldo.qStream.Segment(func() {
		ldo.d.AddOutlet(ldo)
		defer ldo.d.RemoveOutlet(ldo)
		for v := range ldo.inValStream {
			if route.CtxDone(ctx) {
				break
			}
			if ldo.pkc.Contains(v.ChannelConfigID) {
				ldo.oValStream <- v
			}
		}
	}, streamq.WithSegmentName("cluster.chanstream.localDeltaOutlet"))
}

// |||| INLET ||||

type localDeltaInlet struct {
	dr        telem.DataRate
	ctx       context.Context
	qExec     query.Execute
	valStream chan *models.ChannelSample
	errC      chan error
	pkc       model.PKChain
}

func (ldi *localDeltaInlet) Next() <-chan *models.ChannelSample {
	return ldi.valStream
}

func (ldi *localDeltaInlet) Errors() <-chan error {
	return ldi.errC
}

func (ldi *localDeltaInlet) Update(dCtx route.DeltaContext[*models.ChannelSample, outletContext]) {
	ldi.pkc = parsePKC(dCtx)
	ldi.valStream = make(chan *models.ChannelSample, len(ldi.pkc))
}

func (ldi *localDeltaInlet) start() {
	t := time.NewTicker(ldi.dr.Period().ToDuration())
	defer t.Stop()
	for range t.C {
		var samples []*models.ChannelSample
		if len(ldi.pkc) == 0 {
			continue
		}
		if err := streamq.NewTSRetrieve().
			Model(&samples).
			BindExec(ldi.qExec).
			WherePKs(ldi.pkc).
			Exec(ldi.ctx); err != nil {
			ldi.errC <- err
			continue
		}
		for _, s := range samples {
			ldi.valStream <- s
		}
	}
}
