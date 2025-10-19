package svgtransform

import "github.com/abzicht/svgocode/svgocode/math64"

type TransformCommandType string

const (
	TransformCmdMatrix     = TransformCommandType("matrix")
	TransformCmdTranslate  = TransformCommandType("translate")
	TransformCmdTranslateX = TransformCommandType("translateX")
	TransformCmdTranslateY = TransformCommandType("translateY")
	TransformCmdScale      = TransformCommandType("scale")
	TransformCmdScaleX     = TransformCommandType("scaleX")
	TransformCmdScaleY     = TransformCommandType("scaleY")
	TransformCmdRotate     = TransformCommandType("rotate")
	TransformCmdSkew       = TransformCommandType("skew")
	TransformCmdSkewX      = TransformCommandType("skewX")
	TransformCmdSkewY      = TransformCommandType("skewY")
)

type Transform interface {
	TransformCommandType() TransformCommandType
	Apply(p ...math64.VectorF2) []math64.VectorF2
}

type TransformChain []Transform

func (tc TransformChain) Apply(p ...math64.VectorF2) []math64.VectorF2 {
	for _, t := range tc {
		p = t.Apply(p...)
	}
	return p
}
