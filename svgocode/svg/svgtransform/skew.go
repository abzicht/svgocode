package svgtransform

import (
	"github.com/abzicht/svgocode/svgocode/math64"
)

type Skew struct {
	Factor math64.VectorT2[math64.AngRad]
}

func NewSkew(factor math64.VectorT2[math64.AngDeg]) *Skew {
	s := new(Skew)
	s.Factor.X = factor.X.Rad()
	s.Factor.Y = factor.Y.Rad()
	return s
}

func (s *Skew) Apply(p ...math64.VectorF2) []math64.VectorF2 {
	for i, _ := range p {
		p[i] = s.ToMatrix().ApplyP(p[i])
	}
	return p
}

func (s *Skew) ToMatrix() *TransformMatrix {
	xtan := s.Factor.X.Tan()
	ytan := s.Factor.Y.Tan()
	m := math64.MatrixF4Identity()
	m[4] = ytan
	m[1] = xtan
	return NewTransformMatrix(m)
}
