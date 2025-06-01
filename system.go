package mss

import (
	"context"
	"log/slog"

	"github.com/DenisGoldiner/merec"
)

// System is a simulation of the mass service system.
type System[In, Out any] struct {
	source   sourcer[In]
	pre      preprocessor[In, Out]
	post     postprocessor[Out]
	poolSize int
	logger   *slog.Logger
}

// NewSystem initiates the mass service System with
func NewSystem[In, Out any](
	logger *slog.Logger,
	sourcer sourcer[In],
	preprocessor preprocessor[In, Out],
	postprocessor postprocessor[Out],
	poolSize int,
) System[In, Out] {
	return System[In, Out]{
		logger:   logger,
		source:   sourcer,
		pre:      preprocessor,
		post:     postprocessor,
		poolSize: poolSize,
	}
}

// Run starts the System simulation. The processing happens asynchronously.
// The function returns the outputs chan at the beginning.
// If the client does not pop results from the outCh, it will block the processing at some point.
func (s System[In, Out]) Run(
	ctx context.Context,
	call merec.Call[In, Out],
) <-chan merec.Result[Out] {
	inCh := s.source.feed()
	outCh := s.run(ctx, call, inCh)

	return outCh
}

// run is the heart of the task processing engine.
func (s System[In, Out]) run(
	ctx context.Context,
	ef merec.Call[In, Out],
	inCh <-chan In,
) <-chan merec.Result[Out] {
	// TODO: think about the size of the buffer for the postProcessCh
	postProcessCh := make(chan merec.Result[Out], cap(inCh))
	outCh := s.postprocess(postProcessCh)

	go func() {
		globalQueueCh := s.pre.inQueue(inCh, postProcessCh)
		workerResCh, _ := merec.RunWorkerPool(ctx, globalQueueCh, ef, s.poolSize, 0)

		for res := range workerResCh {
			postProcessCh <- res
		}

		close(postProcessCh)
	}()

	return outCh
}

func (s System[In, Out]) postprocess(
	postProcessCh <-chan merec.Result[Out],
) <-chan merec.Result[Out] {
	outCh := make(chan merec.Result[Out], cap(postProcessCh))

	go func() {
		for res := range postProcessCh {
			if res.Err() != nil {
				s.post.err(res.Err())
				outCh <- res

				continue
			}

			s.post.ok(res.Value())
			outCh <- res
		}

		close(outCh)
	}()

	return outCh
}
