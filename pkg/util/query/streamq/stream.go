package streamq

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

type Stream struct {
	SegCount int
	Errors   chan error
	Ctx      context.Context
}

func (s *Stream) Segment(f func()) {
	s.SegCount++
	go f()
}

const (
	streamOptKey query.OptKey = "goExec"
	errBuffSize               = 10
)

func NewStreamOpt(ctx context.Context, p *query.Pack) *Stream {
	errors := make(chan error, errBuffSize)
	s := &Stream{Errors: errors, Ctx: ctx}
	BindStreamOpt(p, s)
	return s
}

func BindStreamOpt(p *query.Pack, s *Stream) {
	p.SetOpt(streamOptKey, s)
}

func StreamOpt(p *query.Pack) (*Stream, bool) {
	opt, ok := p.RetrieveOpt(streamOptKey)
	if !ok {
		return &Stream{}, false
	}
	return opt.(*Stream), true
}
