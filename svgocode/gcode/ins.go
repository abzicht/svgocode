package gcode

import (
	"fmt"

	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/conf"
	"github.com/abzicht/svgocode/svgocode/math64"
)

//Instructions

type Ins struct {
	runtime *conf.RuntimeConfig
}

func NewIns(runtConf *conf.RuntimeConfig) *Ins {
	ins := new(Ins)
	ins.runtime = runtConf
	return ins
}

func (ins *Ins) convUnitF2(v math64.VectorF2) math64.VectorF2 {
	x := math64.LengthConvert(v.X, ins.runtime.SvgUnit, ins.runtime.PlotterUnit)
	y := math64.LengthConvert(v.Y, ins.runtime.SvgUnit, ins.runtime.PlotterUnit)
	return math64.VectorF2{X: x, Y: y}
}

func (ins *Ins) convUnitF3(v math64.VectorF3) math64.VectorF3 {
	x := math64.LengthConvert(v.X, ins.runtime.SvgUnit, ins.runtime.PlotterUnit)
	y := math64.LengthConvert(v.Y, ins.runtime.SvgUnit, ins.runtime.PlotterUnit)
	z := math64.LengthConvert(v.Z, ins.runtime.SvgUnit, ins.runtime.PlotterUnit)
	return math64.VectorF3{X: x, Y: y, Z: z}
}

// Add a comment line
func (ins *Ins) AddComment(g *Gcode, comment string) *Gcode {
	g.AppendCode(fmt.Sprintf("; %s", comment))
	return g
}

// Set unit
func (ins *Ins) SetUnit(g *Gcode, u math64.UnitLength) *Gcode {
	var gcmd string = "G21"
	switch u {
	case math64.UnitMM:
		gcmd = "G21 ; Setting unit (mm)"
	case math64.UnitIN:
		gcmd = "G20 ; Setting unit (in)"
	default:
		llog.Panicf("Unsupported unit type (%s)", u)
	}
	g.AppendCode(gcmd)
	return g
}

// Set extrusion speed for a given mode (G0/G1)
func (ins *Ins) SetExtrusion(g *Gcode, extSpeed math64.Speed, forDrawing bool) *Gcode {
	var gcmd string = "G0"
	if forDrawing {
		gcmd = "G1"
	}
	g.AppendCode(fmt.Sprintf("%s E%f ; Setting Extrusion", gcmd, extSpeed))
	return g
}

// Set the default speed for a given move mode (G0/G1)
func (ins *Ins) SetSpeed(g *Gcode, speed math64.Speed, forDrawing bool) *Gcode {
	var gcmd string = "G0"
	if forDrawing {
		gcmd = "G1"
	}
	g.AppendCode(fmt.Sprintf("%s F%f ; Setting Speed", gcmd, speed))
	return g
}

// Retract to preconfigured height using Move
func (ins *Ins) Retract(g *Gcode) *Gcode {
	//g.AppendCode(fmt.Sprintf("; Retracting"))
	//target := math64.VectorF3{X: g.EndCoord.X, Y: g.EndCoord.Y, Z: ins.runtime.Plotter.RetractHeight}
	//return ins.move(g, target, ins.runtime.Plotter.RetractSpeed, false)
	g.EndCoord.Z = ins.runtime.Plotter.RetractHeight
	g.BoundsMin = g.BoundsMin.Min(g.EndCoord)
	g.BoundsMax = g.BoundsMax.Max(g.EndCoord)
	g.AppendCode(fmt.Sprintf("G0 Z%f F%f ; Retracting", math64.LengthConvert(ins.runtime.Plotter.RetractHeight, math64.UnitMM, ins.runtime.PlotterUnit), math64.SpeedConvert(ins.runtime.Plotter.RetractSpeed, math64.UnitMM, ins.runtime.PlotterUnit)))
	return g
}

// Lower pen to draw height
func (ins *Ins) DrawPos(g *Gcode) *Gcode {
	g.EndCoord.Z = ins.runtime.Plotter.DrawHeight
	g.BoundsMin = g.BoundsMin.Min(g.EndCoord)
	g.BoundsMax = g.BoundsMax.Max(g.EndCoord)
	g.AppendCode(fmt.Sprintf("G1 Z%f F%f ; Lowering", math64.LengthConvert(ins.runtime.Plotter.DrawHeight, math64.UnitMM, ins.runtime.PlotterUnit), math64.SpeedConvert(ins.runtime.Plotter.DrawSpeed, math64.UnitMM, ins.runtime.PlotterUnit)))
	return g
}

// MoveRetracted at retract height to given position.
func (ins *Ins) MoveRetracted(g *Gcode, target math64.VectorF2) *Gcode {
	g.EndCoord.X = target.X
	g.EndCoord.Y = target.Y
	g.EndCoord.Z = ins.runtime.Plotter.RetractHeight
	return ins.move(g, g.EndCoord, ins.runtime.Plotter.RetractSpeed, false)
}

// Draw at the given draw height and speed.
func (ins *Ins) Draw(g *Gcode, target math64.VectorF2) *Gcode {
	return ins.move(g, math64.VectorF3{X: target.X, Y: target.Y, Z: ins.runtime.Plotter.DrawHeight}, ins.runtime.Plotter.DrawSpeed, true)
}

// Move to given position with given speed. Not configured for drawing
func (ins *Ins) Move(g *Gcode, target math64.VectorF3, speed math64.Speed) *Gcode {
	return ins.move(g, target, speed, false)
}

// Move to given position with given speed. Position and speed will be
// converted to plotter's units
// Updates boundary and end coordinate information
func (ins *Ins) move(g *Gcode, target math64.VectorF3, speed math64.Speed, isDrawing bool) *Gcode {
	g.EndCoord = target
	if g.Code.NumLines() == 0 {
		g.BoundsMin = target
		g.BoundsMax = target
	} else {
		g.BoundsMin = g.BoundsMin.Min(target)
		g.BoundsMax = g.BoundsMax.Max(target)
	}
	targetPlUnit := ins.convUnitF3(target)
	speedPlUnit := math64.SpeedConvert(speed, math64.UnitMM, ins.runtime.PlotterUnit)
	if isDrawing {
		g.AppendCode(fmt.Sprintf("G1 X%f Y%f Z%f F%f ; Drawing", targetPlUnit.X, targetPlUnit.Y, targetPlUnit.Z, speedPlUnit))
	} else {
		g.AppendCode(fmt.Sprintf("G0 X%f Y%f Z%f F%f ; Moving", targetPlUnit.X, targetPlUnit.Y, targetPlUnit.Z, speedPlUnit))
	}
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
	centerOffsetPlUnit := ins.convUnitF2(centerOffset)
	g.AppendCode(fmt.Sprintf("%s I%f J%f F%f", gcmd, centerOffsetPlUnit.X, centerOffsetPlUnit.Y, math64.SpeedConvert(ins.runtime.Plotter.DrawSpeed, math64.UnitMM, ins.runtime.PlotterUnit)))
	return g
}
