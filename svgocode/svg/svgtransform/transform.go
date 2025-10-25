package svgtransform

import (
	"github.com/abzicht/svgocode/svgocode/math64"
)

/* Transform parser & handler. Tries to follow
https://www.w3.org/TR/css-transforms-1/.
Offers functionality for summarizing "chains" of transform operations in a
single 4x4 matrix.
*/

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
	TransformCmdMirror     = TransformCommandType("mirror") // Not a standard command
)

type Transform interface {
	Apply(p ...math64.VectorF2) []math64.VectorF2
	ToMatrix() *TransformMatrix
}

type TransformChain []Transform

func (tc TransformChain) Apply(p ...math64.VectorF2) []math64.VectorF2 {
	for _, t := range tc {
		p = t.Apply(p...)
	}
	return p
}

func (tc TransformChain) ToMatrix() *TransformMatrix {
	tMat := NewTransformMatrix(math64.MatrixF4Identity())
	for i, _ := range tc {
		tMat = tMat.Product(tc[i].ToMatrix())
	}
	return tMat
}
