package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/tsquery"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"time"
)

type LocalStorage struct {
	store storage.Storage
}

func (s *LocalStorage) create(ctx context.Context, p *query.Pack) error {
	goExecOpt, ok := tsquery.RetrieveGoExecOpt(p)
	if !ok {
		panic("chanstream queries must be run using goexec")
	}
	for {
		sample, sampleOK := p.Model().ChanRecv()
		if !sampleOK {
			break
		}
		if err := tsquery.NewCreate().Model(sample).Exec(ctx); err != nil {
			goExecOpt.Errors <- err
		}
	}
	return nil
}

func (s *LocalStorage) retrieve(ctx context.Context, p *query.Pack) error {

}

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
	samples, err := lr.retrieveSamples()
	if err != nil {
		lr.relayError(err)
	}
	lr.relaySamples(samples)
	lr.closeDoneQueries()
}

func (lr *localRelay) retrieveSamples() (*model.Reflect, error) {
	samples := model.NewReflect(&[]*models.ChannelSample{})
	return samples, tsquery.NewRetrieve().Model(samples).BindExec(lr.qe).WherePKs(lr.pkc()).Exec(lr.ctx)
}

func (lr *localRelay) relaySamples(samples *model.Reflect) {
	samples.ForEach(func(rfl *model.Reflect, i int) {
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
		o, ok := tsquery.RetrieveGoExecOpt(p)
		if !ok {
			panic("queries must use GoExec")
		}
		select {
		case o.Errors <- err:
		default:
		}
	}
}

func (lr *localRelay) closeDoneQueries() {
	var newPS []*query.Pack
	for _, p := range lr.pc {
		o, ok := tsquery.RetrieveGoExecOpt(p)
		if !ok {
			panic("queries must use GoExec")
		}
		select {
		case <-o.Done:
		default:
			newPS = append(newPS, p)
		}
	}
	lr.pc = newPS
}
