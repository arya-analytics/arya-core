package chanstream

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/model/filter"
	"github.com/arya-analytics/aryacore/pkg/util/query"
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
	dr             telem.DataRate
	_addSend       chan send
	sends          map[send]bool
	_addReceive    chan receive
	receives       map[receive]bool
	_removeSend    chan send
	_removeReceive chan receive
}

func (r *relay) start() {
	t := time.NewTicker(r.dr.Period().ToDuration())
	defer t.Stop()
	for {
		select {
		case snd := <-r._addSend:
			r.processAddSend(snd)
		case rcv := <-r._addReceive:
			r.processAddReceive(rcv)
		case rc := <-r._removeReceive:
			r.processRemoveReceive(rc)
		case s := <-r._removeSend:
			r.processRemoveSend(s)
		case <-t.C:
			r.exec()
		}
	}
}

// |||| SEND MANAGEMENT ||||

func (r *relay) addSend(snd send) {
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

func (r *relay) addReceive(rcv receive) {
	r._addReceive <- rcv
}

func (r *relay) removeReceive(rcv receive) {
	r._removeReceive <- rcv
}

func (r *relay) processAddReceive(rcv receive) {
	r.receives[rcv] = true
}

func (r *relay) processRemoveReceive(rcv receive) {
	delete(r.receives, rcv)
}

// |||| EXECUTION ||||

func (r *relay) exec() {
	samples, err := r.retrieve()
	r.relay(samples)
	if err != nil {
		r.relayError(err)
	}
}

func (r *relay) retrieve() (samples []*models.ChannelSample, err error) {
	for rcv := range r.receives {
		s, err := rcv.retrieve()
		if err != nil {
			return nil, err
		}
		samples = append(samples, s...)
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
	filter.Exec(query.NewRetrieve().Model(samples).WherePKs(s.cfg().pks).Pack(), &os)
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
