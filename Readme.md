# SVGOCODE - Yet Another SVG to GCODE Converter

Does exactly that. Parses SVG as via XML decoder and tries to figure out, what
the image should look like in GCODE such that, e.g. penplotters can draw a
picture.


# Next Steps

* Add option to switch between penplotting and extrusion.
* Implement unit tests.
* Create multithreaded iter.Seq for converting SVG files.
* Create ordering method based on start- and end-coordinates.
