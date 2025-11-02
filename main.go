package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode"
	"github.com/abzicht/svgocode/svgocode/conf"
	"github.com/abzicht/svgocode/svgocode/conv"
	"github.com/abzicht/svgocode/svgocode/gcode"
	"github.com/abzicht/svgocode/svgocode/ordering"
	"github.com/abzicht/svgocode/svgocode/svg"
	"github.com/jessevdk/go-flags"
)

func main() {
	// 1. parse flags
	// 2. read svg
	// 3. convert to gcode
	// 4. write gcode
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
		// User wants to see a template for plotter profiles
		fmt.Println(conf.PlotterConfigLongerLK5ProDefault().YAML(4))
		return
	}
	var plotterConfig *conf.PlotterConfig
	if len(f.PlotterConfigFile) > 0 {
		// Read profile from file
		file_, err := os.Open(f.PlotterConfigFile)
		if err != nil {
			log.Fatal(err.Error())
		}
		plotterConfig, err = conf.InitPlotterConfig(file_)
		if err != nil {
			log.Fatal(err.Error())
		}
		err = file_.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		// or assume default profile.
		llog.Warn("No plotter configuration specified. Using default configuration for LONGER LK5 PRO 3D printer.\n")
		plotterConfig = conf.PlotterConfigLongerLK5ProDefault()
	}
	var reader io.Reader = os.Stdin
	if len(f.SvgFile) > 0 {
		// Read from file (instead of STDIN)
		if _, err := os.Stat(f.SvgFile); errors.Is(err, os.ErrNotExist) {
			llog.Panicf("File %s does not exist", f.SvgFile)
		}
		fi, err := os.Open(f.SvgFile)
		defer func() {
			if err := fi.Close(); err != nil {
				llog.Panicf("Failed to close file %s: %s", f.SvgFile, err.Error())
			}
		}()
		if err != nil {
			llog.Panicf("Failed to open file %s: %s", f.SvgFile, err.Error())
		}
		reader = fi
	}

	var parsed_svg svg.SVG
	decoder := svg.NewDecoder(reader)
	// Decode the SVG
	err = decoder.Decode(&parsed_svg)
	if err != nil {
		llog.Panic(err.Error())
	}
	// Convert to *gcode.Gcode
	gcode_ := svgocode.Svg2Gcode(&parsed_svg, plotterConfig, conv.NewDirect(plotterConfig), ordering.ParseOrdering(ordering.OrderingAlg(f.Ordering)))

	var writer io.Writer = os.Stdout
	if len(f.GcodeFile) > 0 {
		// Write to file (instead of STDOUT)
		fi, err := os.Create(f.GcodeFile)
		defer func() {
			if err := fi.Close(); err != nil {
				llog.Panicf("Failed to close file %s: %s", f.GcodeFile, err.Error())
			}
		}()
		if err != nil {
			llog.Panicf("Failed to open file %s: %s", f.GcodeFile, err.Error())
		}
		writer = fi
	}
	llog.Debug(f.GcodeFile)
	// Encode gcode
	encoder := gcode.NewEncoder(writer)
	if err = encoder.Encode(gcode_); err != nil {
		llog.Panic(err.Error())
	}
	if err = encoder.Close(); err != nil {
		llog.Panic(err.Error())
	}
	// Fin
}
