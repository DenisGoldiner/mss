package mss

// ChanSourcer is a sourcer interface implementation that uses en external channel as a source for input tasks.
type ChanSourcer[In any] struct {
	inputs chan In
}

// NewChanSourcer is a constructor for the ChanSourcer.
func NewChanSourcer[In any](inputs chan In) ChanSourcer[In] {
	return ChanSourcer[In]{inputs: inputs}
}

func (s ChanSourcer[In]) feed() <-chan In {
	return s.inputs
}
