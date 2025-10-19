package path

import (
	"fmt"
	"math"
	"strings"

	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/gcode"
	"github.com/abzicht/svgocode/svgocode/math64"
	"github.com/abzicht/svgocode/svgocode/plotter"
)

// This code was also produced by AI. Explain that to a Victorian peasant.
// Not fully revised yet, but will do.

func PathCommandsToGcode(commands []Command, g *gcode.Gcode, plotterConf *plotter.PlotterConfig) *gcode.Gcode {
	//TODO add gcode bounds
	var b strings.Builder

	var (
		zUp   = plotterConf.RetractHeight // pen/tool up height
		zDown = plotterConf.DrawHeight    // pen/tool down height
		steps = 20.0                      // number of line segments to approximate curves
	)

	current := Point{0, 0}
	start := Point{0, 0}
	penDown := false

	// Header
	fmt.Fprintf(&b, "G0 Z%.3f\n", zUp)

	if len(commands) == 0 {
		llog.Panic("No commands to convert\n")
	}
	for _, cmd := range commands {
		switch cmd.Type {
		case CmdMoveTo:
			if penDown {
				fmt.Fprintf(&b, "G0 Z%.3f\n", zUp)
				penDown = false
			}
			for _, p := range cmd.Points {
				x, y := p.X, p.Y
				if cmd.Relative {
					x += current.X
					y += current.Y
				}
				fmt.Fprintf(&b, "G0 X%.3f Y%.3f\n", x, y)
				current = Point{x, y}
				start = current
			}

		case CmdLineTo:
			if !penDown {
				fmt.Fprintf(&b, "G1 Z%.3f\n", zDown)
				penDown = true
			}
			for _, p := range cmd.Points {
				x, y := p.X, p.Y
				if cmd.Relative {
					x += current.X
					y += current.Y
				}
				fmt.Fprintf(&b, "G1 X%.3f Y%.3f\n", x, y)
				current = Point{x, y}
			}

		case CmdHLineTo:
			if !penDown {
				fmt.Fprintf(&b, "G1 Z%.3f\n", zDown)
				penDown = true
			}
			for _, cx := range cmd.Coordinates {
				x := float64(cx)
				if cmd.Relative {
					x += current.X
				}
				fmt.Fprintf(&b, "G1 X%.3f Y%.3f\n", x, current.Y)
				current.X = x
			}

		case CmdVLineTo:
			if !penDown {
				fmt.Fprintf(&b, "G1 Z%.3f\n", zDown)
				penDown = true
			}
			for _, cy := range cmd.Coordinates {
				y := float64(cy)
				if cmd.Relative {
					y += current.Y
				}
				fmt.Fprintf(&b, "G1 X%.3f Y%.3f\n", current.X, y)
				current.Y = y
			}

		case CmdCurveTo, CmdSmoothCurveTo: // Cubic Bézier curves
			if !penDown {
				fmt.Fprintf(&b, "G1 Z%.3f\n", zDown)
				penDown = true
			}
			for i := 0; i+2 < len(cmd.Points); i += 3 {
				p1 := cmd.Points[i]
				p2 := cmd.Points[i+1]
				p3 := cmd.Points[i+2]
				if cmd.Relative {
					p1.X += current.X
					p1.Y += current.Y
					p2.X += current.X
					p2.Y += current.Y
					p3.X += current.X
					p3.Y += current.Y
				}
				for t := 0.0; t <= 1.0; t += 1.0 / steps {
					x := cubicBezier(t, current.X, p1.X, p2.X, p3.X)
					y := cubicBezier(t, current.Y, p1.Y, p2.Y, p3.Y)
					fmt.Fprintf(&b, "G1 X%.3f Y%.3f\n", x, y)
				}
				current = p3
			}

		case CmdQuadraticBezierTo, CmdSmoothQuadraticBezierTo: // Quadratic Bézier
			if !penDown {
				fmt.Fprintf(&b, "G1 Z%.3f\n", zDown)
				penDown = true
			}
			for i := 0; i+1 < len(cmd.Points); i += 2 {
				p1 := cmd.Points[i]
				p2 := cmd.Points[i+1]
				if cmd.Relative {
					p1.X += current.X
					p1.Y += current.Y
					p2.X += current.X
					p2.Y += current.Y
				}
				for t := 0.0; t <= 1.0; t += 1.0 / steps {
					x := quadraticBezier(t, current.X, p1.X, p2.X)
					y := quadraticBezier(t, current.Y, p1.Y, p2.Y)
					fmt.Fprintf(&b, "G1 X%.3f Y%.3f\n", x, y)
				}
				current = p2
			}

		case CmdEllipticalArc: // Elliptical arc (approximated)
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
				approximateArc(&b, current, to, a, steps)
				current = to
			}

		case CmdClosePath:
			if penDown {
				fmt.Fprintf(&b, "G1 X%.3f Y%.3f\n", start.X, start.Y)
				penDown = false
			}
			current = start

		default:
			fmt.Fprintf(&b, "; Unsupported command: %s\n", cmd.Type)
		}
	}
	{
		g.StartCoord = math64.VectorF3{X: math64.Float(start.X), Y: math64.Float(start.Y), Z: plotterConf.DrawHeight}
		g.EndCoord = math64.VectorF3{X: math64.Float(current.X), Y: math64.Float(current.Y), Z: plotterConf.DrawHeight}
	}
	g.AppendCode(b.String())
	return g
}

// Cubic Bézier interpolation
func cubicBezier(t, p0, p1, p2, p3 float64) float64 {
	u := 1 - t
	return math.Pow(u, 3)*p0 +
		3*math.Pow(u, 2)*t*p1 +
		3*u*math.Pow(t, 2)*p2 +
		math.Pow(t, 3)*p3
}

// Quadratic Bézier interpolation
func quadraticBezier(t, p0, p1, p2 float64) float64 {
	u := 1 - t
	return u*u*p0 + 2*u*t*p1 + t*t*p2
}

// Approximate elliptical arc with line segments
func approximateArc(b *strings.Builder, from, to Point, a EllipticalArcArg, steps float64) {
	// For simplicity: approximate as a simple ellipse arc from 'from' to 'to'
	rx := float64(a.Rx)
	ry := float64(a.Ry)
	startAngle := 0.0
	endAngle := math.Pi / 2 // placeholder; full arc solving is complex

	for i := 0; i <= int(steps); i++ {
		t := float64(i) / float64(steps)
		angle := startAngle + t*(endAngle-startAngle)
		x := from.X + rx*math.Cos(angle)
		y := from.Y + ry*math.Sin(angle)
		fmt.Fprintf(b, "G1 X%.3f Y%.3f\n", x, y)
	}
	fmt.Fprintf(b, "G1 X%.3f Y%.3f\n", to.X, to.Y)
}
