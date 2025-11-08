package conf

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
	GcodePrefix string `yaml:"gprefix"`
	GcodeSuffix string `yaml:"gsuffix"`
	// Unit that is used in PlotterConfig's variables and that will be used for gcode ('mm' or 'in')
	UnitLength    math64.UnitLength `yaml:"length-unit"`
	Plate         Plate             `yaml:"plate"`
	DrawHeight    math64.Float      `yaml:"drawing-height"`
	RetractHeight math64.Float      `yaml:"retract-height"`
	DrawSpeed     math64.Speed      `yaml:"draw-speed"`
	RetractSpeed  math64.Speed      `yaml:"retract-speed"`
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

func (p *PlotterConfig) convUnitF2(v math64.VectorF2, unit math64.UnitLength) math64.VectorF2 {
	x := math64.LengthConvert(v.X, p.UnitLength, unit)
	y := math64.LengthConvert(v.Y, p.UnitLength, unit)
	return math64.VectorF2{X: x, Y: y}
}

// Create a list of transform commands for the given configuration.
// Transform matrix is scaled to the given transformUnit
func (p *PlotterConfig) Transform(transformUnit math64.UnitLength) svgtransform.TransformChain {
	chain := svgtransform.TransformChain{}
	if p.MirrorX || p.MirrorY {
		chain = append(chain, svgtransform.NewMirror(p.MirrorX, p.MirrorY, p.convUnitF2(p.Plate.Center, transformUnit)))
	}
	if !p.PenOffset.Equal(math64.VectorF2{X: 0, Y: 0}) {
		chain = append(chain, svgtransform.NewTranslate(p.convUnitF2(p.PenOffset, transformUnit)))
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
