# SVGOCODE - Yet Another SVG to GCODE Converter

Parses SVG and tries to figure out, what the image should look like in GCODE
such that, e.g. penplotters can draw a picture.

## Features

* Supports the following SVG elements: `svg`, `g`, `a`, `defs`, `use`, `path`,
`line`, `rect`, `circle`, `ellipse`, `polygon`, `polyline`.
  + `use` resolves to its referenced element.
  + `path` commands are fully covered.
* Supports the `transform` attribute and all of its functions (`matrix`, `translate`, `translateX`, `translateY`, `scale`, `scaleX`, `scaleY`, `skew`, `skewX`, `skewY`, `rotate`).
* Translates to GCODE based on customizable plotter / printer profiles.
* Offers algorithms for minimizing travel distance in-between draw operations.
* Defines interfaces for easily extending SVGOCODE with custom converters,
  ordering algorithms, etc.

## Installation

This will compile `build/svgocode` and install the binary in `GOBIN`:
```bash
make install
```

## Use & I/O

SVGOCODE expects SVG-formatted input from `STDIN` (or from files provided with
`-s`). It writes GCODE to `STDOUT` (or to files provided with `-g`), all
auxiliary information is written to `STDERR`.

To convert SVG files to GCODE files, use `svgocode` as follows:

```bash
# a) STDIN / STDOUT:
cat drawing.svg | svgocode > drawing.gcode
# b) via flags:
svgocode -s drawing.svg -g drawing.gcode
```

## Development

* Use `make run` to build and run `svgocode`.
* Use `make dev` to build and run `svgocode` with debugging information.
* Use `make gdb` to build and run `svgocode` inside of `gdb`.

## Configuration

SVGOCODE is configured with a YAML-encoded plotter profile that supplies information such as
the plotter's dimensions, pen offset, drawing speed, etc.

The default profile can be obtained via `svgocode --plotter-config-template`.
It provides an exemplary configuration for the Longer LK5 Pro 3D printer.

Custom profiles are supplied via `svgocode --plotter-config=file.yml`.

Here is a description of the default configuration parameters:

```yaml
# Note: all length units are in mm, speed is in mm/s
gprefix: "gcode that is placed at the start of the output"
gsuffix: "gcode that is placed at the end   of the output"
plate:
    center: # Center coordinates of the plotter's base plate
        "x": 150
        "y": 150
    min:    # Minimum value of each axis
        "x": 0
        "y": 0
        "z": 0
    max:    # Maximum value of each axis
        "x": 300
        "y": 300
        "z": 200
drawing-height: 20  # Absolute position on the Z-axis where drawing takes place
retract-height: 23  # Absolute position on the Z-axis where movement without drawing takes place
draw-speed: 2000    # Speed with that drawing is performed
retract-speed: 4000 # Speed with that movement without drawing is performed
remove-comments: false # Strip all comments (starting with ";") from GCODE.
mirror-x-axis: false # Mirror all X values around the plate's center X-axis.
mirror-y-axis: true  # Mirror all Y values around the plate's center Y-axis.
pen-offset:          # Offset with that the pen is mounted on the printer/plotter.
    "x": 47
    "y": 30
```

## Library

SVGOCODE can be easily used as library, both for parsing SVG and for GCODE
conversion. `main.go` may give you a hint on how to use the library in your own
projects:

```go
	var parsed_svg svg.SVG
	decoder := svg.NewDecoder(READER)
	err = decoder.Decode(&parsed_svg)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	gcode := svgocode.Svg2Gcode( # The all-in-one converter
            &parsed_svg,
            plotter.PlotterConfigLongerLK5ProDefault(), # Plotter profile
            convs.NewDirect(plotterConfig), # Conversion implementation
            ordering.NewGreedy() # Ordering method
        )
	fmt.Println(gcode.String())
```

## Troubleshoot & Disclaimer

So far, SVGOCODE is the creation of one person
([@abzicht](https://github.com/abzicht)). The SVG parser was written by one
person, just as the conversion methods were.

As such, SVGOCODE
can do a whole lot, but it is also quite limited. E.g., SVGOCODE cannot yet

* convert `text`/`tspan` elements,
* account for `transform-origin`/`transform-box`,
* guarantee correctness of complex `transform` hierarchies,
* work with embedded `svg` elements,
* convert `sodipodi` / Inkscape attributes,
* resolve non-local hyperlinks (`href` that do not point to elements inside the
  provided `svg` structure), or
* work with strange length units (`pt`, `px`).

It is, therefore, recommended to

* convert `text`/`tspan` elements to `path` elements manually (e.g., via
  Inkscape),
* consider `mirror-x-axis`/`mirror-y-axis` (cf.
  [Configuration](#Configuration)), if GCODE appears flipped,
* switch to faster algorithms (`svgocode --ordering-algorithm=greedy`), if
  input is large and processing takes too long,
* double-check the SVG's units (must be `mm`, `cm`, or `in`), and
* in general, have a good look at the SVG, e.g., via Inkscape's XML Editor.

Finally, it is crucial to assess the produced GCODE before use. E.g., open the
GCODE in CURA and activate the `Travels` checkbox in `Color Scheme > Line
Type`.

<b>Disclaimer: Running GCODE on expensive machines can cause expensive noises.
The authors of SVGOCODE are not responsible for any damage caused from
GCODE-based mishaps! The produced GCODE is in ASCII and can be easily assessed
before its use. Do so!</b>

## Contributing

Some SVG feature is not yet supported? Please create an _issue_
or DIY and submit a pull request.

## Next Steps

* Add option to switch between penplotting and extrusion.
* Implement unit tests.
