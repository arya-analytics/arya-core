package streamq

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"sync"
)

type Stream struct {
	mu       sync.Mutex
	Segments map[Segment]bool
	Errors   chan error
	Ctx      context.Context
}

func (s *Stream) Segment(f func(), opts ...SegmentOpt) {
	seg := newSegment(opts...)
	s.mu.Lock()
	s.Segments[seg] = true
	s.mu.Unlock()
	go func() {
		defer func() {
			s.mu.Lock()
			s.Segments[seg] = false
			s.mu.Unlock()
		}()
		f()
	}()
}

type Segment struct {
	name string
}

type SegmentOpt func(s Segment) Segment

func newSegment(opts ...SegmentOpt) Segment {
	s := Segment{}
	for _, opt := range opts {
		s = opt(s)
	}
	return s
}

func WithSegmentName(name string) SegmentOpt {
	return func(s Segment) Segment {
		s.name = name
		return s
	}
}

const (
	streamOptKey query.OptKey = "goExec"
	errBuffSize               = 10
)

func NewStreamOpt(ctx context.Context, p *query.Pack) *Stream {
	errors := make(chan error, errBuffSize)
	s := &Stream{Errors: errors, Ctx: ctx, Segments: make(map[Segment]bool)}
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
