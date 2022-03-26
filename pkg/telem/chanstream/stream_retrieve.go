package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type StreamRetrieve struct {
	rel    *relay
	stream chan *models.ChannelSample
	pkc    model.PKChain
}

func newStreamRetrieve(rel *relay) *StreamRetrieve {
	return &StreamRetrieve{rel: rel, stream: make(chan *models.ChannelSample, 1)}
}

func (s *StreamRetrieve) Start(ctx context.Context) chan *models.ChannelSample {
	s.listen()
	return s.stream
}

func (s *StreamRetrieve) listen() {
	s.rel.addSend(s)
}

func (s *StreamRetrieve) WherePKC(pkc model.PKChain) *StreamRetrieve {
	s.pkc = pkc
	return s
}

func (s *StreamRetrieve) send() chan *models.ChannelSample {
	return s.stream
}

func (s *StreamRetrieve) cfg() sendConfig {
	return sendConfig{pks: s.pkc}
}
