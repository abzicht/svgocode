package svg

import (
	"fmt"
	"regexp"

	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/math64"
	"github.com/abzicht/svgocode/svgocode/svg/svgtransform"
)

type SvgId string

type SVGCoreAttributes struct {
	Id    SvgId  `xml:"id,attr"`
	Class string `xml:"class,attr"`
	Style string `xml:"style,attr"`
}

func (s SVGCoreAttributes) ID() SvgId {
	return s.Id
}

type SVGLinkAttributes struct {
	Href  string `xml:"href,attr"`
	XHref string `xml:"xlink:href,attr"`
}

// Return href, xlink:href, or an empty string.
func (s SVGLinkAttributes) GetHref() string {
	if len(s.Href) > 0 {
		return s.Href
	}
	return s.XHref
}

type SVGPresentationLocation struct {
	X math64.Float `xml:"x,attr"`
	Y math64.Float `xml:"y,attr"`
}

type SVGPresentation interface {
	Transform() svgtransform.TransformChain
	AppendTransform(function string, inFront bool)
}

type SVGPresentationTransform struct {
	TransformStr    string `xml:"transform,attr"`
	TransformOrigin string `xml:"transform-origin,attr"`
	TransformBox    string `xml:"transform-box,attr"`
}

func (spt *SVGPresentationTransform) Transform() svgtransform.TransformChain {
	if len(spt.TransformBox) != 0 || len(spt.TransformOrigin) != 0 {
		llog.Panic("Transformation based on transform-origin / transform-box is not (yet) implemented. Try to make the SVG compatible, e.g., by ungrouping elements")
	}
	return svgtransform.ParseTransform(spt.TransformStr)
}

var transformMatch *regexp.Regexp = regexp.MustCompile(`(?i)\b(matrix|translate|scale|rotate|skew|skewX|skewY)\s*\(`)

func (spt *SVGPresentationTransform) AppendTransform(function string, inFront bool) {
	if transformMatch.MatchString(spt.TransformStr) {
		if inFront {
			spt.TransformStr = fmt.Sprintf("%s,%s", function, spt.TransformStr)
		} else {
			spt.TransformStr = fmt.Sprintf("%s,%s", spt.TransformStr, function)
		}
		return
	}
	// No transform functions present, add the first one.
	spt.TransformStr = function
}
