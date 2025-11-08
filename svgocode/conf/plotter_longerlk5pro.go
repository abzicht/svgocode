package conf

import "github.com/abzicht/svgocode/svgocode/math64"

var gCodePrefix string = `
M106 S0 ;Turn-off fan
M107 ; Turn-off fan
M104 S0 ;Turn-off hotend
M140 S0 ;Turn-off bed
G21 ; metric values (mm)
G90 ; absolute positioning
M82 ; set extruder to absolute mode
M107 ; start with the fan off
G92 E0 ; Reset Extruder
G28 ; Home all axes
G92 E0 ; Reset Extruder
G1 Z40.0 F3000 ; Move Z Axis up little to prevent scratching of Heat Bed

G92 E0
G1 F2000 E0
G0 F4000 E0
`

var gCodeSuffix string = `
G1 Z40.0 ;Raise Z

M140 S0
G1 Z40.0 ;Raise Z more
G1 X0 Y300 ;Present print
M84 X Y E ;Disable all steppers but Z
`

var longerLK5ProYamlPrefix string = `
# SVGOCODE plotter configuration template.
# Adjust values to match the parameters of your device (3D printer, etc.).
# Specify the configuration file via --plotter-config.
`

// Return the default/recommended PlotterConfig for the 3D printer LONGER LK5 PRO
func PlotterConfigLongerLK5ProDefault() *PlotterConfig {
	p := new(PlotterConfig)
	p.GcodePrefix = gCodePrefix
	p.GcodeSuffix = gCodeSuffix
	p.UnitLength = math64.UnitMM
	p.Plate = Plate{
		Center: math64.VectorF2{X: 150, Y: 150},
		Min:    math64.VectorF3{X: 0, Y: 0, Z: 0},
		Max:    math64.VectorF3{X: 300, Y: 300, Z: 200},
	}
	p.DrawHeight = 20.0
	p.RetractHeight = 23.0
	p.DrawSpeed = 2000.0
	p.RetractSpeed = 4000.0
	p.RemoveComments = false
	p.MirrorX = false
	p.MirrorY = true
	p.PenOffset = math64.VectorF2{X: 47, Y: 30}
	p.yamlPrefix = longerLK5ProYamlPrefix
	return p
}
