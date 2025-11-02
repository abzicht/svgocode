package ordering

import (
	"slices"

	"github.com/abzicht/svgocode/svgocode/gcode"
)

// An orderer that changes nothing
type Lifo struct {
}

func NewLifo() *Lifo {
	return new(Lifo)
}

func (l *Lifo) Order(gcodes []*gcode.Gcode) []*gcode.Gcode {
	g2 := slices.Clone(gcodes)
	slices.Reverse(g2)
	return g2
}
