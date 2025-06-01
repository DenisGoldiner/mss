package mss

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/DenisGoldiner/kit/kerr"
	"github.com/DenisGoldiner/merec"
)

const (
	numSizes          = 3
	bufferSize        = 20
	signalsBufferSize = 100
)

// Sizer is an interface to be implemented by tasks that should be grouped by size.
type Sizer interface {
	Size() int
}

// LimitedSizedQueueRejectPreprocessor is an implementation of the preprocessor interface.
// It sorts input tasks by size and plans them according to priorities.
type LimitedSizedQueueRejectPreprocessor[In, Out any] struct {
	logger        *slog.Logger
	incomeIntense int
}

// NewLimitedSizedQueueRejectPreprocessor is a constructor for the LimitedSizedQueueRejectPreprocessor.
func NewLimitedSizedQueueRejectPreprocessor[In, Out any](
	logger *slog.Logger,
	incomeIntense int,
) LimitedSizedQueueRejectPreprocessor[In, Out] {
	return LimitedSizedQueueRejectPreprocessor[In, Out]{logger: logger, incomeIntense: incomeIntense}
}

func (p LimitedSizedQueueRejectPreprocessor[In, Out]) inQueue(
	inCh <-chan In,
	postProcessCh chan<- merec.Result[Out],
) <-chan In {
	sizedChanPool := p.sortTasks(inCh, postProcessCh)
	signalsChanPool := merec.SpawnResChanPool[struct{}](len(sizedChanPool), signalsBufferSize)
	done, globalQueueCh := merec.MergeSignalChanPool(sizedChanPool, signalsChanPool)

	go p.prioritizeTasks(done, sizedChanPool, signalsChanPool)

	return globalQueueCh
}

func (p LimitedSizedQueueRejectPreprocessor[In, Out]) prioritizeTasks(
	done <-chan struct{},
	sizedChanPool []chan In,
	signalsChanPool []chan struct{},
) {
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			p.overflowChance(sizedChanPool, signalsChanPool)
			p.sizeBasedChance(signalsChanPool)
		}
	}
}

func (LimitedSizedQueueRejectPreprocessor[In, Out]) overflowChance(
	sizedChans []chan In,
	sizeSignals []chan struct{},
) {
	for i, sch := range sizedChans {
		if len(sch) < cap(sch) {
			continue
		}

		merec.TrySend(sizeSignals[i], struct{}{})
	}
}

func (LimitedSizedQueueRejectPreprocessor[In, Out]) sizeBasedChance(sizeSignals []chan struct{}) {
	for i := range sizeSignals {
		for j := 0; j < len(sizeSignals)-i; j++ {
			merec.TrySend(sizeSignals[i], struct{}{})
		}
	}
}

func (p LimitedSizedQueueRejectPreprocessor[In, Out]) sortTasks(
	inCh <-chan In,
	postProcessCh chan<- merec.Result[Out],
) []chan In {
	sizedChans := merec.SpawnResChanPool[In](numSizes, bufferSize)

	go func() {
		for in := range inCh {
			p.sortTask(sizedChans, postProcessCh, in)
		}

		for i := range sizedChans {
			close(sizedChans[i])
		}
	}()

	return sizedChans
}

func (LimitedSizedQueueRejectPreprocessor[In, Out]) sortTask(
	sizedChans []chan In,
	postProcessCh chan<- merec.Result[Out],
	in In,
) {
	inSizer, ok := any(in).(Sizer)
	if !ok {
		err := liberr.WrapMsg("the task does not support Size detection", ErrInvalidTask)
		postProcessCh <- merec.ErrorResult[Out](err)

		return
	}

	size := inSizer.Size()

	if size == 0 || size-1 > len(sizedChans) {
		err := liberr.WrapMsg("the task Size is not supported", ErrInvalidTask)
		postProcessCh <- merec.ErrorResult[Out](err)

		return
	}

	select {
	case sizedChans[size-1] <- in:
	default:
		msg := fmt.Sprintf("failed to plan the task Size %d execution", size)
		err := liberr.WrapMsg(msg, ErrGlobalQueueOverflow)
		postProcessCh <- merec.ErrorResult[Out](err)
	}
}
