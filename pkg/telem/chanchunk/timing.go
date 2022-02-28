package chanchunk

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	log "github.com/sirupsen/logrus"
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
	TimingErrorTypeChunkOverlap TimingErrorType = iota + 1
	TimingErrorTypeIncompatibleChunks
)

// |||| VALIDATION |||||

func validateTiming(vCtx CreateValidateContext) error {
	ov := vCtx.nextChunk.Overlap(vCtx.prevChunk)
	if !ov.ChunksCompatible() {
		return TimingError{Type: TimingErrorTypeIncompatibleChunks}
	}
	if ov.IsValid() {
		return TimingError{Type: TimingErrorTypeChunkOverlap}
	}
	return nil
}

// |||| RESOLUTION ||||

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
		return discardOverlap(rCtx)
	default:
		return err
	}
}

func discardOverlap(rCtx CreateResolveContext) error {
	ov := rCtx.nextChunk.Overlap(rCtx.prevChunk)
	switch ov.Type() {
	case telem.OverlapTypeNoneOrInvalid:
		log.Warn("discard overlap received a non-overlapping or invalid chunk to resolve")
		return nil
	case telem.OverlapTypeRightPartial:
		log.Warn("received partially overlapping chunk")
		return ov.RemoveFromSource()
	case telem.OverlapTypeDuplicate:
		log.Warn("received duplicate chunk")
		return ov.RemoveFromSource()
	case telem.OverlapTypeSourceConsume:
		log.Warn("received a consumed chunk")
		return ov.RemoveFromConsumed()
	default:
		return TimingError{
			Type:    TimingErrorTypeChunkOverlap,
			Message: fmt.Sprintf("received unresolveable conflict type %s for conflict policy %s", ov.Type(), rCtx.config.ConflictPolicy),
		}
	}
}
