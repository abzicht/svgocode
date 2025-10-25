package ordering

import (
	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/gcode"
)

// An orderer that performs 2-opt ordering. Perfect results, but slow on huge
// input.
type TwoOpt struct {
}

func NewTwoOpt() *TwoOpt {
	return new(TwoOpt)
}

func (tO *TwoOpt) Order(gcodes []*gcode.Gcode) []*gcode.Gcode {
	if len(gcodes) == 0 {
		return gcodes
	}
	ordered := make([]*gcode.Gcode, len(gcodes))
	if n := copy(ordered, gcodes); n < len(gcodes) {
		llog.Panicf("Failed to copy gcode list for ordering (only copied %d).", n)
	}

	improved := true
	for improved {
		improved = false
		bestDistance := gcode.TotalDistanceInBetween(ordered)
		for i := 1; i < len(gcodes)-2; i++ {
			for j := i + 1; j < len(gcodes)-1; j++ {
				tmpOrdered := make([]*gcode.Gcode, len(ordered))
				if n := copy(tmpOrdered, ordered); n < len(gcodes) {
					llog.Panicf("Failed to copy gcode list for ordering (only copied %d).", n)
				}

				// reverse the section between i and j
				for k := 0; k <= (j-i)/2; k++ {
					tmpOrdered[i+k], tmpOrdered[j-k] = tmpOrdered[j-k], tmpOrdered[i+k]
				}
				currentDistance := gcode.TotalDistanceInBetween(tmpOrdered)
				if currentDistance < bestDistance {
					bestDistance = currentDistance
					ordered = tmpOrdered
					improved = true
				}
			}
		}
	}
	return ordered
}
