package ordering

import "github.com/abzicht/svgocode/svgocode/gcode"

// An orderer that changes nothing
type None struct {
}

func NewNone() *None {
	return new(None)
}

func (n *None) Order(gcodes []*gcode.Gcode) []*gcode.Gcode {
	return gcodes
}
