package ordering

import (
	"cmp"
	"slices"

	"github.com/abzicht/svgocode/svgocode/gcode"
)

// An orderer that sorts by number of GCODE instructions per segment
// (descending).
type NumInstructions struct {
	descending bool
}

func NewNumInstructions(descending bool) *NumInstructions {
	nI := new(NumInstructions)
	nI.descending = descending
	return nI
}

// Time: O(n^2), Space: O(n)
func (nI *NumInstructions) Order(gcodes []*gcode.Gcode) []*gcode.Gcode {
	if len(gcodes) == 0 {
		return gcodes
	}
	var ordered []*gcode.Gcode = slices.Clone(gcodes)

	lenCmp := func(a, b *gcode.Gcode) int {
		return cmp.Compare(a.Code.NumInstructions(), b.Code.NumInstructions())
	}
	slices.SortFunc(ordered, lenCmp)
	if nI.descending {
		slices.Reverse(ordered)
	}
	return ordered
}
