package convs

import (
	"fmt"

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

// Add a line that describes the given type and id of the converted svg object
func (d *Direct) addIdComment(g *gcode.Gcode, type_, id string) {
	if len(id) == 0 {
		return
	}
	d.ins.AddComment(g, fmt.Sprintf("SVG %s (ID: %s)", type_, id))
}

func (d *Direct) Path(p *svg.Path) *gcode.Gcode {
	g := gcode.NewGcode()
	d.addIdComment(g, "Path", p.Id)
	return g
}
func (d *Direct) Line(l *svg.Line) *gcode.Gcode {
	g := gcode.NewGcode()
	d.addIdComment(g, "Line", l.Id)
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
	d.addIdComment(g, "Polygon", p.Id)
	llog.Warn("Polygon not implemented\n")
	return g
}

func (d *Direct) Polyline(p *svg.Polyline) *gcode.Gcode {
	g := gcode.NewGcode()
	d.addIdComment(g, "Polyline", p.Id)
	llog.Warn("Polyline not implemented\n")
	return g
}

func (d *Direct) Circle(c *svg.Circle) *gcode.Gcode {
	g := gcode.NewGcode()
	d.addIdComment(g, "Circle", c.Id)
	llog.Warn("Circle not implemented\n")
	return g
}

func (d *Direct) Ellipse(e *svg.Ellipse) *gcode.Gcode {
	g := gcode.NewGcode()
	d.addIdComment(g, "Ellipse", e.Id)
	llog.Warn("Ellipse not implemented\n")
	return g
}

func (d *Direct) Rect(r *svg.Rect) *gcode.Gcode {
	g := gcode.NewGcode()
	d.addIdComment(g, "Rect", r.Id)
	llog.Warn("Rect not implemented\n")
	return g
}
