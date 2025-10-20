package convs

import (
	"fmt"

	"github.com/abzicht/gogenericfunc/fun"
	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/gcode"
	"github.com/abzicht/svgocode/svgocode/math64"
	"github.com/abzicht/svgocode/svgocode/plotter"
	"github.com/abzicht/svgocode/svgocode/svg"
	"github.com/abzicht/svgocode/svgocode/svg/svgtransform"
)

// Direct conversion to gcode paths, no filling of bodies
type Direct struct {
	plotterConf *plotter.PlotterConfig
	ins         *gcode.Ins
}

func NewDirect(plotterConf *plotter.PlotterConfig) *Direct {
	d := new(Direct)
	d.plotterConf = plotterConf
	d.ins = gcode.NewIns(plotterConf)
	return d
}

// Add a line that describes the given type and id of the converted svg object
func (d *Direct) addIdComment(g *gcode.Gcode, type_, id string) {
	if len(id) == 0 {
		return
	}
	d.ins.AddComment(g, fmt.Sprintf("SVG %s (ID: %s)", type_, id))
}

func (d *Direct) Path(p *svg.Path, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode] {
	g := gcode.NewGcode()
	d.addIdComment(g, "Path", p.Id)
	if len(p.D) == 0 {
		return fun.NewNone[*gcode.Gcode]()
	}
	cmds, err := svg.ParseSVGPath(p.D)
	if err != nil {
		llog.Panicf("Failed to parse SVG path (id %s): %s. path D: '%s'\n", p.Id, err.Error(), p.D)
	}
	g = PathCommandsToGcode(cmds, g, d.plotterConf)
	return fun.NewSome[*gcode.Gcode](g)
}
func (d *Direct) Line(l *svg.Line, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode] {
	g := gcode.NewGcode()
	d.addIdComment(g, "Line", l.Id)
	g.BoundsMin.X = l.X1
	g.BoundsMax.X = l.X1
	g.BoundsMin.Y = l.Y1
	g.BoundsMax.Y = l.X1
	// Actual bounds values will be updated by move operation
	g.StartCoord = math64.VectorF3{X: l.X1, Y: l.Y1, Z: d.plotterConf.DrawHeight}
	d.ins.Move(g, g.StartCoord, d.plotterConf.DrawSpeed)
	d.ins.Draw(g, math64.VectorF2{X: l.X2, Y: l.Y2})
	return fun.NewSome[*gcode.Gcode](g)
}

func (d *Direct) Polygon(p *svg.Polygon, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode] {
	g := gcode.NewGcode()
	d.addIdComment(g, "Polygon", p.Id)
	llog.Warn("Polygon not implemented\n")
	return fun.NewNone[*gcode.Gcode]()
}

func (d *Direct) Polyline(p *svg.Polyline, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode] {
	g := gcode.NewGcode()
	d.addIdComment(g, "Polyline", p.Id)
	llog.Warn("Polyline not implemented\n")
	return fun.NewSome[*gcode.Gcode](g)
}

func (d *Direct) Circle(c *svg.Circle, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode] {
	g := gcode.NewGcode()
	d.addIdComment(g, "Circle", c.Id)
	g.BoundsMin.X = c.CX - c.R
	g.BoundsMin.Y = c.CY - c.R
	g.BoundsMax.X = c.CX + c.R
	g.BoundsMax.Y = c.CY + c.R
	// Start at the top of the circle
	g.StartCoord = math64.VectorF3{X: c.CX, Y: c.CY - c.R, Z: d.plotterConf.DrawHeight}
	d.ins.Move(g, g.StartCoord, d.plotterConf.DrawSpeed)
	d.ins.DrawCircle(g, math64.VectorF2{X: 0, Y: c.R}, c.R, true)
	return fun.NewSome[*gcode.Gcode](g)
}

func (d *Direct) Ellipse(e *svg.Ellipse, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode] {
	g := gcode.NewGcode()
	d.addIdComment(g, "Ellipse", e.Id)
	llog.Warn("Ellipse not implemented\n")
	return fun.NewSome[*gcode.Gcode](g)
}

func (d *Direct) Rect(r *svg.Rect, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode] {
	g := gcode.NewGcode()
	d.addIdComment(g, "Rect", r.Id)
	g.StartCoord = math64.VectorF3{X: r.X, Y: r.Y, Z: d.plotterConf.DrawHeight}
	d.ins.Move(g, g.StartCoord, d.plotterConf.DrawSpeed)
	d.ins.Draw(g, math64.VectorF2{X: r.X + r.Width, Y: r.Y})            // to the right
	d.ins.Draw(g, math64.VectorF2{X: r.X + r.Width, Y: r.Y + r.Height}) // down
	d.ins.Draw(g, math64.VectorF2{X: r.X, Y: r.Y + r.Height})           // to the left
	d.ins.Draw(g, math64.VectorF2{X: r.X, Y: r.Y})                      // up
	//TODO: Add RX/RY corners
	return fun.NewSome[*gcode.Gcode](g)
}
