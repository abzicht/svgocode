package svg

import (
	"encoding/xml"

	"github.com/abzicht/gogenericfunc/fun"
	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/math64"
)

type SVGElement interface {
	//Position() math64.VectorF2
	Children() fun.Option[[]SVGElement]
}

type SVGShapeElement interface {
	//Position() math64.VectorF2
}

type SVGElements struct {
	SVGShapeElements
	SVG       []*SVG      `xml:"svg"`
	Groupings []*Grouping `xml:"g"`
	ALinks    []*ALink    `xml:"a"`
}

func (svgElem *SVGElements) Children() fun.Option[[]SVGElement] {
	var children []SVGElement
	for _, s := range svgElem.SVG {
		children = append(children, s)
	}
	for _, g := range svgElem.Groupings {
		children = append(children, g)
	}
	for _, a := range svgElem.ALinks {
		children = append(children, a)
	}
	for _, p := range svgElem.Paths {
		children = append(children, p)
	}
	for _, l := range svgElem.Lines {
		children = append(children, l)
	}
	for _, r := range svgElem.Rects {
		children = append(children, r)
	}
	for _, c := range svgElem.Circles {
		children = append(children, c)
	}
	for _, e := range svgElem.Ellipses {
		children = append(children, e)
	}
	for _, p := range svgElem.Polygons {
		children = append(children, p)
	}
	for _, p := range svgElem.Polylines {
		children = append(children, p)
	}
	return fun.NewSome[[]SVGElement](children)
}

type SVGShapeElements struct {
	Paths     []*Path     `xml:"path"`
	Lines     []*Line     `xml:"line"`
	Rects     []*Rect     `xml:"rect"`
	Circles   []*Circle   `xml:"circle"`
	Ellipses  []*Ellipse  `xml:"ellipse"`
	Polygons  []*Polygon  `xml:"polygon"`
	Polylines []*Polyline `xml:"polyline"`
}

type SVG struct {
	XMLName xml.Name `xml:"svg"`
	X       string   `xml:"x,attr"`
	Y       string   `xml:"y,attr"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
	SVGCoreAttributes
	SVGElements
}

func (s *SVG) UserUnit() (unit math64.UnitType) {
	defer func() {
		if r := recover(); r != nil {
			llog.Warnf("Failed to determine SVG's unit type based on width/height: '%s'. Assuming millimeters. Verify produced gcode!\n", r)
			unit = math64.UnitMM
		}
	}()
	if len(s.Width) > 0 {
		_, unit = math64.NumberUnit(s.Width)
		return unit
	} else if len(s.Height) > 0 {
		_, unit = math64.NumberUnit(s.Height)
		return unit
	} else {
		llog.Panic("Could not determine SVG's unit type. Assuming millimeters. Verify produced gcode!\n")
	}
	return
}

type Grouping struct {
	SVGCoreAttributes
	SVGElements
}

type ALink struct {
	Href  string `xml:"href,attr"`
	XHref string `xml:"xlink:href,attr"`
	SVGCoreAttributes
	SVGElements
}

type Path struct {
	SVGCoreAttributes
	D string `xml:"d,attr"`
}

func (s *Path) Children() fun.Option[[]SVGElement] {
	return fun.NewNone[[]SVGElement]()
}

type Line struct {
	SVGCoreAttributes
	X1 math64.Float `xml:"x1,attr"`
	Y1 math64.Float `xml:"y1,attr"`
	X2 math64.Float `xml:"x2,attr"`
	Y2 math64.Float `xml:"y2,attr"`
}

func (s *Line) Children() fun.Option[[]SVGElement] {
	return fun.NewNone[[]SVGElement]()
}

type Rect struct {
	SVGCoreAttributes
	X      math64.Float `xml:"x,attr"`
	Y      math64.Float `xml:"y,attr"`
	Width  math64.Float `xml:"width,attr"`
	Height math64.Float `xml:"height,attr"`
	RX     math64.Float `xml:"rx,attr"`
	RY     math64.Float `xml:"ry,attr"`
}

func (r *Rect) Children() fun.Option[[]SVGElement] {
	return fun.NewNone[[]SVGElement]()
}

type Circle struct {
	SVGCoreAttributes
	CX math64.Float `xml:"cx,attr"`
	CY math64.Float `xml:"cy,attr"`
	R  math64.Float `xml:"r,attr"`
}

func (c *Circle) Children() fun.Option[[]SVGElement] {
	return fun.NewNone[[]SVGElement]()
}

type Ellipse struct {
	SVGCoreAttributes
	CX math64.Float `xml:"cx,attr"`
	CY math64.Float `xml:"cy,attr"`
	RX math64.Float `xml:"rx,attr"`
	RY math64.Float `xml:"rY,attr"`
}

func (e *Ellipse) Children() fun.Option[[]SVGElement] {
	return fun.NewNone[[]SVGElement]()
}

type Polygon struct {
	SVGCoreAttributes
	Points string `xml:"points,attr"`
}

func (p *Polygon) Children() fun.Option[[]SVGElement] {
	return fun.NewNone[[]SVGElement]()
}

type Polyline struct {
	SVGCoreAttributes
	Points string `xml:"points,attr"`
}

func (p *Polyline) Children() fun.Option[[]SVGElement] {
	return fun.NewNone[[]SVGElement]()
}

// Returns true, if element is a SVG type that can contain children
// Also returns true, if element does not contain children, but could do so.
// Returns false else
func IsCollection(s SVGElement) bool {
	switch s.(type) {
	case *Grouping:
		return true
	case *ALink:
		return true
	case *SVG:
		return true
	}
	return false
}

// Returns true, iff element can not contain children
func IsLeaf(s SVGElement) bool {
	switch s.(type) {
	case *Path:
		return true
	case *Line:
		return true
	case *Rect:
		return true
	case *Circle:
		return true
	case *Ellipse:
		return true
	case *Polygon:
		return true
	case *Polyline:
		return true
	}
	return false
}
