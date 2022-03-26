package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/model/filter"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/tsquery"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"time"
)

type receive interface {
	retrieve() ([]*models.ChannelSample, error)
}

type sendConfig struct {
	pks model.PKChain
}

type send interface {
	cfg() sendConfig
	send() chan *models.ChannelSample
}

type relay struct {
	qExec          query.Execute
	dr             telem.DataRate
	_addSend       chan send
	sends          map[send]bool
	_removeSend    chan send
	_removeReceive chan receive
	clusterQ       *query.Pack
}

func newRelay(dr telem.DataRate, qExec query.Execute) *relay {
	return &relay{
		qExec:          qExec,
		dr:             dr,
		_addSend:       make(chan send),
		sends:          make(map[send]bool),
		_removeSend:    make(chan send),
		_removeReceive: make(chan receive),
	}

}

func (r *relay) start() {
	t := time.NewTicker(r.dr.Period().ToDuration())
	defer t.Stop()
	for {
		select {
		case snd := <-r._addSend:
			r.processAddSend(snd)
		case s := <-r._removeSend:
			r.processRemoveSend(s)
		case <-t.C:
			r.exec()
		}
	}
}

// |||| SEND MANAGEMENT ||||

func (r *relay) addSend(snd send) {
	c := make(chan *models.ChannelSample, len(snd.cfg().pks))
	rfl := model.NewReflect(&c)
	q := tsquery.NewRetrieve().Model(rfl).WherePKs(snd.cfg().pks).BindExec(r.qExec)
	q.GoExec(context.Background())
	r.clusterQ = q.Pack()
	r._addSend <- snd
}

func (r *relay) removeSend(snd send) {
	r._removeSend <- snd
}

func (r *relay) processAddSend(snd send) {
	r.sends[snd] = true
}

func (r *relay) processRemoveSend(snd send) {
	delete(r.sends, snd)
}

// |||| RECEIVE MANAGEMENT ||||

// |||| EXECUTION ||||

func (r *relay) exec() {
	samples, err := r.retrieve()
	r.relay(samples)
	if err != nil {
		r.relayError(err)
	}
}

func (r *relay) retrieve() (samples []*models.ChannelSample, err error) {
	if r.clusterQ == nil {
		return samples, err
	}
	for {
		s, ok := r.clusterQ.Model().ChanTryRecv()
		if !ok {
			break
		}
		samples = append(samples, s.Pointer().(*models.ChannelSample))
	}
	return samples, nil
}

func (r *relay) relay(samples []*models.ChannelSample) {
	for snd := range r.sends {
		relayToSend(snd, samples)
	}
}

func relayToSend(s send, samples []*models.ChannelSample) {
	var os []*models.ChannelSample
	filter.Exec(query.NewRetrieve().Model(&samples).WherePKs(s.cfg().pks).Pack(), &os)
	for _, o := range os {
		s.send() <- o
	}
}

func (r *relay) pkc() model.PKChain {
	pkc := make(model.PKChain, len(r.sends))
	for snd := range r.sends {
		pkc = append(pkc, snd.cfg().pks...)
	}
	return pkc.Unique()
}

func (r *relay) relayError(err error) {}
