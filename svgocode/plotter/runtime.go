package plotter

import (
	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/math64"
)

type RuntimeConfig struct {
	UnitType math64.UnitType
}

func NewRuntimeConfig() *RuntimeConfig {
	r := new(RuntimeConfig)
	r.UnitType = math64.UnitMM
	return r
}

func (r *RuntimeConfig) SetUnitType(u math64.UnitType) {
	switch u {
	case math64.UnitMM, math64.UnitIN:
		break
	default:
		llog.Panicf("Unsupported unit type (%s). Must be 'mm' or 'in'", u)
	}
	r.UnitType = u
}
