package conf

import "github.com/abzicht/svgocode/svgocode/math64"

var templatePrefix string = `
M106 S0 ;Turn-off fan
M104 S0 ;Turn-off hotend
M140 S0 ;Turn-off bed
G21 ; metric values (mm)
G90 ; absolute positioning
M82 ; set extruder to absolute mode
M107 ; start with the fan off
G92 E0 ; Reset Extruder
G28 ; Home all axes
G1 Z40.0 F3000 ; Move Z Axis up to prevent scratching of Heat Bed

G1 F2000 E0 ; Speed for moves, no extrusion
G0 F4000 E0 ; Speed for drawing, no extrusion
`

var templateSuffix string = `
G1 Z80.0 ;Raise Z
G1 X0 Y100 ;Present print
M84 X Y E ;Disable all steppers but Z
`

var templateYamlPrefix string = `
# SVGOCODE plotter configuration template.
# Adjust values to match the parameters of your device (3D printer, etc.).
# Specify the configuration file via --plotter-config.
`

// Return a generic PlotterConfig template for generic printers.
func PlotterConfigTemplate() *PlotterConfig {
	p := new(PlotterConfig)
	p.GcodePrefix = templatePrefix
	p.GcodeSuffix = templateSuffix
	p.Plate = Plate{
		Center: math64.VectorF2{X: 100, Y: 100},
		Min:    math64.VectorF3{X: 0, Y: 0, Z: 0},
		Max:    math64.VectorF3{X: 200, Y: 200, Z: 200}, // a small plate
	}
	p.DrawHeight = 40.0    // Safety margins: draw at height 4cm
	p.RetractHeight = 50.0 // and move at height 5cm
	p.DrawSpeed = 2000.0
	p.RetractSpeed = 4000.0
	p.RemoveComments = false
	p.MirrorX = false
	p.MirrorY = false
	p.PenOffset = math64.VectorF2{X: 0, Y: 0}
	p.yamlPrefix = templateYamlPrefix
	return p
}
