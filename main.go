package main

import (
	"fmt"
	"log"
	"os"

	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode"
	"github.com/abzicht/svgocode/svgocode/convs"
	"github.com/abzicht/svgocode/svgocode/ordering"
	"github.com/abzicht/svgocode/svgocode/plotter"
	"github.com/abzicht/svgocode/svgocode/svg"
	"github.com/jessevdk/go-flags"
)

func main() {
	/* TODO:
	* if clargs is list of files, convert those
	* else, read from stdin until EOF.
	 */

	var f svgocode.Flags
	err := svgocode.ParseFlags(&f)
	switch err.(type) {
	case *flags.Error:
		switch err.(*flags.Error).Type {
		case flags.ErrHelp:
			return
		default:
			log.Fatal(err.Error())
		}
	case nil:
		break
	default:
		log.Fatal(err.Error())
	}
	if f.PlotterConfigTemplate {
		fmt.Println(plotter.PlotterConfigLongerLK5ProDefault().YAML(4))
		return
	}
	var plotterConfig *plotter.PlotterConfig
	if len(f.PlotterConfigFile) > 0 {
		file_, err := os.Open(f.PlotterConfigFile)
		if err != nil {
			log.Fatal(err.Error())
		}
		plotterConfig, err = plotter.InitPlotterConfig(file_)
		if err != nil {
			log.Fatal(err.Error())
		}
		err = file_.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		llog.Warn("No plotter configuration specified. Using default configuration for LONGER LK5 PRO 3D printer.\n")
		plotterConfig = plotter.PlotterConfigLongerLK5ProDefault()
	}
	var parsed_svg svg.SVG
	decoder := svg.NewDecoder(os.Stdin)
	err = decoder.Decode(&parsed_svg)
	if err != nil {
		llog.Panic(err.Error())
	}
	gcode := svgocode.Svg2Gcode(&parsed_svg, plotterConfig, convs.NewDirect(plotterConfig), ordering.ParseOrdering(ordering.OrderingAlg(f.Ordering)))
	fmt.Println(gcode.String())
}
