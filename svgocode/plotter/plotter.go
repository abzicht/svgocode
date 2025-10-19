package plotter

import (
	"io"

	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/math64"
	"gopkg.in/yaml.v2"
)

type Plate struct {
	// Center coordinates
	Center math64.VectorF2 `yaml:"center"`
	// Minimum coordinates
	Min math64.VectorF3 `yaml:"max"`
	// Maximum coordinates
	Max math64.VectorF3 `yaml:"min"`
}

type PlotterConfig struct {
	GcodePrefix   string       `yaml:"gprefix"`
	GcodeSuffix   string       `yaml:"gsuffix"`
	Plate         Plate        `yaml:"plate"`
	DrawHeight    math64.Float `yaml:"drawing-height"`
	RetractHeight math64.Float `yaml:"retract-height"`
	DrawSpeed     math64.Speed `yaml:"draw-speed"`
	RetractSpeed  math64.Speed `yaml:"retract-speed"`
	// RemoveComments: Strip produced gcode from all comments
	RemoveComments bool `yaml:"remove-comments"`
}

// Read a PlotterConfig struct from a reader in YAML-format and return it.
func InitPlotterConfig(r io.Reader) (*PlotterConfig, error) {
	decoder := yaml.NewDecoder(r)
	p := new(PlotterConfig)
	err := decoder.Decode(p)
	return p, err
}

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
