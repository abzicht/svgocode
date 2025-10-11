package svgocode

import "github.com/jessevdk/go-flags"

type Flags struct {
	PlotterConfigFile string `short:"p" long:"plotter-config" description:"YAML-encoded config file for the plotter that is to be used."`
}

func ParseFlags(f *Flags) error {
	parser := flags.NewParser(f, flags.Default)
	var description string = "SVGOCODE is a(nother) tool for converting SVG files to Gcode"

	parser.LongDescription = description
	_, err := parser.Parse()
	return err
}
