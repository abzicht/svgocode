package plotter

import (
	"io"
	"strings"

	"github.com/abzicht/svgocode/svgocode/math64"
	"github.com/abzicht/svgocode/svgocode/svg/svgtransform"
	"gopkg.in/yaml.v3"
)

type Plate struct {
	// Center coordinates
	Center math64.VectorF2 `yaml:"center"`
	// Minimum coordinates
	Min math64.VectorF3 `yaml:"min"`
	// Maximum coordinates
	Max math64.VectorF3 `yaml:"max"`
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
	RemoveComments bool            `yaml:"remove-comments"`
	MirrorX        bool            `yaml:"mirror-x-axis"`
	MirrorY        bool            `yaml:"mirror-y-axis"`
	PenOffset      math64.VectorF2 `yaml:"pen-offset"` // Pen may not be at [X: 0, Y: 0], but instead mounted with an offset.
	yamlPrefix     string
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
	if !p.PenOffset.Equal(math64.VectorF2{X: 0, Y: 0}) {
		chain = append(chain, svgtransform.NewTranslate(p.PenOffset))
	}
	return chain
}

func (p *PlotterConfig) YAML(indentSpaces int) string {
	var b strings.Builder
	b.WriteString(p.yamlPrefix)
	enc := yaml.NewEncoder(&b)
	enc.SetIndent(indentSpaces)
	enc.Encode(p)
	return b.String()
}
