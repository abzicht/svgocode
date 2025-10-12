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
	return ins.move(g, target, ins.plotterConf.RetractSpeed, false)
}

// Draw at the given draw height and speed.
func (ins *Ins) Draw(g *Gcode, target math64.VectorF2) *Gcode {
	g.AppendCode(fmt.Sprintf("; Drawing to X%f Y%f", target.X, target.Y))
	return ins.move(g, math64.VectorF3{X: target.X, Y: target.Y, Z: ins.plotterConf.DrawHeight}, ins.plotterConf.DrawSpeed, true)
}

// Move to given position with given speed. Not configured for drawing
func (ins *Ins) Move(g *Gcode, target math64.VectorF3, speed math64.Speed) *Gcode {
	return ins.move(g, target, speed, false)
}

// Move to given position with given speed
// Updates boundary and end coordinate information
func (ins *Ins) move(g *Gcode, target math64.VectorF3, speed math64.Speed, isDrawing bool) *Gcode {
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

// Draw a circle with the center being measured by an offset to the current
// position
func (ins *Ins) DrawCircle(g *Gcode, centerOffset math64.VectorF2, radius math64.Float, clockwise bool) *Gcode {
	g.AppendCode(fmt.Sprintf("; Drawing circle around offset X%f Y%f with radius $%f from current position", centerOffset.X, centerOffset.Y, radius))
	if centerOffset.X-radius < g.BoundsMin.X {
		g.BoundsMin.X = centerOffset.X - radius
	}
	if centerOffset.Y-radius < g.BoundsMin.Y {
		g.BoundsMin.Y = centerOffset.Y - radius
	}
	if centerOffset.X+radius > g.BoundsMax.X {
		g.BoundsMax.X = centerOffset.X + radius
	}
	if centerOffset.Y+radius > g.BoundsMax.Y {
		g.BoundsMax.Y = centerOffset.Y + radius
	}
	var gcmd string = "G2"
	if !clockwise {
		gcmd = "G3"
	}
	g.AppendCode(fmt.Sprintf("%s I%f J%f F%f", gcmd, centerOffset.X, centerOffset.Y, ins.plotterConf.DrawSpeed))
	return g
}
