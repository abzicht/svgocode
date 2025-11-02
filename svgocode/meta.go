package svgocode

import (
	"fmt"

	"github.com/abzicht/svgocode/svgocode/conf"
	"github.com/abzicht/svgocode/svgocode/gcode"
)

// Create metadata for gcode output

func GcodeAddStatistics(g *gcode.Gcode, plotterConf *conf.PlotterConfig) *gcode.Gcode {
	if !plotterConf.RemoveComments {
		gmeta := g.CopyMeta()
		ins := gcode.NewIns(plotterConf)
		ins.AddComment(gmeta, "SVGOCODE statistics")
		ins.AddComment(gmeta, fmt.Sprintf("Coordinates (min): %s", g.BoundsMin.String()))
		ins.AddComment(gmeta, fmt.Sprintf("Coordinates (max): %s", g.BoundsMax.String()))
		ins.AddComment(gmeta, fmt.Sprintf("Number of instructions: %d", g.Code.NumInstructions()))
		gmeta.Code.Append(g.Code)
		return gmeta
	}
	return g
}
