package convs

import (
	"fmt"
	"strings"

	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/gcode"
	"github.com/abzicht/svgocode/svgocode/math64"
	"github.com/abzicht/svgocode/svgocode/plotter"
	"github.com/abzicht/svgocode/svgocode/svg"
	"github.com/abzicht/svgocode/svgocode/svg/svgtransform"
)

type directPathContext struct {
	g           *gcode.Gcode
	tMat        *svgtransform.TransformMatrix
	plotterConf *plotter.PlotterConfig
}

// This code was partly produced by AI.

func PathCommandsToGcode(commands []svg.PathCommand, transformChain svgtransform.TransformChain, g *gcode.Gcode, plotterConf *plotter.PlotterConfig) *gcode.Gcode {
	var b strings.Builder

	tMat := transformChain.ToMatrix()
	dCtx := directPathContext{g: g, tMat: tMat, plotterConf: plotterConf}
	var (
		zUp   = plotterConf.RetractHeight // pen/tool up height
		zDown = plotterConf.DrawHeight    // pen/tool down height
		steps = math64.Float(20.0)        // number of line segments to approximate curves
	)

	current := math64.VectorF2{X: 0, Y: 0}
	start := math64.VectorF2{X: 0, Y: 0}
	penDown := false

	// Header
	fmt.Fprintf(&b, "G0 Z%.3f\n", zUp)

	if len(commands) == 0 {
		llog.Panic("No commands to convert\n")
	}
	for _, cmd := range commands {
		switch cmd.Type {
		case svg.CmdMoveTo:
			if penDown {
				fmt.Fprintf(&b, "G0 Z%.3f\n", zUp)
				penDown = false
			}
			for _, p := range cmd.PathPoints {
				x, y := p.X, p.Y
				if cmd.Relative {
					x += current.X
					y += current.Y
				}
				fmt.Fprint(&b, drawPointStr(&dCtx, math64.VectorF2{X: x, Y: y}, true))
				current = math64.VectorF2{X: x, Y: y}
				start = current
			}

		case svg.CmdLineTo:
			if !penDown {
				fmt.Fprintf(&b, "G1 Z%.3f\n", zDown)
				penDown = true
			}
			for _, p := range cmd.PathPoints {
				x, y := p.X, p.Y
				if cmd.Relative {
					x += current.X
					y += current.Y
				}
				fmt.Fprint(&b, drawPointStr(&dCtx, math64.VectorF2{X: x, Y: y}, false))
				current = math64.VectorF2{X: x, Y: y}
			}

		case svg.CmdHLineTo:
			if !penDown {
				fmt.Fprintf(&b, "G1 Z%.3f\n", zDown)
				penDown = true
			}
			for _, cx := range cmd.Coordinates {
				x := cx
				if cmd.Relative {
					x += current.X
				}
				fmt.Fprint(&b, drawPointStr(&dCtx, math64.VectorF2{X: x, Y: current.Y}, false))
				current.X = x
			}

		case svg.CmdVLineTo:
			if !penDown {
				fmt.Fprintf(&b, "G1 Z%.3f\n", zDown)
				penDown = true
			}
			for _, cy := range cmd.Coordinates {
				y := cy
				if cmd.Relative {
					y += current.Y
				}
				fmt.Fprint(&b, drawPointStr(&dCtx, math64.VectorF2{X: current.X, Y: y}, false))
				current.Y = y
			}

		case svg.CmdCurveTo, svg.CmdSmoothCurveTo: // Cubic Bézier curves
			if !penDown {
				fmt.Fprintf(&b, "G1 Z%.3f\n", zDown)
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
					fmt.Fprint(&b, drawPointStr(&dCtx, math64.VectorF2{X: x, Y: y}, false))
				}
				current = p3
			}

		case svg.CmdQuadraticBezierTo, svg.CmdSmoothQuadraticBezierTo: // Quadratic Bézier
			if !penDown {
				fmt.Fprintf(&b, "G1 Z%.3f\n", zDown)
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
					fmt.Fprint(&b, drawPointStr(&dCtx, math64.VectorF2{X: x, Y: y}, false))
				}
				current = p2
			}

		case svg.CmdEllipticalArc: // Elliptical arc (approximated)
			if !penDown {
				fmt.Fprintf(&b, "G1 Z%.3f\n", zDown)
				penDown = true
			}
			for _, a := range cmd.ArcArgs {
				to := a.To
				if cmd.Relative {
					to.X += current.X
					to.Y += current.Y
				}
				for _, arcPoint := range approximateArc(current, to, a, steps) {
					fmt.Fprint(&b, drawPointStr(&dCtx, arcPoint, true))
				}
				current = to
			}

		case svg.CmdClosePath:
			if penDown {
				fmt.Fprint(&b, drawPointStr(&dCtx, math64.VectorF2{X: start.X, Y: start.Y}, false))
				fmt.Fprintf(&b, "G0 Z%.3f\n", zUp)
				penDown = false
			}
			current = start

		default:
			fmt.Fprintf(&b, "; Unsupported command: %s\n", cmd.Type)
		}
	}
	{
		startTransformed := tMat.ApplyP(start)
		currentTransformed := tMat.ApplyP(current)
		g.StartCoord = math64.VectorF3{X: startTransformed.X, Y: startTransformed.Y, Z: plotterConf.DrawHeight}
		g.EndCoord = math64.VectorF3{X: currentTransformed.X, Y: currentTransformed.Y, Z: plotterConf.DrawHeight}
	}
	g.AppendCode(b.String())
	return g
}

// Create a G0/G1 move instruction to a point that is transformed using the
// given matrix. Also updates the provided gcode's Min/Max bounds
func drawPointStr(dCtx *directPathContext, point math64.VectorF2, isDrawing bool) string {
	point = dCtx.tMat.ApplyP(point)
	point3 := math64.VectorF3{X: point.X, Y: point.Y, Z: dCtx.plotterConf.RetractHeight}
	gcmd := "G1"
	if isDrawing {
		gcmd = "G0"
		point3.Z = dCtx.plotterConf.DrawHeight
	}
	dCtx.g.BoundsMin = dCtx.g.BoundsMin.Min(point3)
	dCtx.g.BoundsMax = dCtx.g.BoundsMax.Max(point3)
	return fmt.Sprintf("%s X%.3f Y%.3f\n", gcmd, point.X, point.Y)
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
