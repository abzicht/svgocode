package convs

import (
	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/gcode"
	"github.com/abzicht/svgocode/svgocode/svg"
)

type ConverterI interface {
	Path(p *svg.Path) *gcode.Gcode
	Line(l *svg.Line) *gcode.Gcode
	Rect(c *svg.Rect) *gcode.Gcode
	Circle(c *svg.Circle) *gcode.Gcode
	Ellipse(c *svg.Ellipse) *gcode.Gcode
	Polygon(p *svg.Polygon) *gcode.Gcode
	Polyline(p *svg.Polyline) *gcode.Gcode
}

// Convert using the converter, based on the element's type
func SVGConvert(s svg.SVGShapeElement, conv ConverterI) *gcode.Gcode {
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
