package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/tsquery"
)

type deltaOutlet struct {
	pkc    model.PKChain
	s      chan *models.ChannelSample
	errors chan error
}

func (o deltaOutlet) Send(sample *models.ChannelSample) {
	if o.pkc.Contains(sample.ChannelConfigID) {
		select {
		case o.s <- sample:
		}
	}
}

func (o deltaOutlet) SendError(err error) {
	select {
	case o.errors <- err:
	}
}

type deltaInlet struct {
	ctx    context.Context
	qExec  query.Execute
	s      chan *models.ChannelSample
	errors chan error
}

func (i *deltaInlet) Stream() chan *models.ChannelSample {
	return i.s
}

func (i *deltaInlet) Errors() chan error {
	return i.errors
}

func (i *deltaInlet) Update(outlets []deltaOutlet) {
	var allPKC model.PKChain
	for _, inlet := range outlets {
		allPKC = append(allPKC, inlet.pkc...)
	}
	allPKC = allPKC.Unique()
	i.s = make(chan *models.ChannelSample, len(allPKC))
	goe := tsquery.NewRetrieve().Model(&i.s).BindExec(i.qExec).WherePKs(allPKC).GoExec(i.ctx)
	i.errors = goe.Errors
}

type delta struct {
	inlet     *deltaInlet
	outlets   []deltaOutlet
	addOutlet chan deltaOutlet
}

func (d *delta) start() {
	for {
		select {
		case e := <-d.inlet.Errors():
			d.relayErrors(e)
		case o := <-d.addOutlet:
			d.processAddOutlet(o)
		case q := <-d.inlet.Stream():
			d.relay(q)
		}
	}
}

func (d *delta) processAddOutlet(o deltaOutlet) {
	d.outlets = append(d.outlets, o)
	d.inlet.Update(d.outlets)
}

func (d *delta) relay(s *models.ChannelSample) {
	for _, o := range d.outlets {
		o.Send(s)
	}
}

func (d *delta) relayErrors(e error) {
	for _, o := range d.outlets {
		o.SendError(e)
	}
}
