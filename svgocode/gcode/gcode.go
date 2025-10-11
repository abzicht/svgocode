package gcode

import (
	"slices"
	"strings"

	"github.com/abzicht/svgocode/svgocode/math64"
	"github.com/abzicht/svgocode/svgocode/plotter"
)

type Code struct {
	lines []string
}

func NewCode() *Code {
	return new(Code)
}

func (c *Code) Copy() *Code {
	c2 := NewCode()
	c2.lines = slices.Clone(c.lines)
	return c2
}
func (c *Code) String() string {
	return strings.Join(c.lines, "\n") + "\n"
}

func (c *Code) AppendLines(lines ...string) {
	c.lines = append(c.lines, lines...)
}

func (c *Code) Append(c2 *Code) {
	c.AppendLines(c2.lines...)
}

type Gcode struct {
	Code       *Code
	StartCoord math64.VectorF3
	EndCoord   math64.VectorF3
	BoundsMin  math64.VectorF3
	BoundsMax  math64.VectorF3
}

func NewGcode() *Gcode {
	g := new(Gcode)
	g.Code = NewCode()
	return g
}

// Create a new gcode object from the given code.
func NewGcodeFromString(code string) *Gcode {
	g := NewGcode()
	g.AppendCode(code)
	return g
}

// Create a copy of the gcode, but only of its metadata without the actual code
func (g *Gcode) CopyMeta() *Gcode {
	g2 := NewGcode()
	g2.StartCoord = g.StartCoord
	g2.EndCoord = g.EndCoord
	g2.BoundsMin = g.BoundsMin
	g2.BoundsMax = g.BoundsMax
	return g2
}

// Full copy, including the code.
func (g *Gcode) Copy() *Gcode {
	g2 := g.CopyMeta()
	g2.Code = g.Code.Copy()
	return g2
}

func (g *Gcode) String() string {
	return g.Code.String()
}

// Add newline-separated code
func (g *Gcode) AppendCode(code string) {
	if len(code) == 0 {
		return
	}
	lines := strings.Split(code, "\n")
	g.Code.AppendLines(lines...)
}

// Join two gcodes, merging their boundaries, start/end corrdinates, and code.
func (g *Gcode) AppendSimple(g2 *Gcode) {
	g.EndCoord = g2.EndCoord
	if g.BoundsMin.X > g2.BoundsMin.X {
		g.BoundsMin.X = g2.BoundsMin.X
	}
	if g.BoundsMin.Y > g2.BoundsMin.Y {
		g.BoundsMin.Y = g2.BoundsMin.Y
	}
	if g.BoundsMin.Z > g2.BoundsMin.Z {
		g.BoundsMin.Z = g2.BoundsMin.Z
	}
	if g.BoundsMax.X < g2.BoundsMax.X {
		g.BoundsMax.X = g2.BoundsMax.X
	}
	if g.BoundsMax.Y < g2.BoundsMax.Y {
		g.BoundsMax.Y = g2.BoundsMax.Y
	}
	if g.BoundsMax.Z < g2.BoundsMax.Z {
		g.BoundsMax.Z = g2.BoundsMax.Z
	}
	g.Code.Append(g2.Code)
}

// Join two gcodes, merging their boundaries, start/end corrdinates, and code.
// Adds Retract command in-between both codes.
func (g *Gcode) Append(g2 *Gcode, plotterConf *plotter.PlotterConfig) {
	ins := NewIns(plotterConf)
	g2StartRetracted := math64.VectorF3{X: g2.StartCoord.X, Y: g2.StartCoord.Y, Z: plotterConf.RetractHeight}
	gEndRetracted := math64.VectorF3{X: g.EndCoord.X, Y: g2.EndCoord.Y, Z: plotterConf.RetractHeight}
	if !g.EndCoord.Equal(g2.StartCoord) && !g.EndCoord.Equal(g2StartRetracted) && !gEndRetracted.Equal(g2.StartCoord) {
		ins.Retract(g)
		ins.Move(g, g2StartRetracted, plotterConf.RetractSpeed, false)
	}
	g.EndCoord = g2.EndCoord
	g.Code.Append(g2.Code)
}

func Join(gcodes []*Gcode, plotterConf *plotter.PlotterConfig) *Gcode {
	g := NewGcode()
	for _, g2 := range gcodes {
		g.Append(g2, plotterConf)
	}
	return g
}

// Create gcode, based on the given prefix and an initialization towards the
// start coordinates of the first segment
func NewGcodePrefix(plotterConf *plotter.PlotterConfig, firstSegment *Gcode) *Gcode {
	ins := NewIns(plotterConf)
	g := NewGcodeFromString(plotterConf.GcodePrefix)
	g.AppendCode("--- SVGOCODE START ---")
	target := math64.VectorF3{X: firstSegment.StartCoord.X, Y: firstSegment.StartCoord.Y, Z: plotterConf.RetractHeight}
	g.AppendCode("; Moving to start position of first segment")
	ins.Move(g, target, plotterConf.RetractSpeed, false)
	g.StartCoord = target
	g.EndCoord = target
	g.BoundsMin = target
	g.BoundsMax = target
	return g
}

// Create gcode, based on the given suffix and a retract operation beforehand
func NewGcodeSuffix(plotterConf *plotter.PlotterConfig, lastSegment *Gcode) *Gcode {
	ins := NewIns(plotterConf)
	g := NewGcode()
	target := math64.VectorF3{X: lastSegment.EndCoord.X, Y: lastSegment.EndCoord.Y, Z: plotterConf.RetractHeight}
	g.AppendCode("; SVGOCODE finished, retracting")
	ins.Move(g, target, plotterConf.RetractSpeed, false)
	g.AppendCode("; --- SVGOCODE END ---")
	g.StartCoord = target
	g.EndCoord = target
	g.BoundsMin = target
	g.BoundsMax = target
	g.AppendCode(plotterConf.GcodePrefix)
	return g
}
