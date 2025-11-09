package svg

import (
	"encoding/xml"
	"strings"

	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/forgo"
	"github.com/abzicht/svgocode/svgocode/math64"
	"github.com/abzicht/svgocode/svgocode/svg/svgtransform"
)

type SVGElement interface {
	CloneSVGElement() SVGElement
	SVGPresentation
	ID() SvgId
	Children() []SVGElement
	Root() SVGElement
	SetRoot(SVGElement)
}

type SVGShapeElement interface {
	//Position() math64.VectorF2
}

type SVGRoot struct {
	root SVGElement
}

func (s *SVGRoot) Root() SVGElement {
	return s.root
}

func (s *SVGRoot) SetRoot(root SVGElement) {
	s.root = root
}

type SVGCore struct {
	SVGRoot
}

type SVGElements struct {
	SVGShapeElements
	SVG       []*SVG      `xml:"svg"`
	Groupings []*Grouping `xml:"g"`
	ALinks    []*ALink    `xml:"a"`
	Defs      []*Defs     `xml:"defs"`
	Uses      []*Use      `xml:"use"`
	Texts     []*Text     `xml:"text"`
}

func (s *SVGElements) Clone() *SVGElements {
	s2 := new(SVGElements)
	s2.SVGShapeElements = *s.SVGShapeElements.Clone()
	s2.SVG = forgo.Clone[*SVG](s.SVG)
	s2.Groupings = forgo.Clone[*Grouping](s.Groupings)
	s2.ALinks = forgo.Clone[*ALink](s.ALinks)
	s2.Defs = forgo.Clone[*Defs](s.Defs)
	s2.Uses = forgo.Clone[*Use](s.Uses)
	s2.Texts = forgo.Clone[*Text](s.Texts)
	return s2
}

// Produce a TransformChain with all transform operations contained in a given path
func TransformChainForPath(path []SVGElement) svgtransform.TransformChain {
	chain := svgtransform.TransformChain{}
	for _, element := range path {
		chain = append(chain, element.Transform()...)
	}
	return chain
}

func (svgElem *SVGElements) Children() []SVGElement {
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
	for _, d := range svgElem.Defs {
		children = append(children, d)
	}
	for _, u := range svgElem.Uses {
		children = append(children, u)
	}
	for _, t := range svgElem.Texts {
		children = append(children, t)
	}
	children = append(children, svgElem.SVGShapeElements.Children()...)
	return children
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

func (svgElem *SVGShapeElements) Children() []SVGElement {
	var children []SVGElement
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
	return children
}

func (s *SVGShapeElements) Clone() *SVGShapeElements {
	s2 := new(SVGShapeElements)
	s2.Paths = forgo.Clone[*Path](s.Paths)
	s2.Lines = forgo.Clone[*Line](s.Lines)
	s2.Rects = forgo.Clone[*Rect](s.Rects)
	s2.Circles = forgo.Clone[*Circle](s.Circles)
	s2.Ellipses = forgo.Clone[*Ellipse](s.Ellipses)
	s2.Polygons = forgo.Clone[*Polygon](s.Polygons)
	s2.Polylines = forgo.Clone[*Polyline](s.Polylines)
	return s2
}

type SVG struct {
	XMLName xml.Name `xml:"svg"`
	X       string   `xml:"x,attr"`
	Y       string   `xml:"y,attr"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
	SVGCore
	SVGCoreAttributes
	SVGPresentationTransform
	SVGElements
}

func (s *SVG) Clone() *SVG {
	s2 := new(SVG)
	s2.XMLName = s.XMLName
	s2.X = s.X
	s2.Y = s.Y
	s2.Width = s.Width
	s2.Height = s.Height
	s2.SVGCoreAttributes = s.SVGCoreAttributes
	s2.SVGPresentationTransform = s.SVGPresentationTransform
	s2.SVGElements = *s.SVGElements.Clone()
	return s2
}

func (s *SVG) CloneSVGElement() SVGElement {
	return s.Clone()
}

// Determine the unit defined in the SVG's attributes
func (s *SVG) Unit() (unit math64.UnitLength) {
	defer func() {
		if r := recover(); r != nil {
			llog.Warnf("Failed to determine SVG's unit type based on width/height: '%s'. Assuming millimeters. Verify produced gcode!\n", r)
			unit = math64.UnitMM
		}
	}()
	if len(s.Width) > 0 {
		_, unit = math64.NumberUnit(s.Width)
	} else if len(s.Height) > 0 {
		_, unit = math64.NumberUnit(s.Height)
	} else {
		llog.Warn("Could not determine SVG's unit type. Assuming millimeters. Verify produced gcode!\n")
		unit = math64.UnitMM
	}
	return
}

type Tspan struct {
	SVGCore
	SVGCoreAttributes
	SVGPresentationTransform
	X           math64.Float `xml:"x,attr"`
	Y           math64.Float `xml:"y,attr"`
	InnerString string       `xml:"innerhtml"`
}

func (t *Tspan) Clone() *Tspan {
	t2 := new(Tspan)
	t2.X = t.X
	t2.Y = t.Y
	t2.InnerString = t.InnerString
	t2.SVGCoreAttributes = t.SVGCoreAttributes
	t2.SVGPresentationTransform = t.SVGPresentationTransform
	return t2
}

func (t *Tspan) CloneSVGElement() SVGElement {
	return t.Clone()
}

func (t *Tspan) Children() []SVGElement {
	return []SVGElement{}
}

type Text struct {
	SVGCore
	SVGCoreAttributes
	SVGPresentationTransform
	Tspan *Tspan       `xml:"tspan,attr"`
	X     math64.Float `xml:"x,attr"`
	Y     math64.Float `xml:"y,attr"`
}

func (t *Text) Clone() *Text {
	t2 := new(Text)
	t2.X = t.X
	t2.Y = t.Y
	t2.Tspan = t.Tspan.Clone()
	t2.SVGCoreAttributes = t.SVGCoreAttributes
	t2.SVGPresentationTransform = t.SVGPresentationTransform
	return t2
}

func (t *Text) CloneSVGElement() SVGElement {
	return t.Clone()
}

func (t *Text) Children() []SVGElement {
	if nil == t.Tspan {
		return []SVGElement{}
	}
	return []SVGElement{t.Tspan}
}

type Defs struct {
	SVGCore
	SVGCoreAttributes
	SVGPresentationTransform
	SVGElements
}

func (d *Defs) Clone() *Defs {
	d2 := new(Defs)
	d2.SVGCoreAttributes = d.SVGCoreAttributes
	d2.SVGPresentationTransform = d.SVGPresentationTransform
	d2.SVGElements = *d.SVGElements.Clone()
	return d2
}

func (d *Defs) CloneSVGElement() SVGElement {
	return d.Clone()
}

type Grouping struct {
	SVGCore
	SVGCoreAttributes
	SVGPresentationTransform
	SVGElements
}

func (g *Grouping) Clone() *Grouping {
	g2 := new(Grouping)
	g2.SVGCoreAttributes = g.SVGCoreAttributes
	g2.SVGPresentationTransform = g.SVGPresentationTransform
	g2.SVGElements = *g.SVGElements.Clone()
	return g2
}

func (g *Grouping) CloneSVGElement() SVGElement {
	return g.Clone()
}

type ALink struct {
	SVGCore
	SVGCoreAttributes
	SVGPresentationTransform
	SVGLinkAttributes
	SVGElements
}

func (a *ALink) Clone() *ALink {
	a2 := new(ALink)
	a2.SVGCoreAttributes = a.SVGCoreAttributes
	a2.SVGPresentationTransform = a.SVGPresentationTransform
	a2.SVGLinkAttributes = a.SVGLinkAttributes
	a2.SVGElements = *a.SVGElements.Clone()
	return a2
}

func (a *ALink) CloneSVGElement() SVGElement {
	return a.Clone()
}

type Use struct {
	SVGCore
	SVGCoreAttributes
	SVGPresentationTransform
	SVGLinkAttributes
	SVGElements
	X math64.Float `xml:"x,attr"`
	Y math64.Float `xml:"y,attr"`
}

func (u *Use) Clone() *Use {
	u2 := new(Use)
	u2.SVGCoreAttributes = u.SVGCoreAttributes
	u2.SVGPresentationTransform = u.SVGPresentationTransform
	u2.SVGLinkAttributes = u.SVGLinkAttributes
	u2.SVGElements = *u.SVGElements.Clone()
	return u2
}

func (u *Use) CloneSVGElement() SVGElement {
	return u.Clone()
}

// Return the element referenced by the use element (if its id can be found in
// the given map).
func (u *Use) GetRefElement(sMap SvgIdMap) SVGElement {
	href := strings.TrimLeft(u.GetHref(), " \n\r\t")

	if !strings.HasPrefix(href, "#") {
		llog.Panicf("Unsupported HREF value: '%s'. Value must start with '#', i.e., reference another tag.", u.GetHref())
	}
	el, ok := sMap[SvgId(href[1:])]
	if ok {
		return el
	}
	llog.Panicf("Failed to find element with id '%s' (referenced by 'use' with id '%s')", href, u.Id)
	return nil
}

type Path struct {
	SVGCore
	SVGCoreAttributes
	SVGPresentationTransform
	D string `xml:"d,attr"`
}

func (p *Path) Clone() *Path {
	p2 := new(Path)
	p2.SVGCoreAttributes = p.SVGCoreAttributes
	p2.SVGPresentationTransform = p.SVGPresentationTransform
	p2.D = p.D
	return p2
}

func (p *Path) CloneSVGElement() SVGElement {
	return p.Clone()
}

func (p *Path) Children() []SVGElement {
	return []SVGElement{}
}

type Line struct {
	SVGCore
	SVGCoreAttributes
	SVGPresentationTransform
	X1 math64.Float `xml:"x1,attr"`
	Y1 math64.Float `xml:"y1,attr"`
	X2 math64.Float `xml:"x2,attr"`
	Y2 math64.Float `xml:"y2,attr"`
}

func (l *Line) Clone() *Line {
	l2 := new(Line)
	l2.SVGCoreAttributes = l.SVGCoreAttributes
	l2.SVGPresentationTransform = l.SVGPresentationTransform
	l2.X1 = l.X1
	l2.Y1 = l.Y1
	l2.X2 = l.X2
	l2.Y2 = l.Y2
	return l2
}

func (l *Line) CloneSVGElement() SVGElement {
	return l.Clone()
}

func (s *Line) Children() []SVGElement {
	return []SVGElement{}
}

type Rect struct {
	SVGCore
	SVGCoreAttributes
	SVGPresentationTransform
	X      math64.Float `xml:"x,attr"`
	Y      math64.Float `xml:"y,attr"`
	Width  math64.Float `xml:"width,attr"`
	Height math64.Float `xml:"height,attr"`
	RX     math64.Float `xml:"rx,attr"`
	RY     math64.Float `xml:"ry,attr"`
}

func (r *Rect) Clone() *Rect {
	r2 := new(Rect)
	r2.SVGCoreAttributes = r.SVGCoreAttributes
	r2.SVGPresentationTransform = r.SVGPresentationTransform
	r2.X = r.X
	r2.Y = r.Y
	r2.Width = r.Width
	r2.Height = r.Height
	r2.RX = r.RX
	r2.RY = r.RY
	return r2
}

func (r *Rect) CloneSVGElement() SVGElement {
	return r.Clone()
}

func (r *Rect) Children() []SVGElement {
	return []SVGElement{}
}

type Circle struct {
	SVGCore
	SVGCoreAttributes
	SVGPresentationTransform
	CX math64.Float `xml:"cx,attr"`
	CY math64.Float `xml:"cy,attr"`
	R  math64.Float `xml:"r,attr"`
}

func (c *Circle) Clone() *Circle {
	c2 := new(Circle)
	c2.SVGCoreAttributes = c.SVGCoreAttributes
	c2.SVGPresentationTransform = c.SVGPresentationTransform
	c2.CX = c.CX
	c2.CY = c.CY
	c2.R = c.R
	return c2
}

func (c *Circle) CloneSVGElement() SVGElement {
	return c.Clone()
}

func (c *Circle) Children() []SVGElement {
	return []SVGElement{}
}

type Ellipse struct {
	SVGCore
	SVGCoreAttributes
	SVGPresentationTransform
	CX math64.Float `xml:"cx,attr"`
	CY math64.Float `xml:"cy,attr"`
	RX math64.Float `xml:"rx,attr"`
	RY math64.Float `xml:"ry,attr"`
}

func (e *Ellipse) Clone() *Ellipse {
	e2 := new(Ellipse)
	e2.SVGCoreAttributes = e.SVGCoreAttributes
	e2.SVGPresentationTransform = e.SVGPresentationTransform
	e2.CX = e.CX
	e2.CY = e.CY
	e2.RX = e.RX
	e2.RY = e.RY
	return e2
}

func (e *Ellipse) CloneSVGElement() SVGElement {
	return e.Clone()
}

func (e *Ellipse) Children() []SVGElement {
	return []SVGElement{}
}

type Polygon struct {
	SVGCore
	SVGCoreAttributes
	SVGPresentationTransform
	P string `xml:"points,attr"`
}

func (p *Polygon) Clone() *Polygon {
	p2 := new(Polygon)
	p2.SVGCoreAttributes = p.SVGCoreAttributes
	p2.SVGPresentationTransform = p.SVGPresentationTransform
	p2.P = p.P
	return p2
}

func (p *Polygon) CloneSVGElement() SVGElement {
	return p.Clone()
}

func (p *Polygon) Children() []SVGElement {
	return []SVGElement{}
}

func (p *Polygon) Points() []math64.VectorF2 {
	return ParsePointString(p.P)
}

type Polyline struct {
	SVGCore
	SVGCoreAttributes
	SVGPresentationTransform
	P string `xml:"points,attr"`
}

func (p *Polyline) Clone() *Polyline {
	p2 := new(Polyline)
	p2.SVGCoreAttributes = p.SVGCoreAttributes
	p2.SVGPresentationTransform = p.SVGPresentationTransform
	p2.P = p.P
	return p2
}

func (p *Polyline) CloneSVGElement() SVGElement {
	return p.Clone()
}

func (p *Polyline) Children() []SVGElement {
	return []SVGElement{}
}

func (p *Polyline) Points() []math64.VectorF2 {
	return ParsePointString(p.P)
}

// Returns true, if the element is a SVG type that can contain multiple children
// Also returns true, if element does not contain children, but could do so.
// Returns false else
func IsCollection(s SVGElement) bool {
	switch s.(type) {
	case *Grouping, *ALink, *SVG, *Defs:
		return true
	}
	return false
}

// Returns true, iff element can not contain children
func IsLeaf(s SVGElement) bool {
	switch s.(type) {
	case *Path, *Line, *Rect, *Circle, *Ellipse, *Polygon, *Polyline:
		return true
	}
	return false
}
