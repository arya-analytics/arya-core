package chanchunk

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
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
	TimingErrorOverlap TimingErrorType = iota + 1
	TimingErrorTypeIncompatibleChunks
	TimingErrorTypeNonContiguous
)

// |||| VALIDATION |||||

func validateTiming(vCtx NextChunkValidateContext) error {
	if vCtx.prevChunk == nil {
		return nil
	}
	ov := vCtx.nextChunk.Overlap(vCtx.prevChunk)
	if !ov.ChunksCompatible() {
		return TimingError{Type: TimingErrorTypeIncompatibleChunks}
	}
	if ov.IsValid() {
		return TimingError{Type: TimingErrorOverlap}
	}
	if vCtx.nextChunk.Start() < vCtx.prevChunk.Start() {
		return TimingError{Type: TimingErrorTypeNonContiguous}
	}

	return nil
}

// |||| RESOLUTION ||||

func resolveTiming(sErr error, rCtx NextChunkResolveContext) (bool, error) {
	err, ok := sErr.(TimingError)
	if !ok {
		return false, sErr
	}
	// The only error type we're resolving rn are overlaps
	if err.Type != TimingErrorOverlap {
		return true, err
	}
	return true, resolveChunkOverlap(err, rCtx)
}

func resolveChunkOverlap(err TimingError, rCtx NextChunkResolveContext) error {
	switch rCtx.config.ConflictPolicy {
	case models.ChannelConflictPolicyDiscard:
		return discardOverlap(rCtx)
	default:
		return err
	}
}

func discardOverlap(rCtx NextChunkResolveContext) error {
	ov := rCtx.nextChunk.Overlap(rCtx.prevChunk)
	switch ov.Type() {
	case telem.OverlapTypeNoneOrInvalid:
		return nil
	case telem.OverlapTypeRightPartial:
		return ov.RemoveFromSource()
	case telem.OverlapTypeDuplicate:
		return ov.RemoveFromSource()
	case telem.OverlapTypeDestConsume:
		return ov.RemoveFromConsumed()
	default:
		return TimingError{
			Type:    TimingErrorTypeNonContiguous,
			Message: fmt.Sprintf("received unresolveable conflict type %s for conflict policy %s", ov.Type(), rCtx.config.ConflictPolicy),
		}
	}
}
