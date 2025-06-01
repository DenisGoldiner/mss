package mss

import (
	"github.com/DenisGoldiner/kit/kerr"
	"github.com/DenisGoldiner/merec"
)

// LimitedQueueRejectPreprocessor is an implementation for the preprocessor interface.
// It simulates the classic mass service system with a limited queue and rejects in case of overfill.
type LimitedQueueRejectPreprocessor[In, Out any] struct{}

// NewLimitedQueueRejectPreprocessor is a constructor for the LimitedQueueRejectPreprocessor.
func NewLimitedQueueRejectPreprocessor[In, Out any]() LimitedQueueRejectPreprocessor[In, Out] {
	return LimitedQueueRejectPreprocessor[In, Out]{}
}

// LimitedQueueRejectEngine simulates limited queue and rejects new task if the waiting line is full.
func (LimitedQueueRejectPreprocessor[In, Out]) inQueue(
	inCh <-chan In,
	postProcessCh chan<- merec.Result[Out],
) <-chan In {
	globalQueueCh := make(chan In, cap(inCh))

	go func() {
		for in := range inCh {
			select {
			case globalQueueCh <- in:
			default:
				err := liberr.WrapMsg("failed to plan the task execution", ErrGlobalQueueOverflow)
				postProcessCh <- merec.ErrorResult[Out](err)
			}
		}

		close(globalQueueCh)
	}()

	return globalQueueCh
}
