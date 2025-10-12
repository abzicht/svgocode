package convs

import (
	"github.com/abzicht/gogenericfunc/fun"
	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/gcode"
	"github.com/abzicht/svgocode/svgocode/svg"
)

type ConverterI interface {
	Path(p *svg.Path) fun.Option[*gcode.Gcode]
	Line(l *svg.Line) fun.Option[*gcode.Gcode]
	Rect(c *svg.Rect) fun.Option[*gcode.Gcode]
	Circle(c *svg.Circle) fun.Option[*gcode.Gcode]
	Ellipse(c *svg.Ellipse) fun.Option[*gcode.Gcode]
	Polygon(p *svg.Polygon) fun.Option[*gcode.Gcode]
	Polyline(p *svg.Polyline) fun.Option[*gcode.Gcode]
}

// Convert using the converter, based on the element's type
func SVGConvert(s svg.SVGShapeElement, conv ConverterI) fun.Option[*gcode.Gcode] {
	switch s.(type) {
	case *svg.Path:
		return conv.Path(s.(*svg.Path))
	case *svg.Line:
		return conv.Line(s.(*svg.Line))
	case *svg.Rect:
		return conv.Rect(s.(*svg.Rect))
	case *svg.Circle:
		return conv.Circle(s.(*svg.Circle))
	case *svg.Ellipse:
		return conv.Ellipse(s.(*svg.Ellipse))
	case *svg.Polygon:
		return conv.Polygon(s.(*svg.Polygon))
	case *svg.Polyline:
		return conv.Polyline(s.(*svg.Polyline))
	default:
		llog.Panicf("Unknown SVG object received, cannot convert to gcode. Type: %T\n", s)
		return nil
	}
}
