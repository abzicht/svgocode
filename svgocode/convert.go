package svgocode

import (
	"github.com/abzicht/gogenericfunc/fun"
	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/conf"
	"github.com/abzicht/svgocode/svgocode/conv"
	"github.com/abzicht/svgocode/svgocode/gcode"
	"github.com/abzicht/svgocode/svgocode/math64"
	"github.com/abzicht/svgocode/svgocode/ordering"
	"github.com/abzicht/svgocode/svgocode/svg"
)

// Convert an SVG object to GCODE instructions
func Svg2Gcode(s *svg.SVG, plotterConf *conf.PlotterConfig, converter conv.ConverterI, order ordering.OrderingI) *gcode.Gcode {
	svgUnit := s.Unit()
	runtConf := conf.NewRuntimeConfig(plotterConf, plotterConf.UnitLength, svgUnit)
	converter = conv.WithConfig(converter, conv.NewConvConf(runtConf))

	plotterTransform := runtConf.Plotter.Transform(svgUnit)

	var gcodes []*gcode.Gcode
	for svgElementPath := range svg.PathSeq(s) {
		if len(svgElementPath) == 0 {
			continue
		}
		svgElement := svgElementPath[len(svgElementPath)-1]
		// Convert the SVG objects to individual GCODE segments using the
		// provided converter.
		if svg.IsLeaf(svgElement) {
			transformChain := append(plotterTransform, svg.TransformChainForPath(svgElementPath)...)
			gcodeOpt := conv.SVGConvert(svgElement, transformChain, converter)
			switch gcodeOpt.(type) {
			case fun.Some[*gcode.Gcode]:
				gcodes = append(gcodes, gcodeOpt.GetValue())
				if gcodeOpt.GetValue().BoundsMin.Equal(math64.VectorF3{X: 0, Y: 0, Z: 20.0}) {
					llog.Panic(svgElement.ID())
				}
			case fun.None[*gcode.Gcode]:
			default:
				llog.Panicf("Unknown option type: %T\n", gcodeOpt)
			}
		} else {
			switch svgElement.(type) {
			case *svg.Text:
				llog.Warn("Text elements are not supported and will not be added to GCODE. Please convert text to paths.\n")
			}
		}
	}

	if len(gcodes) < 1 {
		llog.Warn("No GCODE produced\n")
		return gcode.NewGcode()
	}
	if llog.GetLevel() >= llog.LDebug {
		//Only call this function, if we even want to print this info
		llog.Debugf("Non-drawing travel distance before ordering: %.0f%s\n", gcode.TotalDistanceInBetween(gcodes), runtConf.PlotterUnit)
	}
	// Order the gcode segments, e.g., such that travel distance is
	// minimized (depends on the given ordering method)
	gcodes = order.Order(gcodes)
	if llog.GetLevel() >= llog.LDebug {
		totalTravelDist := gcode.TotalDistanceInBetween(gcodes)
		llog.Debugf("Non-drawing travel distance after ordering: %.0f%s\n", totalTravelDist, runtConf.PlotterUnit)
	}

	// Join all instructions
	gcode_joined := gcode.Join(gcodes, runtConf)

	// Add statistics, prefix, and suffix
	gcode_full := GcodeAddSummary(gcode.Join([]*gcode.Gcode{
		NewGcodePrefix(runtConf, gcode_joined),
		gcode_joined,
		NewGcodeSuffix(runtConf, gcode_joined),
	}, runtConf), runtConf)

	WarnBoundariesConditional(runtConf, gcode_full)
	// Remove comments, if they are not desired
	if runtConf.Plotter.RemoveComments {
		gcode_full.Code = gcode_full.Code.RemoveComments()
	}
	return gcode_full
}

func WarnBoundariesConditional(runtConf *conf.RuntimeConfig, g *gcode.Gcode) {
	if !g.BoundsMin.Min(runtConf.Plotter.Plate.Min).Equal(runtConf.Plotter.Plate.Min) {
		llog.Warnf("(Parts of) GCODE lies outside of plotter dimensions. Minimum GCODE position: %s. Minimum plotter coordinates: %s.\n", g.BoundsMin.String(), runtConf.Plotter.Plate.Min.String())
	}
	if !g.BoundsMax.Max(runtConf.Plotter.Plate.Max).Equal(runtConf.Plotter.Plate.Max) {
		llog.Warnf("(Parts of) GCODE lies outside of plotter dimensions. Maximum GCODE position: %s. Maximum plotter coordinates: %s.\n", g.BoundsMax.String(), runtConf.Plotter.Plate.Max.String())
	}
}

// Create gcode for the plotter's gcode prefix
func NewGcodePrefix(runtConf *conf.RuntimeConfig, body *gcode.Gcode) *gcode.Gcode {
	ins := gcode.NewIns(runtConf)
	g := body.CopyMeta()
	g.AppendCode(runtConf.Plotter.GcodePrefix)
	ins.AddComment(g, "--- SVGOCODE START ---")
	ins.SetUnit(g, runtConf.PlotterUnit)
	ins.SetExtrusion(g, 0, true)
	ins.SetExtrusion(g, 0, false)
	ins.SetSpeed(g, runtConf.Plotter.RetractSpeed, false)
	ins.SetSpeed(g, runtConf.Plotter.DrawSpeed, true)
	target := math64.VectorF3{X: body.StartCoord.X, Y: body.StartCoord.Y, Z: runtConf.Plotter.RetractHeight}
	ins.AddComment(g, "Moving to start position of first segment")
	ins.Move(g, target, runtConf.Plotter.RetractSpeed)
	g.StartCoord = target
	g.EndCoord = target
	g.BoundsMin = target
	g.BoundsMax = target
	return g
}

// Create gcode for the plotter's gcode suffix, based on the given gcode body
func NewGcodeSuffix(runtConf *conf.RuntimeConfig, body *gcode.Gcode) *gcode.Gcode {
	ins := gcode.NewIns(runtConf)
	g := gcode.NewGcode()
	g.StartCoord = body.EndCoord
	g.EndCoord = body.EndCoord
	g.BoundsMin = body.EndCoord
	g.BoundsMax = body.EndCoord
	ins.AddComment(g, "SVGOCODE finished, retracting")
	ins.Retract(g)
	ins.AddComment(g, "--- SVGOCODE END ---")
	g.AppendCode(runtConf.Plotter.GcodeSuffix)
	return g
}
