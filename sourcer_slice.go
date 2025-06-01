package mss

// SliceSourcer is an implementation for the sourcer interface.
// It uses the external slice as a source for input tasks.
type SliceSourcer[In any] struct {
	inputs []In
}

// NewSliceSourcer is a constructor for the SliceSourcer.
func NewSliceSourcer[In any](inputs []In) SliceSourcer[In] {
	return SliceSourcer[In]{inputs: inputs}
}

func (s SliceSourcer[In]) feed() <-chan In {
	inputIntense := len(s.inputs)
	inCh := make(chan In, inputIntense)

	go func() {
		for _, in := range s.inputs {
			inCh <- in
		}

		defer close(inCh)
	}()

	return inCh
}
