package plotter

import "github.com/abzicht/svgocode/svgocode/math64"

var gCodePrefix string = `
;FLAVOR:Marlin
;TIME:1375
;Filament used: 1.01345m
;Layer height: 0.2
;MINX:111.389
;MINY:111.557
;MINZ:0.2
;MAXX:188.611
;MAXY:188.443
;MAXZ:1
;TARGET_MACHINE.NAME:LONGER LK5 Pro
;Generated with Cura_SteamEngine 5.10.0
M106 S0 ;Turn-off fan
M107 ; Turn-off fan
M104 S0 ;Turn-off hotend
M140 S0 ;Turn-off bed
; LONGER Start G-code
G21 ; metric values (mm)
G90 ; absolute positioning
M82 ; set extruder to absolute mode
M107 ; start with the fan off
G92 E0 ; Reset Extruder
G28 ; Home all axes
G1 Z40.0 F3000 ; Move Z Axis up little to prevent scratching of Heat Bed
G92 E0 ; Reset Extruder

M82 ;absolute extrusion mode
G92 E0
G1 F2700
;LAYER_COUNT:5
;LAYER:0
;MESH:Untitled.stl
;TYPE:WALL-INNER
G1 F2700 E0
`

var gCodeSuffix string = `
G1 Z40.0 ;Raise Z

M140 S0
; LONGER End G-code
G91 ;Relative positioning
G1 X5 Y5 F3000 ;Wipe out
G1 Z40.0 ;Raise Z more
G90 ;Absolute positioning
G1 X0 Y300 ;Present print
M106 S0 ;Turn-off fan
M104 S0 ;Turn-off hotend
M140 S0 ;Turn-off bed
M84 X Y E ;Disable all steppers but Z

M82 ;absolute extrusion mode
M104 S0
`

// Return the default/recommended PlotterConfig for the 3D printer LONGER LK5 PRO
func PlotterConfigLongerLK5ProDefault() *PlotterConfig {
	p := new(PlotterConfig)
	p.GcodePrefix = gCodePrefix
	p.GcodeSuffix = gCodeSuffix
	p.Plate = Plate{
		Center: math64.VectorF2{X: 150, Y: 150},
		Min:    math64.VectorF3{X: 0, Y: 0, Z: 0},
		Max:    math64.VectorF3{X: 300, Y: 300, Z: 200},
	}
	p.DrawHeight = 2.0
	p.RetractHeight = 4.0
	p.DrawSpeed = 70.0
	p.RetractSpeed = 100.0
	return p
}
