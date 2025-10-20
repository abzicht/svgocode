package svg

import (
	"github.com/abzicht/svgocode/svgocode/svg/svgtransform"
)

type SVGCoreAttributes struct {
	Id    string `xml:"id,attr"`
	Class string `xml:"class,attr"`
	Style string `xml:"style,attr"`
}

type SVGPresentationTransform struct {
	TransformStr string `xml:"transform,attr"`
}

func (spt *SVGPresentationTransform) Transform() svgtransform.TransformChain {
	return svgtransform.ParseTransform(spt.TransformStr)
}
