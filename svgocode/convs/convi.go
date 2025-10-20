package convs

import (
	"github.com/abzicht/gogenericfunc/fun"
	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/gcode"
	"github.com/abzicht/svgocode/svgocode/svg"
	"github.com/abzicht/svgocode/svgocode/svg/svgtransform"
)

type ConverterI interface {
	Path(p *svg.Path, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode]
	Line(l *svg.Line, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode]
	Rect(c *svg.Rect, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode]
	Circle(c *svg.Circle, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode]
	Ellipse(c *svg.Ellipse, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode]
	Polygon(p *svg.Polygon, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode]
	Polyline(p *svg.Polyline, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode]
}

// Convert using the converter, based on the element's type
func SVGConvert(s svg.SVGShapeElement, transformChain svgtransform.TransformChain, conv ConverterI) fun.Option[*gcode.Gcode] {
	switch s.(type) {
	case *svg.Path:
		return conv.Path(s.(*svg.Path), transformChain)
	case *svg.Line:
		return conv.Line(s.(*svg.Line), transformChain)
	case *svg.Rect:
		return conv.Rect(s.(*svg.Rect), transformChain)
	case *svg.Circle:
		return conv.Circle(s.(*svg.Circle), transformChain)
	case *svg.Ellipse:
		return conv.Ellipse(s.(*svg.Ellipse), transformChain)
	case *svg.Polygon:
		return conv.Polygon(s.(*svg.Polygon), transformChain)
	case *svg.Polyline:
		return conv.Polyline(s.(*svg.Polyline), transformChain)
	default:
		llog.Panicf("Unknown SVG object received, cannot convert to gcode. Type: %T\n", s)
		return nil
	}
}
