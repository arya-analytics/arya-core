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

type Error struct {
	Type    ErrorType
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("%s - %s - %s", errKey, e.Type, e.Message)
}

type ErrorType int

//go:generate stringer -type=TimingErrorType
const (
	ErrorTimingOverlap ErrorType = iota + 1
	ErrorTimingIncompatibleChunks
	ErrorTimingNonContiguous
)

// |||| VALIDATION |||||

func validateTiming(vCtx nextChunkContext) error {
	if vCtx.prev == nil {
		return nil
	}
	ov := vCtx.next.Overlap(vCtx.prev)
	if !ov.ChunksCompatible() {
		return Error{Type: ErrorTimingIncompatibleChunks}
	}
	if ov.IsValid() {
		return Error{Type: ErrorTimingOverlap}
	}
	if vCtx.next.Start() < vCtx.prev.Start() {
		return Error{Type: ErrorTimingNonContiguous}
	}

	return nil
}

// |||| RESOLUTION ||||

func resolveTiming(sErr error, nCtx nextChunkContext) (bool, error) {
	err, ok := sErr.(Error)
	if !ok {
		return false, sErr
	}
	// The only error type we're resolving rn are overlaps
	if err.Type != ErrorTimingOverlap {
		return true, err
	}
	return true, resolveChunkOverlap(err, nCtx)
}

func resolveChunkOverlap(err Error, nCtx nextChunkContext) error {
	switch nCtx.cfg.ConflictPolicy {
	case models.ChannelConflictPolicyDiscard:
		return discardOverlap(nCtx)
	default:
		return err
	}
}

func discardOverlap(nCtx nextChunkContext) error {
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
		return Error{
			Type:    ErrorTimingNonContiguous,
			Message: fmt.Sprintf("received unresolveable conflict type %s for conflict policy %s", ov.Type(), nCtx.cfg.ConflictPolicy),
		}
	}
}
