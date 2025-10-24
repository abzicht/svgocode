package svgtransform

import (
	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/math64"
)

type Mirror struct {
	X      bool            // Mirror the X axis
	Y      bool            // Mirror the Y axis
	Center math64.VectorF2 // The point that other points are being mirrored around.
}

func NewMirror(x, y bool, center math64.VectorF2) *Mirror {
	m := new(Mirror)
	m.X = x
	m.Y = y
	m.Center = center
	return m
}

func (m *Mirror) ToMatrix() *TransformMatrix {
	//TODO
	llog.Warn("Matrix for Mirror transform not yet implemented")
	return NewTransformMatrix(math64.MatrixF4Identity())
}

func (m *Mirror) apply(p math64.VectorF2) math64.VectorF2 {
	// p' = c + (c - p)
	if m.X {
		p.X = m.Center.X + (m.Center.X - p.X)
	}
	if m.Y {
		p.Y = m.Center.Y + (m.Center.Y - p.Y)
	}
	return p
}
func (m *Mirror) Apply(p ...math64.VectorF2) []math64.VectorF2 {
	for i, _ := range p {
		p[i] = m.apply(p[i])
	}
	return p
}
