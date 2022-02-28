package chanchunk

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/models"
)

// |||| ERROR ||||

const (
	errKey = "timing"
)

type TimingError struct {
	Type    TimingErrorType
	Message string
}

func (e TimingError) Error() string {
	return fmt.Sprintf("%s - %s - %s", errKey, e.Type, e.Message)
}

type TimingErrorType int

//go:generate stringer -type=TimingErrorType
const (
	TimingErrorTypeChunkDuplicate TimingErrorType = iota + 1
	TimingErrorTypeChunkOverlap
	TimingErrorTypeChunkNonSequential
)

func validateTiming(vCtx CreateValidateContext) error {
	// No overlap, we're good to go
	if vCtx.prevChunk.End() < vCtx.nextChunk.Start() {
		return nil
	}
	// Entire previous chunk is before next chunk
	if vCtx.nextChunk.Start() < vCtx.prevChunk.Start() && vCtx.nextChunk.End() < vCtx.prevChunk.Start() {
		return TimingError{Type: TimingErrorTypeChunkNonSequential}

	}
	// Next chunk is exactly the same as previous chunk
	if vCtx.nextChunk.Start() == vCtx.prevChunk.Start() && vCtx.nextChunk.End() == vCtx.prevChunk.End() {
		return TimingError{Type: TimingErrorTypeChunkDuplicate}
	}
	// If neither of those are true, we have some overlap that we need to resolve
	return TimingError{Type: TimingErrorTypeChunkOverlap}
}

func resolveTiming(sErr error, rCtx CreateResolveContext) (bool, error) {
	err, ok := sErr.(TimingError)
	if !ok {
		return false, sErr
	}
	// The only error type we're resolving rn are overlaps
	if err.Type != TimingErrorTypeChunkOverlap {
		return true, err
	}

	return true, resolveChunkOverlap(err, rCtx)
}

func resolveChunkOverlap(err TimingError, rCtx CreateResolveContext) error {
	switch rCtx.config.ConflictPolicy {
	case models.ChannelConflictPolicyDiscard:
		return discardChunkOverlap(rCtx)
	default:
		return err
	}
}

func discardChunkOverlap(rCtx CreateResolveContext) error {
	//oRange, ok := rCtx.prevChunk.Overlap(rCtx.nextChunk)
	//if !ok {
	//	return TimingError{
	//		Type:    TimingErrorTypeChunkNonSequential,
	//		Message: fmt.Sprintf("Unresolvable overlap %s with span %s", oRange, oRange.Span()),
	//	}
	//}
	return nil
}
