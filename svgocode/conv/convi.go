package conv

import (
	"github.com/abzicht/gogenericfunc/fun"
	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/conf"
	"github.com/abzicht/svgocode/svgocode/gcode"
	"github.com/abzicht/svgocode/svgocode/svg"
	"github.com/abzicht/svgocode/svgocode/svg/svgtransform"
)

type ConverterI interface {
	SetConfig(*ConvConf)
	Path(p *svg.Path, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode]
	Line(l *svg.Line, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode]
	Rect(c *svg.Rect, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode]
	Circle(c *svg.Circle, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode]
	Ellipse(c *svg.Ellipse, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode]
	Polygon(p *svg.Polygon, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode]
	Polyline(p *svg.Polyline, transformChain svgtransform.TransformChain) fun.Option[*gcode.Gcode]
}

// ConvConf: The greatest type name so far
type ConvConf struct {
	plotter *conf.PlotterConfig
	runtime *conf.RuntimeConfig
}

func NewConvConf(plotterConf *conf.PlotterConfig, runtConf *conf.RuntimeConfig) *ConvConf {
	c := new(ConvConf)
	c.plotter = plotterConf
	c.runtime = runtConf
	return c
}

// Equip a given converter with configuration
func WithConfig(converter ConverterI, config *ConvConf) ConverterI {
	converter.SetConfig(config)
	return converter
}

// Convert using the converter, based on the element's type
func SVGConvert(s svg.SVGShapeElement, transformChain svgtransform.TransformChain, converter ConverterI) fun.Option[*gcode.Gcode] {
	switch s.(type) {
	case *svg.Path:
		return converter.Path(s.(*svg.Path), transformChain)
	case *svg.Line:
		return converter.Line(s.(*svg.Line), transformChain)
	case *svg.Rect:
		return converter.Rect(s.(*svg.Rect), transformChain)
	case *svg.Circle:
		return converter.Circle(s.(*svg.Circle), transformChain)
	case *svg.Ellipse:
		return converter.Ellipse(s.(*svg.Ellipse), transformChain)
	case *svg.Polygon:
		return converter.Polygon(s.(*svg.Polygon), transformChain)
	case *svg.Polyline:
		return converter.Polyline(s.(*svg.Polyline), transformChain)
	default:
		llog.Panicf("Unknown SVG object received, cannot convert to gcode. Type: %T\n", s)
		return nil
	}
}
