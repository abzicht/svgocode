package svgocode

import (
	"github.com/abzicht/svgocode/llog"
	"github.com/jessevdk/go-flags"
)

type Flags struct {
	Verbosity             int    `short:"v" long:"verbosity" description:"Verbosity (fatal: 0, error: 1, warn: 2, info: 3, debug: 4)." default:"3"`
	SvgFile               string `short:"s" long:"svg" description:"SVG file to read from (in place of STDIN)"`
	GcodeFile             string `short:"g" long:"gcode" description:"File that GCODE will be written to (in place of STDOUT)"`
	PlotterConfigFile     string `short:"p" long:"plotter-config" description:"YAML-encoded config file for the plotter that is to be used."`
	PlotterConfigTemplate bool   `long:"plotter-config-template" description:"Print an exemplary plotter configuration file in YAML-encoding (cf. flag --plotter-config)"`
	Ordering              string `short:"o" long:"ordering-algoritm" description:"Algorithm for finding a GCODE segment order. Available algorithms: '2opt' (perfect result; for small input), 'greedy' (not perfect; for large input), 'reverse' (draw last elements first), and 'none' (skip ordering)." default:"2opt"`
}

func ParseFlags(f *Flags) error {
	parser := flags.NewParser(f, flags.Default)
	var description string = `SVGOCODE: A(nother) tool for converting SVG files to GCODE.
Copyright (C) 2025 Abzicht <https://github.com/abzicht>.
This program comes with ABSOLUTELY NO WARRANTY; This is free software, and you are welcome to redistribute it under the GNU GPLv3 license (see <https://www.gnu.org/licenses/>).`

	parser.LongDescription = description
	_, err := parser.Parse()
	llog.SetLevel(llog.LogLevel(f.Verbosity))
	return err
}
