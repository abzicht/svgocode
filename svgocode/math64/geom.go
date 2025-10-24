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
