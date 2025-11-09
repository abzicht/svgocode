package conv

import (
	"fmt"
	"math"

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
	pathStart := math64.VectorF2{X: 0, Y: 0}        // The very first point
	pathSegmentStart := math64.VectorF2{X: 0, Y: 0} // The first point since the last drawing began
	penDown := false

	// Header
	dCtx.ins.Retract(g)

	if len(commands) == 0 {
		llog.Panic("No commands to convert\n")
	}
	for i, cmd := range commands {
		switch cmd.Type {
		case svg.CmdMoveTo:
			if penDown {
				dCtx.ins.Retract(g)
				penDown = false
			}
			for j, p := range cmd.PathPoints {
				x, y := p.X, p.Y
				if cmd.Relative {
					x += current.X
					y += current.Y
				}
				dCtx.ins.MoveRetracted(g, dCtx.project(math64.VectorF2{X: x, Y: y}))
				current = math64.VectorF2{X: x, Y: y}
				pathSegmentStart = current
				if 0 == i && 0 == j {
					pathStart = current
				}
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
				dCtx.ins.Draw(g, dCtx.project(math64.VectorF2{X: pathSegmentStart.X, Y: pathSegmentStart.Y}))
				dCtx.ins.Retract(g)
				penDown = false
			}
			current = pathSegmentStart

		default:
			dCtx.ins.AddComment(g, fmt.Sprintf("; Unsupported command: %s\n", cmd.Type))
		}
	}
	{
		startTransformed := dCtx.project(pathStart)
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
	rx := math64.Float(math.Abs(float64(a.R.X)))
	ry := math64.Float(math.Abs(float64(a.R.Y)))
	if rx == 0 || ry == 0 {
		points = append(points, to)
		return points
	}

	// Convert rotation to radians
	//phi := float64(a.XAxis) * math.Pi / 180.0
	phi := math64.AngDeg(a.XAxis).Rad()

	// Step 1: compute (x1', y1')
	dx := (from.X - to.X) / 2.0
	dy := (from.Y - to.Y) / 2.0
	x1p := phi.Cos()*dx + phi.Sin()*dy
	y1p := -phi.Sin()*dx + phi.Cos()*dy

	// Step 2: correct radii if too small
	rx2 := rx * rx
	ry2 := ry * ry
	x1p2 := x1p * x1p
	y1p2 := y1p * y1p

	rCheck := x1p2/rx2 + y1p2/ry2
	if rCheck > 1 {
		scale := rCheck.Sqrt()
		rx *= scale
		ry *= scale
		rx2 = rx * rx
		ry2 = ry * ry
	}

	// Step 3: compute center in transformed coordinates (cx', cy')
	sign := math64.Float(-1.0)
	if a.Large != a.Sweep {
		sign = 1.0
	}

	num := rx2*ry2 - rx2*y1p2 - ry2*x1p2
	den := rx2*y1p2 + ry2*x1p2
	if den == 0 {
		den = 1e-9
	}
	cfac := sign * (num / den).Max(0).Sqrt() // math.Sqrt(math.Max(0, num/den))
	cxp := cfac * (rx * y1p / ry)
	cyp := cfac * (-ry * x1p / rx)

	// Step 4: transform center back to original coordinate system
	cx := phi.Cos()*cxp - phi.Sin()*cyp + (from.X+to.X)/2
	cy := phi.Sin()*cxp + phi.Cos()*cyp + (from.Y+to.Y)/2

	// Step 5: compute start and end angles
	v1x := (x1p - cxp) / rx
	v1y := (y1p - cyp) / ry
	v2x := (-x1p - cxp) / rx
	v2y := (-y1p - cyp) / ry

	startAngle := math64.Atan2(v1y, v1x)
	endAngle := math64.Atan2(v2y, v2x)

	// Step 6: compute delta angle
	delta := endAngle - startAngle
	if !a.Sweep && delta > 0 {
		delta -= 2 * math.Pi
	} else if a.Sweep && delta < 0 {
		delta += 2 * math.Pi
	}

	// Step 7: approximate the arc
	for i := 1; i <= int(steps); i++ {
		t := math64.AngRad(i) / math64.AngRad(steps)
		angle := startAngle + t*delta
		x := cx + rx*angle.Cos()*phi.Cos() - ry*angle.Sin()*phi.Sin()
		y := cy + rx*angle.Cos()*phi.Sin() + ry*angle.Sin()*phi.Cos()
		points = append(points, math64.VectorF2{X: x, Y: y})
	}
	return points
}
