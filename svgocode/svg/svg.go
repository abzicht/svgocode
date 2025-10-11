package svg

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/abzicht/svgocode/svgocode/math64"
)

type SVG struct {
	XMLName   xml.Name    `xml:"svg"`
	Paths     []*Path     `xml:"path"`
	Lines     []*Line     `xml:"line"`
	Rects     []*Rect     `xml:"rect"`
	Circles   []*Circle   `xml:"circle"`
	Ellipses  []*Ellipse  `xml:"ellipse"`
	Polygons  []*Polygon  `xml:"polygon"`
	Polylines []*Polyline `xml:"polyline"`
}

func (s *SVG) GetGraphicsElements() []SVGGraphicsElement {
	var elements []SVGGraphicsElement
	for _, p := range s.Paths {
		elements = append(elements, p)
	}
	for _, l := range s.Lines {
		fmt.Printf("%v\n", l)
		elements = append(elements, l)
	}
	for _, r := range s.Rects {
		elements = append(elements, r)
	}
	for _, c := range s.Circles {
		elements = append(elements, c)
	}
	for _, e := range s.Ellipses {
		elements = append(elements, e)
	}
	for _, p := range s.Polygons {
		elements = append(elements, p)
	}
	for _, p := range s.Polylines {
		elements = append(elements, p)
	}
	return elements
}

type SVGGraphicsElement interface {
	//Position() math64.VectorF2
}

type Path struct {
	D string `xml:"d,attr"`
}

type Line struct {
	X1 math64.Float `xml:"x1,attr"`
	Y1 math64.Float `xml:"y1,attr"`
	X2 math64.Float `xml:"x2,attr"`
	Y2 math64.Float `xml:"y2,attr"`
}

type Polygon struct {
	Points string `xml:"points,attr"`
}

type Polyline struct {
	Points string `xml:"points,attr"`
}

type Circle struct {
	CX math64.Float `xml:"cx,attr"`
	CY math64.Float `xml:"cy,attr"`
	R  math64.Float `xml:"r,attr"`
}

type Ellipse struct {
	CX math64.Float `xml:"cx,attr"`
	CY math64.Float `xml:"cy,attr"`
	RX math64.Float `xml:"rx,attr"`
	RY math64.Float `xml:"rY,attr"`
}

type Rect struct {
	X      math64.Float `xml:"x,attr"`
	Y      math64.Float `xml:"y,attr"`
	Width  math64.Float `xml:"width,attr"`
	Height math64.Float `xml:"height,attr"`
	RX     math64.Float `xml:"rx,attr"`
	RY     math64.Float `xml:"ry,attr"`
}

type Decoder struct {
	xmlDecoder *xml.Decoder
}

func NewDecoder(r io.Reader) *Decoder {
	d := new(Decoder)
	d.xmlDecoder = xml.NewDecoder(r)
	return d
}

func (d *Decoder) Decode(svg *SVG) error {
	return d.xmlDecoder.Decode(svg)
}
