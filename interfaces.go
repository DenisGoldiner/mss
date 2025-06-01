package mss

import "github.com/DenisGoldiner/merec"

type sourcer[In any] interface {
	feed() <-chan In
}

// TODO: could be renamed to queue manager
type preprocessor[In, Out any] interface {
	inQueue(inCh <-chan In, postProcessCh chan<- merec.Result[Out]) <-chan In
}

type postprocessor[Out any] interface {
	ok(Out)
	err(error)
}
