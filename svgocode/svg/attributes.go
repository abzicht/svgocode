package svg

import (
	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/svg/svgtransform"
)

type SVGCoreAttributes struct {
	Id    string `xml:"id,attr"`
	Class string `xml:"class,attr"`
	Style string `xml:"style,attr"`
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
