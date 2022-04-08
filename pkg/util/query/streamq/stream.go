package streamq

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"sync"
)

// Stream represents a query that returns a stream of data as opposed a single value (i.e. a streaming operation
// instead of a unary operation). This is common for time series and event data such as TSRetrieve and TSCreate.
//
// Using a Stream as a Client:
//
// Typically, a stream is returned as the first return value of a query run using .Stream(ctx). For example:
//
//		stream, err := tsquery.NewRetrieve().Stream(ctx)
//
// The error returned as the second value represents an error encountered during pipeline 'segment' assembly (things
// such as validating query parameters, resolving hosts, doing lookups on context items from the data store).
//
// The Stream returned pipes any errors encountered during the actual transportation of query results to
// Stream.Errors. In short, errors encountered during construction of the stream are returned upon completing construction, while
// errors encountered during stream operation will be piped through Stream.Errors.
//
// Stream.Ctx is the same context used to construct the stream. This is useful for canceling the stream.
//
// To cancel a Stream that receives values (a query that retrieves things), DO NOT CLOSE THE CHANNEL. Instead, cancel
// the context associated with the stream. Closing the channel will most likely result in a panic further down the road.
//
// To cancel a Stream that sends values (a query that creates, deletes, or updates things), it's ok to close the channel
// as long as you're sure nothing else is sending values to it. Cancelling the context will work as well.
//
// Using a Stream as a query provider:
//
// A stream can be retrieved using a similar API to other query options. Simply call:
//
// 		stream, ok := streamq.RetrieveStreamOpt(p)
//
// It's also common to require that a stream be used to run a certain query.
//
//		stream, _ := streamq.RetrieveStreamOpt(p, query.RequireOpt())
//
// It's useful to think of a stream of values as a comprised by segments that received values from the previous segment
// and sent values to the next segment (most likely doing some routing, filtering, modification in each stage). These
// segments will involve using goroutines.
//
// For diagnostic and debug reasons, we want to track the quantity and identity of the goroutines used to serve a query.
// Use Stream.Segment to start these goroutines.
//
// For error handling, return errors encountered during construction (i.e. outside of those goroutines) directly instead
// of using Stream.Errors. This is useful for separating error types and maximizing runtime safety. Pipe errors encountered
// within segments to Stream.Errors.
//
// To stop and garbage collect Segments, check for cancellation of the context (such as route.CtxDone()) and break the
// goroutine loop. Close any channels sending values to the next segment. DO NOT CLOSE any channels receiving values
// from the previous segment.
//
type Stream struct {
	mu sync.Mutex
	// Segments is map of the goroutines involved in running the stream. The boolean value indicates whether the goroutine
	// is running or not.
	Segments map[Segment]bool
	// Errors is a channel that receives errors encountered during the stream operation.
	Errors chan error
	// Ctx is the context used to construct the stream.
	Ctx context.Context
}

// Segment adds a goroutine as a segment of the stream. Used for observability purposes.
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

func RetrieveStreamOpt(p *query.Pack, opts ...query.OptRetrieveOpt) (*Stream, bool) {
	opt, ok := p.RetrieveOpt(streamOptKey, opts...)
	if !ok {
		return &Stream{}, false
	}
	return opt.(*Stream), true
}
