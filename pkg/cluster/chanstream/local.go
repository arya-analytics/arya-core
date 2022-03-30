package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/model/filter"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"time"
)

type LocalStorage struct {
	relay *localRelay
	qe    query.Execute
}

func (ls *LocalStorage) exec(ctx context.Context, p *query.Pack) error {
	return query.Switch(ctx, p, query.Ops{
		&streamq.TSCreate{}:   ls.create,
		&streamq.TSRetrieve{}: ls.retrieve,
		&query.Create{}:       ls.qe,
		&query.Delete{}:       ls.qe,
		&query.Retrieve{}:     ls.qe,
		&query.Delete{}:       ls.qe,
	})
}

func (ls *LocalStorage) create(ctx context.Context, p *query.Pack) error {
	stream := stream(p)
	stream.Segment(func() {
		for {
			sample, sampleOK := p.Model().ChanRecv()
			if !sampleOK {
				break
			}
			if err := streamq.NewTSCreate().Model(sample).BindExec(ls.qe).Exec(ctx); err != nil {
				stream.Errors <- err
			}
		}
	})
	return nil
}

func (ls *LocalStorage) retrieve(ctx context.Context, p *query.Pack) error {
	_, ok := streamq.StreamOpt(p)
	if ok {
		ls.relay.add <- p
		return nil
	}
	return ls.qe(ctx, p)
}

// |||| RELAY ||||

type localRelay struct {
	ctx context.Context
	qe  query.Execute
	dr  telem.DataRate
	add chan *query.Pack
	pc  []*query.Pack
}

func (lr *localRelay) start() {
	t := time.NewTicker(lr.dr.Period().ToDuration())
	for {
		select {
		case p := <-lr.add:
			lr.processAdd(p)
		case <-t.C:
			lr.exec()
		}
	}
}

func (lr *localRelay) processAdd(p *query.Pack) {
	if !p.Model().IsChan() {
		panic("local relay can't process non channel queries")
	}
	_, ok := query.PKOpt(p)
	if !ok {
		panic("queries to local relay must use a primary key")
	}
	_, ok = streamq.StreamOpt(p)
	if !ok {
		panic("queries to local relay must use a Exec opt")
	}
	lr.pc = append(lr.pc, p)
}

func (lr *localRelay) pkc() (pkc model.PKChain) {
	for _, p := range lr.pc {
		nPKC, _ := query.PKOpt(p)
		pkc = append(pkc, nPKC...)
	}
	return pkc.Unique()
}

func (lr *localRelay) exec() {
	lr.filterDoneQueries()
	s, err := lr.retrieveSamples()
	if err != nil {
		lr.relayError(err)
	}
	lr.relaySamples(s)
}

func (lr *localRelay) retrieveSamples() (samples []*models.ChannelSample, err error) {
	return samples, streamq.NewTSRetrieve().Model(&samples).BindExec(lr.qe).WherePKs(lr.pkc()).Exec(lr.ctx)
}

func (lr *localRelay) relaySamples(samples []*models.ChannelSample) {
	for _, p := range lr.pc {
		relay(p, samples)
	}
}

func relay(p *query.Pack, samples []*models.ChannelSample) {
	pkc, _ := query.PKOpt(p)
	os := model.NewReflect(&[]*models.ChannelSample{})
	filter.Exec(query.NewRetrieve().Model(&samples).WherePKs(pkc).Pack(), os)
	os.ForEach(func(rfl *model.Reflect, i int) { p.Model().ChanTrySend(rfl) })
}

func (lr *localRelay) relayError(err error) {
	for _, p := range lr.pc {
		o, _ := streamq.StreamOpt(p)
		select {
		case o.Errors <- err:
		default:
		}
	}
}

func (lr *localRelay) filterDoneQueries() (newPS []*query.Pack) {
	for _, p := range lr.pc {
		o, _ := streamq.StreamOpt(p)
		select {
		case <-o.Ctx.Done():
			//log.Infof("relay: query %s done", p.String())
		default:
			newPS = append(newPS, p)
		}
	}
	return newPS
}
