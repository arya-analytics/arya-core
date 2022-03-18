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

func validateTiming(vCtx NextChunkContext) error {
	if vCtx.prev == nil {
		return nil
	}
	ov := vCtx.next.Overlap(vCtx.prev)
	if !ov.ChunksCompatible() {
		return TimingError{Type: TimingErrorTypeIncompatibleChunks}
	}
	if ov.IsValid() {
		return TimingError{Type: TimingErrorOverlap}
	}
	if vCtx.next.Start() < vCtx.prev.Start() {
		return TimingError{Type: TimingErrorTypeNonContiguous}
	}

	return nil
}

// |||| RESOLUTION ||||

func resolveTiming(sErr error, nCtx NextChunkContext) (bool, error) {
	err, ok := sErr.(TimingError)
	if !ok {
		return false, sErr
	}
	// The only error type we're resolving rn are overlaps
	if err.Type != TimingErrorOverlap {
		return true, err
	}
	return true, resolveChunkOverlap(err, nCtx)
}

func resolveChunkOverlap(err TimingError, nCtx NextChunkContext) error {
	switch nCtx.cfg.ConflictPolicy {
	case models.ChannelConflictPolicyDiscard:
		return discardOverlap(nCtx)
	default:
		return err
	}
}

func discardOverlap(nCtx NextChunkContext) error {
	ov := nCtx.next.Overlap(nCtx.prev)
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
			Message: fmt.Sprintf("received unresolveable conflict type %s for conflict policy %s", ov.Type(), nCtx.cfg.ConflictPolicy),
		}
	}
}
