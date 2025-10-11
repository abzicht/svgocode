package convs

import (
	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/gcode"
	"github.com/abzicht/svgocode/svgocode/math64"
	"github.com/abzicht/svgocode/svgocode/plotter"
	"github.com/abzicht/svgocode/svgocode/svg"
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

func (d *Direct) Path(p *svg.Path) *gcode.Gcode {
	g := gcode.NewGcode()
	return g
}
func (d *Direct) Line(l *svg.Line) *gcode.Gcode {
	g := gcode.NewGcode()
	g.BoundsMin.X = l.X1
	g.BoundsMax.X = l.X1
	g.BoundsMin.Y = l.Y1
	g.BoundsMax.Y = l.X1
	// Actual bounds values will be updated by move operation
	g.StartCoord = math64.VectorF3{X: l.X1, Y: l.Y1, Z: d.plotterConf.DrawHeight}
	d.ins.Move(g, g.StartCoord, d.plotterConf.DrawSpeed, false)
	d.ins.Draw(g, math64.VectorF2{X: l.X2, Y: l.Y2})
	return g
}

func (d *Direct) Polygon(p *svg.Polygon) *gcode.Gcode {
	g := gcode.NewGcode()
	llog.Warn("Polygon not implemented\n")
	return g
}

func (d *Direct) Polyline(p *svg.Polyline) *gcode.Gcode {
	g := gcode.NewGcode()
	llog.Warn("Polyline not implemented\n")
	return g
}

func (d *Direct) Circle(c *svg.Circle) *gcode.Gcode {
	g := gcode.NewGcode()
	llog.Warn("Circle not implemented\n")
	return g
}

func (d *Direct) Ellipse(c *svg.Ellipse) *gcode.Gcode {
	g := gcode.NewGcode()
	llog.Warn("Ellipse not implemented\n")
	return g
}

func (d *Direct) Rect(c *svg.Rect) *gcode.Gcode {
	g := gcode.NewGcode()
	llog.Warn("Rect not implemented\n")
	return g
}
