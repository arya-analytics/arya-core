package chanstream

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query/tsquery"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	"time"
)

type receive interface{}

type sendConfig struct {
	pks model.PKChain
}

type send interface {
	cfg() sendConfig
}

type config struct {
	dr telem.DataRate
}

type relay struct {
	cfg         config
	_addSend    chan send
	sends       map[send]bool
	_addReceive chan receive
	receives    map[receive]bool
}

func (r *relay) start() {
	t := time.NewTicker(r.cfg.dr.Period().ToDuration())
	defer t.Stop()
	for {
		select {
		case snd := <-r._addSend:
			r.processAddSend(snd)
		case rcv := <-r._addReceive:
			r.processAddReceive(rcv)
		case <-t.C:
			r.exec()
		}
	}
}

// |||| SEND MANAGEMENT ||||

func (r *relay) addSend(snd send) {
	r._addSend <- snd
}

func (r *relay) processAddSend(snd send) {
	r.sends[snd] = true
}

// |||| RECEIVE MANAGEMENT ||||

func (r *relay) addReceive(rcv receive) {
	r._addReceive <- rcv
}

func (r *relay) processAddReceive(rcv receive) {
	r.receives[rcv] = true
}

// |||| EXECUTION ||||

func (r *relay) exec() {
	pkc := r.parseSendPKC()
	samples := make(chan *models.ChannelSample, len(pkc))
	errs := make(chan error)
	errs := tsquery.NewRetrieve().Model(samples).WherePKs(pkc).GoExec(ctx, errs)
}

func (r *relay) parseSendPKC() model.PKChain {
	pkc := model.NewPKChain([]uuid.UUID{})
	for k := range r.sends {
		pkc = append(pkc, k.cfg().pks...)
	}
	return pkc.Unique()
}
