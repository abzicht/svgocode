package conv

import (
	"fmt"
	"strings"

	"github.com/abzicht/gogenericfunc/fun"
	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/gcode"
	"github.com/abzicht/svgocode/svgocode/math64"
	"github.com/abzicht/svgocode/svgocode/svg"
	"github.com/abzicht/svgocode/svgocode/svg/svgtransform"
)

// Direct conversion to gcode paths, no filling of bodies
type Direct struct {
	conf *ConvConf
	ins  *gcode.Ins
}

// Do not use the returned *Direct before calling SetConfig on it.
func NewDirect() *Direct {
	d := new(Direct)
	return d
}

func (d *Direct) SetConfig(config *ConvConf) {
	d.conf = config
	d.ins = gcode.NewIns(d.conf.runtime)
}

// Add a line that describes the given type and id of the converted svg object
func (d *Direct) addIdComment(g *gcode.Gcode, type_ string, id svg.SvgId) {
	if len(id) == 0 {
		return
	}
	d.ins.AddComment(g, fmt.Sprintf("SVG %s (ID: %s)", type_, id))
}

func (d *Direct) PathStr(g *gcode.Gcode, pathStr string, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode] {
	if len(pathStr) == 0 {
		return fun.NewNone[*gcode.Gcode]()
	}
	cmds, err := svg.ParseSVGPath(pathStr)
	if err != nil {
		llog.Panicf("Failed to parse SVG path: %s. Path string: '%s'\n", err.Error(), pathStr)
	}
	g = PathCommandsToGcode(cmds, transformChain, g, d.conf.runtime, d.ins)
	return fun.NewSome[*gcode.Gcode](g)
}

func (d *Direct) Path(p *svg.Path, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode] {
	g := gcode.NewGcode()
	if len(p.D) == 0 {
		return fun.NewNone[*gcode.Gcode]()
	}
	d.addIdComment(g, "Path", p.Id)
	return d.PathStr(g, p.D, transformChain)
}

func (d *Direct) Line(l *svg.Line, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode] {
	g := gcode.NewGcode()
	tMatrix := transformChain.ToMatrix()
	p1 := tMatrix.ApplyP(math64.VectorF2{X: l.X1, Y: l.Y1})
	p2 := tMatrix.ApplyP(math64.VectorF2{X: l.X2, Y: l.Y2})
	d.addIdComment(g, "Line", l.Id)
	g.BoundsMin.X = p1.X
	g.BoundsMin.Y = p1.Y
	g.BoundsMax.X = p1.X
	g.BoundsMax.Y = p1.Y
	// Actual bounds values will be updated by move operation
	g.StartCoord = math64.VectorF3{X: p1.X, Y: p1.Y, Z: d.conf.runtime.Plotter.DrawHeight}
	d.ins.Move(g, g.StartCoord, d.conf.runtime.Plotter.DrawSpeed)
	d.ins.Draw(g, math64.VectorF2{X: p2.X, Y: p2.Y})
	return fun.NewSome[*gcode.Gcode](g)
}

func (d *Direct) Polygon(p *svg.Polygon, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode] {
	g := gcode.NewGcode()
	d.addIdComment(g, "Polygon", p.Id)
	return d.PathStr(g, svg.PointsToPathStr(p.Points(), true), transformChain)
}

func (d *Direct) Polyline(p *svg.Polyline, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode] {
	g := gcode.NewGcode()
	d.addIdComment(g, "Polyline", p.Id)
	return d.PathStr(g, svg.PointsToPathStr(p.Points(), false), transformChain)
}

func (d *Direct) Circle(c *svg.Circle, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode] {
	g := gcode.NewGcode()
	d.addIdComment(g, "Circle", c.Id)
	if len(transformChain) != 0 {
		// We transform the circle, better let the path converter figure that one out
		pathStr := fmt.Sprintf(
			"M %g %g A %g %g 0 1 0 %g %g A %g %g 0 1 0 %g %g Z",
			c.CX-c.R, c.CY, c.R, c.R, c.CX+c.R, c.CY, c.R, c.R, c.CX-c.R, c.CY)
		return d.PathStr(g, pathStr, transformChain)
	}
	// No transformations, direct gcode call to draw a circle.
	g.BoundsMin.X = c.CX - c.R
	g.BoundsMin.Y = c.CY - c.R
	g.BoundsMax.X = c.CX + c.R
	g.BoundsMax.Y = c.CY + c.R
	// Start at the top of the circle
	g.StartCoord = math64.VectorF3{X: c.CX, Y: c.CY - c.R, Z: d.conf.runtime.Plotter.DrawHeight}
	d.ins.Move(g, g.StartCoord, d.conf.runtime.Plotter.DrawSpeed)
	d.ins.DrawCircle(g, math64.VectorF2{X: 0, Y: c.R}, c.R, true)
	return fun.NewSome[*gcode.Gcode](g)
}

func (d *Direct) Ellipse(e *svg.Ellipse, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode] {
	g := gcode.NewGcode()
	d.addIdComment(g, "Ellipse", e.Id)
	// No need to use math here, just convert the ellipse to a SVG path and let the
	// path parser do the work.
	pathStr := fmt.Sprintf(
		"M %g %g A %g %g 0 1 0 %g %g A %g %g 0 1 0 %g %g Z",
		e.CX-e.RX, e.CY, e.RX, e.RY, e.CX+e.RX, e.CY, e.RX, e.RY, e.CX-e.RX, e.CY)
	return d.PathStr(g, pathStr, transformChain)
}

func (d *Direct) Rect(r *svg.Rect, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode] {
	if r.Width <= 0 || r.Height <= 0 {
		return fun.NewNone[*gcode.Gcode]()
	}

	// Rounded corners are complicated, let the path parser do the work.

	g := gcode.NewGcode()
	d.addIdComment(g, "Rect", r.Id)
	// Apply SVG rounding rules
	rx := r.RX
	ry := r.RY
	if rx < 0 {
		rx = 0
	}
	if ry < 0 {
		ry = 0
	}
	if rx == 0 && ry > 0 {
		rx = ry
	}
	if ry == 0 && rx > 0 {
		ry = rx
	}
	if rx > r.Width/2 {
		rx = r.Width / 2
	}
	if ry > r.Height/2 {
		ry = r.Height / 2
	}

	var pathStr string
	// No rounded corners
	if rx == 0 && ry == 0 {
		pathStr = fmt.Sprintf("M %g,%g H %g V %g H %g Z",
			r.X, r.Y,
			r.X+r.Width,
			r.Y+r.Height,
			r.X,
		)
	} else {
		// With rounded corners (arcs)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("M %g,%g ", r.X+rx, r.Y))
		sb.WriteString(fmt.Sprintf("H %g ", r.X+r.Width-rx))
		sb.WriteString(fmt.Sprintf("A %g,%g 0 0 1 %g,%g ", rx, ry, r.X+r.Width, r.Y+ry))
		sb.WriteString(fmt.Sprintf("V %g ", r.Y+r.Height-ry))
		sb.WriteString(fmt.Sprintf("A %g,%g 0 0 1 %g,%g ", rx, ry, r.X+r.Width-rx, r.Y+r.Height))
		sb.WriteString(fmt.Sprintf("H %g ", r.X+rx))
		sb.WriteString(fmt.Sprintf("A %g,%g 0 0 1 %g,%g ", rx, ry, r.X, r.Y+r.Height-ry))
		sb.WriteString(fmt.Sprintf("V %g ", r.Y+ry))
		sb.WriteString(fmt.Sprintf("A %g,%g 0 0 1 %g,%g Z", rx, ry, r.X+rx, r.Y))
		pathStr = sb.String()
	}
	return d.PathStr(g, pathStr, transformChain)
}
