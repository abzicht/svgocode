package math64

import "math"

type AngRad Float // Radians
type AngDeg Float // Degree

/**
* FAQ
* Q: Why so many types that are, essentially, float64?
* A: So that nobody attempts to, e.g., math.Cos(degrees)
 */

func (r AngRad) Deg() AngDeg {
	return AngDeg(r * 180 / math.Pi)
}

func (d AngDeg) Rad() AngRad {
	return AngRad(d * math.Pi / 180)
}

func (r AngRad) Sin() Float {
	return Float(math.Sin(float64(r)))
}
func (r AngRad) Cos() Float {
	return Float(math.Cos(float64(r)))
}
func (r AngRad) Tan() Float {
	return Float(math.Tan(float64(r)))
}

func (v VectorF2) DistEuclid(v2 VectorF2) Float {
	x := v.X - v2.X
	y := v.Y - v2.Y
	return Float(math.Hypot(float64(x), float64(y)))
}

func (v VectorF2) DistManhattan(v2 VectorF2) Float {
	return Float(math.Abs(float64(v.X-v2.X)) + math.Abs(float64(v.Y-v2.Y)))
}

func (v VectorF3) DistEuclid(v2 VectorF3) Float {
	x := float64(v.X - v2.X)
	y := float64(v.Y - v2.Y)
	z := float64(v.Z - v2.Z)
	return Float(math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2) + math.Pow(z, 2)))
}

func (v VectorF3) DistManhattan(v2 VectorF3) Float {
	return Float(math.Abs(float64(v.X-v2.X)) + math.Abs(float64(v.Y-v2.Y)) + math.Abs(float64(v.Z-v2.Z)))
}
