package svgtransform

import (
	"github.com/abzicht/svgocode/svgocode/math64"
)

type TransformMatrix struct {
	M math64.MatrixF4
}

func NewTransformMatrix(m math64.MatrixF4) *TransformMatrix {
	tMat := new(TransformMatrix)
	tMat.M = m
	return tMat
}

// Create a new matrix that is the product of the two given matrices
func (tMat *TransformMatrix) Product(tMat2 *TransformMatrix) *TransformMatrix {
	return NewTransformMatrix(tMat.M.MProduct(tMat2.M))
}

func (tMat *TransformMatrix) ToMatrix() *TransformMatrix {
	// Don't return tMat, who knows what the return value will be used for
	return NewTransformMatrix(tMat.M)
}

func (tMat *TransformMatrix) ApplyP(p math64.VectorF2) math64.VectorF2 {
	v4 := tMat.M.VProduct(math64.VectorF4{X: p.X, Y: p.Y, Z: 1, W: 1})
	return math64.VectorF2{X: v4.X, Y: v4.Y}
}

func (tMat *TransformMatrix) Apply(points ...math64.VectorF2) []math64.VectorF2 {
	points2 := make([]math64.VectorF2, len(points))
	for _, p := range points {
		points2 = append(points2, tMat.ApplyP(p))
	}
	return points2
}
