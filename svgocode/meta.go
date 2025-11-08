package svgocode

import (
	"fmt"

	"github.com/abzicht/svgocode/svgocode/conf"
	"github.com/abzicht/svgocode/svgocode/gcode"
)

// Create metadata for gcode output

func GcodeAddSummary(g *gcode.Gcode, runtConf *conf.RuntimeConfig) *gcode.Gcode {
	if !runtConf.Plotter.RemoveComments {
		gmeta := g.CopyMeta()
		ins := gcode.NewIns(runtConf)
		ins.AddComment(gmeta, "SVGOCODE Summary")
		ins.AddComment(gmeta, fmt.Sprintf("Unit: %s", runtConf.PlotterUnit))
		ins.AddComment(gmeta, fmt.Sprintf("Coordinates (min): %s", g.BoundsMin.String()))
		ins.AddComment(gmeta, fmt.Sprintf("Coordinates (max): %s", g.BoundsMax.String()))
		ins.AddComment(gmeta, fmt.Sprintf("Number of instructions: %d", g.Code.NumInstructions()))
		gmeta.Code.Append(g.Code)
		return gmeta
	}
	return g
}
