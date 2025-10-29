package svgocode

import (
	"github.com/abzicht/gogenericfunc/fun"
	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/convs"
	"github.com/abzicht/svgocode/svgocode/gcode"
	"github.com/abzicht/svgocode/svgocode/math64"
	"github.com/abzicht/svgocode/svgocode/ordering"
	"github.com/abzicht/svgocode/svgocode/plotter"
	"github.com/abzicht/svgocode/svgocode/svg"
)

// Convert an SVG object to gcode instructions
func Svg2Gcode(s *svg.SVG, plotterConf *plotter.PlotterConfig, conv convs.ConverterI, order ordering.OrderingI) *gcode.Gcode {
	runtConf := plotter.NewRuntimeConfig()
	runtConf.SetUnitType(s.UserUnit())
	plotterTransform := plotterConf.Transform()

	var gcodes []*gcode.Gcode
	for svgElementPath := range svg.PathSeq(s, true) {
		if len(svgElementPath) == 0 {
			continue
		}
		svgElement := svgElementPath[len(svgElementPath)-1]
		// First, convert the svg objects to individual gcode segments using the
		// provided converter.
		if svg.IsLeaf(svgElement) {
			transformChain := append(plotterTransform, svg.TransformChainForPath(svgElementPath)...)
			gcodeOpt := convs.SVGConvert(svgElement, transformChain, conv)
			switch gcodeOpt.(type) {
			case fun.Some[*gcode.Gcode]:
				gcodes = append(gcodes, gcodeOpt.GetValue())
			case fun.None[*gcode.Gcode]:
			default:
				llog.Panicf("Unknown option type: %T\n", gcodeOpt)
			}
		} else {
			switch svgElement.(type) {
			case *svg.Text:
				llog.Warn("Text elements are not supported and will not be added to gcode. Please convert text to paths.\n")
			}
		}
	}
	// Then, order the gcode segments, e.g., such that travel distance is
	// minimized (depends on the given ordering method)

	if llog.GetLevel() >= llog.LDebug {
		//Only call this function, if we even want to print this info
		llog.Debugf("Non-drawing travel distance before ordering: %.0f%s\n", gcode.TotalDistanceInBetween(gcodes), runtConf.UnitType)
	}
	gcodes = order.Order(gcodes)
	if llog.GetLevel() >= llog.LDebug {
		//Only call this function, if we even want to print this info
		llog.Debugf("Non-drawing travel distance after ordering: %.0f%s\n", gcode.TotalDistanceInBetween(gcodes), runtConf.UnitType)
	}

	// Finally, add prefix and suffix
	if len(gcodes) < 1 {
		llog.Warn("No gcode produced\n")
		return gcode.NewGcode()
	}
	// Join all instructions
	gcode_joined := gcode.Join(gcodes, plotterConf)

	// Add statistics, prefix, and suffix
	gcode_full := GcodeAddStatistics(gcode.Join([]*gcode.Gcode{
		NewGcodePrefix(plotterConf, runtConf, gcode_joined),
		gcode_joined,
		NewGcodeSuffix(plotterConf, gcode_joined),
	}, plotterConf), plotterConf)
	// Remove comments, if they are not desired
	if plotterConf.RemoveComments {
		gcode_full.Code = gcode_full.Code.RemoveComments()
	}
	return gcode_full
}

// Create gcode for the plotter's gcode prefix
func NewGcodePrefix(plotterConf *plotter.PlotterConfig, runtConf *plotter.RuntimeConfig, body *gcode.Gcode) *gcode.Gcode {
	ins := gcode.NewIns(plotterConf)
	g := body.CopyMeta()
	g.AppendCode(plotterConf.GcodePrefix)
	ins.AddComment(g, "--- SVGOCODE START ---")
	ins.SetUnit(g, runtConf.UnitType)
	ins.SetExtrusion(g, 0, true)
	ins.SetExtrusion(g, 0, false)
	ins.SetSpeed(g, plotterConf.RetractSpeed, false)
	ins.SetSpeed(g, plotterConf.DrawSpeed, true)
	target := math64.VectorF3{X: body.StartCoord.X, Y: body.StartCoord.Y, Z: plotterConf.RetractHeight}
	ins.AddComment(g, "Moving to start position of first segment")
	ins.Move(g, target, plotterConf.RetractSpeed)
	g.StartCoord = target
	g.EndCoord = target
	g.BoundsMin = target
	g.BoundsMax = target
	return g
}

// Create gcode for the plotter's gcode suffix, based on the given gcode body
func NewGcodeSuffix(plotterConf *plotter.PlotterConfig, body *gcode.Gcode) *gcode.Gcode {
	ins := gcode.NewIns(plotterConf)
	g := gcode.NewGcode()
	g.StartCoord = body.EndCoord
	g.EndCoord = body.EndCoord
	g.BoundsMin = body.EndCoord
	g.BoundsMax = body.EndCoord
	ins.AddComment(g, "SVGOCODE finished, retracting")
	ins.Retract(g)
	ins.AddComment(g, "--- SVGOCODE END ---")
	g.AppendCode(plotterConf.GcodeSuffix)
	return g
}
