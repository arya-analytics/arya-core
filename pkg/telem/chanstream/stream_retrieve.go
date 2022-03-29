package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type StreamRetrieve struct {
	delta  *delta
	stream chan *models.ChannelSample
	errors chan error
	pkc    model.PKChain
}

func newStreamRetrieve(delta *delta) *StreamRetrieve {
	return &StreamRetrieve{
		delta:  delta,
		stream: make(chan *models.ChannelSample, 1),
		errors: make(chan error, 1),
	}
}

func (s *StreamRetrieve) Start(ctx context.Context) chan *models.ChannelSample {
	s.listen()
	return s.stream
}

func (s *StreamRetrieve) WherePKC(pkc model.PKChain) *StreamRetrieve {
	s.pkc = pkc
	return s
}

func (s *StreamRetrieve) Errors() chan error {
	return s.errors
}

func (s *StreamRetrieve) listen() {
	s.delta.addOutlet <- deltaOutlet{s: s.stream, pkc: s.pkc, errors: s.errors}
}
