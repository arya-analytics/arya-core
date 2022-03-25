package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/tsquery"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"time"
)

type LocalStorage struct {
	relay *localRelay
	qe    query.Execute
}

func (s *LocalStorage) exec(ctx context.Context, p *query.Pack) error {
	return query.Switch(ctx, p, query.Ops{
		&tsquery.Create{}:   s.create,
		&tsquery.Retrieve{}: s.retrieve,
		&query.Create{}:     s.qe,
		&query.Delete{}:     s.qe,
		&query.Retrieve{}:   s.qe,
		&query.Delete{}:     s.qe,
	})
}

func (s *LocalStorage) create(ctx context.Context, p *query.Pack) error {
	goe, ok := tsquery.RetrieveGoExecOpt(p)
	if !ok {
		panic("chanstream queries must be run using goexec")
	}
	for {
		sample, sampleOK := p.Model().ChanRecv()
		if !sampleOK {
			break
		}
		if err := tsquery.NewCreate().Model(sample).BindExec(s.qe).Exec(ctx); err != nil {
			goe.Errors <- err
		}
	}
	return nil
}

func (s *LocalStorage) retrieve(ctx context.Context, p *query.Pack) error {
	_, ok := tsquery.RetrieveGoExecOpt(p)
	if ok {
		s.relay.add <- p
		return nil
	}
	return s.qe(ctx, p)
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
	_, ok = tsquery.RetrieveGoExecOpt(p)
	if !ok {
		panic("queries to local relay must use a GoExec opt")
	}
	lr.pc = append(lr.pc, p)
}

func (lr *localRelay) pkc() (pkc model.PKChain) {
	for _, p := range lr.pc {
		nPKC, _ := query.PKOpt(p)
		pkc = append(pkc, nPKC...)
	}
	return pkc
}

func (lr *localRelay) exec() {
	lr.filterDoneQueries()
	s, err := lr.retrieveSamples()
	if err != nil {
		lr.relayError(err)
	}
	lr.relaySamples(s)
}

func (lr *localRelay) retrieveSamples() (*model.Reflect, error) {
	samples := model.NewReflect(&[]*models.ChannelSample{})
	return samples, tsquery.NewRetrieve().Model(samples).BindExec(lr.qe).WherePKs(lr.pkc()).Exec(lr.ctx)
}

func (lr *localRelay) relaySamples(samples *model.Reflect) {
	samples.ForEach(func(rfl *model.Reflect, _ int) {
		for _, p := range lr.pc {
			pkc, _ := query.PKOpt(p)
			if pkc.Contains(rfl.PK()) {
				p.Model().ChanTrySend(rfl)
			}
		}
	})
}

func (lr *localRelay) relayError(err error) {
	for _, p := range lr.pc {
		o, _ := tsquery.RetrieveGoExecOpt(p)
		select {
		case o.Errors <- err:
		default:
		}
	}
}

func (lr *localRelay) filterDoneQueries() (newPS []*query.Pack) {
	for _, p := range lr.pc {
		o, _ := tsquery.RetrieveGoExecOpt(p)
		select {
		case <-o.Done:
		default:
			newPS = append(newPS, p)
		}
	}
	return newPS
}
