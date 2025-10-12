package svgocode

import (
	"github.com/abzicht/gogenericfunc/fun"
	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/convs"
	"github.com/abzicht/svgocode/svgocode/gcode"
	"github.com/abzicht/svgocode/svgocode/ordering"
	"github.com/abzicht/svgocode/svgocode/plotter"
	"github.com/abzicht/svgocode/svgocode/svg"
)

// Convert an SVG object to gcode instructions
func Svg2Gcode(s *svg.SVG, plotterConf *plotter.PlotterConfig, conv convs.ConverterI, order ordering.OrderingI) *gcode.Gcode {
	var gcodes []*gcode.Gcode
	s.UserUnit()
	for svgElement := range svg.Seq(s) {
		// First, convert the svg objects to individual gcode segments using the
		// provided converter.
		if svg.IsLeaf(svgElement) {
			gcodeOpt := convs.SVGConvert(svgElement, conv)
			switch gcodeOpt.(type) {
			case fun.Some[*gcode.Gcode]:
				gcodes = append(gcodes, gcodeOpt.GetValue())
			case fun.None[*gcode.Gcode]:
			default:
				llog.Panicf("Unknown option type: %T\n", gcodeOpt)
			}
		}
	}
	// Then, order the gcode segments, e.g., such that travel distance is
	// minimized (depends on the given ordering method)
	gcodes = order.Order(gcodes)

	// Finally, add prefix and suffix
	if len(gcodes) < 1 {
		llog.Warn("No gcode produced\n")
		return gcode.NewGcode()
	}
	gcodes = append([]*gcode.Gcode{gcode.NewGcodePrefix(plotterConf, gcodes[0])}, gcodes...)
	gcodes = append(gcodes, gcode.NewGcodeSuffix(plotterConf, gcodes[len(gcodes)-1]))

	// And join them together
	return gcode.Join(gcodes, plotterConf)
}
