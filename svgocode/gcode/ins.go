package gcode

import (
	"fmt"

	"github.com/abzicht/svgocode/svgocode/math64"
	"github.com/abzicht/svgocode/svgocode/plotter"
)

//Instructions

type Ins struct {
	plotterConf *plotter.PlotterConfig
}

func NewIns(plotterConf *plotter.PlotterConfig) *Ins {
	ins := new(Ins)
	ins.plotterConf = plotterConf
	return ins
}

// Add a comment line
func (ins *Ins) AddComment(g *Gcode, comment string) *Gcode {
	g.AppendCode(fmt.Sprintf("; %s", comment))
	return g
}

// Retract to preconfigured height using Move
func (ins *Ins) Retract(g *Gcode) *Gcode {
	g.AppendCode(fmt.Sprintf("; Retracting"))
	target := math64.VectorF3{X: g.EndCoord.X, Y: g.EndCoord.Y, Z: ins.plotterConf.RetractHeight}
	return ins.Move(g, target, ins.plotterConf.RetractSpeed, false)
}

// Draw at the given draw height and speed.
func (ins *Ins) Draw(g *Gcode, target math64.VectorF2) *Gcode {
	g.AppendCode(fmt.Sprintf("; Drawing to X%f Y%f", target.X, target.Y))
	return ins.Move(g, math64.VectorF3{X: target.X, Y: target.Y, Z: ins.plotterConf.DrawHeight}, ins.plotterConf.DrawSpeed, true)
}

// Move to given position with given speed
// Updates boundary and end coordinate information
func (ins *Ins) Move(g *Gcode, target math64.VectorF3, speed math64.Speed, isDrawing bool) *Gcode {
	g.EndCoord = target
	if target.X < g.BoundsMin.X {
		g.BoundsMin.X = target.X
	}
	if target.Y < g.BoundsMin.Y {
		g.BoundsMin.Y = target.Y
	}
	if target.Z < g.BoundsMin.Z {
		g.BoundsMin.Z = target.Z
	}
	if target.X > g.BoundsMax.X {
		g.BoundsMax.X = target.X
	}
	if target.Y > g.BoundsMax.Y {
		g.BoundsMax.Y = target.Y
	}
	if target.Z > g.BoundsMax.Z {
		g.BoundsMax.Z = target.Z
	}
	var gcmd string = "G0"
	if isDrawing {
		gcmd = "G1"
	}
	g.AppendCode(fmt.Sprintf("%s X%f Y%f Z%f F%f", gcmd, target.X, target.Y, target.Z, speed))
	return g
}
