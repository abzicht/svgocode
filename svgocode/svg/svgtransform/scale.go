package svgtransform

import "github.com/abzicht/svgocode/svgocode/math64"

type Scale struct {
	Factor math64.VectorF2 // X/Y axis scale
}

func NewScale(factor math64.VectorF2) *Scale {
	s := new(Scale)
	s.Factor = factor
	return s
}

func (s *Scale) apply(p math64.VectorF2) math64.VectorF2 {
	return math64.VectorF2{X: p.X * s.Factor.X, Y: p.Y * s.Factor.Y}
}

func (s *Scale) Apply(p ...math64.VectorF2) []math64.VectorF2 {
	for i, _ := range p {
		p[i] = s.ToMatrix().ApplyP(p[i])
	}
	return p
}

func (s *Scale) ToMatrix() *TransformMatrix {
	m := math64.MatrixF4Identity()
	m[0] = s.Factor.X
	m[5] = s.Factor.Y
	return NewTransformMatrix(m)
}
