package ordering

import (
	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/gcode"
)

type OrderingAlg string

const (
	OrderingAlgTwoOpt = OrderingAlg("2opt")
	OrderingAlgGreedy = OrderingAlg("greedy")
	OrderingAlgNone   = OrderingAlg("none")
	OrderingAlgLifo   = OrderingAlg("reverse")
)

type OrderingI interface {
	Order([]*gcode.Gcode) []*gcode.Gcode
}

func ParseOrdering(alg OrderingAlg) OrderingI {
	switch alg {
	case OrderingAlg(""): // Default is 2opt
		fallthrough
	case OrderingAlgTwoOpt:
		return NewTwoOpt()
	case OrderingAlgGreedy:
		return NewGreedy()
	case OrderingAlgNone:
		return NewNone()
	case OrderingAlgLifo:
		return NewLifo()
	default:
		llog.Panicf("Unknown ordering algorithm: '%s'", alg)
		return nil
	}
}
