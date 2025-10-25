package plotter

import (
	"io"

	"github.com/abzicht/svgocode/svgocode/math64"
	"github.com/abzicht/svgocode/svgocode/svg/svgtransform"
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
	MirrorX        bool `yaml:"mirror-x-axis"`
	MirrorY        bool `yaml:"mirror-y-axis"`
}

// Read a PlotterConfig struct from a reader in YAML-format and return it.
func InitPlotterConfig(r io.Reader) (*PlotterConfig, error) {
	decoder := yaml.NewDecoder(r)
	p := new(PlotterConfig)
	err := decoder.Decode(p)
	return p, err
}

// Create a list of transform commands for the given configuration.
func (p *PlotterConfig) Transform() svgtransform.TransformChain {
	chain := svgtransform.TransformChain{}
	if p.MirrorX || p.MirrorY {
		chain = append(chain, svgtransform.NewMirror(p.MirrorX, p.MirrorY, p.Plate.Center))
	}
	return chain
}
