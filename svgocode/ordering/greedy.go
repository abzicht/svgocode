package ordering

import (
	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/gcode"
	"github.com/abzicht/svgocode/svgocode/math64"
)

// An orderer that performs greedy optimization.
// Works for big input, but results are most certainly not perfect.
type Greedy struct {
}

func NewGreedy() *Greedy {
	return new(Greedy)
}

func remove(g []*gcode.Gcode, i int) []*gcode.Gcode {
	g[i] = g[len(g)-1]
	return g[:len(g)-1]
}

// Time: O(n^2), Space: O(n)
func (gr *Greedy) Order(gcodes []*gcode.Gcode) []*gcode.Gcode {
	if len(gcodes) == 0 {
		return gcodes
	}
	var ordered []*gcode.Gcode

	candidates := make([]*gcode.Gcode, len(gcodes))
	if n := copy(candidates, gcodes); n < len(gcodes) {
		llog.Panicf("Failed to copy gcode list for ordering (only copied %d).", n)
	}

	current := candidates[0]
	ordered = append(ordered, current)
	candidates = candidates[1:]

	for len(candidates) > 0 {
		var bestIndex int = 0
		var bestDist math64.Float = current.EndCoord.DistEuclid(candidates[bestIndex].StartCoord)
		for j, candidate := range candidates {
			dist := current.EndCoord.DistEuclid(candidate.StartCoord)
			if bestDist > dist {
				bestIndex = j
				bestDist = dist
			}
		}
		current = candidates[bestIndex]
		ordered = append(ordered, current)
		candidates = remove(candidates, bestIndex)
	}
	return ordered
}
