package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	"github.com/arya-analytics/aryacore/pkg/util/route"
)

type streamRetrieve struct {
	delta *route.Delta[*models.ChannelSample, outletContext]
}

func newStreamRetrieve(delta *route.Delta[*models.ChannelSample, outletContext]) *streamRetrieve {
	return &streamRetrieve{delta: delta}
}

func (sr *streamRetrieve) exec(ctx context.Context, p *query.Pack) error {
	s := stream(p)
	pkc, ok := query.PKOpt(p)
	if !ok {
		panic("no pk")
	}
	d := &deltaOutlet{qStream: s, pkc: pkc, valStream: make(chan *models.ChannelSample), model: p.Model()}
	d.Start(ctx)
	sr.delta.AddOutlet(d)
	return nil
}

// || DELTA OUTLET IMPL ||

type deltaOutlet struct {
	pkc       model.PKChain
	model     *model.Reflect
	qStream   *streamq.Stream
	valStream chan *models.ChannelSample
}

func (o *deltaOutlet) SendError() chan<- error {
	return o.qStream.Errors
}

func (o *deltaOutlet) Send() chan<- *models.ChannelSample {
	return o.valStream
}

func (o *deltaOutlet) Context() outletContext {
	return outletContext{pkc: o.pkc}
}

func (o *deltaOutlet) Start(ctx context.Context) {
	o.qStream.Segment(func() {
		for v := range o.valStream {
			select {
			case <-ctx.Done():
				return
			default:
				r := model.NewReflect(v)
				if o.pkc.Contains(r.PK()) {
					o.model.ChanSend(r)
				}
			}

		}
	})
}

type outletContext struct {
	pkc model.PKChain
}

// |||| INLET ||||

type deltaInlet struct {
	cancel    context.CancelFunc
	qExec     query.Execute
	qStream   *streamq.Stream
	valStream chan *models.ChannelSample
}

func (i *deltaInlet) Next() <-chan *models.ChannelSample {
	return i.valStream
}

func (i *deltaInlet) Errors() <-chan error {
	return i.qStream.Errors
}

func (i *deltaInlet) Update(dCtx route.DeltaContext[*models.ChannelSample, outletContext]) {
	pkc := parsePKC(dCtx)
	i.valStream = make(chan *models.ChannelSample, len(pkc))
	ctx, cancel := context.WithCancel(context.Background())
	pQStream, err := streamq.NewTSRetrieve().Model(&i.valStream).BindExec(i.qExec).WherePKs(pkc).Stream(ctx)
	if err != nil {
		i.qStream.Errors <- err
		cancel()
		return
	}
	if i.cancel != nil {
		i.cancel()
	}
	i.cancel = cancel
	i.qStream = pQStream
}

func parsePKC(dCtx route.DeltaContext[*models.ChannelSample, outletContext]) (pkc model.PKChain) {
	for o := range dCtx.Outlets {
		pkc = append(pkc, o.Context().pkc...)
	}
	return pkc.Unique()
}

func stream(p *query.Pack) *streamq.Stream {
	s, ok := streamq.StreamOpt(p)
	if !ok {
		panic("qStream not found on query")
	}
	return s
}
