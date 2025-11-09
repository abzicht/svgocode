package conv

import (
	"fmt"

	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/conf"
	"github.com/abzicht/svgocode/svgocode/gcode"
	"github.com/abzicht/svgocode/svgocode/math64"
	"github.com/abzicht/svgocode/svgocode/svg"
	"github.com/abzicht/svgocode/svgocode/svg/svgtransform"
)

type directPathContext struct {
	g       *gcode.Gcode
	tMat    *svgtransform.TransformMatrix
	runtime *conf.RuntimeConfig
	ins     *gcode.Ins
}

// Project point into gcode space by applying svg transformations and unit
// conversion
func (d *directPathContext) project(p math64.VectorF2) math64.VectorF2 {
	p2 := d.tMat.ApplyP(p)
	x := math64.LengthConvert(p2.X, d.runtime.SvgUnit, d.runtime.PlotterUnit)
	y := math64.LengthConvert(p2.Y, d.runtime.SvgUnit, d.runtime.PlotterUnit)
	return math64.VectorF2{X: x, Y: y}
}

// Convert a slice of svg path commands into gcode, applying a transform chain
// and svg-to-plotter-unit conversion to all instructions.
func PathCommandsToGcode(commands []svg.PathCommand, transformChain svgtransform.TransformChain, g *gcode.Gcode, runtConf *conf.RuntimeConfig, ins *gcode.Ins) *gcode.Gcode {

	tMat := transformChain.ToMatrix()
	dCtx := directPathContext{g: g, tMat: tMat, runtime: runtConf, ins: ins}
	var (
		steps = math64.Float(20.0) // number of line segments to approximate curves
	)

	current := math64.VectorF2{X: 0, Y: 0}
	start := math64.VectorF2{X: 0, Y: 0}
	penDown := false

	// Header
	dCtx.ins.Retract(g)

	if len(commands) == 0 {
		llog.Panic("No commands to convert\n")
	}
	for _, cmd := range commands {
		switch cmd.Type {
		case svg.CmdMoveTo:
			if penDown {
				dCtx.ins.Retract(g)
				penDown = false
			}
			for _, p := range cmd.PathPoints {
				x, y := p.X, p.Y
				if cmd.Relative {
					x += current.X
					y += current.Y
				}
				dCtx.ins.MoveRetracted(g, dCtx.project(math64.VectorF2{X: x, Y: y}))
				current = math64.VectorF2{X: x, Y: y}
				start = current
			}

		case svg.CmdLineTo:
			if !penDown {
				dCtx.ins.DrawPos(g)
				penDown = true
			}
			for _, p := range cmd.PathPoints {
				x, y := p.X, p.Y
				if cmd.Relative {
					x += current.X
					y += current.Y
				}
				dCtx.ins.Draw(g, dCtx.project(math64.VectorF2{X: x, Y: y}))
				current = math64.VectorF2{X: x, Y: y}
			}

		case svg.CmdHLineTo:
			if !penDown {
				dCtx.ins.DrawPos(g)
				penDown = true
			}
			for _, cx := range cmd.Coordinates {
				x := cx
				if cmd.Relative {
					x += current.X
				}
				dCtx.ins.Draw(g, dCtx.project(math64.VectorF2{X: x, Y: current.Y}))
				current.X = x
			}

		case svg.CmdVLineTo:
			if !penDown {
				dCtx.ins.DrawPos(g)
				penDown = true
			}
			for _, cy := range cmd.Coordinates {
				y := cy
				if cmd.Relative {
					y += current.Y
				}
				dCtx.ins.Draw(g, dCtx.project(math64.VectorF2{X: current.X, Y: y}))
				current.Y = y
			}

		case svg.CmdCurveTo, svg.CmdSmoothCurveTo: // Cubic Bézier curves
			if !penDown {
				dCtx.ins.DrawPos(g)
				penDown = true
			}
			for i := 0; i+2 < len(cmd.PathPoints); i += 3 {
				p1 := cmd.PathPoints[i]
				p2 := cmd.PathPoints[i+1]
				p3 := cmd.PathPoints[i+2]
				if cmd.Relative {
					p1 = p1.Add(current)
					p2 = p2.Add(current)
					p3 = p3.Add(current)
				}
				for t := math64.Float(0.0); t <= 1.0; t += 1.0 / steps {
					x := cubicBezier(math64.Float(t), current.X, p1.X, p2.X, p3.X)
					y := cubicBezier(math64.Float(t), current.Y, p1.Y, p2.Y, p3.Y)
					dCtx.ins.Draw(g, dCtx.project(math64.VectorF2{X: x, Y: y}))
				}
				current = p3
			}

		case svg.CmdQuadraticBezierTo, svg.CmdSmoothQuadraticBezierTo: // Quadratic Bézier
			if !penDown {
				dCtx.ins.DrawPos(g)
				penDown = true
			}
			for i := 0; i+1 < len(cmd.PathPoints); i += 2 {
				p1 := cmd.PathPoints[i]
				p2 := cmd.PathPoints[i+1]
				if cmd.Relative {
					p1.X += current.X
					p1.Y += current.Y
					p2.X += current.X
					p2.Y += current.Y
				}
				for t := math64.Float(0.0); t <= 1.0; t += 1.0 / steps {
					x := quadraticBezier(math64.Float(t), current.X, p1.X, p2.X)
					y := quadraticBezier(math64.Float(t), current.Y, p1.Y, p2.Y)
					dCtx.ins.Draw(g, dCtx.project(math64.VectorF2{X: x, Y: y}))
				}
				current = p2
			}

		case svg.CmdEllipticalArc: // Elliptical arc (approximated)
			if !penDown {
				dCtx.ins.DrawPos(g)
				penDown = true
			}
			for _, a := range cmd.ArcArgs {
				to := a.To
				if cmd.Relative {
					to.X += current.X
					to.Y += current.Y
				}
				for _, arcPoint := range approximateArc(current, to, a, steps) {
					dCtx.ins.Draw(g, dCtx.project(arcPoint))
				}
				current = to
			}

		case svg.CmdClosePath:
			if penDown {
				dCtx.ins.Draw(g, dCtx.project(math64.VectorF2{X: start.X, Y: start.Y}))
				dCtx.ins.Retract(g)
				penDown = false
			}
			current = start

		default:
			dCtx.ins.AddComment(g, fmt.Sprintf("; Unsupported command: %s\n", cmd.Type))
		}
	}
	{
		startTransformed := dCtx.project(start)
		currentTransformed := dCtx.project(current)
		g.StartCoord = math64.VectorF3{X: startTransformed.X, Y: startTransformed.Y, Z: runtConf.Plotter.DrawHeight}
		g.EndCoord = math64.VectorF3{X: currentTransformed.X, Y: currentTransformed.Y, Z: runtConf.Plotter.DrawHeight}
	}
	return g
}

// Cubic Bézier interpolation
func cubicBezier(t, p0, p1, p2, p3 math64.Float) math64.Float {
	u := 1 - t
	return u.Pow(3)*p0 +
		3*u.Pow(2)*t*p1 +
		3*u*t.Pow(2)*p2 +
		t.Pow(3)*p3
}

// Quadratic Bézier interpolation
func quadraticBezier(t, p0, p1, p2 math64.Float) math64.Float {
	u := 1 - t
	return u*u*p0 + 2*u*t*p1 + t*t*p2
}

// Approximate elliptical arc with line segments
func approximateArc(from, to math64.VectorF2, a svg.EllipticalArcArg, steps math64.Float) []math64.VectorF2 {
	var points []math64.VectorF2
	// For simplicity: approximate as a simple ellipse arc from 'from' to 'to'
	var startAngle math64.AngRad = 0.0
	endAngle := math64.AngDeg(1).Rad() // placeholder; full arc solving is complex

	for i := 0; i <= int(steps); i++ {
		t := math64.AngRad(i) / math64.AngRad(steps)
		angle := startAngle + t*(endAngle-startAngle)
		x := from.X + a.R.X*angle.Cos()
		y := from.Y + a.R.Y*angle.Sin()
		points = append(points, math64.VectorF2{X: x, Y: y})
	}
	points = append(points, to)
	return points
}
