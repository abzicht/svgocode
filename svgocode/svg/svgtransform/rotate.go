package svgtransform

import (
	"github.com/abzicht/svgocode/svgocode/math64"
)

type Rotate struct {
	AngRad math64.AngRad   // By which degree
	Point  math64.VectorF2 // Around a point
}

func NewRotate(degrees math64.AngDeg, point math64.VectorF2) *Rotate {
	r := new(Rotate)
	r.AngRad = degrees.Rad()
	r.Point = point
	return r
}

func (t *Rotate) Apply(p ...math64.VectorF2) []math64.VectorF2 {
	for i, _ := range p {
		p[i] = t.ToMatrix().ApplyP(p[i])
	}
	return p
}

func (r *Rotate) ToMatrix() *TransformMatrix {
	c := r.AngRad.Cos()
	s := r.AngRad.Sin()
	m := math64.MatrixF4Identity()
	m[0] = c
	m[4] = s
	m[1] = -s
	m[5] = c
	m[2] = r.Point.X - c*r.Point.X + s*r.Point.Y
	m[6] = r.Point.Y - c*r.Point.Y - s*r.Point.X
	return NewTransformMatrix(m)
}
