package svgtransform

import "github.com/abzicht/svgocode/svgocode/math64"

type TransformMatrix struct {
	M math64.MatrixF3
}

func NewTransformMatrix(m math64.MatrixF3) *TransformMatrix {
	tMat := new(TransformMatrix)
	tMat.M = m
	return tMat
}

func (tMat *TransformMatrix) ToMatrix() *TransformMatrix {
	// Don't return tMat, who knows what the return value will be used for
	return NewTransformMatrix(tMat.M)
}

func (tMat *TransformMatrix) ApplyP(p math64.VectorF2) math64.VectorF2 {
	return tMat.M.VProductF2(p)
}

func (tMat *TransformMatrix) Apply(points ...math64.VectorF2) []math64.VectorF2 {
	points2 := make([]math64.VectorF2, len(points))
	for _, p := range points {
		points2 = append(points2, tMat.ApplyP(p))
	}
	return points2
}
