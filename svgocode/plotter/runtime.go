package plotter

import (
	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/math64"
)

type RuntimeConfig struct {
	UnitLength math64.UnitLength
}

func NewRuntimeConfig() *RuntimeConfig {
	r := new(RuntimeConfig)
	r.UnitLength = math64.UnitMM
	return r
}

func (r *RuntimeConfig) SetUnitLength(u math64.UnitLength) {
	switch u {
	case math64.UnitMM, math64.UnitIN:
		break
	default:
		llog.Panicf("Unsupported unit type (%s). Must be 'mm' or 'in'", u)
	}
	r.UnitLength = u
}
