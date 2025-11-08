package conf

import (
	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/math64"
)

type RuntimeConfig struct {
	Plotter     *PlotterConfig
	PlotterUnit math64.UnitLength
	SvgUnit     math64.UnitLength
}

func NewRuntimeConfig(plotter *PlotterConfig, plotterUnit, svgUnit math64.UnitLength) *RuntimeConfig {
	r := new(RuntimeConfig)
	if nil == plotter {
		llog.Panicf("Pointer to plotter configuration is nil! We really need a config, please do better")
	}
	r.Plotter = plotter
	r.SetPlotterUnit(plotterUnit)
	r.SetSvgUnit(svgUnit)
	return r
}

func (r *RuntimeConfig) SetPlotterUnit(u math64.UnitLength) {
	switch u {
	case math64.UnitMM, math64.UnitIN:
		break
	default:
		llog.Panicf("Unsupported unit (%s). Must be 'mm' or 'in'", u)
	}
	r.PlotterUnit = u
}

func (r *RuntimeConfig) SetSvgUnit(u math64.UnitLength) {
	switch u {
	case math64.UnitCM, math64.UnitMM, math64.UnitIN:
		break
	default:
		llog.Panicf("Unsupported unit (%s). Must be 'cm', 'mm', or 'in'", u)
	}
	r.SvgUnit = u
}
